package nats

import (
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/stepanbukhtii/easy-tools/elog"
)

func NewConnection(url string) (*nats.Conn, error) {
	return nats.Connect(
		url,
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
			if err != nil {
				slog.With(elog.Err(err)).Error("nats disconnected")
			}
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) { slog.Info("nats reconnected") }),
	)
}
