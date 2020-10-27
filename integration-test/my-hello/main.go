package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"os"
)

func main(){
	conn, _ := tls.Dial("tcp", ":1965", &tls.Config{InsecureSkipVerify: true})
	fmt.Fprint(conn, "gemini://foo.bar/integration-test/resources/hello.txt\r\n")
	for {
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		status, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(status)
	}
}
