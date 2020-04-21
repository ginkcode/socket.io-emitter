package util

import "github.com/go-redis/redis/v7"

type RedisWrapper struct {
	*redis.Client
	*redis.ClusterClient
	IsCluster bool
}

func (w *RedisWrapper) PubSubNumSub(channels ...string) *redis.StringIntMapCmd {
	if w.IsCluster {
		return w.ClusterClient.PubSubNumSub(channels...)
	}
	return w.Client.PubSubNumSub(channels...)
}

func (w *RedisWrapper) PubSubChannels(pattern string) *redis.StringSliceCmd {
	if w.IsCluster {
		return w.ClusterClient.PubSubChannels(pattern)
	}
	return w.Client.PubSubChannels(pattern)
}

func (w *RedisWrapper) Publish(channel string, message interface{}) *redis.IntCmd {
	if w.IsCluster {
		return w.ClusterClient.Publish(channel, message)
	}
	return w.Client.Publish(channel, message)
}

func (w *RedisWrapper) Subscribe(channels ...string) *redis.PubSub {
	if w.IsCluster {
		return w.ClusterClient.Subscribe(channels...)
	}
	return w.Client.Subscribe(channels...)
}
