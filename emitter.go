package ioemitter

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/vmihailenco/msgpack/v4"

	"github.com/ginkcode/socket.io-emitter/util"
)

var (
	flags = make([]string, 0)
	uuid  = "emitter"
)

type Emitter struct {
	Prefix    string
	Namespace string
	Channel   string
	Rooms     []string
	Flags     []string
	*util.RedisWrapper
}

func NewEmitter(prefix, namespace string, wrapper *util.RedisWrapper) *Emitter {
	e := Emitter{
		Prefix:       prefix,
		Namespace:    namespace,
		Flags:        flags,
		RedisWrapper: wrapper,
	}
	return &e
}

func (e *Emitter) To(room string) *Emitter {
	if e.Rooms == nil {
		e.Rooms = []string{room}
		return e
	}
	if !find(e.Rooms, room) {
		e.Rooms = append(e.Rooms, room)
	}
	return e
}

func (e *Emitter) Emit(data ...interface{}) error {
	if util.IsZero(e.RedisWrapper) {
		return errors.New("must set redis wrapper client first")
	}

	raw, err := encode(e, data...)
	if err != nil {
		return err
	}
	if e.Rooms != nil && len(e.Rooms) == 1 {
		e.Channel = fmt.Sprintf("%s#%s#%s#", e.Prefix, e.Namespace, e.Rooms[0])
	} else {
		e.Channel = fmt.Sprintf("%s#%s#", e.Prefix, e.Namespace)
	}
	defer func() {
		e.Rooms = nil
	}()
	cmd := e.RedisWrapper.Publish(e.Channel, raw)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

func (e *Emitter) Broadcast(data interface{}) error {
	e.Rooms = []string{}
	e.Channel = fmt.Sprintf("%s#%s#", e.Prefix, e.Namespace)
	raw, err := encode(e, "broadcast", data)
	if err != nil {
		return err
	}
	cmd := e.RedisWrapper.Publish(e.Channel, raw)
	if err := cmd.Err(); err != nil {
		return err
	}
	return nil
}

func encode(e *Emitter, data ...interface{}) ([]byte, error) {
	packet := map[string]interface{}{
		"type": "2",
		"data": data,
		"nsp":  e.Namespace,
	}
	msg := []interface{}{
		uuid,
		packet,
		map[string]interface{}{
			"rooms": e.Rooms,
			"flags": flags,
		},
	}
	return msgpack.Marshal(msg)
}

func find(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}

func NewWrapper(client *redis.Client) *util.RedisWrapper {
	return &util.RedisWrapper{
		Client: client,
	}
}

func NewWrapperWithCluster(cluster *redis.ClusterClient) *util.RedisWrapper {
	return &util.RedisWrapper{
		ClusterClient: cluster,
		IsCluster:     true,
	}
}

func NewWrapperWithOptions(opts interface{}) (*util.RedisWrapper, error) {
	w := new(util.RedisWrapper)
	switch opts.(type) {
	case *redis.Options:
		o := opts.(*redis.Options)
		w.Client = redis.NewClient(o)
	case *redis.ClusterOptions:
		o := opts.(*redis.ClusterOptions)
		w.ClusterClient = redis.NewClusterClient(o)
		w.IsCluster = true
	default:
		return nil, errors.New("invalid options, must be Options or ClusterOptions pointer")
	}
	return w, nil
}
