-- name: CreateStatement :exec
INSERT INTO STATEMENTS(TXN_TYPE, AMOUNT, TAG, DESCRIPTION) VALUES(?, ?, ?, ?) RETURNING *;

-- name: GetStatementsByTag :many
SELECT * FROM STATEMENTS WHERE TAG = ?;

-- name: GetAllStatementsPaginated :many
SELECT *
FROM STATEMENTS
ORDER BY CREATED_AT DESC
LIMIT 10 OFFSET ?;

-- name: GetStatementsCount :one
SELECT COUNT(*) FROM STATEMENTS;

-- name: GetAllStatements :many
SELECT * FROM STATEMENTS;

-- name: GetStatementsLimit :many
SELECT * FROM STATEMENTS LIMIT ?;

-- name: GetCurrentMonthBalance :one
SELECT 
    SUM(CASE WHEN TXN_TYPE = 'INCOME' THEN AMOUNT ELSE -AMOUNT END) AS monthly_net_balance
FROM STATEMENTS
WHERE strftime('%Y-%m', CREATED_AT) = strftime('%Y-%m', 'now');

-- name: GetCurrentMonthCredit :one
SELECT 
    SUM(AMOUNT) AS monthly_net_balance
FROM STATEMENTS
WHERE strftime('%Y-%m', CREATED_AT) = strftime('%Y-%m', 'now') AND TXN_TYPE = 'credit';

-- name: GetCurrentMonthDebit :one
SELECT 
    SUM(AMOUNT) AS monthly_net_balance
FROM STATEMENTS
WHERE strftime('%Y-%m', CREATED_AT) = strftime('%Y-%m', 'now') AND TXN_TYPE = 'debit';

-- name: GetStatementsByCategory :many
SELECT SUM(AMOUNT), TAG FROM STATEMENTS GROUP BY TAG;
