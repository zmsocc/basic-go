package main

type List interface {
	Add(idx int, val any) error // 正常来说，接口里面定义的应该都是公有的方法
	Append(val any)
	Delete(idx int) (any, error)
	//toSlice() ([]any, error) // 开头小写，全局不可用
}
