我的sb通过代码`sc := reflect.SliceOf(sb)`转成了sc  
接下来就要通过reflect.MakeSlice()把sc变成真正能用的slice类型。
mySlice := reflect.MakeSlice(sc,0,10)
