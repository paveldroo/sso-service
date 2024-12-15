package models

type User struct {
	ID       int    `db:"id"`
	Email    string `db:"email"`
	PassHash string `db:"pass_hash"`
}
