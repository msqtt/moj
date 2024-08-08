package listener

type Listener interface {
	OnEvent(any) error
}
