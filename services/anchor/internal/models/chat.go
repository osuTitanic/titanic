package models

type Message interface {
	Sender() string
	Content() string
	Target() string
	SenderId() int32
}

type Channel interface {
	Name() string
	Topic() string
	Owner() string
	UserCount() int16
}
