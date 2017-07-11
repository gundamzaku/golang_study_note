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
	b := (*[20]byte)(unsafe.Pointer(n.Info()))
	for _, v := range b {
		fmt.Println(string(v))
	}

}
```
`b := (*[20]byte)(unsafe.Pointer(n.Info()))`是为了拿到完整个的字节，否则只能拿到指针的头，因为Go不支持指针位移，所以暂时只能用这个方式取值。  
打印出来的结果是：  
b=`&[0 0 7 42 115 116 114 105 110 103 0 0 7 42 117 105 110 116 49 54]`
让我们对照着之前的newName()方法来看，第一位是b[0]是根据exported、lan(tag)和pkgPath来生成的。  
目前三个位数都不存在，所以为0。  

b[1] = uint8(len(val) >> 8)，最后得出来还是0，虽然不知道这段的意思。

b[2] = b[2] = uint8(len(val))

