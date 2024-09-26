-- name: GenImageCreate :one
insert into generated_images (url) values ($1) returning id;
