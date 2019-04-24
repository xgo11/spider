package core

type IQueue interface {
	Name() string
	Put(message ...string) error
	Pop(count ...int) []string
	Size() int
	Limit() int
}
