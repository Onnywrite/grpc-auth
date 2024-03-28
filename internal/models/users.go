package models

// TODO: add validation
type User struct {
	Login    string
	Email    string
	Phone    string
	Password string
}

type SavedUser struct {
	Id    int64
	Login string
	Email string
	Phone string
}

func (u *User) Saved(id int64) *SavedUser {
	return &SavedUser{
		Id:    id,
		Login: u.Login,
		Email: u.Email,
		Phone: u.Phone,
	}
}
