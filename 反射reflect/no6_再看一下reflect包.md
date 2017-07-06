
在reflect包中，定义了一组常量
```go
type Kind uint
const (
	Invalid Kind = iota
	Bool
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	Array
	Chan
	Func
	Interface
	Map
	Ptr
	Slice
	String
	Struct
	UnsafePointer
)
```
这组常量都是Kind类型（即uint类型），iota值为0，这个是在builtin包中定义的一个常量  
`builtin包是go的预声明定义,包括go语言中常用的各种类型和方法声明,包括变量和常量两部分.`  
`const iota = 0 // Untyped int.`  
按照go语言特殊的定义，从Invalid向下，就是不停的++的过程，Invalid为0，Bool为1，Int为2……以此类推。  
可是我在`fmt.Println(reflect.Int)`的时候，显示的result为`int`，而不是数字，这是怎么回事？
问题出在下面一段代码的定义上面
```go
var kindNames = []string{
	Invalid:       "invalid",
	Bool:          "bool",
	Int:           "int",
	Int8:          "int8",
	Int16:         "int16",
	Int32:         "int32",
	Int64:         "int64",
	Uint:          "uint",
	Uint8:         "uint8",
	Uint16:        "uint16",
	Uint32:        "uint32",
	Uint64:        "uint64",
	Uintptr:       "uintptr",
	Float32:       "float32",
	Float64:       "float64",
	Complex64:     "complex64",
	Complex128:    "complex128",
	Array:         "array",
	Chan:          "chan",
	Func:          "func",
	Interface:     "interface",
	Map:           "map",
	Ptr:           "ptr",
	Slice:         "slice",
	String:        "string",
	Struct:        "struct",
	UnsafePointer: "unsafe.Pointer",
}
```
在value.go的文件里面，go将kindNames和最上面的常量串在了一起。只是单纯的定义了这个kindNames并没有实际的意思
```go
func (k Kind) String() string {
	if int(k) < len(kindNames) {
		return kindNames[k]
	}
	return "kind" + strconv.Itoa(int(k))
}
```
真正将两者串起来的是这个String()的隐函数。之所以称之为隐函数，是因为它控制了fmt.println()
我现在`fmt.println(Int8)`，实际上是传了一个值8到String()里面去，代码判定8小于kindNames的总长度，因此，将`Uint8:"uint8"`返回
如果我传一个超长的值，const (UnsafePointer)后面我再加一个常量如Unknow，这个常量在KindNames中并没有，然后再试着println  
`result:kind27`,解决了`return "kind" + strconv.Itoa(int(k))`的方法
`将整数转换为十进制字符串形式（即：FormatInt(i, 10) 的简写）`,这样，这组常量的定义算是明白了。  

```go
// ChanDir represents a channel type's direction.
type ChanDir int

const (
	RecvDir ChanDir             = 1 << iota // <-chan
	SendDir                                 // chan<-
	BothDir = RecvDir | SendDir             // chan
)
```
看到chan，就可以断定这组常量和go的通道有关系了。
很明显，这组常量和
```go
func (d ChanDir) String() string {
	switch d {
	case SendDir:
		return "chan<-"
	case RecvDir:
		return "<-chan"
	case BothDir:
		return "chan"
	}
	return "ChanDir" + strconv.Itoa(int(d))
}
```
也对应上了，功能效果与上面的一样。  
目前并不知道用在哪里，先跳过。
再有一组
```go
type SelectDir int
const (
	_             SelectDir = iota
	SelectSend              // case Chan <- Send
	SelectRecv              // case <-Chan:
	SelectDefault           // default
)
```
这一次，这组常量老老实实地用了int，没有做任何转义。  
这三组就是reflect里面对外的三种类型的常量了，我们都可以直接调用。
分别是Kind，ChanDir和SelectDir类型。  

接下来，再看一下reflect的自定义类型，一共十一个
分别是：
```go
type Type interface {}
type ChanDir int
type Kind uint
type Method struct {}
type Value struct {}
type SelectCase struct {}
type SelectDir int
type SliceHeader struct {}
type StringHeader struct {}
type StructField struct {}
type StructTag string
type ValueError struct {}
```

