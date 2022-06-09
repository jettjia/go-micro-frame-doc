package project

//go:generate mockgen -source=./db.go -destination ./mock/mock_db.go -package mock
type DB interface {
	Get(key string) (int, error)
}

func GetFromDB(db DB, key string) int {
	if value, err := db.Get(key); err == nil {
		return value
	}

	return -1
}
