学习完了reflect.Append()方法之后，其实我还有不少的问题。  
首先一个问题就是，reflect.Append()是为一个切片添加内容，那我为什么不用原生的Slice切片来添加呢？原生的应该更简单吧。  

```go
func main() {
	slice := []string{"Mon","Tues","Wed","Thur","Fri"}
	newSlice := append(slice,"sat","Sun");
	fmt.Println(newSlice)
}
```
同样的方法名，完全实现了一样的功能。那它和`reflect.Append()`有什么区别呢？  
append()方法在builtin包中的builtin.go文件中被声明，
`func append(slice []Type, elems ...Type) []Type`

可惜我找遍了整个src目录，都没有找到append()的实现方法，看网上有提到`总的意思是我在源码找到的 builtin.go 只是一个文档，并不是源代码，这些函数是经过特殊处理的，我们不能直接去使用他们。`，按这个说法的话，这个方法对我们而言，是不透明的了。这大概也是反射和原生的一个区别。
不过看注释的时候，注意到了有趣的一点说明：  
```go
// As a special case, it is legal to append a string to a byte slice, like this:
作为一个特别的安全，将字符串追加到byte类型的切片中是合法的。
// slice = append([]byte("hello "), "world"...)
```
很有意思的小技艺，能记的话就记一下。  

既然原生的方法没得看了，那我把目光再返回到reflect中。
`网上有人论证过，原生的append()性能远大于reflect.append()`

在所有的方法中，我又看到一个叫`reflect.AppendSlice()`的方法。看这名字，似乎和slice又有点关系，那到底是做什么的呢？

