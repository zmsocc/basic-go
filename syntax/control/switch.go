package main

func Switch(status int) {
	switch status {
	case 0:
		println("初始化")
	case 1:
		println("运行中")
	default: // default 有没有都可以，没有default的话就会直接跳出switch语句
		println("未知状态")
	}
}

func SwitchBool(age int) {
	switch {
	case age >= 18:
		println("adult")
	case age > 12:
		println("teenager")
	}
}
