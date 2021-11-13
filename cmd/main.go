package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/feitian124/goapi/config"
	"github.com/feitian124/goapi/datasource"
)

func main() {
	// 定义命令行参数方式1
	var dsn string
	flag.StringVar(&dsn, "name", "张三", "姓名")

	// 解析命令行参数
	flag.Parse()
	fmt.Println(dsn)

	c, err := config.New()
	if err != nil {
		os.Exit(-1)
	}

	s, err := datasource.Analyze(c.DSN)
	if err != nil {
		os.Exit(-1)
	}
	fmt.Printf("%+v", s)
}
