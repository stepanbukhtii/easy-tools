package kafka

type Event interface {
	Topic() string
}
