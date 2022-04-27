package contract

const IDKey = "eehe:id"

type IDService interface {
	NewID() string
}
