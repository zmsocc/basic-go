package domain

import "time"

// User 领域对象， 是 DDD 中的entity
// BO(business object)
type User struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time
}

//type Address struct {
//}
