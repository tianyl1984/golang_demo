package templateDemo

import (
	"fmt"
	"html/template"
	"os"
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

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}
func m1() {
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
