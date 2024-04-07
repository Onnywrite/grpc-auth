package models

import "time"

type User struct {
	Login    string  `db:"login" validate:"ne=me,gt=0,max=30"`
	Email    *string `db:"email" validate:"email,max=255"`
	Phone    *string `db:"phone" validate:"e164"`
	Password string  `db:"password" validate:"lte=72,gte=8" secret:"1"`
}

type SavedUser struct {
	Id        int64      `db:"user_id"`
	Login     string     `db:"login"`
	Email     *string    `db:"email"`
	Phone     *string    `db:"phone"`
	Password  string     `db:"password"`
	DeletedAt *time.Time `db:"deleted_at"`
}

func (u *SavedUser) IsDeleted() bool {
	return u.DeletedAt != nil
}
