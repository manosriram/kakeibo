package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manosriram/kakeibo/internal/bot"
	"github.com/manosriram/kakeibo/internal/handlers"
	"github.com/manosriram/kakeibo/internal/utils"
	"github.com/manosriram/kakeibo/sqlc/db"
	"github.com/qdrant/go-client/qdrant"

	_ "embed"

	_ "modernc.org/sqlite"
)

//go:embed sqlc/schema.sql
var schema string

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func InitDB(path string) (*db.Queries, error) {
	conn, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	conn.Exec(schema)

	return db.New(conn), nil
}

func InjectDb(db *db.Queries) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	}
}

func importSqliteToCsv(d *db.Queries) error {
	statements, err := d.GetAllStatements(context.Background())
	fmt.Println(statements)
	if err != nil {
		return err
	}

	// wd, err := os.Getwd()
	// if err != nil {
	// return err
	// }
	f, err := os.OpenFile("/data/spends.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, stmt := range statements {
		_, err = fmt.Fprintf(f, "%s,%s,INR %v,%s,%s\n", stmt.Tag.String, stmt.TxnType.String, stmt.Amount.Int64, stmt.Description.String, stmt.CreatedAt.Time)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	e := echo.New()
	e.Static("/", "static")

	qdrantClient, err := qdrant.NewClient(&qdrant.Config{
		Host: os.Getenv("QDRANT_URL"),
		Port: 6334,
	})

	collectionExists, err := qdrantClient.CollectionExists(context.Background(), "kakeibo_knowledge_base")
	if err != nil {
		log.Fatalf("Error checking qdrant collection status: %s\n", err.Error())
	}
	defer qdrantClient.Close()
	if !collectionExists {
		err = qdrantClient.CreateCollection(context.Background(), &qdrant.CreateCollection{
			CollectionName: "kakeibo_knowledge_base",
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     1536, // OpenAI embedding dimension
				Distance: qdrant.Distance_Cosine,
			}),
		})
		if err != nil {
			log.Fatalf("Error creating qdrant collection: %s\n", err.Error())
		}
	}

	tmpl := template.New("").Funcs(template.FuncMap{
		"formatDate": utils.FormatDateTime,
	})

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting wd")
	}

	tmpl, err = tmpl.ParseGlob(wd + "/templates/*.html")
	if err != nil {
		panic(err)
	}

	e.Renderer = &Template{
		templates: tmpl,
	}

	q, err := InitDB("/data/kakeibo.db")
	// q, err := InitDB("./kakeibo.db")
	if err != nil {
		log.Fatalf("Error starting sqlite db")
	}

	if os.Getenv("KAKEIBO_SYNC_SQLITE_TO_CSV") == "1" {
		go importSqliteToCsv(q)
	}

	e.Use(InjectDb(q))

	go bot.StartTelegramBot(q)

	// Middleware
	e.Use(middleware.RequestLogger()) // use the default RequestLogger middleware with slog logger
	e.Use(middleware.Recover())       // recover panics as errors for proper error handling

	e.GET("/", handlers.HomeHandler)
	e.GET("/transactions", handlers.GetAllTransactionsAPI)
	e.POST("/api/transaction", handlers.CreateTransactionAPI)
	e.GET("/api/health", handlers.HealthAPI)
	e.GET("/api/insights", handlers.QueryRagAPI)

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
