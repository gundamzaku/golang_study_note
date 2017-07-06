### GO中的包  

`下面的代码都是基于go1.5来写的，实际上现在我安装了go1.8以后，发现reflect中的代码已经重新改过结构了。详细的可见第（6）篇`

`type.go`  
`reflect.TypeOf`方法  
```
var i int = 123
fmt.Println(reflect.TypeOf(i))
```
出来的结果是`int`

看一下源代码做了什么操作
```go
func TypeOf(i interface{}) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}
```
首先传入的是一个`interface{}`,空interface(`interface{}`)不包含任何的method，正因为如此，所有的类型都实现了空interface。  
顺着这一点，我决定新建一个包，将reflect包中的文件精简出来，看一下TypeOf到底是怎么样运作的。  
在reflact包中，最最重要的就是`type.go`和`value.go`两个文件。  
现在我新建一个新的包，命名为selfReflact包，并新建一个`type.go`文件。  

通过TypeOf，我扒出了另外一个方法toType，然后又找到emptyInterface的结构，具体的精简代码可以参考我的后面的源码，此不复述。  
值得注意的是`unsafe.Pointer()`这个方法，这是go的另外一个包unsafe里的内容，在这里不做深入，仅了解它的意思。  
顾名思义，这是一个不安全的方法，他在此处的功能就是强制进行指针类型的转换，而这在平时的go语言中是不被允许的，前面看过一篇文章，是讲unsafe这个包的不实用性，实际上可以看出，在go语言包库中还是有一定用处的。  
了解到这一点以后，就可以看出，上面的语句就是把我们传进来的i这个变量的地址强行转成了emptyInterface结构并做成了一个指针对象。  

```go
type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
}
```
emptyInterface由typ和word两个元素组成，word是一个unsafe.Pointer类型，目前他到底代表什么意思？我不得而知。  
而typ则是一个rtype的结构  

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
在rtype又有一个uncommonType的结构体  
```go
type uncommonType struct {
	name    *string  // name of type
	pkgPath *string  // import path; nil for built-in types like int, string
	methods []method // methods associated with type
}
```
这一来一去，事情就变得复杂了。他们到底代表什么意思？虽然源码中还有英文的注释，但是究竟他们之间是如何运作起来的，我还是不清楚。  
目前可以肯定的是`eface := *(*emptyInterface)(unsafe.Pointer(&i))`这段代码执行以后。  
现在fmt.println(eface.typ.uncommonType.name)`(eface.typ.name)`也可以，他会自动跳过uncommonType，形成一个缩写。  
我决定单独试一下这段代码  
于是我把这段代码单独拎了出来  
```go
package main

import (
	"unsafe"
	"fmt"
)

type typeAlg struct {
	hash func(unsafe.Pointer, uintptr) uintptr
	equal func(unsafe.Pointer, unsafe.Pointer) bool
}

type spcommonType struct {
	name    *string  // name of type
}
type supertype struct {
	size          uintptr
	ptrdata       uintptr
	hash          uint32         // hash of type; avoids computation in hash tables
	alg           *typeAlg       // algorithm table
	gcdata        *byte          // garbage collection data
	string        *string        // string form; unnecessary but undeniably useful
	*spcommonType                // (relatively) uncommon fields
}

type toConvInterface struct {
	typ  *supertype
	word unsafe.Pointer
}

func main()  {
	var i interface{} = 900
	eface := *(*toConvInterface)(unsafe.Pointer(&i))
	fmt.Println(eface.typ.name)
}
```
现在执行，发现打印出来的是0x113b90，怎么回事？一看这个就是个内存地址了，`fmt.Println(*eface.typ.name)`是我忘记加上指针了。  
现在可以成功的显示`int`  
我试着想再次精简代码，发现删除任何一个地方，都会使程序无法执行，于是就到此为至。  
目前暂时不想去研究这块的深理，故且我认为这是unsafe.Pointer进行强制转类型以后的某种特点的格式吧。在这段DEMO中我有意将方法的名字全部重新定义了一下，发现仍然可以执行，那显然是和结构有关系，而和方法名称无关。  

