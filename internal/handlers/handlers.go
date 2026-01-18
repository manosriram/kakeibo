package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/manosriram/kakeibo/sqlc/db"
)

type CreateTxn struct {
	Description string `json:"description"`
}

func CreateStatement(d *db.Queries, description string) error {
	err := d.CreateStatement(context.Background(), db.CreateStatementParams{
		TxnType:     sql.NullString{String: "online", Valid: true},
		Amount:      sql.NullInt64{Int64: 1000, Valid: true},
		Tag:         sql.NullString{String: "food", Valid: true},
		Description: sql.NullString{String: description, Valid: true},
	})
	return err
}

func HomeHandler(c echo.Context) error {
	db := c.Get("db").(*db.Queries)
	txns, err := db.GetAllStatements(context.Background())
	if err != nil {
	}

	return c.Render(http.StatusOK, "index.html", map[string]any{
		"statements": txns,
	})
}

func CreateTransactionAPI(c echo.Context) error {
	d := c.Get("db").(*db.Queries)

	txn := new(CreateTxn)
	if err := c.Bind(txn); err != nil {
		fmt.Println(err)
	}
	// txnType := c.FormValue("")
	// amount := c.FormValue("")
	// tag := c.FormValue("")
	err := CreateStatement(d, txn.Description)
	if err != nil {
		fmt.Println("err = ", err.Error())
	}

	return c.JSON(200, map[string]any{
		"message": "Transaction created",
	})
}

func GetAllTransactionsAPI(c echo.Context) error {
	db := c.Get("db").(*db.Queries)
	txns, _ := db.GetAllStatements(context.Background())

	return c.JSON(200, map[string]any{
		"statements": txns,
	})
}
