package club

type clientState = string

const (
	inClub  clientState = "IN_CLUB"
	atTable clientState = "TABLE"
	waiting clientState = "WAITING"
	gone    clientState = "GONE"
)

type client struct {
	name  string
	state clientState
	table *table
}
