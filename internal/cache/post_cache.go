package cache

import (
	"context"
	"crud_service/internal/domain"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	postsAllKey string = "posts:all"
	postKeyFmt  string = "posts:%d"
	postsTTL           = 60 * time.Second
)

type PostCache struct {
	Client *redis.Client
}

func NewPostCache(client *redis.Client) *PostCache {
	return &PostCache{Client: client}
}

func (c *PostCache) GetAll(ctx context.Context) ([]*domain.Post, error) {
	data, err := c.Client.Get(ctx, postsAllKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("PostCache.GetAll: %w", err)
	}

	var posts []*domain.Post

	if err := json.Unmarshal(data, &posts); err != nil {
		return nil, fmt.Errorf("PostCache.GetAll unmarshal: %w", err)
	}

	return posts, nil
}

func (c *PostCache) SetAll(ctx context.Context, posts []*domain.Post) error {
	data, err := json.Marshal(&posts)
	if err != nil {
		return fmt.Errorf("PostCache.SetAll Marshal: %w", err)
	}

	if err := c.Client.Set(ctx, postsAllKey, data, postsTTL).Err(); err != nil {
		return fmt.Errorf("PostCache.SetAll set: %w", err)
	}

	return nil
}

func (c *PostCache) GetById(ctx context.Context, id int) (*domain.Post, error) {
	data, err := c.Client.Get(ctx, fmt.Sprintf(postKeyFmt, id)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("PostCache.GetById: %w", err)
	}

	var post *domain.Post

	if err := json.Unmarshal(data, &post); err != nil {
		return nil, fmt.Errorf("PostCache.GetById unmarshal: %w", err)
	}

	return post, nil
}

func (c *PostCache) SetById(ctx context.Context, post *domain.Post) error {
	data, err := json.Marshal(&post)
	if err != nil {
		return fmt.Errorf("PostCache.SetById marshal: %w", err)
	}

	if err := c.Client.Set(ctx, fmt.Sprintf(postKeyFmt, post.Id), data, postsTTL).Err(); err != nil {
		return fmt.Errorf("PostCache.SetById set: %w", err)
	}

	return nil
}

// CREATE, UPDATE, DELETE for post
func (c *PostCache) Invalidate(ctx context.Context, id int) error {
	if err := c.Client.Del(ctx, postsAllKey, fmt.Sprintf(postKeyFmt, id)).Err(); err != nil {
		return fmt.Errorf("PostCache.SetById del: %w", err)
	}

	return nil
}
