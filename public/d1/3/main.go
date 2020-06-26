package main

import (
	"fmt"
	"math"
)

func main() {
	//math包 Abs绝对值，Pow（x,y）以x为底的y次幂。 Sqrt（x） 求平方跟
	a := -2.0
	fmt.Println(math.Abs(a))
	b := math.Pow(a, 2)
	fmt.Println(b)

	fmt.Println(math.Sqrt(b))
}