最后是他的方法，一共21个，有点多。下面一个一个来看到底是干什么用的吧。  
首先我先看之前提到过的reflect.typeof()方法，这是获得变量类型的方法。我看了一下源代码，却发现和我之前使用1.5的源代码是完全不一样的。  
这让我感觉有些意外，是因为我从1.5升级到1.8的关系，还是因为mac机和win机系统不一样的关系？（我是mac和win混着用的）  
这样一来，我要重新梳理一下。  
可以对比一下代码的不同之处  
首先看rtype类型  
原先的rtype
```go
type rtype struct {
	size          uintptr
	ptrdata       uintptr
	hash          uint32         // hash of type; avoids computation in hash tables
	alg           *typeAlg       // algorithm table
	gcdata        *byte          // garbage collection data
	string        *string        // string form; unnecessary but undeniably useful
	*uncommonType                // (relatively) uncommon fields
	kind          uint8

}
```
现在的rtype
```go
type rtype struct {
	size       uintptr
	ptrdata    uintptr
	hash       uint32   // hash of type; avoids computation in hash tables
	tflag      tflag    // extra type information flags
	align      uint8    // alignment of variable with this type
	fieldAlign uint8    // alignment of struct field with this type
	kind       uint8    // enumeration for C
	alg        *typeAlg // algorithm table
	gcdata     *byte    // garbage collection data
	str        nameOff  // string form
	ptrToThis  typeOff  // type for pointer to this type, may be zero
}
```
关键的`*uncommonType`没有了
在看一下原先的string()方法
```go
func (t *rtype) String() string { return *t.string }
```
现在的string()方法
```go
func (t *rtype) String() string {
	s := t.nameOff(t.str).name()
	if t.tflag&tflagExtraStar != 0 {
		return s[1:]
	}
	return s
}
```
怎么一下子复杂了这么多？
直接从String()入手，首先，它产生了一个t.str的变量。
我拿
```go
func main() {
	var val string = "hello";
	fmt.Println(reflect.TypeOf(val))
}
```
来做例子吧

`package runtime主要是与go的runtime系统进行互动操作，例如控制goroutine的函数等。它也包含reflect package所需的低等级信息。`

按顺序执行，到`s := t.nameOff(t.str).name()`的时候，我们将s打印出来，发现是`string`，看来这个nameOff()就是获取变量的信息的。

`小技巧：修改go/src包中的代码，用go build -a 编译才可以生效，否则不行`

之前我是用的抠代码的战术，现在不用了，实在是太费力，如果这样用下去，我要看到牛年马月了，所以现在直接修改src包的代码。
我再看看t.nameOff(t.str)到底做了什么？
t.str得到的值是5424，这是什么意思？或者说这代表的是字符串？为了释疑，我又试着将`var val string = "hello";`的类型换成了整形，得到939；换成了布尔类型得到了3224。
目前我无法确认它们的意思，不过似乎是每一种类型会对应一种数字。先标记一下，以后说不定会有地方碰到。

`return name{(*byte)(resolveNameOff(unsafe.Pointer(t), int32(off)))}`
nameoff()返回的是一个结构体，经过resolveNameOff返回一个无类型的指针，然后又将这个指针转换化`(*byte)`类型。

`resolveNameOff方法的实现非常复杂，并且在runtime包的type.go方法里面，因此这里略过。`

type name struct {}是一个对象，并且下面有五个方法继承
```go
func (n name) data(off int) *byte {}
func (n name) isExported() bool {}
func (n name) nameLen() int {}
func (n name) tagLen() int {}
func (n name) name() (s string) {}
func (n name) tag() (s string) {}
func (n name) pkgPath() string {}
```
看到这里，有点晕，但是毫无疑问的是，这个name的对象，保存的就是变量的一些基础信息  
再来看resolveNameOff()这个方法，在type.go文件里，这个方法并没有实现，只有一段定义。
`func resolveNameOff(ptrInModule unsafe.Pointer, off int32) unsafe.Pointer`
这就奇怪了，那它到底是怎么实现的？在程序里面，不可能存在只有定义没有实现但是可以执行的代码吧。
还好上面有一段注释`mplemented in the runtime package.`
根据注释，在runtime包里面成功到找到了一个type.go的文件。  
至于为什么会有两个type.go，我不知道，不过目前可以肯定的是，resolveNameOff()在runtime包里面实现了。因此再看一下这个具体的resolveNameOff代码，相当复杂，应该调用了不少底层的东西。
```go
func resolveNameOff(ptrInModule unsafe.Pointer, off nameOff) name {
	if off == 0 {
		return name{}
	}
	base := uintptr(ptrInModule)
	for md := &firstmoduledata; md != nil; md = md.next {
		if base >= md.types && base < md.etypes {
			res := md.types + uintptr(off)
			if res > md.etypes {
				println("runtime: nameOff", hex(off), "out of range", hex(md.types), "-", hex(md.etypes))
				throw("runtime: name offset out of range")
			}
			return name{(*byte)(unsafe.Pointer(res))}
		}
	}
	//后面省略，还有一长串代码…………
}
```
经过断点，发现在代码执行到一半的时候就返回了`return name{(*byte)(unsafe.Pointer(res))}`，松了一口气，要知道下面还有一大串的代码。  
这段代码第一步是把rtype的指针地址转了过来，经过uintptr的转化，变成一个整形。  

