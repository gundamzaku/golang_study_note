这是根据网上的资料整理出的一个代码demo，也可以适当的了解一下unsafe的用法
代码结构
>src  
>>main  
>>>dan.go  
>>>main.go  

```go
dan.go
package main
import "fmt"

type dan struct  {

	name string
	age int8
}

func (info *dan)getName() string  {
	fmt.Println("getName:"+info.name)
	return info.name
}

func (info *dan)getAge() int8  {
	fmt.Println("getAge:",info.age)
	return info.age
}

func (info *dan)setName(name string)  {
	fmt.Println("setName:"+name)
	info.name = name
}

func (info *dan)setAge(age int8){
	fmt.Println("setAge:",age)
	info.age = age
}
```
这是一个class对象,dan就是我，我的个人信息

```go
package main
import "unsafe"

func main()  {
	info := &dan{}
	info.setName("liudan")
	info.setAge(18)

	var i *string = (*string)(unsafe.Pointer(info))
	*i = string("batman's father")


	var j *int8 = (*int8)(unsafe.Pointer(uintptr(unsafe.Pointer(info)) 
		+ uintptr(unsafe.Sizeof(string(0)))))
	*j = int8(20)
	info.getName()
	info.getAge()
}

result:
setName:liudan
setAge: 18
getName:batman's friend
getAge: 20
```
其实这段代码非常容易理解，就是绕过了dan对象的方法，直接用unsafe的指针去修改对象中的值。  
比较费解的一段应该是`uintptr(unsafe.Sizeof(string(0)))`  
其实就是在struct里面，每一个定义的变量/对象都有一个地址，而且是按顺序排列的，比如小明站在了位置1上面，那小李就只能站在位置1+1上面了。 
我去掉这一段试试  
```
setName:liudan
setAge: 18
getName:recoveryＸＸＸbad
getAge: 18
```
可以看到，名字不对了，说明小明和小李重叠了。  

不过如果只是这样的话，我用`info.name = "batman's enemy"`不也就搞定了么？何必如此大费章折？  

正当我困惑之际，突然想到一个问题。在Java体系中，对象是分private和public和protect的。  
我这里怎么都是public了。在一想，是我目录结构的问题，go把这两个文件当成是一个包的文件了，自然都是public了。  
因此，我要重新调整一下目录结构  
>src  
>>dan  
>>>dan.go  
>>main  
>>>main.go  

把`dan.go`单独拎出来做为一个独立包

现在的代码如下
```go
package dan
import "fmt"

type Dan struct  {

	name 	string
	age 	int8
	desc 	string
}

func (info *Dan)GetName() string  {
	fmt.Println("getName:"+info.name)
	return info.name
}

func (info *Dan)GetAge() int8  {
	fmt.Println("getAge:",info.age)
	return info.age
}

func (info *Dan)GetDesc() string  {
	fmt.Println("getDesc:",info.desc)
	return info.desc
}

func (info *Dan)SetName(name string)  {
	fmt.Println("setName:"+name)
	info.name = name
}

func (info *Dan)SetAge(age int8){
	fmt.Println("setAge:",age)
	info.age = age
}
```
要注意方法名的大小写问题，因为不习惯，已经在这上面吃过无数次亏了。  
`只有首字母大小才能被外界访问，成为公共变量`  

```go
package main
import (
	"unsafe"
	"dan"
)

func main()  {
	info := &dan.Dan{}
	info.SetName("liudan")
	info.SetAge(18)

	info.GetName()

	var i *string = (*string)(unsafe.Pointer(info))
	*i = string("batman's friend")

	var j *int64 = (*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(info))
		+uintptr(unsafe.Sizeof(string(0)))))
	*j = int64(18)

	var k *string = (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(info))
		+uintptr(unsafe.Sizeof(string(0)))
		+uintptr(unsafe.Sizeof(int64(0)))))
	*k = string("i will conqure the world")

	info.GetName()
	info.GetAge()
	info.GetDesc()

}

result:
setName:liudan
setAge: 18
getName:liudan
getName:batman's friend
getAge: 18
getDesc: i will conqure the world
```

这和之前的代码有什么不同呢？  
现在再试一下`info.name = "batman's enemy"`看看，对不起，不好意思，现在它已经成为私有变量，你无法修改。  
但是unsafe的方式可以改啊！  
这就是unsafe的妙处，绕开了go的限制，也因此不受go的推荐。  

uintptr到现在我还没有完全理解这个概念，下一次要学一下了。  

```
在Go语言中，还有一个专门用于存储内存地址的类型 uintptr。而 uintptr 类型与 int 类型和 uint 类型一样，也属于数组类型。它的值是一个能够保存一个指针类型值（简称指针值）的位模式形式。
```


