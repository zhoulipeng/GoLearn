package main

import "fmt"

type Base struct {
	Name string
}
func (base *Base) Foo() {
	fmt.Println("foo", base.Name)
}
func (base *Base) Bar() {
	fmt.Println("bar", base.Name)
}
type Vertex struct {
	X, Y int
	*Base
}

func (v *Vertex) Bar() {
	v.Base.Bar()
	fmt.Println("extend bar", v.Name)
}

func (v *Vertex) Bar1() {
	fmt.Println("extend bar1", v.Name)
	v.Base.Bar()
}

var (
	v1 = Vertex{1, 2, &Base{"name1"}}  // 创建一个 Vertex 类型的结构体
	v2 = Vertex{X: 1}  // Y:0 被隐式地赋予
	v3 = Vertex{}      // X:0 Y:0
	p  = &Vertex{1, 2, &Base{"name2"}} // 创建一个 *Vertex 类型的结构体（指针）
)

func main() {
	fmt.Println(v1, p, v2, v3)
	v1.Base.Bar()
	v1.Bar()
	v1.Bar1()
	p.Base.Bar()
	p.Bar()
	p.Bar1()
}
