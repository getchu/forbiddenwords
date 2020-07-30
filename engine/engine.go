package engine

type Engine interface {
	Add(string) error
	Find(string) string
	Len() int
	IsExists(string) bool
}

//空结构体，内存为0
var empty struct{}
