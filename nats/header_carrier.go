package nats

import (
	"github.com/nats-io/nats.go"
)

type HeadersCarrier nats.Header

func (c HeadersCarrier) Get(key string) string {
	return nats.Header(c).Get(key)
}

func (c HeadersCarrier) Set(key, val string) {
	nats.Header(c).Set(key, val)
}

func (c HeadersCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}
