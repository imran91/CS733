package main

import ( 
		 "net"
		 "fmt" 
     "bufio"
     "strings"
		 )


func handleConnection(conn net.Conn){

  for{
    // will listen for message to process ending in newline (\n)     
    message, _ := bufio.NewReader(conn).ReadString('\n')     
    // output message received     
    fmt.Print("Message Received:", string(message))     
    // sample process for string received     
    newmessage := strings.ToUpper(message)    
     // send new string back to client     
    conn.Write([]byte(newmessage + "\n"))   
  } 
}
var ln net.Listener
func main() {
  


fmt.Println("Starting server...")

  // listen on all interfaces
//for {
    ln, list_err := net.Listen("tcp", ":8081")
    if list_err!= nil {
            fmt.Println(list_err)
    } /*else {
      		  	break
    	   	  }
  }
*/
// run loop forever (or until ctrl-c)   
for {     
	// accept connection on port
  conn, con_err := ln.Accept()
  if con_err!= nil {
  		fmt.Println(con_err)
      continue
  	}
  go handleConnection(conn)
}

}
