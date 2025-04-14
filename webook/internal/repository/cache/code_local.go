package cache

//import (
//	"context"
//	lru "github.com/hashicorp/golang-lru"
//	"sync"
//	"time"
//)
//
//// 技术选型考虑的点
//// 1.功能性：功能是否能够完全覆盖你的需求。
//// 2.社区和支持度：社区是否活跃，文档是否齐全，
//// 	 以及百度（搜索引擎）能不能搜索到你需要的各种信息，有没有帮你踩过坑
//// 3.非功能性：易用性（用户友好度，学习曲线要平滑），
//// 	 扩展性（如果开源软件的某些功能需要定制，框架是否支持定制，以及定制的难度高不高）
////   性能（追求性能的公司，往往有能力自研）
//
//// LocalCodeCache 本地缓存实现
//type LocalCodeCache struct {
//	cache *lru.Cache
//	// 普通锁、或者说写锁
//	lock sync.Mutex
//	// 读写锁
//	rwlock     sync.RWMutex
//	expiration time.Duration
//}
//
//func NewLocalCodeCache(c *lru.Cache, expiration time.Duration) *LocalCodeCache {
//	return &LocalCodeCache{
//		cache:      c,
//		expiration: expiration,
//	}
//}
//
//func (l *LocalCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
//	l.lock.Lock()
//	defer l.lock.Unlock()
//	// 这里可以考虑用读写锁来优化，但是效果不会很好
//	// 因为你可以预期，大部分时候是要走到写锁里面的
//
//	// 我选用的本地缓存，很不幸的是，没有获得过期时间的接口，所以都是自己维持了一个过期的接口
//	firstPageKey := l.
//	now := time.Now()
//	val, ok := l.cache.Get(firstPageKey)
//	if !ok {
//		// 说明没有验证码
//		l.cache.Add(firstPageKey, codeItem{})
//	}
//}
