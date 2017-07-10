在`func SliceOf(t Type) Type {}`中，我们可以看到它较多地用到一种新的知识，叫做'cache'。  

首先，从传参来看，他要求我们传入的是经过reflect.typeof()转化后的变量类型（比如string,int)，而他返回的是经过再次转化后，变成了切片的变量类型（仍然是Type，但是打印的时候可以看出，变成了[]string,[]int了）
