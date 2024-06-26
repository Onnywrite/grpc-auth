package models

import "time"

type User struct {
	Login    *string `db:"login" validate:"gte=3,max=30"`
	Email    *string `db:"email" validate:"omitempty,email,max=255"`
	Phone    *string `db:"phone" validate:"omitempty,e164"`
	Password string  `db:"password" validate:"required,lte=72,gte=8" secret:"1"`
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
