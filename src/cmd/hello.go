//go入口函数main包中的main方法
//不同的文件可以有相同的package
package main

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"session"
	"strconv"
	"strings"
	"subpackage"
)

var sessionManager *session.SessionManager

//go执行顺序：import -> const -> var -> init() -> main()
func init() {
	fmt.Println("this is in main")
	subpackage.M()
	subpackage.M2()
	fmt.Println("----------------------")
	var err error
	//不能使用:=，:=会创建局部变量
	sessionManager, err = session.NewManager("memory", "sessionId", 300)
	if err != nil {
		fmt.Println(err)
	}
	go sessionManager.GC()
}

//使用别名
//import std "fmt"

/*
go语言保留字
break default func interface select
case defer go map struct
chan else goto package switch
const fallthrough if range type
continue for import return var 

*/

func m1() {
	fmt.Println(math.Pi)
}

func m2(a int) {
	fmt.Println(a)
}

//多返回值
func add(a, b int) (int, int) {
	return a + b, a * b
}

//命名返回值参数，return后不需要指定值，会返回c、d
func swap(a, b string) (c, d string) {
	c = b
	d = a
	return
}

//不定长参数，必须是最后一个参数，参数为slice的值拷贝，直接传递slice也是拷贝，拷贝的地址
func m3(a ...int) {
	for i := 0; i < len(a); i++ {
		fmt.Printf("%d", a[i])
	}
	fmt.Printf("\n")
	fmt.Printf("%v\n", a)
}

//函数作为值、类型
func m3a() {
	//函数作为值
	a := m1
	//m2参数类型和m1不同，不能再赋值给a
	//a = m2
	a()
	a = func() {
		fmt.Printf("匿名函数:%v\n", a)
	}
	a()
	n := closure(1)(2)
	fmt.Println(n)
	//函数作为类型
	type testMethod func(int) bool
	c := func(num int, m testMethod) {
		if m(num) {
			fmt.Println("aaaaaaaaaa")
		} else {
			fmt.Println("bbbbbbbbbb")
		}
	}
	b := func(num int) bool {
		return num%2 != 0
	}
	c(1, b)
}

//返回值为函数
func closure(a int) func(int) int {
	return func(b int) int {
		//此处的a与closure方法传进来的a为同一个地址，不是值的拷贝
		fmt.Printf("%d\n", a)
		return a + b
	}
}

//defer，函数执行结束后（return后）执行，后进先出的调用规则，
func m3b() {
	defer fmt.Println("a")
	defer fmt.Println("b")
	fmt.Println("c")
	a := 1
	defer func(a int) {
		fmt.Println("defer A", a)
	}(a) //此时a已经拷贝
	a = 2
	defer func() {
		fmt.Println("defer B", a) //直接使用外部的a
	}()
	a = 3
	func() {
		//defer需要提前放到堆栈中，否则defer没有定义就返回了
		defer func() {
			if recover() != nil {
				fmt.Println("start recover")
			}
		}()
		m3c()
	}()
	fmt.Println("after painc")
}

//错误处理，panic可以在任何地方引发，recover只能在defer中使用。panic和recover都是内置函数
func m3c() {
	//panic发生时，直接返回调用初，继续执行painc，直到遇到defer中的recover
	panic("painc in m3c()")
	fmt.Println("this never print")
}

//传值，传指针
func m3d() {
	a := 1
	m3d1(a)
	fmt.Println(a)
	m3d2(&a)
	fmt.Println(a)
}

func m3d1(a int) {
	a = 2
}

func m3d2(ap *int) {
	*ap = 2
}

//全局变量
var aa int = 100
var bb string = "bb string"

//系统推断类型
var cc = false

//全局变量不能使用最简方式
//c1 := true

//int8 int32 int64 float32（小数点后7位） float64（小数点后15位）
//byte 和 int8 类型一样
var (
	dd bool = true
	ee int8 = 2
	ff uint = 23
	gg byte = 10
)

//类型别名
type (
	字符串 string
)

var aString 字符串 = "类型别名测试"

