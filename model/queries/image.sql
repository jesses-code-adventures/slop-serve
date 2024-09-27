-- name: GenImageCreate :one
insert into generated_images (url, user_id) values ($1, $2) returning id;
