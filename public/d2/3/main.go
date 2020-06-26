package main

import (
	"fmt"
	"os"
)

func main() {
	//多个变量一起声明不用麻烦那么多
	s, sep := "", ""
	//忽略第一个index的参数，用range来取值
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = "-"
	}
	fmt.Println(s)
}
