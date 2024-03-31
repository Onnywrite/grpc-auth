package models

type User struct {
	Login    string  `validation:"ne=me,gt=0,max=30"`
	Email    *string `validation:"email,max=255"`
	Phone    *string `validation:"e164"`
	Password string  `validation:"lte=72"`
}

type SavedUser struct {
	Id       int64 `db:"user_id"`
	Login    string
	Email    *string
	Phone    *string
	Password string
}