接下来，再回到原来的程序中去。我在main方法里面调用了TypeOf(i)方法  
```go
package main


import (
	"fmt"
	"selfReflect"
)

func main()  {
	var i int = 123
	fmt.Println(selfReflect.TypeOf(i))    //int

}
```
理论上来讲，应该也显示出`Int`才对吧。事实上并没有，而是打印出了一串我看不懂的东西`&{8 0 4149441018 0x191630 0x135098 0x114040 0xbfec8 0}`  
这是什么回事，typeof方法和reflect中的源代码方法比较了一下，并没有什么问题  
```go
func TypeOf(i interface{}) Type {
	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}

func toType(t *rtype) Type {
	if t == nil {
		return nil
	}
	return t
}

```
那问题会出在哪里？？  

经过无法用言语来形容的缜密排查之后，我终于灵感大开，发现一段不起眼的代码  
`func (t *rtype) String() string { return *t.string }`
难道是……  
我把这段代码和上面的代码合并了一下  
```go
func toType(t *rtype) Type {
	if t == nil {
		return nil
	}
	return *t.string
}
```
return的内容作了调整，意外的发现成功了，现在是成功打印出`int`  

那么那段代码到底做了什么事情呢？从字面上来看，他是rtype的一个字方法（go用struct来做面向对象真的让人很不习惯）  
然后它重写了t。  
这样说有点拗口，或者说他重定义了打印t时，t显示的内容。  
这点需要我再起一段代码验证一下。  

```go
package main

import (
	"fmt"
)

type supertype struct {
	string        *string
}

func (t *supertype) SetT() {
	var s string = "hello go"
	t.string = &s
}

func (t *supertype) String() string {
	return *t.string
}

func main()  {
	t := supertype{}
	t.SetT()
	s:=&t
	fmt.Println(s)
}
```
然后运行一下，咦……打印出了"hello go"，这么看来还真被我给蒙对了，String()方法是隐性的。  

这样一来，所有的问题和疑惑就完全解开啦〜可以安心收工了。  

```go
package selfReflect
import (
	"unsafe"
)

type Type interface {
	//String() string

}

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

type typeAlg struct {
	hash func(unsafe.Pointer, uintptr) uintptr
	equal func(unsafe.Pointer, unsafe.Pointer) bool
}

type method struct {
	name    *string        // name of method
	pkgPath *string        // nil for exported Names; otherwise import path
	mtyp    *rtype         // method type (without receiver)
	typ     *rtype         // .(*FuncType) underneath (with receiver)
	ifn     unsafe.Pointer // fn used in interface call (one-word receiver)
	tfn     unsafe.Pointer // fn used for normal method call
}

type uncommonType struct {
	name    *string  // name of type
	pkgPath *string  // import path; nil for built-in types like int, string
	methods []method // methods associated with type
}

type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
}

func (t *rtype) String() string { return *t.string }

//空interface(interface{})不包含任何的method，正因为如此，所有的类型都实现了空interface。
func TypeOf(i interface{}) Type {

	//在ＧＯ中有内存安全性的限制，于是下面的转换是不被允许的
	//两个不同指针类型的值，例如 int64和 float64。
	//指针类型和uintptr的值。
	//unsafe.Pointer，我们可以打破Go类型和内存安全性，并使上面的转换成为可能
	//这里是把i强制转换成了emptyInterface指针
	//unsafe里只有几个定义没有具体实现？
	//fmt.Println(unsafe.Pointer(&i))

	eface := *(*emptyInterface)(unsafe.Pointer(&i))
	return toType(eface.typ)
}

func toType(t *rtype) Type {
	if t == nil {
		return nil
	}
	return t
}
```
main.go
```go
package main


import (
	"fmt"
	"selfReflect"
)

func main()  {
	var i int = 123
	fmt.Println(selfReflect.TypeOf(i))    //int

}
```
