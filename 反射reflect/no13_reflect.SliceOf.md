## func SliceOf(t Type) Type {}

在`func SliceOf(t Type) Type {}`中，我们可以看到它较多地用到一种新的知识，叫做'cache'。  

首先，从传参来看，他要求我们传入的是经过reflect.typeof()转化后的变量类型（比如string,int)，而他返回的是经过再次转化后，变成了切片的变量类型（仍然是Type，但是打印的时候可以看出，变成了[]string,[]int了）

我们以一段代码为例：
```go
func main()  {
	var my string = "hello"
	sb := reflect.TypeOf(my)
	fmt.Println(sb)
	sc := reflect.SliceOf(sb)
	fmt.Println(sc)
}
result:
string
[]string
```

可以看到返回结果的变化，并且，继续用这段代码来解释reflect.SliceOf(sb)的作用。  
老规矩，第一行代码中，拿到了sb的rtype集合。  
`typ := t.(*rtype)`  
接下来，就用到了cache的结构 
`ckey := cacheKey{Slice, typ, nil, 0}`
```go
type cacheKey struct {
	kind  Kind
	t1    *rtype
	t2    *rtype
	extra uintptr
}
```
可以看出，cacheKey存放了{变量类型,sb的属性，空，0}  
接下来，我们可以通过cacheGet(ckey)的方法得到一个新的Type类型变量。这里恐怕变得复杂起来了。先看一下cacheGet()的方法。  
```go
func cacheGet(k cacheKey) Type {
	lookupCache.RLock()
	t := lookupCache.m[k]
	lookupCache.RUnlock()
	if t != nil {
		return t
	}

	lookupCache.Lock()
	t = lookupCache.m[k]
	if t != nil {
		lookupCache.Unlock()
		return t
	}

	if lookupCache.m == nil {
		lookupCache.m = make(map[cacheKey]*rtype)
	}

	return nil
}
```
可以看到，用到了锁的机制。  
首先用到了`lookupCache.RLock()`   
```go
// The lookupCache caches ArrayOf, ChanOf, MapOf and SliceOf lookups.
var lookupCache struct {
	sync.RWMutex
	m map[cacheKey]*rtype
}
```
lookupCache这是一个结构体，他包含了  
sync.RWMutex  
还有一个被称为m的集合，程序执行到目前，m是不应该存在值的，可以先忽略。而sync.RWMutex则是sync包中的rwmutex.go里面的内容了。真是一环扣一环，可是sync暂时还不在这次的学习内容里面。先看一下他的定义吧，首先，嗯，这是一个对象。其次，他和读写锁有关系。  

先看一段定义`golang中sync包实现了两种锁Mutex （互斥锁）和RWMutex（读写锁），其中RWMutex是基于Mutex实现的，只读锁的实现使用类似引用计数器的功能．`  
再看一段定义：
```
RWMutex提供四个方法：
func (*RWMutex) Lock //写锁定
func (*RWMutex) Unlock //写解锁
func (*RWMutex) RLock //读锁定
func (*RWMutex) RUnlock //读解锁
```
这么看看，就突然感觉简单了不少，lookupCache.RLock()，把读锁定了，我读了，其它人就不能读了，要等我读完。  

然后我从lookupCache的m集合中根据传入的cacheKey来取值，然后把读锁给解开。

不过前面已经说了，现在是空的，取不到。  

所以又转入下一步，lookupCache.Lock()，把写锁定，现在只有我能写，其它人都等着。  

然后我再从lookupCache的m集合中去取一次（双保险？）确保我没取到数据。  

安全以后，我把sb的变量属性写进map集合里面，`lookupCache.m = make(map[cacheKey]*rtype)`

收工。至于为什么要做这一步，我现在还不知道。继续看代码。  

回到SliceOf(t Type)方法  
```go
if slice := cacheGet(ckey); slice != nil {
	return slice
}
```
取到值，直接返回，因为已经有了。没取到，继续走下去。  

`s := "[]" + typ.String()`，硬拼出一个字符串，果然够霸道，也就是说我前面传入的string类型，这里转为字符串，然后前面再加上"[]"，就凑成了一个"[]string"的字符串了。  

调用typesByString()方法，又是一个较长的方法，先看一下实现了什么吧。

`for _, tt := range typesByString(s) {}`

`*range 用来遍历数组和切片的时候返回索引和元素值`

英文盲强行再看注释
```
// typesByString returns the subslice of typelinks() whose elements have
// the given string representation.
typesByString返回typelinks()的子切片？若这个元素具有给定的字符串表示
// It may be empty (no known types with that string) or may have
// multiple elements (multiple types with that string).
它可能为空（那个字符串具有未知类型）或为多个元素（那个字符串具有多类型）
```
好吧，看不懂。。 先看看我的sb传进去是什么结果。  
[1/1]0xc042004030  
好像是意味着只有一个元素。  

先看一下typesByString()方法中的第一行`sections, offset := typelinks()`  
typelinks()，完了，又是一个在go代码里面找不到的内置方法  
只能看注释：  
```go
// typelinks is implemented in package runtime.
在runtime包中实现（奇怪，我怎么没找到，难道是func typelinksinit() {}方法？）
// It returns a slice of the sections in each module,
返回切片，基于每个模块的部件
// and a slice of *rtype offsets in each module.
和每个模块的rtype地址
//
// The types in each module are sorted by string. That is, the first
每个模块的类型被存成字符
// two linked types of the first module are:
//
//	d0 := sections[0]
//	t1 := (*rtype)(add(d0, offset[0][0]))
//	t2 := (*rtype)(add(d0, offset[0][1]))
//
// and
//
//	t1.String() < t2.String()
//
得，这个例子看得我一脸懵逼
// Note that strings are not unique identifiers for types:
// there can be more than one with a given string.
// Only types we might want to look up are included:
// pointers, channels, maps, slices, and arrays.
func typelinks() (sections []unsafe.Pointer, offset [][]int32)
```

