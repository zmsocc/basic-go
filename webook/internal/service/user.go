package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gitee.com/zmsoc/gogogo/webook/internal/domain"
	"gitee.com/zmsoc/gogogo/webook/internal/repository"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicateEmail
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	// 先找用户
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	// 再比较密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 考虑加密放在哪里
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)

	// 然后就是存起来
	err = svc.repo.Create(ctx, u)
	if err != nil {
		return err
	}
	// Redis 不知道怎么处理这个 u，所以要转化
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	// 要求 u 的 id 不为 0
	err = svc.redis.Set(ctx, fmt.Sprintf("user:info:%d", u.Id), val, time.Minute*30)
	return err
}

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	// 第一个念头是
	val, err := svc.redis.Get(ctx, fmt.Sprintf("user:info:%d", id)).Result()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal([]byte(val), &u)
	if err != nil {
		return u, err
	}
	// 接下来就是从数据库里面查找
}
