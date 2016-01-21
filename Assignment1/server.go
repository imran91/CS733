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
  "errors"
)

type command_meta struct {
	operation string
	filename  string
	numbytes  int
	expiry    int
	version   int64
  data      string
  create_time int
}

//map manager to hold files 
var kvmanager = struct{
    sync.RWMutex
    kv map[string]*command_meta
}{kv: make(map[string]*command_meta)}

//parse the command and return it into command_meta object
func parser(cmd string) (command_meta,error) {
	cobj := command_meta{"", "", 0, 0, 0,"",0}
	command_split := strings.Split(cmd, " ")
	cobj.operation = command_split[0]
  err := errors.New("ERR_CMD_ERR\r\n")
  num_err := errors.New("ERR_NON_NUMERIC_CONTENT")
  no_of_param := len(command_split)
  
	switch command_split[0] {
	case "write":
    if no_of_param == 3 || no_of_param == 4 {
       cobj.filename = command_split[1]
      nbytes, err1 := strconv.Atoi(command_split[2])
      if err1 != nil {
        log.Println(err1)
        fmt.Println("nbytes error")
        return cobj, num_err
      }
      cobj.numbytes = nbytes  

      if len(command_split) == 4 {
      exp, err2 := strconv.Atoi(command_split[3])
      if err2 != nil {
        log.Println(err2)
        fmt.Println("exp error")
        return cobj, num_err
      }
      cobj.expiry = exp
      }
      cobj.create_time = int(time.Now().Unix())

    } else {
      fmt.Println("more parameters or less parameters")
      return cobj, err
    }
	
  case "read":
    if no_of_param == 2 {
        cobj.filename = command_split[1]
        return cobj ,nil  
      } else {
        return cobj, err
      }
  
  case "cas":
    if no_of_param == 5 || no_of_param == 6 {
          cobj.filename = command_split[1]
          ver, err1 := strconv.Atoi(command_split[2])
          if err1 != nil {
            log.Println(err1)
            return cobj, num_err
          }
          cobj.version = int64(ver)
          nbytes, err2 := strconv.Atoi(command_split[3])
          if err2 != nil {
            log.Println(err2)
            return cobj, num_err
          }
          cobj.numbytes = nbytes
          if len(command_split) == 5 {
            exp, err3 := strconv.Atoi(command_split[4])
            if err3 != nil {
              log.Println(err3)
              return cobj, num_err
            }
            cobj.expiry = exp
          }
            cobj.create_time = int(time.Now().Unix())
          return cobj, nil
      } else {
        return cobj, err
      }
		
	case "delete":
    if no_of_param == 2 {
      cobj.filename = command_split[1]
      return cobj, nil
    } else {
      return cobj,err
    }
	}//end of switch case
	return cobj, nil
}

func write(parsed_comm command_meta,ver int64,conn net.Conn,addr net.Addr) {
      //fmt.Println("time_new:",parsed_comm.expiry)

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

}

