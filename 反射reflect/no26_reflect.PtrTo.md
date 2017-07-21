## func PtrTo(t Type) Type {}

PtrTo是一个超简单的方法。  
```go
// PtrTo returns the pointer type with element t.
// For example, if t represents type Foo, PtrTo(t) represents *Foo.
func PtrTo(t Type) Type {
	return t.(*rtype).ptrTo()
}
```
看注释就明白，也就是你传个变量过来，它把这个变量变成了指针形式。  
不过我仍然没有弄清楚这个的应用场景到底在什么地方，我照着官方的说明写了一个例子，大概已经能很好的诠释reflect.PtrTo（）的功能了。  
代码如下，主要是转成指针，然后再还原回来：  
```go
func main(){
	var s string
	s = "hello"

	ns:=reflect.PtrTo(reflect.TypeOf(s))
	fmt.Println(ns)

	v:=reflect.New(ns)

	var p *string
	var val string = "hello"
	p = &val

	v.Elem().Set(reflect.ValueOf(p))
	fmt.Println(*(*string)(unsafe.Pointer(v.Elem().Pointer())))
}
result:
*string
hello
```
