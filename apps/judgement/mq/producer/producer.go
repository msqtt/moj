package producer

type Producer interface {
	Send(message any) error
	Close()
}
