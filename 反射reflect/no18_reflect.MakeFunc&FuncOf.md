
## func FuncOf(in, out []Type, variadic bool) Type {}
传入in,out，均为[]Type类型的切片，variadic为布尔
返回Type变量

在学习reflect.MakeFunc()之前，根据以往的惯例，有Make必有Of，所以很理所当然地我们找到了reflect.FuncOf()方法。  
reflect.FuncOf()的源代码很长，在此处就不帖出来了。比较难理解的是，之前的xxOf()都是让传一个Type类型，现在这里是[]Type。

怎么理解呢，我们先想办法写一个例子看一下比较好。  
