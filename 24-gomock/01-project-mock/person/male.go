package person

//go:generate mockgen -source=male.go -destination=../mock/male_mock.go -package=mock
type Male interface {
	Get(id int64) error
}
