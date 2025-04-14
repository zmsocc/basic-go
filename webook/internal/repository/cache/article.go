package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee.com/zmsoc/gogogo/webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, author int64, arts []domain.Article) error
	DelFirstPage(ctx context.Context, author int64) error
	Set(ctx context.Context, id int64) error

	// SetPub 正常来说，创作者和读者的 Redis 集群要分开，因为读者是一个核心中的核心
	SetPub(ctx context.Context, article domain.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
}

type RedisArticleCache struct {
	client redis.Cmdable
}

func (r *RedisArticleCache) Set(ctx context.Context, art domain.Article, id int64) error {
	data, err := json.Marshal(art)
	if err != nil {
		return err
	}
	// 过期时间要短，你的预测效果越不好，就要越短
	return r.client.Set(ctx, r.key(id), data, time.Minute).Err()
}

func (r *RedisArticleCache) GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedisArticleCache) SetFirstPage(ctx context.Context, author int64, arts []domain.Article) error {
	for i := 0; i < len(arts); i++ {
		arts[i].Content = arts[i].Abstract()
	}
	data, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.firstPageKey(author), data, time.Minute*10).Err()
}

func (r *RedisArticleCache) key(id int64) string {
	return fmt.Sprintf("article:%d", id)
}

func (r *RedisArticleCache) firstPageKey(uid int64) string {
	return fmt.Sprintf("article:first_page:%d", uid)
}

func (r *RedisArticleCache) DelFirstPage(ctx context.Context, author int64) error {
	//TODO implement me
	panic("implement me")
}
