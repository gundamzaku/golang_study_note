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