好吧，我表示看了也是白看，还不如想想到底是怎么实现的吧。
typelinks返回的是一个无类型的sections地址，一个int32的多维数组。我多写几种例子来测一下有什么区别吧。  

先用原来传入的sb看一下。  
```
sections, offset := typelinks()
println(sections)
println(offset)
result:
[1/1]0xc042050018
[1/1]0xc04203e400
```
咦，不是一个指针一个int吗，怎么两个都是指针了，好吧，我表示很晕，不过在注释上又写着`and a slice of *rtype offsets in each module`，那应该其实返回的就是rtype的地址吧。  

我把我的sb的目标变量改成`var my = [3]string{"a","b","c"}`，仍然是单一的。  
```
result:
[1/1]0xc042050018
[1/1]0xc04203e400
```
很可惜，因为自己的一知半解，在这块上面无法获知更多的信息，谈了几种方式，最后的结果都是[1/1]，或许是用在别的地方的吧。  

我在假设这个方法不透明的情况下，继续执行下去。  
```go
for offsI, offs := range offset {
	//奇怪，offs变成了`[716/716]0x4bf140`,offsI则是0 	
	section := sections[offsI] //section等于0x489420

	i, j := 0, len(offs)	//i是0，j是716，难道要做716次循环？
	for i < j {
		//还好，只执行了9次，每次都要折半，分别是358,537,448,493,471,482,488,491,492
		h := i + (j-i)/2 // avoid overflow when computing h
		/*
		似乎搞清了这里取出来的是什么东西了。
		358=[16]uint32 大于等于s
		537=[]strconv.leftCheat
		448=[8192]uint16 大于等于s
		493=[][32]*runtime._defer
		471=[8]uint32 大于等于s
		482=[]*runtime._type 大于等于s
		488=[]*runtime.mspan 大于等于s
		491=[]*runtime.timer 大于等于s
		492=[]*sync.Pool
		然而我还是不明白这是什么意思
		*/
		if !(rtypeOff(section, offs[h]).String() >= s) { //rtypeOff()方法是做了地址的偏移
			i = h + 1 // preserves f(i-1) == false
		} else {
			j = h // preserves f(j) == true
		}
	}
	//最后i与j均为492，小于716（len），继续做加法，难道要做300多次？
	for j := i; j < len(offs); j++ {
		typ := rtypeOff(section, offs[j])
		if typ.String() != s {//typ.string()为[]*sync.Pool,与s不相等，直接退出（break）
			break
		}
		//ret是一个[]*rtype，存放切片的，把这里的[]*sync.Pool存进去然后返回。
		ret = append(ret, typ)
	}
}
return ret
```

看完了，感觉非常晕，也完全没有看懂，事实上最后返回的ret也是一个空的定义。  
所以在`func SliceOf(t Type) Type {}`中这整个一段都没有执行。
```go
for _, tt := range typesByString(s) {
	slice := (*sliceType)(unsafe.Pointer(tt))
	if slice.elem == typ {
		return cachePut(ckey, tt)
	}
}
```	
我在目前为止，可能还是无法完全理解这其中所包含的原理，不过从上面打出的sync.Pool来看，有可能和Go的pool包有关系。  

`Go 1.3 的sync包中加入一个新特性：Pool。这个类设计的目的是用来保存和复用临时对象，以减少内存分配，降低CG压力。`
更多的解释  
```
我们可以把sync.Pool类型值看作是存放可被重复使用的值的容器。此类容器是自动伸缩的、高效的，同时也是并发安全的。为了描述方便，我们也会把sync.Pool类型的值称为临时对象池，而把存于其中的值称为对象值。
```
在此处先做一个引申，暂且跳过，或许在不久的学习过程中，就会重新碰到并了解其实际的应用价值。  

还是回到原来的方法。  

```go
// Make a slice type.
//产生一个空的指针？
var islice interface{} = ([]unsafe.Pointer)(nil)
//然后把这个指针声明成slice类型
prototype := *(**sliceType)(unsafe.Pointer(&islice))
slice := *prototype//分配给slice?
slice.tflag = 0
//先是newName()，再是resolveReflectName(),都是相当重的方法，很难理解。
slice.str = resolveReflectName(newName(s, "", "", false))
//用了一种fnv1的hash生成方式
//FNV能快速hash大量数据并保持较小的冲突率，它的高度分散使它适用于hash一些非常相近的字符串，比如URL，hostname，文件名，text，IP地址等。
slice.hash = fnv1(typ.hash, '[')
slice.elem = typ
slice.ptrToThis = 0
```
真是没想到创建一个变量需要这么多的操作和内容，呼呼，就此打住了。如果深入下去真的没有底了。  
最后`return cachePut(ckey, &slice.rtype)`先存入缓存，再返回。(记得解锁）  
到这里，reflect.sliceOf()算是完了，可是我对这个方法的具体流程仍然处于一知半解的状态，非常遗憾。  
