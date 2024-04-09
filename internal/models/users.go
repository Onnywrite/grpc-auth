package models

import "time"

type User struct {
	Login    *string `db:"login" validate:"gte=3,max=30"`
	Email    *string `db:"email" validate:"omitempty,email,max=255"`
	Phone    *string `db:"phone" validate:"omitempty,e164"`
	Password string  `db:"password" validate:"required,lte=72,gte=8" secret:"1"`
}

func (u *User) Idendifier() *UserIdentifier {
	var identifier UserIdentifier

	switch {
	case u.Login != nil:
		identifier = UserIdentifier{Key: "login", Value: *u.Login}
	case u.Email != nil:
		identifier = UserIdentifier{Key: "email", Value: *u.Email}
	case u.Phone != nil:
		identifier = UserIdentifier{Key: "phone", Value: *u.Phone}
	default:
		identifier = UserIdentifier{Key: "login", Value: ""}
	}

	return &identifier
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