//for，go语言循环只有for，没有while
func m4() {

	sum := 0
	for i := 0; i <= 100; i++ {
		sum += i
	}
	fmt.Println(sum)

	sum2 := 1
	for sum2 < 100 {
		sum2++
	}
	fmt.Println(sum2)

	var a1 int = 1
	for {
		a1++
		if a1 > 10 {
			break
		}
	}
	fmt.Println(a1)
}

//if 没有括号
func m5() {
	if aa < 100 {
		fmt.Println("aa小于100")
	} else {
		fmt.Println("aa大于或等于100")
	}

	//可以先初始化一个值，作用域为if...else块
	if tt := 100; tt < 100 {
		fmt.Println("tt小于100")
	} else {
		fmt.Println("tt大于或等于100")
	}
	//tt作用域为if...else块中，此处不能使用
	//fmt.Println(tt)
}

//首字母大写，其他包可以调用
type Student struct {
	name  string
	email string
	int   //匿名字段使用类型作为字段名称，相当于int int
}

type Human struct {
	name string
	age  int
}

//自定义类型
type Skills []string

type Stu struct {
	Human     //继承Humen的字段和方法
	stuNumber string
	Skills
}

//为struct定义method：func 接受者 methodName(参数) 返回值{}
//接收者为指针，可以修改原值
func (s *Stu) say() int {
	fmt.Println("I am student:" + s.name)
	s.age = 100 //无需*s.age = 100，当使用指针访问字段(指针并没有字段)，会自动处理
	return 1
}

//接收者不是指针，对原值修改无效
func (h Human) doJob() {
	h.age = 9999 //对原值无效
	fmt.Println("do my job:" + h.name)
}

//重载Stu继承自humen的方法
func (s Stu) doJob() {
	fmt.Println("do study:" + s.name)
}

//struct使用
func m6() {
	//声明:按顺序提供值
	var stu = Student{"zhangsan", "zhangsan@126.com", 12}
	fmt.Println(stu)
	stu.name = "lisi"
	fmt.Println(stu)
	stu = Student{name: "lisi", email: "lisi@126.com"}
	fmt.Println(stu)
	stu.int = 100 //修改匿名字段值
	fmt.Println(stu)
	//使用new声明
	// stu2 := new(Student)
	// fmt.Println(stu2)
	//匿名struct
	var t = struct {
		name string
	}{name: "lisi"}
	fmt.Println(t)
}

//继承
func m7() {
	//声明
	var stu = Stu{stuNumber: "lisi", Human: Human{name: "lisi", age: 12}}
	fmt.Println(stu)
	//访问继承的字段
	stu.Human.age = 13
	stu.name = "zhangsan"
	fmt.Println(stu)
	//访问自定义类型
	stu.Skills = []string{"aaa", "bbb"}
	fmt.Println(stu)
	//调用struct的method，当method接受者为指针时，不需要&stu.say()，go语言自动转为指针
	stu.say()
	fmt.Println(stu)
	stu.doJob()
	fmt.Println(stu)
}

//定义interface
type Men interface {
	doJob()
}

//interface继承
type Person interface {
	Men
	personSay()
}

//定义空接口，任何类型都能赋予空接口类型
type emptyInterface interface{}

//interface使用
func m7a() {
	h := Human{name: "zhangsan"}
	s := Stu{}
	s.name = "lisi"
	//声明m为Men接口类型
	var m Men
	m = h //Humen实现了接口m的方法，所以m能存储Humen类型
	m.doJob()
	m = s
	m.doJob()
	//判断类型if方式
	if val, ok := m.(Human); ok {
		fmt.Printf("m is Human:%v\n", val)
	}
	if val, ok := m.(Stu); ok {
		fmt.Printf("m is Stu:%v\n", val)
	}
	//判断类型switch方式
	switch val := m.(type) {
	case Human:
		fmt.Printf("m is Human:%v\n", val)
	case Stu:
		fmt.Printf("m is Stu:%v\n", val)
	}
}

