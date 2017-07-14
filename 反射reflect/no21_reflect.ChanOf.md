## func ChanOf(dir ChanDir, t Type) Type {}

到这里的时候，Make系列也只剩下一个Chan了，并且随着复杂度的提高，这里越来越多的只能讲使用上的技巧，而非原理。 

既然有MakeChan()，必然就有ChanOf(),先看一下ChanOf()的使用方式。  

首先，它要我们传入一个ChanDir的参数，这又是什么？  

在reflect包的type.go文件里面，有这么一个定义：  
`type ChanDir int`
