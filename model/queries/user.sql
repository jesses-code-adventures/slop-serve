-- name: UserCreate :one
INSERT INTO users (first_name, last_name, email, password_hash)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: HashedPasswordGet :one
select users.password_hash, users.id user_id
from users
where users.email = $1;
