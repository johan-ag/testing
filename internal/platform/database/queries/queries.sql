-- Users
-- name: SaveUser :execresult
INSERT INTO `users` (
    `name`, `age`
) VALUES ( ?, ? );

-- name: FindUser :one
SELECT * FROM `users` WHERE `id` = ? ;  

-- name: FindUserByParams :many
SELECT * FROM `users` WHERE `name` = ? AND `age` = ?;