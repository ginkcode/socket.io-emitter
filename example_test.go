package ioemitter

import (
	"testing"

	"github.com/go-redis/redis/v7"
)

func TestConnectAndEmit(t *testing.T) {
	opts := &redis.Options{Addr: "localhost:6379"}
	wrapper, err := NewWrapperWithOptions(opts)
	if err != nil {
		t.Error("Can't create wrapper")
	}
	emitter := NewEmitter("development", "/", wrapper)
	if err := emitter.To("chat").Emit("Hello world"); err != nil {
		t.Error("Can't send message to room")
	}

	if err := emitter.Broadcast("Hello World"); err != nil {
		t.Error("Can't broadcast message")
	}
}
