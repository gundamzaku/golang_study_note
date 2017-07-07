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

我在网上找了一篇文章，关于Golang内存管理源码剖析，里面有提到  
`Golang 的内存管理基于 tcmalloc`  
还有专门的一篇文章（可惜现在被墙了，无法访问）  
http://goog-perftools.sourceforge.net/doc/tcmalloc.html
全英文的，以后慢慢研究。目前如果深入下去，恐怕后面的内容都不用看了。  

这里只简单的讲一下   
go的数组结构有三种：  
mcache: per-P cache，可以认为是 local cache。  
mcentral: 全局 cache，mcache 不够用的时候向 mcentral 申请。  
mheap: 当 mcentral 也不够用的时候，通过 mheap 向操作系统申请。  

第一种就是Small objects的存放地址

`我们知道每个 Gorontine 的运行都是绑定到一个 P 上面，mcache 是每个 P(Processer) 的 cache。这么做的好处是分配内存时不需要加锁`

大的就直接送到mheap里面去了。

我们先看这mallocgc()的方法定义  
`func mallocgc(size uintptr, typ *_type, needzero bool) unsafe.Pointer {}`
调用
`mallocgc(typ.size*uintptr(n), typ, true)`  
传入尺寸，这个尺寸做了一个乘法，`typ.size*uintptr(n)`切片变量的size是16，系统给扩展的大小是10，结果是160。至于为什么是160，现在我不知道。  
传入了切片的属性，传入了一个布尔的true（needzero)  
在代码块里做了一个if  
```go
if gcBlackenEnabled != 0 {
	// Charge the current user G for this allocation.
	assistG = getg()
	if assistG.m.curg != nil {
		assistG = assistG.m.curg
	}
	// Charge the allocation against the G. We'll account
	// for internal fragmentation at the end of mallocgc.
	assistG.gcAssistBytes -= int64(size)

	if assistG.gcAssistBytes < 0 {
		// This G is in debt. Assist the GC to correct
		// this before allocating. This must happen
		// before disabling preemption.
		gcAssistAlloc(assistG)
	}
}
```	
gcBlackenEnabled在这里是160，因此触发下面的代码块。  
assistG = getg()  
g是什么？似乎前面有到过关于G的介绍`G: Goroutine 执行的上下文环境。`并不是特别的理解。  
在代码里，也仅仅是在runtime包的stubs.go里面找到方法的定义。没有找到实现。  
不过好在还有注释  
```
// getg returns the pointer to the current g.
返回当前g的指针
// The compiler rewrites calls to this function into instructions
这编汇重写了这个方法到指示的调用
// that fetch the g directly (from TLS or from the dedicated register).
直接到达g
```
g是一个结构体，在runtime包中的runtime2.go文件里有定义type g struct {}，里面的内容太多，现在暂时不想看。先跳过吧 
在代码中能够看到`var x unsafe.Pointer`段，这个x就是最后还返回的变量，目前是无类型指针。 
下面有一段`if size <= maxSmallSize {}`这就好理解了，maxSmallSize是一个常量，表示32kb，size就是切片尺寸，小于32kb，就执行代码块里的内容。否则就写入heap里面去。  
好了，mallocgc()的内容暂时就讲到这里打住，否则也会越讲越乱，等到有机会真正接触到go的内核的时候，再回过头来研究一下。  
我们现在需要知道的是mallocgc()返回的是一个指针数据（Data）  
```go
// sliceHeader is a safe version of SliceHeader used within this package.
type sliceHeader struct {
	Data unsafe.Pointer
	Len  int
	Cap  int
}
```
在sliceHeader结构体可以看出，第一个就是mallocgo()返回的Data。这应该是go语言中对slice切片的一种结构的定义。  
之后，就将他嵌入到`Value{typ.common(), unsafe.Pointer(&s), flagIndir | flag(Slice)}`的结构体中，一并返回。  
再用Copy(t, s)的方式，我的天啊，一个Copy()也是一段非常复杂的代码。  
先不管，我们看一下Copy中的t，t就是前面讲的那一堆乱七八糟的内容最后返回的一个sliceHeader{}，s则是老的切片。  

```go
// Copy copies the contents of src into dst until either
// dst has been filled or src has been exhausted.
// It returns the number of elements copied.
// Dst and src each must have kind Slice or Array, and
// dst and src must have the same element type.
func Copy(dst, src Value) int {
```
