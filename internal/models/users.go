package models

// TODO: add validation
type User struct {
	Login    string
	Email    string
	Phone    string
	Password string
}

type SavedUser struct {
	Id    int64 `db:"user_id"`
	Login string
	Email *string // to enable nullability
	Phone *string // to enable nullability
}
