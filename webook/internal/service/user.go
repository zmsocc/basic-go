package service

import (
	"context"
	"errors"
	"gitee.com/zmsoc/gogogo/webook/internal/domain"
	"gitee.com/zmsoc/gogogo/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService interface {
	Login(ctx context.Context, email string, password string) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, u domain.User) error
}

type userService struct {
	repo repository.UserRepository
}

// NewUserService 我用的人只管用，一点都不关心如何初始化
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, u domain.User) error {
	panic("implement me")
}

func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
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

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	// 考虑加密放在哪里
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	// 要判断，有没有这个用户
	// 这个叫做快路径
	if err != repository.ErrUserNotFound {
		// nil 会进来这里
		// 不为 ErrUserNotFound 的也进来这里
		return u, err
	}
	// 在系统资源不足，触发降级之后，不执行慢路径了
	//if ctx.Value("降级") == "true" {
	//	return domain.User{}, errors.New("系统降级了")
	//}
	// 这个叫做慢路径
	// 你明确知道，没有这个用户
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	// 这里会遇到主从延迟的问题
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	// 在系统内部，基本上都是用 id 的
	// 有些人的系统比较复杂，有一个 GUID (Global Unique Id)
	return svc.repo.FindById(ctx, id)
}

func PathDownGrade(ctx context.Context, quick, slow func()) {
	quick()
	if ctx.Value("降级") == "true" {
		return
	}
	slow()
}

//func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
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
