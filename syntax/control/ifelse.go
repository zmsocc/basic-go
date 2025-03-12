package main

func IfOnly(age int) {
	if age > 18 {
		println("成年了")
	}
}

func IfElse(age int) {
	if age > 18 {
		println("成年了")
	} else {
		println("未成年")
	}
}

func IfElseIf(age int) {
	if age > 18 {
		println("成年了")
	} else if age > 12 {
		println("青年")
	} else {
		println("小孩子")
	}
}

func IfNewVariable(start int, end int) string {
	if distance := end - start; distance > 100 {
		return "太远了"
	} else if distance > 60 {
		return "有点远"
	} else {
		return "还行"
	}
}
