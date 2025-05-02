package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"main_service/pkg/models"
	"strconv"
)

var RedisClient *redis.Client

func InitRedisClient(addr string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})
}

type RedisPublisher struct{}

type Message interface{}

func (p *RedisPublisher) PublishNewDocument(msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal redis message: %w", err)
	}

	channel := fmt.Sprintf("notifications:message:%d", msg)
	return RedisClient.Publish(context.Background(), channel, data).Err()
}

func (p *RedisPublisher) Loginlinkchin(log models.AuditLog) error {
	data, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("failed to marshal redis message: %w", err)
	}

	channel := fmt.Sprintf("notifications:message:%s", strconv.FormatUint(uint64(log.ID), 10))
	return RedisClient.Publish(context.Background(), channel, data).Err()
}
func (p *RedisPublisher) PublishIservice(changes Message) error {
	data, err := json.Marshal(changes)
	if err != nil {
		return fmt.Errorf("failed to marshal redis message: %w", err)
	}

	channel := "AI"
	return RedisClient.Publish(context.Background(), channel, data).Err()
}