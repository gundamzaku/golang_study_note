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
