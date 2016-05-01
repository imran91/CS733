package main

import (
	//"encoding/gob"
	 "fmt"
	"time"
	"errors"
	"strconv"
	"testing"
	//	"reflect"
)

var raft_cluster = []NetConfig{
	{	
		Id: 101, 
		Host: "127.0.0.1", 
		Port: 5001,
	},
	{
		Id: 102, 
		Host: "127.0.0.1", 
		Port: 5002,
	},
	{
		Id: 103, 
		Host: "127.0.0.1", 
		Port: 5003,
	},
	{
		Id: 104, 
		Host: "127.0.0.1", 
		Port: 5004,
	},
	{
		Id: 105, 
		Host: "127.0.0.1", 
		Port: 5005,
	},
}


func retConfig(selfId int) Config {
	var	config Config

	config.cluster = raft_cluster
	config.Id = selfId
	config.LogDir = strconv.Itoa(config.Id)
	config.ElectionTimeout = 500
	config.HeartbeatTimeout = 100
	return config
}

func createRaft(index int) Node{
	var config Config
	config = retConfig(index)
	node := New(config)
	return node
}

func createRafts() ([]Node) {
	var nodes []Node
	for i,_ :=range raft_cluster {
		raft_node := createRaft(101+i)
		nodes = append(nodes, raft_node)
	}
	return nodes
}

func getLeader(rafts []Node) (Node,error){
	var id int
	var leader_id int
	var status bool
	//fmt.Println("Inside getLeader")
	for _, temp := range rafts {
		id,status = temp.Id()
		if !status {
			continue
		}

		leader_id,status = temp.LeaderId()
		if !status {
			continue
		}

		if id == leader_id {
			return temp,nil
		}
	}
	return nil,errors.New("LEADER_NOT_FOUND")
}
var rafts []Node
var ok error
var ldr Node

func TestBasic(t *testing.T) {

	rafts = createRafts() // array of []raft.Node
	data := "foo"
	
	ok = errors.New("err")
	for ok != nil {
		ldr, ok = getLeader(rafts)
	}
	res := ldr.Append([]byte(data))
	if !res {
		t.Error("basic: append on shutdown leader")
	}

	//time.Sleep(2 * time.Second)

	for _, rn := range rafts {
		ch, _ := rn.CommitChannel()
		//node := rn.(*RaftNode)
		select {
		case c := <- ch:
			//fmt.Println("data=====",node.sm.id,string(c.Data))
			if c.Err != nil {
				t.Error(c.Err)
			} else if string(c.Data) != data {
				t.Error("basic: Different data")
			}
		}
	}
}


func TestAppend(t *testing.T) {
	var data []string
	data = append(data,"Imran")
	data = append(data,"Pathan")
	data = append(data,"IIT powai")

	for i:=0;i<3;i++ {
		ok := errors.New("err")
		for ok != nil {
			ldr, ok = getLeader(rafts)	
		}
		//id,_:=ldr.Id()
		//fmt.Println("leader is=",id)
		
		res := ldr.Append([]byte(data[i]))
		if !res {
			t.Error("append on shutdown leader")
		}
		//time.Sleep(time.Second * 2)
	}

	for _,testData := range data {
		for _, rn := range rafts {
			ch, _ := rn.CommitChannel()
			node := rn.(*RaftNode)
			select {
			case c := <- ch:
				fmt.Println("data=====",node.sm.id,string(c.Data))
				if c.Err != nil {
					t.Error(c.Err)
				} else if string(c.Data) != testData {
					t.Error("basic: Different data")
				}
			}
		}
	}
}

func TestLeaderFailure(t *testing.T){
	var data string
	data = "Test me if you can Leader"
	ok = errors.New("err")
	for ok != nil {
		ldr, ok = getLeader(rafts)
	}
	fmt.Println("old leader",ldr)
	old_ldr_id,_ :=ldr.Id()
	ldr.Shutdown()
	time.Sleep(time.Second * 5)
	ok = errors.New("err")
	for ok != nil {
		ldr, ok = getLeader(rafts)
	}
	fmt.Println("new leader",ldr)
	res := ldr.Append([]byte(data))
	if !res {
		t.Error("basic: append on shutdown leader")
	}

	for _, rn := range rafts {
		ch, _ := rn.CommitChannel()
		node := rn.(*RaftNode)
		if node.sm.id != old_ldr_id {
					select {
						case c := <- ch:
							fmt.Println("data=====",node.sm.id,string(c.Data))
							if c.Err != nil {
								t.Error(c.Err)
							} else if string(c.Data) != data {
								t.Error("basic: Different data")
							}
					}
			}
	}
	for _, node := range rafts {
		node.Shutdown()
	}
}
