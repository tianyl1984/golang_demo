package xmlDemo

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

func StartDemo() {
	m1()
}

func m1() {
	file, err := os.Open("..\\xmlDemo\\school.xml") //打开文件
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
