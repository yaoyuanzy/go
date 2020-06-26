package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
)

//要求前端的网页打印出10的平方，用math.Pow函数完成
func hand(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hello,world,from Go!")
	val := math.Pow(10, 2)
	fmt.Fprintln(w, val)
	fmt.Fprintf(w, "%v\n", val)
	//不是自己做出来的，都尼玛记不住Pow咋用看来还是得多复习啊
}

//the entrance
func main() {
	log.Println("Starting http server on localhost:8080")
	http.HandleFunc("/", hand)
	//访问根目录的时候调用hand的方法在前端显示文字
	http.ListenAndServe(":8080", nil)
}
