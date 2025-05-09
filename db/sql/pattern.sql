-- name: SelectAllPattern :many
select "id", "name", "desc", graph_data from oteldemo.pattern;

-- name: InsertPattern :one
insert into oteldemo.pattern (name, "desc", graph_data)
values (sqlc.arg(name), sqlc.arg(description), sqlc.arg(graph_data))
returning id;