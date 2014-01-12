package reflectDemo

import (
	"fmt"
	"reflect"
)

type Student struct {
	Name string
	Age  int
}

func (stu Student) Hello(str string) {
	fmt.Println(stu.Name, str)
}

func m1() {
	var stu = Student{"zhangsan", 12}
	t := reflect.TypeOf(stu)
	if k := t.Kind(); k != reflect.Struct {
		fmt.Println("不是struct")
		return
	}
	fmt.Println("type:", t.Name())
	v := reflect.ValueOf(stu)
	//获取字段，若获取不可导出的字段值就会报错
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		val := v.Field(i).Interface()
		fmt.Println(f.Name, ":", f.Type, ":", val)
	}
	//获取方法
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		fmt.Println(m.Name, m.Type)
	}
	//修改
	v2 := reflect.ValueOf(&stu)
	if v2.Kind() == reflect.Ptr && !v2.Elem().CanSet() {
		fmt.Println("不能修改")
		return
	} else {
		v2 = v2.Elem()
	}
	//判断是否找到字段
	//v2.FieldByName("aaa").IsValid()
	if f2 := v2.FieldByName("Name"); f2.Kind() == reflect.String {
		f2.SetString("lisi")
		fmt.Println(stu)
	}

	//动态调用方法
	v3 := reflect.ValueOf(stu)
	m := v3.MethodByName("Hello")
	if !m.IsValid() {
		fmt.Println("方法:", m.IsValid())
		return
	}
	args := []reflect.Value{reflect.ValueOf("aaaa")}
	m.Call(args)
}

func StartDemo() {
	m1()
}
