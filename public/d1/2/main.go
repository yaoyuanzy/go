package main

/*c10k goroute的示例 轻量级线程，内存占用少
CSP MPG反正我不懂 https://www.cnblogs.com/sunsky303/p/9115530.html
*/
import (
	"fmt"
	"log"
	"net"
	"time"
)

func hand(conn net.Conn) {
	//当前时间用string的打出来写到链接中然后关闭
	fmt.Fprintf(conn, "%s", time.Now().String())
	conn.Close()
}
func main() {
	log.Println("Start listen on tcp 8080")
	// 本地tcp 8080上起来
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go hand(conn)
	}
}
