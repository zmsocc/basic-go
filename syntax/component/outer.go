package main

// 组合不是继承，没有多态

type Inner struct {
}

func (i Inner) DoSmoething() {
	println("这是Inner")
}

func (i Inner) SayHello() {
	println("hello", i.Name())
}

func (i Inner) Name() string {
	return "inner"
}

type Outer struct {
	Inner
}

func (o Outer) Name() string {
	return "outer"
}

type OuterV1 struct {
	Inner
}

func (o OuterV1) DoSmoething() {
	println("这是OuterV1")
}

type OuterPtr struct {
	*Inner
}

type OOOOuter struct {
	Outer
}

func UseInner() {
	var o Outer
	o.DoSmoething() // 组合了就可以调用里面的方法,常用这种

	var op *OuterPtr
	op.DoSmoething()

	// 初始化
	o1 := Outer{
		Inner: Inner{},
	}
	op1 := OuterPtr{
		Inner: &Inner{},
	}
	o1.DoSmoething()
	op1.DoSmoething()
}

func main() {
	var o1 OuterV1
	o1.DoSmoething()
	o1.Inner.DoSmoething()

	var o Outer
	o.SayHello() // 输出hello inner
}