//map使用
func m8() {
	//map的key必须支持==或!=，不能是map、slice、函数。
	//make(map[keyType]valType,cap)，cap可省略
	var mymap map[string]Student = make(map[string]Student, 10)
	mymap["aaa"] = Student{"zhangsan", "zhangsan@126.com", 12}
	mymap["bbb"] = Student{"lisi", "lisi@126.com", 13}
	fmt.Println(mymap["bbb"].email)
	fmt.Printf("len:%d,%v\n", len(mymap), mymap)
	//直接赋值初始化
	mymap2 := map[int]string{1: "one", 2: "two", 3: "three", 4: "four", 5: "five"}
	fmt.Printf("%+v\n", mymap2)
	fmt.Printf("%v\n", mymap2)
	//使用delete删除map中的key
	delete(mymap, "aaa")
	fmt.Printf("%v\n", mymap)
	//多级map需全部初始化才能使用
	map1 := make(map[int]map[int]string)
	//必须初始化才能使用
	map1[1] = make(map[int]string)
	map1[1][1] = "str"
	fmt.Println(map1[1][1])
	//判断key-value是否存在
	_, ok := map1[2][1]
	fmt.Printf("%v\n", ok)
	//使用range迭代
	for key, val := range mymap {
		//key、val均是拷贝，不会影响map中的值
		val.name = "newName"
		fmt.Printf("key:%v,value:%v", key, mymap[key])
	}
}

//数组使用
func m9() {
	a := [5]int{1, 2, 3, 4, 5}
	fmt.Println(a)
	//长度作为类型的一部分，长度不同的数组不是一种类型。
	var b [3]int = [3]int{}
	//索引赋值方式
	b = [3]int{0: 1, 1: 2}
	fmt.Println(b)
	//使用...推断类型
	c := [...]string{2: "aaa"}
	fmt.Println(c)
	//指向数组的指针
	d := [...]int{1, 2, 0}
	var dp *[3]int = &d
	fmt.Println(dp)
	//使用new直接返回数组的指针
	dp = new([3]int)
	fmt.Println(dp)
	//指针数组
	e, f := "eeee", "ffff"
	g := [...]*string{&e, &f}
	fmt.Println(g)
	//数组可以用==和!=比较是否相等(必须类型相同)，大小不能比较
	fmt.Println(b == d)
	//可以使用下标直接操作数组内元素，指向数组的指针同样可以使用下标操作数组元素
	fmt.Println(a[1])
	fmt.Println(dp[2])
	//多维数组
	var h [2][3]int = [2][3]int{{0, 1, 2}, {7, 8, 9}}
	fmt.Println(h)
	//修改数组元素
	a1 := [2]int{1, 2}
	fmt.Printf("a1=%v\n", a1)
	a1[0] = 2
	fmt.Printf("a1=%v\n", a1)
}

//slice
/*
本身不是数组，指向底层数组；
引用类型

*/
func m9a() {
	var a = [2]int{1, 2} //创建数组
	var s = []int{1, 2}  //创建slice
	// s = append(a, 11, 1, 1)//数组不能append
	s = append(s, 111, 1, 1, 1)
	fmt.Printf("%v%v\n", a, s)
	//创建不指定个数，和数组区别
	var s1 []int
	fmt.Println(s1)
	//使用数组创建，a1[开始索引,终止索引]，开始索引、终止索引可以省略
	a1 := [10]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	s1 = a1[2:4] //[2,3]
	fmt.Println(s1)
	//slice分配连续内存，使用数组创建，len为(终止索引-开始索引)，cap为(数组len-开始索引)。
	fmt.Printf("len:%d,cap:%d\n", len(s1), cap(s1))
	//不能读取超过len的内容，尽管没有超过cap
	//fmt.Printf("%d\n", s1[2])//index out of range
	//reslice，可以在cap长度内reslice，不能超过cap
	s2 := s1[1:5]
	fmt.Println(s2) //[3,4,5,6]
	//使用make声明slice，make([]类型,length,cap)，不指定cap
	s1 = make([]int, 3, 6)
	fmt.Printf("len:%d,cap:%d\n", len(s1), cap(s1))
	//append
	fmt.Printf("%v,%p\n", s1, s1)
	s2 = append(s1, 1, 2, 3) //追加元素，没有超过cap，返回原来的slice
	fmt.Printf("%v,%p\n", s2, s2)
	s2 = append(s1, 1, 2, 3, 4, 5) //追加元素，超过cap，创建新的slice返回
	fmt.Printf("%v,%p\n", s2, s2)
	//slice为引用类型
	a2 := [5]int{1, 2, 3, 4, 5}
	s1 = a2[0:3]
	s2 = a2[2:]
	fmt.Printf("%v,%v,%v\n", a2, s1, s2)
	s1[2] = 9 //修改后会导致a2,s2都修改。注意append操作可能创建新的slice，a2、s2将不会改变
	fmt.Printf("%v,%v,%v\n", a2, s1, s2)
	//copy(to,from)，copy操作只会copy最小len个元素，to的len不变
	//copy不会修改地址，只会修改值
	s1 = []int{1, 2, 3, 4, 5, 6}
	s2 = []int{7, 8, 9}
	fmt.Printf("%p,%p\n", s1, s2)
	copy(s1, s2)
	fmt.Printf("%p,%p\n", s1, s2)
	fmt.Printf("%v,%v\n", s1, s2)
	//迭代
	for i, j := range s1 {
		fmt.Printf("%d:%d\n", i, j)
	}
}

