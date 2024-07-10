package queue

type EventQueue interface {
	EnQueue(event any) error
}
