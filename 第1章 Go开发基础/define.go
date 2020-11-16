声明变量的一般形式是使用 var 关键字：
var name type
或者 
name :=""
num :=123
array :=

var arrInt= []int{2,4,6}


var arrString = []string{"2","4","6"}

 //1、创建
    var g1 map[int]string //默认值为nil
    g2 := map[int]string{}
    g3 := make(map[int]string)
    g4 := make(map[int]string, 10)
    fmt.Println(g1, g2, g3, g4, len(g4))

    //2、初始化
    var m1 map[int]string = map[int]string{1: "ck_god", 2: "god_girl"}
    m2 := map[int]string{1: "ck_god", 2: "god_girl"}
    fmt.Println(m1, m2)