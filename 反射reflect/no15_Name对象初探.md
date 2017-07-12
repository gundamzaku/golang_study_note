## Name对象初探

这一篇开头，先来了解一下Name这个对象的具体定义。  

name是一个对象，他声明了一个变量`bytes *byte`，是比特类型的指针。在上一节中，我们用rs:=newName("[]string","", "", false)就是创建了这么一个比特数据。

之前也讲到了，他是一个变量的一些定义参数，最主要的有三个属性`tag, pkgPath, exported`

先拉一下name对象里的方法：
```go
func (n name) data(off int) *byte {}
func (n name) isExported() bool {}
func (n name) nameLen() int {}
func (n name) tagLen() int {}
func (n name) name() (s string) {}
func (n name) tag() (s string) {}
func (n name) pkgPath() string {}
```

一共有七个方法，而name的生成则有：  
```go
func newName(n, tag, pkgPath string, exported bool) name {}
```
可以操作，这个在前面已经讲过。

那么name的读取呢？在rptye{}中我们可以找到一个str对应nameOff类型，而nameOff其实是一个方法，返回得正是name对象。  
```go
func (t *rtype) nameOff(off nameOff) name {}
```

现在来做一个试验。  

因为在reflect包中，一些特有的方法是私有的，所以我把源码做了一点改造。  

首先：
type Type interface {}
在Type接口中增加一个方法描述
Info() name //返回一个name类型  

接着，我在实现这个接口的type name struct {}下面增加这个方法的实现：  
```go
func (n name) Info() *byte {
	var B *byte
	B = n.bytes
	return B
}
```
用大写的Info，确保我在外部能够调用。

最后我用代码实现  
```
func main()  {
	var val string
	val = "hell world"

	rs:=reflect.TypeOf(val)
	n:= rs.Info()
	b := (*[9]byte)(unsafe.Pointer(n.Info()))
	for _, v := range b {
		fmt.Println(string(v))
	}

}
```
`b := (*[9]byte)(unsafe.Pointer(n.Info()))`是为了拿到完整个的字节，否则只能拿到指针的头，因为Go不支持指针位移，所以暂时只能用这个方式取值。  
打印出来的结果是：  
b=`&[0 0 7 42 115 116 114 105 110 103]`
让我们对照着之前的newName()方法来看，第一位是b[0]是根据exported、lan(tag)和pkgPath来生成的。  
目前三个位数都不存在，所以为0。  

`*n is string`

b[1] = uint8(len(n) >> 8)，最后得出来还是0，虽然不知道这段的意思。

b[2] = uint8(len(n))
这里是7，实际上送进去的类型等于`*string`，正好7位  
后面的[b3:]-9是英文的描述，我们用string()转义一下可以看到，正好是`*string`  

这就是对之前newName()的知识的一些补充，可是有几个问题仍然悬而未解，`tag, pkgPath, exported`到底是什么意思？  
实现上我模拟了好几个变量定义，最后的结果b[0]均是0，完全没有突破。  
失之东隅，收之桑榆，尽管在变量的声明上面没有找到发现，不过在type struce{}定义中，我们可以看到tag和pkgPath的意思。  
```go
package main

import (
	"reflect"
	"fmt"
)

func main()  {
	s := struct{
		name string "this is tags"
	}{"dan"}

	t := reflect.TypeOf(s)
	fmt.Printf("Name: %q, PkgPath: %q\n", t.Name(), t.PkgPath())
	fmt.Printf("Name: %q, PkgPath: %q, Tag: %q\n", t.Field(0).Name, t.Field(0).PkgPath, t.Field(0).Tag)
}

result:
Name: "", PkgPath: ""
Name: "name", PkgPath: "main", Tag: "this is tags"
```

可以看到，在分析结构体内部的变量的时候，tags是对变量的标签(注：似乎只有结构体里的变量有这个属性)，而pkgPath则是这个变量对应在哪个包的包名。  
可惜的是，我没有办法用我之前的办法，把这两个定义在字体中体现出来。这也许是结构体特有的属性，exported根据字面上的意思，结果现在的包的变量有大小写的区别，首字大写的变量方法才有可能被其它的包引用，而小写的变量方法则是私有的。我想，可能exported就和这个有关系吧。  

话休絮烦，这一块的内容就打此为此了，暂时我也无法再深入下去，或许将来的某一天，我还会再次碰到。
