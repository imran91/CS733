package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestMain(t *testing.T) {
	go serverMain() // launch the server as a goroutine.
	time.Sleep(1 * time.Second)
}

// Simple serial check of getting and setting
func TestTCPSimple(t *testing.T) {
	name := "hi.txt"
	contents := "bye"
	exptime := 300000
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err.Error()) // report error through testing framework
	}
	scanner := bufio.NewScanner(conn)
	// Write a file
	fmt.Fprintf(conn, "write %v %v %v\r\n%v\r\n", name, len(contents), exptime, contents)
	scanner.Scan()                  // read first line
	resp := scanner.Text()          // extract the text from the buffer
	arr := strings.Split(resp, " ") // split into OK and <version>
	expect(t, arr[0], "OK")
	ver, err := strconv.Atoi(arr[1]) // parse version as number
	if err != nil {
		t.Error("Non-numeric version found")
	}
	version := int64(ver)

	fmt.Fprintf(conn, "read %v\r\n", name) // try a read now
	scanner.Scan()

	arr = strings.Split(scanner.Text(), " ")
	expect(t, arr[0], "CONTENTS")
	expect(t, arr[1], fmt.Sprintf("%v", version)) // expect only accepts strings, convert int version to string
	expect(t, arr[2], fmt.Sprintf("%v", len(contents)))
	scanner.Scan()
	expect(t, contents, scanner.Text())
	conn.Close()
}

// Useful testing function
func expect(t *testing.T, a string, b string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}

func TestExpiry(t *testing.T) {
	name := "hi.txt"
	contents := "bye"
	exptime := 5
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err.Error()) // report error through testing framework
	}
	scanner := bufio.NewScanner(conn)
	// Write a file
	fmt.Fprintf(conn, "write %v %v %v\r\n%v\r\n", name, len(contents), exptime, contents)
	scanner.Scan()                  // read first line
	resp := scanner.Text()          // extract the text from the buffer
	arr := strings.Split(resp, " ") // split into OK and <version>
	expect(t, arr[0], "OK")
	ver, err := strconv.Atoi(arr[1]) // parse version as number
	if err != nil {
		t.Error("Non-numeric version %v found", ver)
	}
	time.Sleep(5 * time.Second) //wait till file content gets expired

	fmt.Fprintf(conn, "read %v\r\n", name) // try a read now
	scanner.Scan()

	rep := scanner.Text()
	expect(t, rep, "ERR_FILE_NOT_FOUND")
	conn.Close()
}

func TestConcurrentWrite(t *testing.T) {
	N := 100 //Number of clients
	k := 10  //Number of operations per thread

	ack := make(chan bool, N)

	for i := 0; i < N; i++ {
		go CreateClient(t, k, ack, i)
		time.Sleep(time.Nanosecond) //	For some reason this is necessary
	}

	tick := time.Tick(20 * time.Second)

	for i := 0; i < N; i++ {
		select {
		case <-tick:
			{
				t.Error("Timeout", N-i, "threads were unable to complete write operation")
				break
			}
		case <-ack:
		}
	}
}

func CreateClient(t *testing.T, k int, ack chan<- bool, id int) {
	name := "hi.txt"
	contents := "bye"
	exptime := 500
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err, id)
		ack <- true
		return
	}

	scanner := bufio.NewScanner(conn)
	for j := 1; j < k; j++ {
		// Write a file
		fmt.Fprintf(conn, "write %v %v %v\r\n%v\r\n", name, len(contents), exptime, contents)
		scanner.Scan()                  // read first line
		resp := scanner.Text()          // extract the text from the buffer
		arr := strings.Split(resp, " ") // split into OK and <version>
		expect(t, arr[0], "OK")
		ver, err := strconv.Atoi(arr[1]) // parse version as number
		if err != nil {
			t.Error("Non-numeric version %v found", ver)
		}
	}
	conn.Close()
	ack <- true
}

func TestMiscellaneous(t *testing.T) {
	name := "hi.txt"
	contents := "bye"
	exptime := 5
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Error(err.Error()) // report error through testing framework
	}

	scanner := bufio.NewScanner(conn)
	// Write a file
	fmt.Fprintf(conn, "write %v %v %v\r\n%v\r\n", name, len(contents), exptime, contents)
	scanner.Scan()                  // read first line
	resp := scanner.Text()          // extract the text from the buffer
	arr := strings.Split(resp, " ") // split into OK and <version>
	expect(t, arr[0], "OK")
	ver, err := strconv.Atoi(arr[1]) // parse version as number
	if err != nil {
		t.Error("Non-numeric version %v found", ver)
	}
	version := int64(ver)

	//checking correctness of the cas command
	fmt.Fprintf(conn, "cas %v %v %v %v\r\n%v\r\n", name, version, len(contents), exptime, contents)
	scanner.Scan()                 // read first line
	resp = scanner.Text()          // extract the text from the buffer
	arr = strings.Split(resp, " ") // split into OK and <version>
	expect(t, arr[0], "OK")
	ver, err = strconv.Atoi(arr[1]) // parse version as number
	if err != nil {
		t.Error("Non-numeric version %v found", ver)
	}
	version = version + 1
	expect(t, arr[1], strconv.Itoa(int(version)))

	//checking correctness of delete operation
	fmt.Fprintf(conn, "delete %v\r\n", name)
	scanner.Scan()        // read first line
	resp = scanner.Text() // extract the text from the buffer
	arr = strings.Split(resp, " ")
	expect(t, arr[0], "OK")

	//checking read operation correctness
	fmt.Fprintf(conn, "read %v\r\n", name)
	scanner.Scan()        // read first line
	resp = scanner.Text() // extract the text from the buffer
	arr = strings.Split(resp, " ")
	expect(t, arr[0], "ERR_FILE_NOT_FOUND")

	//checking invalid no of arguments in cas command
	fmt.Fprintf(conn, "cas %v %v \r\n%v\r\n", name, version, contents)
	scanner.Scan()        // read first line
	resp = scanner.Text() // extract the text from the buffer
	arr = strings.Split(resp, " ")
	expect(t, arr[0], "ERR_CMD_ERR")

}