//switch使用
/*
可以使用任何类型或表达式作为条件语句
不需要显示使用break，一旦条件满足自动终止，如希望继续执行检查，需使用fallthrough语句
switch后可以没有条件表达式，在case后使用表达式
*/
func m10() {
	a := 4
	switch a {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3, 4, 5:
		fmt.Println("three four five")
	default:
		fmt.Println("Error")
	}

	//定义变量
	switch b := 0; b {
	case 0:
		fmt.Println(b)
	default:
		fmt.Println("b != 0")
	}
	//switch中定义的变量作用域在switch
	//fmt.Println(b)

	switch c := 4; {
	case c < 4:
		fmt.Println("c < 4")
	case c > 3:
		fmt.Println("c > 3")
		fallthrough //使用fallthrough，继续执行不跳出
	case c == 4:
		fmt.Println("c == 4")
		fallthrough
	default:
		fmt.Println("default")
	}
}

//强制类型转换
func m11() {
	var a float32 = 1.3
	var b int = int(a)
	fmt.Println(b)
	b = 65
	c := string(b)
	fmt.Println(c)
	//数字转为字符串
	fmt.Println(strconv.Itoa(b))
	//字符串转为数字
	d, _ := strconv.Atoi("123")
	fmt.Println(d)
}

//常量
const PI float32 = 3.1415

//常量组
//常量组中不提供初始化值的，表示使用上行的表达式
//iota是常量计数器，遇到const置零，每定义一个常量自动加一
const (
	one = iota
	two
	three
)

func m12() {
	fmt.Println(PI)
	fmt.Println(one)
	fmt.Println(two)
	fmt.Println(three)
}

//运算符
//从左至右结合
/*
优先级（高-->低）
^ !
* / % 

*/
/*
位运算符
 6: 0110
11: 1011
&   0010  与运算
|   1111  或运算
^   1101  非
&^  0100  如果第二个运算数为1就会把第一个数改为0，否则不变
*/
func m13() {
	//取反
	fmt.Println(!true)
	//左移，右移
	fmt.Println(1 << 10)
}

//指针不能运算，操作符&取变量地址，*通过指针间接访问目标对象，默认值为nil，不是NULL
func m14() {
	var a int = 1
	var ap *int = &a
	fmt.Println(ap)
	fmt.Println(*ap)
}

//goto break continue
func m15() {
	i := 0
	//LABEL内容先执行一遍
LABEL:
	{
		i++
		fmt.Println("in label")
	}
	for i == 1 {
		fmt.Println("i == 1")
		goto LABEL
	}
	fmt.Println("end")
}

//字符串使用
func m16() {
	s := "aaa" +
		"bbb"
	s = `
	     aaa
	     bbb
	     ccc
	     `
	fmt.Println(s)
	s = "abcdefghij"
	fmt.Println(s[1:])
}

//错误类型
func m17() {
	err := errors.New("this is a error")
	fmt.Println(err)
}

