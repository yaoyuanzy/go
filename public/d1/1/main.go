//package
package main

//std lib
import (
	"fmt"
	"log"
	"net/http"
)

//一个简单的http服务器，感觉有点意思
func hand(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello,world,from Go!")
}

//the entrance
func main() {
	log.Println("Starting http server on localhost:8080")
	http.HandleFunc("/", hand)
	//访问根目录的时候调用hand的方法在前端显示文字
	http.ListenAndServe(":8080", nil)
}
