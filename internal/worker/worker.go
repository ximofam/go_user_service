package worker

var GlobalPool Pool

func InitGlobalPool() {
	GlobalPool = NewPoolV1(3, 3)

	GlobalPool.Run()
}