然后在循环中的firstmoduledata是一个变量定义
var firstmoduledata moduledata  // linker symbol
而moduledata又是一个很长的结构体，感觉么玩下去要疯啊……
```go
type moduledata struct {
	pclntable    []byte
	ftab         []functab
	filetab      []uint32
	findfunctab  uintptr
	minpc, maxpc uintptr

	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr
	types, etypes         uintptr

	textsectmap []textsect
	typelinks   []int32 // offsets from types
	itablinks   []*itab

	ptab []ptabEntry

	pluginpath string
	pkghashes  []modulehash

	modulename   string
	modulehashes []modulehash

	gcdatamask, gcbssmask bitvector

	typemap map[typeOff]*_type // offset to *_rtype in previous module

	next *moduledata
}
看注释的意思：
moduledata records information about the layout of the executable image
moduledata记录了关于可执行的映象的层的信息
moduledata is stored in read-only memory; none of the pointers here are visible to the garbage collector.
moduledata是存在只读内存里面的，对垃圾回收器来说没有指针可见
好吧……再次表示看不太懂。
```
看看这段`for md := &firstmoduledata; md != nil; md = md.next {}`都做了什么  
moduledata是通过指针链的关系串起来的，因此这里的for就是判断`moduledata里面的next *moduledata`存不存在，存在的话一直遍历。  
很明显，这里不存在，只遍历了一次。接下来，base就要同里面的属性进行对比。

在函数中主要用到了md.types和md.etypes两个变量，我以String类型作为参数打印了一下这两个值的信息  
md.types:4754432  
md.etypes:4972988  
接下来，我再用int参数和boole参数再测试一下  
int类型  
md.types:4754432  
md.etypes:4972956  
boole类型  
md.types:4754432  
md.etypes:4972956  

base这个变量目前是4807776,介于两个数据之间，也使得`if base >= md.types && base < md.etypes {}`这个条件成立。
说明了md.types是变量类型的下限，md.etypes是变量类型的上限，变量类型这个信息是不会超过两者之外的。  
等一下，根据前面学到的知识base存放的是指针的整形，那么md.types和md.etypes岂不一样？

最后还原为： 0x488c60-0x4be23d
`Go 有指针，但是没有指针运算。你不能用指针变量遍历字符串的各个字节。`

`return name{(*byte)(unsafe.Pointer(res))}`
返回的name{}是一个对象，要注意的是，这个name{}对象仍然是reflect包中的type.go文件里的，而非runtime包里面的，尽管runtime包中也存在  
我们把焦点定位在这个name{}对象中来。  
```go
func (n name) name() (s string) {

	if n.bytes == nil {
		return
	}
	b := (*[4]byte)(unsafe.Pointer(n.bytes))

	hdr := (*stringHeader)(unsafe.Pointer(&s))
	hdr.Data = unsafe.Pointer(&b[3])
	hdr.Len = int(b[1])<<8 | int(b[2])
	return s
}
```
获得变量的类型的name是通过这个方法
`b := (*[4]byte)(unsafe.Pointer(n.bytes))`是把n.bytes指针转换成了一个4位数组的指针。  
然后声明出一个string类型的变量s，接着把hdr指到了s的指针地址上面  
然后分别赋值Data和Len，然后返回，这时候println(s)发现result就是string，，这是什么鬼？  
我把这段代码单独抽出来测试
然而抽出来的代码无法执行，走到`for next := &firstmoduledata; next != nil; next = next.next {}`的时候，一直拿不到firstmoduledata的值。
md.types和md.etypes取到的全是0。一点用也没。  
折腾了我半天，突然前面在厕所小便的时候灵光一闪。  
把type.go放在runtime里面应该不是没有理由的。  
runtime可以理解成go的虚拟机不是？  
虚拟机在运行的时候，有些变量在内存里是动态的不是？  
所以type.go里面的内容，很可能是跟着虚拟机走，动态分配的，不能抽出来。  
为了验证这个道理，我在runtime/type.go里面加入一段代码  
```go
func GetFm() {
	println(*&firstmoduledata.types)
	println(*&firstmoduledata.etypes)
	//return &firstmoduledata.types
}
```
然后在我的main()方法中调用，主要是为了看数据。  
果然，这两个值一直在变，这两个值是指针，也就是说地址一直在变，于是我把这两个变动的值直接写死到我代码中（不一定每次成功，我调了几次）  
```
string
4757024
4982027
*string
```
终于出现了我要的数值。

大体上我也明白了go是怎么操作反射去取一些变量或方法的信息的了。  

不过很可惜，这篇学习比较惨，因为我没有能深入地了解到go的本质上面，一些概念仍然一知半解。同时我也发现，越是往低层的东西，其实我们可以修改的越是有限。  

下面是我抽出来的代码，里面有一些冗余，但是不想改了   

