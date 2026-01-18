package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/manosriram/kakeibo/internal/llm"
	"github.com/manosriram/kakeibo/sqlc/db"
)

type StatementFromLLM struct {
	Tag         string `json:"tag"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

type CreateTxn struct {
	Description string `json:"description"`
}

func CreateStatement(d *db.Queries, description string) error {
	result, err := llm.NewOpenAI(description).Call()
	if err != nil {
		fmt.Println("error calling openai ", err)
	} else {
		// fmt.Println("res = ", result)
	}

	result = strings.TrimLeft(result, "```json")
	result = strings.Trim(result, "```")

	var statements []StatementFromLLM
	json.Unmarshal([]byte(result), &statements)
	fmt.Println("statements = ", statements)

	for _, statement := range statements {
		err = d.CreateStatement(context.Background(), db.CreateStatementParams{
			// TxnType:     sql.NullString{String: "online", Valid: true},
			Amount:      sql.NullInt64{Int64: statement.Amount, Valid: true},
			Tag:         sql.NullString{String: statement.Tag, Valid: true},
			Description: sql.NullString{String: statement.Description, Valid: true},
		})
		if err != nil {
			return err
		}
	}

	return nil
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

	result, err := llm.NewOpenAI(txn.Description).Call()
	if err != nil {
		fmt.Println("error calling openai ", err)
	} else {
		fmt.Println("res = ", result)
	}

	// txnType := c.FormValue("")
	// amount := c.FormValue("")
	// tag := c.FormValue("")
	err = CreateStatement(d, txn.Description)
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
