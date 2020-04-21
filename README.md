## Socket.io Emitter

Library supports sending message from non-socket.io process via redis. (Mostly configured by socket.io and socket.io-redis in NodeJs)

### Example
`go get -u github.com/ginkcode/socket.io-emitter`

```.go
import "github.com/ginkcode/socket.io-emitter"
import "github.com/go-redis/redis/v7"

func Test() {
    opts := &redis.Options{Addr: "localhost:6379"}
    rdWrapper, err := ioemitter.NewWrapperWithOptions(opts)
    if err != nil {
        panic("Can't create Redis wrapper")
    }
    emitter := ioemitter.NewEmitter("prefix", "/namespace", rdWrapper)
    if err := emitter.To("chat").Emit("Hello world"); err != nil {
        panic("Can't send message to chat room")
    }
    
    if err := emitter.Broadcast("Hello World"); err != nil {
        panic("Can't broadcast message")
    }
}

```

## Notes

- `emitter` use its own Rooms to combine and pack the message before emitting and then reset it.
Let's create new instance or check mutex lock to avoid race condition.
- Can reuse redis connection by passing `redis.Client` or `redis.ClusterClient` to `NewWrapper()` or `NewWrapperWithCluster()`
