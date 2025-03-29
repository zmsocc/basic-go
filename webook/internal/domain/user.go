package domain

import "time"

// User 领域对象， 是 DDD 中的entity
// BO(business object)
type User struct {
	Id       int64
	Email    string
	Password string
	Phone    string
	Ctime    time.Time
	Nickname string
	Birthday time.Time
	AboutMe  string
	// 不要组合，万一将来有同名字段，dingdingInfo
	WechatInfo WechatInfo
}

//type Address struct {
//}
