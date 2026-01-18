-- name: CreateStatement :exec
INSERT INTO STATEMENTS(TXN_TYPE, AMOUNT, TAG, DESCRIPTION) VALUES(?, ?, ?, ?) RETURNING *;

-- name: GetStatementsByTag :many
SELECT * FROM STATEMENTS WHERE TAG = ?;

-- name: GetAllStatements :many
SELECT * FROM STATEMENTS;

-- name: GetStatementsLimit :many
SELECT * FROM STATEMENTS LIMIT ?;
