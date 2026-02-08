package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/manosriram/kakeibo/internal/llm"
	"github.com/manosriram/kakeibo/internal/rag"
	"github.com/manosriram/kakeibo/sqlc/db"
)

type StatementFromLLM struct {
	Tag             string `json:"tag"`
	Amount          int64  `json:"amount"`
	Description     string `json:"description"`
	TransactionType string `json:"txn_type"`
}

type CreateTxn struct {
	Description string `json:"description"`
}

func CreateStatement(d *db.Queries, description string) error {
	// wd, err := os.Getwd()
	// if err != nil {
	// return err
	// }
	f, err := os.OpenFile("/app/data/spends.csv", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	result, err := llm.NewClaude(description).Call()
	if err != nil {
		fmt.Println("error calling claude ", err)
	} else {
	}

	result = strings.TrimLeft(result, "```json")
	result = strings.Trim(result, "```")

	var statements []StatementFromLLM
	json.Unmarshal([]byte(result), &statements)
	fmt.Println("statements = ", statements)

	for _, statement := range statements {
		stmt, err := d.CreateStatement(context.Background(), db.CreateStatementParams{
			TxnType:     sql.NullString{String: statement.TransactionType, Valid: true},
			Amount:      sql.NullInt64{Int64: statement.Amount, Valid: true},
			Tag:         sql.NullString{String: statement.Tag, Valid: true},
			Description: sql.NullString{String: statement.Description, Valid: true},
		})
		if err != nil {
			return err
		}

		_, err = fmt.Fprintf(f, "%s,%s,INR %v,%s,%s\n", stmt.Tag.String, stmt.TxnType.String, stmt.Amount.Int64, stmt.Description.String, stmt.CreatedAt.Time)
		if err != nil {
			return err
		}
	}

	return nil
}

func HealthHandler(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"message": "Health ok",
	})
}

func HomeHandler(c echo.Context) error {
	db := c.Get("db").(*db.Queries)
	page := c.QueryParam("page")
	if page == "" {
		page = "1"
	}
	p, _ := strconv.Atoi(page)

	// offset = (page_number - 1) * page_size
	offset := (p - 1) * 10

	count, err := db.GetStatementsCount(context.Background())
	if err != nil {
		return err
	}

	next := p
	if int64(p+1) <= count {
		next = p + 1
	}

	prev := 0
	if p-1 >= 0 && int64(p-1) < count {
		prev = p - 1
	}

	txns, err := db.GetAllStatementsPaginated(context.Background(), int64(offset))
	if err != nil {
		return err
	}

	totalPages := count / 2

	cr, err := db.GetCurrentMonthCredit(context.Background())
	credit := cr.Float64

	de, err := db.GetCurrentMonthDebit(context.Background())
	debit := de.Float64

	netSavings := credit - debit
	var savingsPerc string
	if credit == 0 {
		savingsPerc = "0"
	} else {
		savingsPerc = fmt.Sprintf("%.2f", (netSavings/credit)*100)
	}

	statementsByTag, _ := db.GetStatementsByCategory(context.Background())

	return c.Render(http.StatusOK, "index.html", map[string]any{
		"statements":        txns,
		"page":              p,
		"pageCount":         totalPages,
		"next":              next,
		"prev":              prev,
		"credit":            credit,
		"debit":             debit,
		"savings":           netSavings,
		"savingsPercentage": savingsPerc,
		"spendByTag":        statementsByTag,
		"totalEntries":      count,
		"hasNextPage":       p < int(totalPages)-1,
	})
}

func QueryRagAPI(c echo.Context) error {
	query := c.QueryParam("query")
	r := rag.NewRAG()
	if query != "" {
		answer, err := r.Query(query)
		if err == nil {
			return c.JSON(200, map[string]any{
				"llmResponse": answer,
			})
		} else {
			return c.JSON(400, map[string]any{
				"error": err.Error(),
			})
		}
	} else {
		return c.JSON(400, map[string]any{
			"message": "Query should not be empty",
		})
	}
}

func HealthAPI(c echo.Context) error {
	return c.JSON(200, map[string]any{
		"health": "ok",
	})
}

func CreateTransactionAPI(c echo.Context) error {
	d := c.Get("db").(*db.Queries)

	txn := new(CreateTxn)
	if err := c.Bind(txn); err != nil {
		fmt.Println(err)
	}

	result, err := llm.NewClaude(txn.Description).Call()
	if err != nil {
		fmt.Println("error calling claude ", err)
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
