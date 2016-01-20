package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"io"
  "sync"
  "time"
)

type command_meta struct {
	operation string
	filename  string
	numbytes  int
	expiry    int
	version   int64
  data string
}

//map manager to hold files 
var kvmanager = struct{
    sync.RWMutex
    kv map[string]*command_meta
}{kv: make(map[string]*command_meta)}

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

    if len(command_split) == 4 {
    exp, err2 := strconv.Atoi(command_split[3])
		if err2 != nil {
			log.Println(err2)
			break
		}
		cobj.expiry = exp
  }
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
		cobj.version = int64(ver)
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
			cobj.expiry = exp
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
  var ver int64 = 1
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
      t := parsed_comm.expiry
      fmt.Println("t:",t)
      if t != 0 {
          fmt.Println("time_before:",time.Now().Unix())
          t += int(time.Now().Unix())
        }
      parsed_comm.expiry = t
      fmt.Println("time_new:",parsed_comm.expiry)

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

        kvmanager.RLock()
        obj, present := kvmanager.kv[parsed_comm.filename]
        kvmanager.RUnlock()
        
      if present {
        ver = obj.version
        ver = ver + 1
        parsed_comm.version = ver
        } else {
        parsed_comm.version = ver
        ver = ver + 1
      }

        ref := &(parsed_comm)
        kvmanager.Lock()
        kvmanager.kv[parsed_comm.filename] = ref
        kvmanager.Unlock()

        kvmanager.RLock()
        conn.Write([]byte( "OK " + strconv.Itoa(int(kvmanager.kv[parsed_comm.filename].version)) +"\r\n"))
        kvmanager.RUnlock()      
        break

    case "read":
        kvmanager.RLock()
        obj, present := kvmanager.kv[parsed_comm.filename]
        if present {
          kvmanager.RUnlock()
          t := obj.expiry
          fmt.Println("read expiry time is:",t)
          if t != 0 {
            t = obj.expiry - int(time.Now().Unix()) // remaining time
            fmt.Println("new expiry time is:",t)
          }
          if t < 0 {
           fmt.Println("less than zero time is:",t)
            t = 0
          }

          //kvmanager.RUnlock()
          _, err := conn.Write([]byte( "CONTENTS " + strconv.Itoa(int(obj.version)) + " " + strconv.Itoa(obj.numbytes) + " " + 
                strconv.Itoa(t) + " " + "\r\n" + obj.data +"\r\n"))
          if err != nil {
               fmt.Println(err)
          }
        } else {
          kvmanager.RUnlock()
          _, err := conn.Write([]byte("ERRNOTFOUND\r\n"))
          fmt.Println(err)
        }
      break
      
    case "cas":
        kvmanager.RLock()
        obj, present := kvmanager.kv[parsed_comm.filename]
        if present {
            kvmanager.RUnlock()
            if obj.version == parsed_comm.version {
              t := parsed_comm.expiry
              if t != 0 {
              t += int(time.Now().Unix())
              }
              parsed_comm.expiry = t

              buffer := make([]byte, parsed_comm.numbytes)
               _, err1 := io.ReadFull(reader, buffer)
              if err1 != nil {
                //Read error
                log.Println("ERROR reading data:", addr, err1)
                break
              }     

              extra_data, err2 := reader.ReadString('\n')
              if err2 != nil {
                //Read error
                log.Println("ERROR reading post-data:", addr, err2)
                break
              }     

              parsed_comm.data = string(buffer)
              ver = obj.version
              ver = ver + 1
              parsed_comm.version = ver
              
              if (strings.TrimRight(extra_data,"\r\n") != "") || (len(parsed_comm.data) != parsed_comm.numbytes) {
                log.Println("ERROR IN CMD")
                continue
              }

              kvmanager.Lock()
              fmt.Println("hello")
              kvmanager.kv[parsed_comm.filename].expiry = parsed_comm.expiry
              kvmanager.kv[parsed_comm.filename].version = parsed_comm.version
              kvmanager.kv[parsed_comm.filename].data = parsed_comm.data
              kvmanager.kv[parsed_comm.filename].numbytes = parsed_comm.numbytes
              kvmanager.Unlock()

              fmt.Println("filename:", parsed_comm.filename)
              fmt.Println(" numbytes:", parsed_comm.numbytes)
              fmt.Println(" expiry: ", parsed_comm.expiry)
              fmt.Println(" version: ", parsed_comm.version)
              fmt.Println(" data: ", parsed_comm.data)
              kvmanager.RLock()
              conn.Write([]byte( "Ok " + strconv.Itoa(int(kvmanager.kv[parsed_comm.filename].version)) +"\r\n"))
              kvmanager.RUnlock()
            } else {
              _, err := conn.Write([]byte("VER NOT FOUND\r\n"))
              fmt.Println(err)
            }
        } else {
          kvmanager.RUnlock()
          _, err := conn.Write([]byte("FILE NOT FOUND\r\n"))
          fmt.Println(err)
        }
      break

    case "delete":
      kvmanager.RLock()
        _, present := kvmanager.kv[parsed_comm.filename]
      kvmanager.RUnlock()
        if present {
          kvmanager.Lock()
          delete(kvmanager.kv,parsed_comm.filename)
          kvmanager.Unlock()
        _, err := conn.Write([]byte("OK\r\n"))
        if err != nil {
               fmt.Println(err)
          }
        } else {
          _, err := conn.Write([]byte("ERRNOTFOUND\r\n"))
          fmt.Println(err)
        }
      break
    }
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
