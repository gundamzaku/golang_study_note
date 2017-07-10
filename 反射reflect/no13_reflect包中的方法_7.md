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
