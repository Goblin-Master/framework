package domain

type User struct {
	Id   int64
	Name string
	Addr Addr
}

// 值对象，只能整体替换
type Addr struct {
	Id       int64
	Province string
	City     string
	District string
}
