package status

type Handler interface {
	Handle(interface{}) bool
}

type Publisher interface {
	AddHandler(Handler)
	Publish(interface{})
}