先看注释：  
```
// AppendSlice appends a slice t to a slice s and returns the resulting slice.
AppendSlice就将一个slice追加到一个slice中并且返回slice结果的方法
// The slices s and t must have the same element type.
这个slices中的s和t必须是同样的类型元素
```
这么一看就懂了，之前的append()是将一个任意的元素类型(string,int之类）添加到slice当中。而这一个AppendSlice则是将slice添加到slice当中。  
`注意：这个方法只允许两个传入参数，而不是之前的append()一样可以有多个传入参数`

我写一段代码进行测试，看上去并没有什么新意。  
```go
slice := []string{"Mon","Tues","Wed","Thur","Fri"}

sliceExt := []string{"sat","sun"}

newSlice := reflect.AppendSlice(slice,sliceExt)
fmt.Println(newSlice)

result:
.\mian.go:39: cannot use slice (type []string) as type reflect.Value in argument to reflect.AppendSlice
.\mian.go:39: cannot use sliceExt (type []string) as type reflect.Value in argument to reflect.AppendSlice
```

报错了……呃……和之前一样，传入的参数必须是valueof()转化过的。  
最后结果就是：  
```
[Mon Tues Wed Thur Fri sat sun]
```

再看一眼这个方法的源代码：  
```
func AppendSlice(s, t Value) Value {
	s.mustBe(Slice)
	t.mustBe(Slice)
	typesMustMatch("reflect.AppendSlice", s.Type().Elem(), t.Type().Elem())
	s, i0, i1 := grow(s, t.Len())
	Copy(s.Slice(i0, i1), t)
	return s
}
```
两个mustBe来检查变量的类型，grow是划分出一个新的slice对象，然后把要追加的slice拷贝到s.Slice(i0, i1)之中，并且返回。  
当中有一段
`typesMustMatch("reflect.AppendSlice", s.Type().Elem(), t.Type().Elem())`
这是干什么的呢？  
```go
func typesMustMatch(what string, t1, t2 Type) {
	if t1 != t2 {
		panic(what + ": " + t1.String() + " != " + t2.String())
	}
}
```
注意：what string, t1, t2 Type，这里t1没有指定类型，是Go的一种写未能，其实是表示t1,t2是一种类型。  
这个方法非常简单，是验证s和t是否是同一种类型，可是上面已经有过mustbe了，这里是不是有点多此一举了？  
其实不是，这里的验证类型，主要是验证slice元素内部的数据类型。  
比如我slice{里面的元素是string}，sliceExt{里面的元素是int}  
就会收到报错信息`panic: reflect.AppendSlice: string != int`  

好了，这样一来`reflect.AppendSlice()`方法的用处我们也知道了。
到目前为止，我们大抵上了解了反射中的四个方法：  
```
reflect.Append()
reflect.ValueOf()
reflect.AppendSlice()
reflect.typeof()
```

节奏一下子快了很多，实际上，在append()的方法当中，我们一直有用到一个Copy()的方法，其实这个方法也是可以对外使用的，是公有方法之一。即reflect.Copy()。  
那么，既然如此，我们就趁热打铁，看一下Copy()方法的使用。在上面的代码中，可以看到`Copy(s.Slice(i0, i1), t)`  
而这个方法的实现则是`func Copy(dst, src Value) int {}`  
dst和src都是Value数据类型。至于Copy()的定义，之前贴过，这里再贴一下，反正也不要钱。  
```go
// Copy copies the contents of src into dst until either
复制src的内容到dst
// dst has been filled or src has been exhausted.
直到dst被填满或者src用完
// It returns the number of elements copied.
返回复制元素的个数
// Dst and src each must have kind Slice or Array, and
dst和src必须符合slice或array
// dst and src must have the same element type.
dst和src也必须是同样的元素类型
func Copy(dst, src Value) int {}
```
完整代码：(根据前面的slice的实例来分析一下）

首先，dst接受到的是我传入第一个参数是s.Slice(i0, i1)，i0为5，i1为7，即s这个切片在5-7区间的空间。  
第二个参数是t，就是新的slice，是准备填到5-7这个空间去的。  
```go
func Copy(dst, src Value) int {
	dk := dst.kind()  //开始验证类型，23,妥妥的slice类型
	if dk != Array && dk != Slice {//既然是slice了，必然通过
		panic(&ValueError{"reflect.Copy", dk})
	}
	if dk == Array {//不是数组，这一步跳过（关于数组，这里先不讲，后面再讲）
		dst.mustBeAssignable()
	}
	dst.mustBeExported()//必须是输出的？
	//if f records that the value was obtained using an unexported field.
	//如果f记录的值是非输出领域得到的？不太理解
	
	//下面与dk雷同
	sk := src.kind()
	if sk != Array && sk != Slice {
		panic(&ValueError{"reflect.Copy", sk})
	}
	src.mustBeExported()
	
	//验证dst和src内部的元素类型，必须一样
	de := dst.typ.Elem()
	se := src.typ.Elem()
	typesMustMatch("reflect.Copy", de, se)

	var ds, ss sliceHeader
	if dk == Array {
		ds.Data = dst.ptr
		ds.Len = dst.Len()
		ds.Cap = ds.Len
	} else {
		//ds指到sliceHeader（故且先这么理解吧）
		ds = *(*sliceHeader)(dst.ptr)
	}
	if sk == Array {
		ss.Data = src.ptr
		ss.Len = src.Len()
		ss.Cap = ss.Len
	} else {
		//ss指到sliceHeader（故且先这么理解吧）
		ss = *(*sliceHeader)(src.ptr)
	}

	return typedslicecopy(de.common(), ds, ss)
}
```

最后，调用了下面这个方法，返回结果。
`func typedslicecopy(elemType *rtype, dst, src sliceHeader) int`
同样的，这个方法在runtime包中的mbarrier.go中有实现的方法，可以参考一下，这里就不多讲了，反正讲半天也是一知半解。  
返回的是int类型，告诉结果是成功还是失败。

以上是我拿之前的slice的例子进行的调试，看Copy()的源码时可以看到是支持Array数组的，那么为了单独测试一下reflect.Copy()的用法，我就找一个Array的例子吧。  

```go
var a = [5]int {1, 2, 3, 4, 5}
var b = [5]int {6, 7, 8, 9, 0}
c := reflect.Copy(reflect.ValueOf(a),reflect.ValueOf(b))
fmt.Println(a)
fmt.Println(b)
fmt.Println(c)

result:
panic: reflect: reflect.Copy using unaddressable value
```

执行了一下……唔，报错。  
错误出在`func (f flag) mustBeAssignable() {}`方法之中。
```go
if f&flagAddr == 0 {
	panic("reflect: " + methodName() + " using unaddressable value")
}
```
这是怎么回事，似乎是f&flagAddr变成了0？
f和flagAddr分别是145和256，两者的与运算确实是0。