```go
runtime/type.go里面改造，追加下面两个方法
func GetTypes() uintptr{
	return *&firstmoduledata.types
}
func GetETypes() uintptr{
	return *&firstmoduledata.etypes
}
```

```go
package main

import (
	"fmt"
	"unsafe"
	"runtime"
)

type name struct {
	bytes *byte
}
type stringHeader struct {
	Data unsafe.Pointer
	Len  int
}

func (n name) name() (s string) {
	if n.bytes == nil {
		return
	}
	b := (*[4]byte)(unsafe.Pointer(n.bytes))
	hdr := (*stringHeader)(unsafe.Pointer(&s))
	hdr.Data = unsafe.Pointer(&b[3])
	hdr.Len = int(b[1]) << 8 | int(b[2])
	println(s)
	return s
}

type nameOff int32
type functab struct {
	entry   uintptr
	funcoff uintptr
}
type textsect struct {
	vaddr    uintptr // prelinked section vaddr
	length   uintptr // section length
	baseaddr uintptr // relocated section address
}

type itab struct {
	inter  *interfacetype
	_type  *_type
	link   *itab
	bad    int32
	inhash int32      // has this itab been added to hash?
	fun    [1]uintptr // variable sized
}

type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod
}
type imethod struct {
	name nameOff
	ityp typeOff
}
type ptabEntry struct {
	name nameOff
	typ  typeOff
}
type modulehash struct {
	modulename   string
	linktimehash string
	runtimehash  *string
}
type bitvector struct {
	n        int32 // # of bits
	bytedata *uint8
}
type moduledata struct {
	pclntable    []byte
	ftab         []functab
	filetab      []uint32
	findfunctab  uintptr
	minpc, maxpc uintptr
	text, etext           uintptr
	noptrdata, enoptrdata uintptr
	data, edata           uintptr
	bss, ebss             uintptr
	noptrbss, enoptrbss   uintptr
	end, gcdata, gcbss    uintptr
	types, etypes         uintptr
	textsectmap []textsect
	typelinks   []int32 // offsets from types
	itablinks   []*itab
	ptab []ptabEntry
	pluginpath string
	pkghashes  []modulehash
	modulename   string
	modulehashes []modulehash
	gcdatamask, gcbssmask bitvector
	typemap map[typeOff]*_type // offset to *_rtype in previous module
	next *moduledata
}
type typeAlg struct {
	hash  func(unsafe.Pointer, uintptr) uintptr
	equal func(unsafe.Pointer, unsafe.Pointer) bool
}

var fm moduledata  // linker symbol
type tflag uint8
type typeOff int32 // offset to an *rtype
/*
type rtype struct {
	size       uintptr
	ptrdata    uintptr
	hash       uint32   // hash of type; avoids computation in hash tables
	tflag      tflag    // extra type information flags
	align      uint8    // alignment of variable with this type
	fieldAlign uint8    // alignment of struct field with this type
	kind       uint8    // enumeration for C
	alg        *typeAlg // algorithm table
	gcdata     *byte    // garbage collection data
	str        nameOff  // string form
	ptrToThis  typeOff  // type for pointer to this type, may be zero
}*/
type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldalign uint8
	kind       uint8
	alg        *typeAlg
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}
/*
func (t *rtype) nameOff(off nameOff) name {
	return name{(*byte)(resolveNameOff(unsafe.Pointer(t), off))}
}*/

func (t *_type) nameOff(off nameOff) name {
	return rOff(unsafe.Pointer(t), off)
}

func rOff(ptrInModule unsafe.Pointer, off nameOff) name {
	if off == 0 {
		return name{}
	}
	base := uintptr(ptrInModule)
	//println(*&fm.etypes)
	getTypes := runtime.GetTypes()
	GetETypes := runtime.GetETypes()
	for md := &fm; md != nil; md = md.next {
		//这里是写死的数据，仅作参考
		md.types = getTypes;
		md.etypes = GetETypes;

		if base >= md.types && base < md.etypes {
			res := md.types + uintptr(off)
			if res > md.etypes {
				println("error")
			}
			return name{(*byte)(unsafe.Pointer(res))}
		}
	}
	return name{}
}

type Type interface{}
type emptyInterface struct {
	typ  *_type
	word unsafe.Pointer
}

func toType(t *_type) Type {
	if t == nil {
		return nil
	}
	return t
}

func TypeOf(i interface{}) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}

const (
	tflagUncommon tflag = 1 << 0
	tflagExtraStar tflag = 1 << 1
	tflagNamed tflag = 1 << 2
)

func (t *_type) String() string {
	s := t.nameOff(t.str).name()
	if t.tflag & tflagExtraStar != 0 {
		return s[1:]
	}
	return s
}

func main() {
	var s string = "hello"
	//var s int = 100
	fmt.Println(TypeOf(s))
}
```
