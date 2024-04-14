package session

type State = string

var (
	NoSession State = "no"
	Active    State = "active"
	Suspended State = "suspended"
	NeedPin   State = "need_pin"
)

var ActivePeriodSeconds float64 = 200

type UserInfo struct {
	ID        string
	HashPin   string
	HashToken string
}
