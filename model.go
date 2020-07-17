package testutils

type Model interface {
	GetID() string
	Equals(interface{}) bool
}