//make new 区别
func m18() {
	//make只能用于内建类型map、slice、channel的内存分配，new可以用于各种类型的内存分配
	//make(T,args)返回T类型
	//new(T,args)返回*T类型（指针）
	var s []int
	fmt.Println(s)
	s = make([]int, 3)
	fmt.Println(s)

}

//反射
func m19() {
	var stu = Student{name: "zhangsan"}
	t := reflect.TypeOf(stu)
	fmt.Println(t)
	v := reflect.ValueOf(stu)
	fmt.Println(v)
}

//并发
func m20() {
	go func() {
		for i := 0; i < 5; i++ {
			runtime.Gosched() //cpu把时间片让给其他goroutines
			fmt.Println("Hello")
		}
	}()
	fmt.Println("World") //当前goroutines执行
}

//简单http服务器
func m21() {
	//监听地址及回调
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm() //解析参数，默认不解析参数
		fmt.Println(r.Form)
		fmt.Println("Path:" + r.URL.Path)
		fmt.Println("Scheme:" + r.URL.Scheme)
		for k, v := range r.Form {
			fmt.Println("key:" + k)
			fmt.Println("value:" + strings.Join(v, ""))
		}
		fmt.Fprintln(w, "OK!")
	})
	startListent()
}

func startListent() {
	fmt.Println("start...9999")
	//监听端口
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Printf("error:%v\n", err)
		log.Fatal(err)
	}
}

//使用html模板
func m22() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "<html><body><a href='/login'>登录</a></body></html>")
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		printHeader(r)
		printRequest(r)
		fmt.Println("method:", r.Method)
		if r.Method == "GET" {
			t, _ := template.ParseFiles("login.gtpl")
			t.Execute(w, nil)
		}
		if r.Method == "POST" {
			r.ParseForm()
			fmt.Println("username:", r.Form["username"])
			if len(strings.Join(r.Form["username"], "")) == 0 {
				fmt.Println("empty username")
			}
			//设置cookie
			cookie := http.Cookie{Name: "username", Value: string([]byte(strings.Join(r.Form["username"], "")))}
			http.SetCookie(w, &cookie)
			//输出到客户端
			// fmt.Fprint(w, strings.Join(r.Form["username"], ""))
			//html转义输出到客户端
			template.HTMLEscape(w, []byte(strings.Join(r.Form["username"], "")))
		}
	})
	startListent()
}

//session使用
func m22Session() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		session := sessionManager.SessionStart(w, r)
		name := session.Get("username")
		if name != nil {
			fmt.Fprintln(w, "<html><body><a href='/login'>登录</a>", name, "</body></html>")
		} else {
			fmt.Fprintln(w, "<html><body><a href='/login'>登录</a></body></html>")
		}
	})
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		// printHeader(r)
		// printRequest(r)
		// fmt.Println("method:", r.Method)
		if r.Method == "GET" {
			t, _ := template.ParseFiles("login.gtpl")
			t.Execute(w, nil)
		}
		if r.Method == "POST" {
			r.ParseForm()
			fmt.Println("username:", r.Form["username"])
			if len(strings.Join(r.Form["username"], "")) == 0 {
				fmt.Println("empty username")
			}
			//设置session
			session := sessionManager.SessionStart(w, r)
			session.Set("username", string([]byte(strings.Join(r.Form["username"], ""))))
			//输出到客户端
			// fmt.Fprint(w, strings.Join(r.Form["username"], ""))
			//html转义输出到客户端
			template.HTMLEscape(w, []byte(strings.Join(r.Form["username"], "")))
		}
	})
	startListent()
}

//打印请求header
func printHeader(r *http.Request) {
	fmt.Println("-----------------Header Start----------------")
	for key, valArr := range r.Header {
		fmt.Printf("%v:%v\n", key, strings.Join(valArr, " "))
	}
	fmt.Println("-----------------Header End------------------")
}

//打印请求参数
func printRequest(r *http.Request) {
	fmt.Println("-----------------Request Start----------------")
	r.ParseForm()
	for k, v := range r.Form {
		fmt.Printf("%v:%v\n", k, strings.Join(v, " "))
	}
	fmt.Println("-----------------Request End------------------")
}

