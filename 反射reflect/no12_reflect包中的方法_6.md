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

