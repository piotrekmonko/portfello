-- name: UsersList :many
SELECT * FROM users ORDER BY display_name;

-- name: RolesListForUser :many
SELECT * FROM roles WHERE user_id = $1;