//访问数据库
/*
CREATE TABLE `userinfo` (
	`id` INT(10) NOT NULL AUTO_INCREMENT,
	`name` VARCHAR(50) NULL DEFAULT NULL,
	`age` INT NULL DEFAULT NULL,
	`email` VARCHAR(50) NULL DEFAULT NULL,
	PRIMARY KEY (`id`)
)
*/
func m23() {
	db, err := sql.Open("mysql", "root:tyl6632460@tcp(localhost:3306)/mytest1?charset=utf8")
	checkError(err)
	// m23Select(db)
	// m23Insert(db)
	//m23Update(db)
	m23Delete(db)
}

func m23Delete(db *sql.DB) {
	fmt.Println("<---------start delete--------->")
	stmt, err := db.Prepare("delete from userinfo where id = ?")
	checkError(err)
	res, err := stmt.Exec(1)
	rows, err := res.RowsAffected()
	fmt.Println("delete:", rows)
	fmt.Println("<---------start delete--------->")
}

func m23Update(db *sql.DB) {
	fmt.Println("<---------start update--------->")
	stmt, err := db.Prepare("update userinfo set email = ? where id = ?")
	checkError(err)
	rs, err := stmt.Exec("lisi@163.com", 3)
	checkError(err)
	fmt.Println(rs.RowsAffected())
	fmt.Println("<---------end update----------->")
}

func m23Insert(db *sql.DB) {
	fmt.Println("<---------start insert--------->")
	stmt, err := db.Prepare("insert into userinfo(name,age,email) value(?,?,?)")
	checkError(err)
	res, err := stmt.Exec("zhangsan", 12, "zhangsan@126.com")
	checkError(err)
	id, err := res.LastInsertId()
	checkError(err)
	fmt.Println("insert:", id)
	fmt.Println("<---------end insert----------->")
}

func m23Select(db *sql.DB) {
	fmt.Println("<---------start select--------->")
	rows, err := db.Query("select * from userinfo")
	for rows.Next() {
		var id int
		var name string
		var age int
		var email string
		err = rows.Scan(&id, &name, &age, &email)
		checkError(err)
		fmt.Printf("ID:%d,name:%v,age:%d,email:%v", id, name, age, email)
	}
	fmt.Println("<---------end select----------->")
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

//使用struct tag指定对应xml中的内容。struct tag可以通过反射获取到
type Schools struct {
	XMLName  xml.Name `xml:"schools"`
	Location string   `xml:"location,attr"`
	Schools  []School `xml:"school"`
	Desc     string   `xml:",innerxml"`
}

type School struct {
	XMLName      xml.Name `xml:"school"`
	Name         string   `xml:"name"`
	StudentCount int      `xml:"studentCount"`
}

//处理xml
func m24() {
	file, err := os.Open("E:\\workspace\\goproject\\src\\cmd\\school.xml") //打开文件
	if err != nil {
		fmt.Printf("error0:%v", err)
		return
	}
	defer file.Close() //是否需要在错误之前
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("error1:%v", err)
		return
	}
	val := Schools{}
	err = xml.Unmarshal(data, &val)
	if err != nil {
		fmt.Printf("error2:%v", err)
	}
	fmt.Println(val)
	fmt.Println("---------生成xml--------")
	schools := Schools{Location: "南京"}
	schools.Schools = append(schools.Schools, School{Name: "南京中学", StudentCount: 100})
	schools.Schools = append(schools.Schools, School{Name: "南京中学2", StudentCount: 200})
	output, err := xml.MarshalIndent(schools, "	", "	")
	if err != nil {
		fmt.Printf("error3:%v", err)
	}
	fmt.Println(string(output))
}

type Grade struct {
	Num      int
	Name     string
	Eclasses []Eclass
}

type GradeSlice struct {
	Grades []Grade
}

type Eclass struct {
	Name  string
	Email string
}

