package entity

type User struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	Email    string `db:"email"`
	FullName string `db:"full_name"`
}

type GetUser struct {
	Username string `db:"username"`
	FullName string `db:"full_name"`
}
