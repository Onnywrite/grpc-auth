package models

type Service struct {
	OwnerId int64 `db:"owner_fk"`
	Name    string
}

type SavedService struct {
	Id      int64 `db:"service_id"`
	OwnerId int64 `db:"owner_fk"`
	Name    string
}
