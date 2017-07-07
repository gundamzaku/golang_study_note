承接上一篇，再回到`unsafe_NewArray(typ.Elem().(*rtype), cap)`中来  
unsafe_NewArray()这个方法在value.go中有声明  
`func unsafe_NewArray(*rtype, int) unsafe.Pointer`

但是却没有具体的实现！经过前面几章的训练，很容易就能猜到，一定又是放在runtime包里头的value.go里面了。  
看一下名字我就猜到这个方法的作用，不安全的……嗯创建一个新的数组。  

不过这一次，很明显我想错了。这方法是在runtime包里面的malloc.go里面。  
```go
//go:linkname reflect_unsafe_NewArray reflect.unsafe_NewArray
func reflect_unsafe_NewArray(typ *_type, n int) unsafe.Pointer {
	return newarray(typ, n)
}
```
看这一段代码的描述，reflect_unsafe_NewArray与reflect.unsafe_NewArray做了链接。
