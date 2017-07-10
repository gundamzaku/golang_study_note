前面不知不觉又说了一大堆的废话，忘记了初心了。  

我们要确定的一点是，Value()这是一个对象，这个对象有很多关于变量属性的信息。这个对象是通过reflact.valueof()方法来产生的，并且不会对原值造成影响，只是追加了一个flag。

具体产生Value的方法是在valueof()里面的unpackEface()。  
在new Value()的时候，共创建三个参数。t, e.word, f  
t是从e.typ的赋值，而e.typ则是将变量强行转成无参指针，然后重定义为`*emptyInterface`类型。  

在
```go
type emptyInterface struct {
	typ  *rtype
	word unsafe.Pointer
}
```
里面，word就是这个变量对应的地址，现在是无参指针类型。  
typ就是rtype结构，这在之前都全部讲过。  

至于他们具体是怎么对应的，这是Go的内部约定，目前不了解底层的我还不是特别清晰。  

接着，e.word我们也基本知道是什么东西了，然后就是f，f作为flag，在变成Value类型的时候，是会被重写的。  
重写的规则也很简单，变量类型所对应的数字，和常量flagIndir做一个或的操作，当然，这是建立在ifaceIndir(t)的基础之上（必须为true）  
ifaceIndir(t)做了什么？他将变量类型所对应的数字和常量kindDirectIface进行了与的操作，并且等于0的话，认为是true。

这里面一下子就涉及到两个常量flagIndir和kindDirectIface，这两个常量具体什么用，代码里没有明确说明，不过也可以猜到一点点。  

比如这个kindDirectIface，他的值是32，想必是约定kind的边界（1-28，预留4位），不能超过这个范围，否则就是违规。  

flagIndir应该是下标，他的是值是128，往往和变量类型所对应的数字相加得到新的flag值。

我挑几个常用的来展示一下：

int: 130  
string: 152  
array: 145  
slice: 151  
boolen: 129  
应该都是它们特定的标志。  

基本上，到这里应该有个定论，就是所有的变量都可以重定义为interface{}类型，而所有的变量，都可以转成Value类型。

复习完这些知识以后，我们继续我们的代码之旅……  

趁热打铁，我将目标盯上了reflect.MakeSlice()这个方法，其实在reflect里面，带Make字眼的方法一共有四个。
```go
reflect.MakeSlice()
reflect.MakeChan()
reflect.MakeFunc()
reflect.MakeMap()
```

既然是叫Make了，那肯定就是制造/创建，Slice是Go的切片，Chan是Go的管道，Func是Go的方法，Map是Go的集合。  
其实看字面的意思我们已经可以知道这四个方法是干什么用的了，但是具体是怎么用的呢？先从reflect.MakeSlice()看起吧。  

先看方法：  
```go
func MakeSlice(typ Type, len, cap int) Value {
	if typ.Kind() != Slice {
		panic("reflect.MakeSlice of non-slice type")
	}
	if len < 0 {
		panic("reflect.MakeSlice: negative len")
	}
	if cap < 0 {
		panic("reflect.MakeSlice: negative cap")
	}
	if len > cap {
		panic("reflect.MakeSlice: len > cap")
	}

	s := sliceHeader{unsafe_NewArray(typ.Elem().(*rtype), cap), len, cap}
	return Value{typ.common(), unsafe.Pointer(&s), flagIndir | flag(Slice)}
}
```

抛去前面的验证，下面的两个调用的方法我们在之前的reflect.AppendSlice()里面已经全部接触过了，都不难理解。这里让我们传入三个参数`typ Type, len, cap int`

len是长度，cap是容量，都是int类型。比较费解的是这个typ，是Type对象。Type对象可是在type.go文件里面的一个类，具体我要怎么把它生成出来？  
似乎也不难，只要`var typ = new(reflect.Type)`一下不就可以了？  
好像一切都很简单一样，于是我试着用这个方式去MakeSlice一下。  
```go
func main()  {

	var typ = new(reflect.Type)
	mySlice := reflect.MakeSlice(typ,10,5)
	fmt.Println(mySlice)
}

result:
cannot use typ (type *reflect.Type) as type reflect.Type in argument to reflect.MakeSlice:
*reflect.Type is pointer to interface, not interface
```

果然没有这么一帆风顺的事情，看错误的提示，似乎和指针有关系。  
我用`fmt.Println(reflect.TypeOf(typ))`打印了一下，果然是`*reflect.Type`类型，而非是方法传递指明要用的Type类型。那怎么解决？  

我在网上找到了一种方案：  
```go
type s struct {

}

func main()  {
	sc := reflect.SliceOf(reflect.TypeOf(&s{}))
	fmt.Println(sc)

	mySlice := reflect.MakeSlice(sc,5,10)
	fmt.Println(mySlice)
}

result:
*main.s
[]*main.s
[<nil> <nil> <nil> <nil> <nil>]
```
可以看出，首先，要创建一个type（我怎么就没想到的。。）,然后用reflect.TypeOf(&s{})转换成Value，然后再用reflect.SliceOf(X)转换，最后就变成了可以传入reflect.MakeSlice()的参数。  

这里有几点困惑的地方，首先，为什么要传入引用（&s[]）；其次，为什么要用reflect.TypeOf()，再次，为什么要用reflect.SliceOf()  

其实我不传引用，这段程序同样可以执行。  
```go
func main()  {

	sb := reflect.TypeOf(s{})
	fmt.Println(sb)
	sc := reflect.SliceOf(sb)
	fmt.Println(sc)

	mySlice := reflect.MakeSlice(sc,5,10)
	fmt.Println(mySlice)
}
result:
main.s
[]main.s
[{} {} {} {} {}]
```
