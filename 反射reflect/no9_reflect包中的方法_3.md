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

```go
// newarray allocates an array of n elements of type typ.
func newarray(typ *_type, n int) unsafe.Pointer {
	if n < 0 || uintptr(n) > maxSliceCap(typ.size) {
		panic(plainError("runtime: allocation size out of range"))
	}
	return mallocgc(typ.size*uintptr(n), typ, true)
}
```

最后指到mallocgc()方法上面，这方法一看就晕了，太长。根据以前对C知识的了解，malloc，肯定是开内存去了。  
先看一下注释吧。  
```
// Allocate an object of size bytes.
分配一个比特大小的目标
// Small objects are allocated from the per-P cache's free lists.
小目标被分配到per-P 缓存的空闲列表中
// Large objects (> 32 kB) are allocated straight from the heap.
大目标（大小32KB）的，被分配到直接堆中？
```
这下搞大了，直接渗透到go的底层来了。
