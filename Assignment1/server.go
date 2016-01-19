package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"io"
)

type command_meta struct {
	operation string
	filename  string
	numbytes  int
	expiry    int64
	version   int
  data string
}

//parse the command and return it into command_meta object
func parser(cmd string) command_meta {
	cobj := command_meta{"", "", 0, 0, 0,""}
	command_split := strings.Split(cmd, " ")
	cobj.operation = command_split[0]

	switch command_split[0] {
	case "write":
		cobj.filename = command_split[1]
		nbytes, err1 := strconv.Atoi(command_split[2])
		if err1 != nil {
			log.Println(err1)
			break
		}
		cobj.numbytes = nbytes
		exp, err2 := strconv.Atoi(command_split[3])
		if err2 != nil {
			log.Println(err2)
			break
		}
		cobj.expiry = int64(exp)
		return cobj

	case "read":
		cobj.filename = command_split[1]
		return cobj
	case "cas":
		cobj.filename = command_split[1]
		ver, err1 := strconv.Atoi(command_split[2])
		if err1 != nil {
			log.Println(err1)
			break
		}
		cobj.version = ver
		nbytes, err2 := strconv.Atoi(command_split[3])
		if err2 != nil {
			log.Println(err2)
			break
		}
		cobj.numbytes = nbytes
		if len(command_split) == 5 {
			exp, err3 := strconv.Atoi(command_split[4])
			if err3 != nil {
				log.Println(err3)
				break
			}
			cobj.expiry = int64(exp)
		}
		return cobj
	case "delete":
		cobj.filename = command_split[1]
		return cobj
	}
	return cobj
}

func handleConnection(conn net.Conn) {
	addr := conn.RemoteAddr()
	log.Println(addr, "connected.")
	reader := bufio.NewReader(conn)
	for {

		message, err := reader.ReadString('\n')
		if err != nil {
			//log.Println("ERROR reading:")
			continue
		}
		message = strings.TrimRight(message, "\r\n")
		parsed_comm := parser(message)
    switch parsed_comm.operation {
    case "write":
        buffer := make([]byte, parsed_comm.numbytes)
        _, err1 := io.ReadFull(reader, buffer)
        if (err1) != nil {
          //Read error
          log.Println("ERROR reading data:", addr, err1)
          break
        }

        extra_data, err2 := reader.ReadString('\n')
        if (err2) != nil {
          //Read error
          log.Println("ERROR reading post-data:", addr, err2)
          break
        }
        parsed_comm.data = string(buffer)
        if (strings.TrimRight(extra_data,"\r\n") != "") || (len(parsed_comm.data) != parsed_comm.numbytes) {
          log.Println("ERROR IN CMD")
          continue
        }
        break;
    }

    
    fmt.Println("filename:", parsed_comm.filename)
    fmt.Println(" numbytes:", parsed_comm.numbytes)
    fmt.Println(" expiry: ", parsed_comm.expiry)
    fmt.Println(" version: ", parsed_comm.version)
    fmt.Println(" data: ", parsed_comm.data)
		/*cmd, e := parser(message)
		  if e != nil {
		    //write(writer, addr, e.Error())
		    log.Println("ERROR parsing:", addr, e)
		    } else {

		    }*/
		// output message received

		// sample process for string received
		//newmessage := strings.ToUpper(message)
		// send new string back to client
		//conn.Write([]byte(newmessage + "\n"))
	}

	conn.Close() //used for clean up actions
}

func serverMain() {
	var ln net.Listener
	var list_err error
	fmt.Println("Starting server...")

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}

	for {
		ln, list_err = net.ListenTCP("tcp", tcpAddr)
		if list_err != nil {
			log.Println(err)
		} else {
			break
		}
	}

	defer ln.Close()

	// run loop forever (or until ctrl-c)
	for {
		// accept connection on port
		conn, con_err := ln.Accept()
		if con_err != nil {
			fmt.Println("Hello", con_err)
			continue
		}
		go handleConnection(conn)
	}
}

func main() {
	serverMain()
}
