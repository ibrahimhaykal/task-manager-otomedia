package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/ibrahimhaykal/otomedia-task-manager/internal/model"
	"github.com/redis/go-redis/v9"
)

const taskListKeyPrefix = "tasks:list:"

type TaskCache interface {
	GetTaskList(ctx context.Context, key string) (*model.TaskListResponse, error)
	SetTaskList(ctx context.Context, key string, data *model.TaskListResponse, ttl time.Duration) error
	InvalidateTaskLists(ctx context.Context) error
}

type redisTaskCache struct {
	client *redis.Client
}

func NewRedisTaskCache(client *redis.Client) *redisTaskCache {
	return &redisTaskCache{client: client}
}

func BuildTaskListKey(filter model.TaskFilter) string {
	params := url.Values{}
	params.Set("status", filter.Status)
	params.Set("keyword", filter.Keyword)
	params.Set("assignee", filter.Assignee)
	params.Set("page", strconv.Itoa(filter.Page))
	params.Set("limit", strconv.Itoa(filter.Limit))
	params.Set("sort", filter.Sort)
	return taskListKeyPrefix + params.Encode()
}

func (c *redisTaskCache) GetTaskList(ctx context.Context, key string) (*model.TaskListResponse, error) {
	rawValue, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("redis get task list: %w", err)
	}

	var cachedResponse model.TaskListResponse
	if err := json.Unmarshal([]byte(rawValue), &cachedResponse); err != nil {
		return nil, fmt.Errorf("unmarshal cached task list: %w", err)
	}

	return &cachedResponse, nil
}

func (c *redisTaskCache) SetTaskList(ctx context.Context, key string, data *model.TaskListResponse, ttl time.Duration) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal task list: %w", err)
	}

	if err := c.client.Set(ctx, key, payload, ttl).Err(); err != nil {
		return fmt.Errorf("redis set task list: %w", err)
	}

	return nil
}

func (c *redisTaskCache) InvalidateTaskLists(ctx context.Context) error {
	var cursor uint64
	pattern := taskListKeyPrefix + "*"

	for {
		keys, nextCursor, err := c.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("redis scan task list keys: %w", err)
		}

		if len(keys) > 0 {
			if err := c.client.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("redis delete task list keys: %w", err)
			}
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return nil
}
