package rabbitmq

type Event interface {
	Exchange() string
	RoutingKey() string
}
