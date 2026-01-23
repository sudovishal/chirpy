# chirpy_boot.dev

# Goose Migration Up
goose postgres <connection_string> up

example:
`goose postgres "postgres://vishals:@localhost:5432/chirpy" up`

# Goose Migration Down
goose postgres <connection_string> down

example:
`goose postgres "postgres://vishals:@localhost:5432/chirpy" down`

# SQLC Generate

`sqlc generate`
