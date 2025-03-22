package service

import (
	"context"
	"errors"
	"gitee.com/zmsoc/gogogo/webook/internal/domain"
	"gitee.com/zmsoc/gogogo/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
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
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error {
	u, err := svc.repo.FindByPhone(ctx, phone)
	// 要判断，有没有这个用户
	if err != repository.ErrUserNotFound {
		// nil 会进来这里
		// 不为 ErrUserNotFound 的也进来这里
		return u, err
	}

	// 你明确知道，没有这个用户
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	return u, err
})

func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, id)
	return u, err

}

//func (svc *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
//// 第一个念头是
//val, err := svc.redis.Get(ctx, fmt.Sprintf("user:info:%d", id)).Result()
//if err != nil {
//	return domain.User{}, err
//}
//var u domain.User
//err = json.Unmarshal([]byte(val), &u)
//if err != nil {
//	return u, err
//}
//// 接下来就是从数据库里面查找
//}
