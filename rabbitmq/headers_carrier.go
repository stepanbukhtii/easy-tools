package rabbitmq

import "github.com/rabbitmq/amqp091-go"

type HeadersCarrier amqp091.Table

func (c HeadersCarrier) Get(key string) string {
	if v, ok := c[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func (c HeadersCarrier) Set(key string, value string) {
	c[key] = value
}

func (c HeadersCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}
