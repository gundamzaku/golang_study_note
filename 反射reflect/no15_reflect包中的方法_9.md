这一篇开头，先来了解一下Name这个对象的具体定义。  

name是一个对象，他声明了一个变量`bytes *byte`，是比特类型的指针。在上一节中，我们用rs:=newName("[]string","", "", false)就是创建了这么一个比特数据。

之前也讲到了，他是一个变量的一些定义参数，最主要的有三个属性`tag, pkgPath, exported`
