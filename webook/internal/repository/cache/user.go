package cache

import (
	"context"
	"gitee.com/zmsoc/gogogo/webook/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	// 传单机 Redis 可以
	// 传 cluster 的 Redis 也可以
	cmd *redis.Client
}

// A 用到了 B， B 一定是接口
// A 用到了 B， B 一定是 A 的字段
// A 用到了 B， A 绝对不初始化 B， 而是外面注入
func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client: client,
	}
}

func (u *UserCache) GetUser(ctx context.Context, id int64) (domain.User, error) {

}
