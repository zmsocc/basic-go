package lock

import "sync"

// ZLockDemo
// 优先使用 RWMutex, 优先加读锁
// 常用的并发手段，用读写锁来优化读锁
type ZLockDemo struct {
	lock sync.Mutex
}

func (l *ZLockDemo) PanicDemo() {
	l.lock.Lock()
	// 在中间 panic 了，无法释放锁
	panic("hello")
	l.lock.Unlock()
}

func (l *ZLockDemo) LockDemo() {
	l.lock.Lock()
	defer l.lock.Unlock()
}
