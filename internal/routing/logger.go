package routing

//go:generate mockgen -destination=logger_mock_test.go -package routing . Logger

type Logger interface {
	Debug(s string)
	Info(s string)
	Warn(s string)
	Error(s string)
}
