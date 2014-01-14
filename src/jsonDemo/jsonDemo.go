package jsonDemo

import (
	"encoding/json"
	"fmt"
)

func StartDemo() {
	m1()
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

func m1() {
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
