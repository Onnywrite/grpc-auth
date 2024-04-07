package models

type Service struct {
	OwnerId int64  `db:"owner_fk" validate:"gte=0"`
	Name    string `db:"name" validate:"gte=1,alphanumunicode"`
}

type SavedService struct {
	Id      int64  `db:"service_id"`
	OwnerId int64  `db:"owner_fk"`
	Name    string `db:"name"`
}
