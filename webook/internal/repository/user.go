package repository

import (
	"context"
	"database/sql"
	"gitee.com/zmsoc/gogogo/webook/internal/domain"
	"gitee.com/zmsoc/gogogo/webook/internal/repository/cache"
	"gitee.com/zmsoc/gogogo/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

//var ErrUserDuplicateEmailV1 = fmt.Errorf("%w 邮箱错误", dao.ErrUserDuplicateEmail)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
}

type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, c cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: c,
	}
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	// 先从 cache 里面找
	// 再从 dao 里面找
	// 找到了回写 cache
	u, err := r.cache.Get(ctx, id)
	if err != nil {
		// 必然是有数据
		return domain.User{}, err
	}
	// 没这个数据
	//if err == cache.ErrKeyNotExist {
	//	去数据库里面加载
	//}
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	u = r.entityToDomain(ue)
	err = r.cache.Set(ctx, u)
	if err != nil {
		// 打日志，做监控
	}
	return u, err
}

func (r *CachedUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			// 我确实有手机号
			Valid: u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			// 我确实有手机号
			Valid: u.Phone != "",
		},
		Password: u.Password,
		Ctime:    u.Ctime.UnixMilli(),
	}
}

func (r *CachedUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}
