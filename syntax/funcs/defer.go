package main

func Defer() {
	defer func() {
		println("第一个defer")
	}()
	defer func() {
		println("第二个defer")
	}()
}

func DeferClosure() {
	i := 0
	defer func() {
		println(i)
	}()
	i = 1
}

func DeferClosure1() {
	i := 0
	defer func(i int) {
		println(i)
	}(i)
	i = 1
}

func DeferReturn() int {
	a := 0
	defer func() {
		a = 1
	}()
	return a
}

func DeferReturn1() (a int) {
	a = 0
	defer func() { // defer 里面可以改带返回值名字的参数
		a = 1
	}()
	return a
}

func DeferClosureLoopV1() {
	for i := 0; i < 10; i++ {
		defer func() {
			println(i)
		}()
	}
}

func DeferClosureLoopV2() {
	for i := 0; i < 10; i++ {
		defer func(val int) {
			println(val)
		}(i)
	}
}

func DeferClosureLoopV3() {
	for i := 0; i < 10; i++ {
		j := i
		defer func() {
			println(j)
		}()
	}
}
