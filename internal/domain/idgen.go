package domain

type IDGenerator interface {
	GenerateID() (string, error)
}
