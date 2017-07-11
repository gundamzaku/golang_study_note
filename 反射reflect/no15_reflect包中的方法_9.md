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
