package main

var Cfg = ConfigInit()

func main() {
	FlagInit()
	LogInit()
	Serve()
}
