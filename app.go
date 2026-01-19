package main

import (
	"database/sql"
	"errors"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manosriram/kakeibo/internal/handlers"
	"github.com/manosriram/kakeibo/internal/utils"
	"github.com/manosriram/kakeibo/sqlc/db"

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

func main() {
	e := echo.New()

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

	q, err := InitDB(wd + "/kakeibo.db")
	if err != nil {
		log.Fatalf("Error starting sqlite db")
	}

	e.Use(InjectDb(q))

	// go bot.StartTelegramBot(q)

	// Middleware
	e.Use(middleware.RequestLogger()) // use the default RequestLogger middleware with slog logger
	e.Use(middleware.Recover())       // recover panics as errors for proper error handling

	e.GET("/", handlers.HomeHandler)
	e.GET("/transactions", handlers.GetAllTransactionsAPI)
	e.POST("/api/transaction", handlers.CreateTransactionAPI)

	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}
