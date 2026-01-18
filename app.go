package main

import (
	"database/sql"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/manosriram/kakeibo/internal/handlers"
	"github.com/manosriram/kakeibo/sqlc/db"

	// db "github.com/manosriram/kakeibo/sqlc/db"

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

func initDB(path string) (*db.Queries, error) {
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

func getDaySuffix(day int) string {
	if day >= 11 && day <= 13 {
		return "th"
	}
	switch day % 10 {
	case 1:
		return "st"
	case 2:
		return "nd"
	case 3:
		return "rd"
	default:
		return "th"
	}
}

func formatDateTime(t time.Time) string {
	t = t.Local()
	day := t.Day()
	suffix := getDaySuffix(day)

	// Format hour for 12-hour clock
	hour := t.Hour()
	minute := t.Minute()
	ampm := "am"

	if hour >= 12 {
		ampm = "pm"
		if hour > 12 {
			hour -= 12
		}
	}
	if hour == 0 {
		hour = 12
	}

	return fmt.Sprintf("%d%s %s, %d %d:%02d%s",
		day, suffix, t.Month().String(), t.Year(), hour, minute, ampm)
}

func main() {
	e := echo.New()

	tmpl := template.New("").Funcs(template.FuncMap{
		"formatDate": formatDateTime,
	})
	tmpl, err := tmpl.ParseGlob("templates/*.html")
	if err != nil {
		panic(err)
	}

	e.Renderer = &Template{
		templates: tmpl,
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting wd")
	}

	q, err := initDB(wd + "/kakeibo.db")
	if err != nil {
		log.Fatalf("Error starting sqlite db")
	}

	e.Use(InjectDb(q))

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
