package main

import (
	//"bufio"
	"fmt"
	"net"
	/*"strconv"
	"strings"
	"os"
	"testing"*/)

func main() {
	name := "hi.txt"
	contents := "abcde\r\nhij\r\n"
	contents1 := "imran\r\n"
	//exptime := 300
	exptime1 := 1000
	exptime2 := 5000
	ver1 := 1
	// connect to this socket
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error connecting")
	}

	//scanner := bufio.NewScanner(conn)
	fmt.Fprintf(conn, "write %v %v %v\r\n%v\r\ncas %v %v %v %v\r\n%v\r\n", name, len(contents), exptime1,contents,name,ver1,len(contents1),exptime2,contents1)
	//fmt.Fprintf(conn, "write %v %v %v\r\n%v\r\n", name, len(contents), exptime, contents)
	/*scanner.Scan() // read first line
	  resp := scanner.Text() // extract the text from the buffer
	  arr := strings.Split(resp, " ") // split into OK and <version>
	  //expect(t, arr[0], "OK")
	  ver, err := strconv.Atoi(arr[1]) // parse version as number
	  if err != nil {
	  		fmt.Println("Non-numeric version found")
	  	}
	  	version := int64(ver)
	  fmt.Println("%v %v ",arr[0],version)
	  for {     // read in input from stdin
	  		reader := bufio.NewReader(os.Stdin)
	  		fmt.Print("Text to send: ")
	  		text, _ := reader.ReadString('\n')
	  		// send to socket
	  		fmt.Fprintf(conn, text + "\n")
	  		// listen for reply
	  		message, _ := bufio.NewReader(conn).ReadString('\n')
	  		fmt.Print("Message from server: "+message)
	  	}
	*/
}
