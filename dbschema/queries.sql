-- name: HistoryList :many
SELECT * FROM history ORDER BY id;

-- name: HistoryInsert :exec
INSERT INTO history (id, namespace, reference, event, email, created_at) VALUES ($1, $2, $3, $4, $5, $6);

-- name: WalletsByUser :many
SELECT * FROM wallet WHERE user_id = $1 ORDER BY wallet.created_at;

-- name: WalletInsert :exec
INSERT INTO wallet (id, user_id, balance, currency, created_at) VALUES ($1, $2, $3, $4, $5);

-- name: WalletUpdateBalance :exec
UPDATE wallet SET balance = $1 WHERE id = $2;

-- name: ExpenseInsert :exec
INSERT INTO expense (id, wallet_id, amount, description, created_at) VALUES ($1, $2, $3, $4, $5);

-- name: ExpenseListByWallet :many
SELECT * FROM expense WHERE wallet_id = $1 ORDER BY id;

-- name: ExpenseListByWalletByUser :many
SELECT * FROM expense WHERE wallet_id = $1 AND wallet_id IN (
    SELECT id FROM wallet WHERE user_id = $2
) 
ORDER BY id;

-- name: LocalUserInsert :exec
INSERT INTO local_user (email, display_name, roles, created_at, pwdhash) VALUES ($1, $2, $3, $4, $5);

-- name: LocalUserUpdate :exec
UPDATE local_user SET roles = $1 WHERE email = $2;

-- name: LocalUserSetPass :exec
UPDATE local_user SET pwdhash = $1 WHERE email = $2;

-- name: LocalUserGetByEmail :one
SELECT * FROM local_user WHERE email = $1;

-- name: LocalUserList :many
SELECT * FROM local_user ORDER BY email;
