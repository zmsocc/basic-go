package main

func YourName(name string, aliases ...string) {
	// aliases 是一个切片
}

func CallYourName() {
	YourName("大明")
	YourName("大明", "小明")
	YourName("大明", "小明", "拉拉")
	aliases := []string{"小明", "拉拉"}
	YourName("大明", aliases...)
}
