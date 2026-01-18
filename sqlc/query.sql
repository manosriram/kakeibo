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
