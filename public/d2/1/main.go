package main

import (
	"fmt"
	"math"
)

const (
	a = 3
	b = 4
)

func main() {

	c := math.Sqrt(a*a + b*b)
	fmt.Printf("%.1f", c)
	//简单的勾股定理，c格式化一个小数点
}
