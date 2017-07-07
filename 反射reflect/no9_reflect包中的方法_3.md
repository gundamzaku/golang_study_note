承接上一篇，再回到`unsafe_NewArray(typ.Elem().(*rtype), cap)`中来  
unsafe_NewArray()这个方法在value.go中有声明  
`func unsafe_NewArray(*rtype, int) unsafe.Pointer`
