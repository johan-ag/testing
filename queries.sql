-- Users
-- name: SaveUser :execresult
INSERT INTO `users` (
    `name`, `age`, `random`
) VALUES ( ?, ?, ? );

-- name: FindUser :one
SELECT * FROM `users` WHERE `id` = ? ;  

-- Books
-- name: SaveBook :execresult
INSERT INTO `books` (
    `title`, `author`
) VALUES ( ?, ? );

-- name: FindBook :one
SELECT * FROM `books` WHERE `id` = ? ;  