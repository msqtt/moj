package listener

type Listener interface {
	OnEvent(event any) error
}