//处理json
func m25() {
	var str = `{"Grades":[{"Num":1,"Name":"1年级"},{"Num":2,"Name":"2年级"}]}`
	var grades GradeSlice
	//知道json结构
	err := json.Unmarshal([]byte(str), &grades)
	if err != nil {
		fmt.Println(err)
		return
	}
	//json结构未知
	fmt.Println(grades)
	var f interface{}
	err = json.Unmarshal([]byte(str), &f)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(f)
	//生成json
	var grades2 GradeSlice
	grades2.Grades = append(grades2.Grades, Grade{Num: 1, Name: "1年级"})
	grades2.Grades = append(grades2.Grades, Grade{Num: 2, Name: "2年级"})
	b, err := json.Marshal(grades2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))
}

//模板的使用
func m26() {
	//创建模板
	t1 := template.New("t1")
	t1, err := t1.Parse("序号:{{.Num}},名称:{{.Name}}\n")
	checkError(err)
	var g = Grade{Num: 1, Name: "1年级"}
	t1.Execute(os.Stdout, g) //输出到标准输出流
	g.Eclasses = append(g.Eclasses, Eclass{Name: "1班", Email: "1b@123.com"})
	g.Eclasses = append(g.Eclasses, Eclass{Name: "2班", Email: "2b@123.com"})
	g.Eclasses = append(g.Eclasses, Eclass{Name: "3班", Email: "3b@123.com"})
	g.Eclasses = append(g.Eclasses, Eclass{Name: "4班", Email: "4b@123.com"})
	t2, err2 := template.New("t2").Parse(`
		Grade:{{.Name}}
		{{range .Eclasses}}
		Name:{{.Name}} Email:{{.Email}}{{end}}
		`)
	checkError(err2)
	t2.Execute(os.Stdout, g)
	//模板函数
	t3 := template.New("t3")
	t3.Funcs(template.FuncMap{"emailDeal": func(args ...interface{}) string {
		s, ok := args[0].(string)
		if !ok {
			return "error"
		}
		return s + "!Deal"
	}})
	t3, err3 := t3.Parse(`
		Grade:{{.Name}}
		{{range .Eclasses}}
		Name:{{.Name}} Email:{{.Email | emailDeal}}{{end}}
		`)
	checkError(err3)
	t3.Execute(os.Stdout, g)
	//嵌套模板
	t4, err4 := template.ParseFiles("header.gtpl", "body.gtpl", "foot.gtpl")
	checkError(err4)
	fmt.Println("------------------")
	t4.ExecuteTemplate(os.Stdout, "header", nil)
	fmt.Println("------------------")
	t4.ExecuteTemplate(os.Stdout, "body", nil)
	fmt.Println("------------------")
	t4.ExecuteTemplate(os.Stdout, "foot", nil)
}

//文件操作
func m27() {
	//文件夹操作
	path := "aaa/bbb"
	fi, err := os.Stat("aaa/bbb")
	if err != nil {
		if os.IsExist(err) {
			fmt.Println("文件夹存在")
		} else {
			fmt.Println("文件夹不存在")
			os.MkdirAll(path, 0777)
			fmt.Println("已创建文件夹")
		}
	} else {
		if fi.IsDir() {
			fmt.Println("文件夹存在！！")
		} else {
			fmt.Println(path, "不是目录！！")
		}
		os.RemoveAll(path)
		fmt.Println("已删除", path, "！！")
	}
	f, err := os.Create("aaa.txt") //已存在会覆盖
	checkError(err)
	fmt.Println("已创建文件：", f.Name())
	f.WriteString("测试文件！！！！！！")
	defer f.Close()
	f2, err := os.Open("login.gtpl")
	defer f2.Close()
	checkError(err)
	buf := make([]byte, 1024)
	for {
		n, _ := f2.Read(buf)
		if 0 == n {
			break
		}
		fmt.Println(string(buf[:n]))
	}
}

func main() {
	// m1()
	// m2(123)
	// m3(1, 2, 3)
	// m3a()
	// m3b()
	// m3d()
	// m4()
	// m5()
	// m6()
	// m7()
	// m7a()
	// m8()
	// m9()
	// m9a()
	// m10()
	// m11()
	// m12()
	// m13()
	// m14()
	// m15()
	// m16()
	// m17()
	// m18()
	// m19()
	// m20()
	// m21()
	// m22()
	// m22Session()
	// m23()
	// m24()
	// m25()
	// m26()
	m27()
}
