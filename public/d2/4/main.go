package main

import (
	"fmt"
	"os"
	"strings"
)

//跟之前的几个一样通过strings包实现输入参数的组和
func main() {
	fmt.Println(strings.Join(os.Args[1:], "-"))
}