func handleConnection(conn net.Conn) {
	addr := conn.RemoteAddr()
	//log.Println(addr, "connected.")
	reader := bufio.NewReader(conn)
  var ver int64 = 1
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			//log.Println("ERROR reading:")
			break
		}

		message = strings.TrimRight(message, "\r\n")
		parsed_comm, err := parser(message)
    if(err != nil) {
      s := err.Error()
      conn.Write([]byte(s+"\r\n"))
      conn.Close()
    } else {
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
          log.Println("ERR_IN_CMD\r\n")
          conn.Write([]byte("ERR_IN_CMD\r\n"))
          //conn.Close()
          break
        }        
        write(parsed_comm,ver,conn,addr)
        break

    case "read":
        rem := 0
        kvmanager.RLock()
        obj, present := kvmanager.kv[parsed_comm.filename]
        if present {
          kvmanager.RUnlock()
          t := obj.expiry
          
          if t != 0 && t > 0{
            elapsed := int(time.Now().Unix())-obj.create_time
            rem = t - elapsed
            //fmt.Println("rem time is:",rem)
            if rem <= t && rem > 0 {
              //log.Println("new expiry time is:",rem)
            } else {
              //delete expired file 
              kvmanager.Lock()
              delete(kvmanager.kv,parsed_comm.filename)
              kvmanager.Unlock()
              conn.Write([]byte("ERR_FILE_NOT_FOUND\r\n"))
              break 
            }   
          }

          _, err := conn.Write([]byte( "CONTENTS " + strconv.Itoa(int(obj.version)) + " " + strconv.Itoa(obj.numbytes) + " " + 
                strconv.Itoa(rem) + " " + "\r\n" + obj.data +"\r\n"))
          if err != nil {
               fmt.Println(err)
          }
        } else {
          kvmanager.RUnlock()
          _, err := conn.Write([]byte("ERR_FILE_NOT_FOUND\r\n"))
          if err != nil {
               fmt.Println(err)
          }
        }
      break
      
    case "cas":
        rem := 0
        kvmanager.RLock()
        obj, present := kvmanager.kv[parsed_comm.filename]
        if present {
            kvmanager.RUnlock()
            if obj.version == parsed_comm.version {
              t := parsed_comm.expiry
              if t != 0 {
                  elapsed := int(time.Now().Unix())-obj.create_time
                  rem = t - elapsed
                  //fmt.Println("rem time is:",rem)
                  if rem <= t && rem > 0 {
                      kvmanager.Lock()
                      kvmanager.kv[parsed_comm.filename].create_time = int(time.Now().Unix())
                      kvmanager.Unlock()
                     // fmt.Println("new expiry time is:",rem)
                  } else {
                  //delete expired file 
                    kvmanager.Lock()
                    delete(kvmanager.kv,parsed_comm.filename)
                    kvmanager.Unlock()
                    conn.Write([]byte("ERR_FILE_NOT_FOUND\r\n"))
                    break 
                  }
              }

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
                conn.Write([]byte("ERR_IN_CMD\r\n"))
                //conn.Close()
              }

              kvmanager.Lock()
              kvmanager.kv[parsed_comm.filename].expiry = parsed_comm.expiry
              kvmanager.kv[parsed_comm.filename].version = parsed_comm.version
              kvmanager.kv[parsed_comm.filename].data = parsed_comm.data
              kvmanager.kv[parsed_comm.filename].numbytes = parsed_comm.numbytes
              kvmanager.Unlock()

              kvmanager.RLock()
              conn.Write([]byte( "OK " + strconv.Itoa(int(kvmanager.kv[parsed_comm.filename].version)) +"\r\n"))
              kvmanager.RUnlock()
            } else {
              kvmanager.RLock()
              conn.Write([]byte( "ERR_VERSION " + strconv.Itoa(int(kvmanager.kv[parsed_comm.filename].version)) +"\r\n"))
              kvmanager.RUnlock()
            }
        } else {
          kvmanager.RUnlock()
          _, err := conn.Write([]byte("ERR_FILE_NOT_FOUND\r\n"))
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
          _, err := conn.Write([]byte("ERR_FILE_NOT_FOUND\r\n"))
          if err != nil {
               fmt.Println(err)    
          }
        }
      break
    }
    }    
    
	}
	conn.Close() //used for clean up actions
}

func serverMain() {
	var ln net.Listener
	var list_err error
	//fmt.Println("Starting server...")

	tcpAddr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}

	for {
		ln, list_err = net.ListenTCP("tcp", tcpAddr)
		if list_err != nil {
			log.Println(list_err)
		} else {

			break
		}
	}
	defer ln.Close()
	for {
		// accept connection on port
		conn, con_err := ln.Accept()
		if con_err != nil {
			log.Println(con_err)
			continue
		}
		go handleConnection(conn)
	}
}

func main() {
	serverMain()
}
