package main

import "fmt"

func NewUser() {
	// 初始化结构体
	u := User{}
	fmt.Printf("%v \n", u)
	fmt.Printf("%+v \n", u)

	// up 是一个指针
	up := &User{}
	fmt.Printf("up: %+v \n", up)
	up2 := new(User)
	fmt.Printf("up2: %+v \n", up2)

	u4 := User{Name: "Tom", Age: 0}
	fmt.Printf("u4: %+v \n", u4)

	u4.Name = "zzz"
	u4.Age = 24
	fmt.Printf("u4: %+v \n", u4)

	var up3 *User
	// nil 上访问字段，或方法,会报错
	println(up3.FirstName)
	println(up3)
}

type User struct {
	Name      string
	FirstName string
	Age       int
}

func (u User) ChangeName(name string) {
	fmt.Printf("change name 中 u 的地址：%p \n", &u)
	u.Name = name
}

func (u *User) ChangeAge(age int) {
	fmt.Printf("change age 中 u 的地址：%p \n", u)
	u.Age = age
}

func ChangeUser() {
	u1 := User{Name: "Tom", Age: 18}
	fmt.Printf("change name 中 u 的地址：%p \n", &u1)
	u1.ChangeName("mmm")
	u1.ChangeAge(100)
	fmt.Printf("u1: %+v \n", u1)

	up1 := &User{}
	up1.ChangeName("mmm")
	up1.ChangeAge(99)
	fmt.Printf("up1: %+v \n", up1)
}

type Fish struct {
}

func (f Fish) Swim() {
	fmt.Printf("fish 在游")
}

// FakeFish 衍生类型就是一个全新的类型
type FakeFish Fish

func UseFish() {
	f1 := Fish{}
	f1.Swim()
	f2 := FakeFish(f1)
	//f2.Swim()
	println(f1)
	println(f2)
}
