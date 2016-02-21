package main

import (
	"fmt"
	"testing"
	"strconv"
//	"reflect"
)

func exampleInitialise(sm []StateMachine){
	
	sm[0].id = 1 //leader
	sm[1].id = 2 //peer1
	sm[2].id = 3 //peer2

	sm[0].state = 3 //leader
	sm[1].state = 1 //follower
	sm[2].state = 1 //follower
	
	for i:=0; i<len(sm); i++ {
		sm[i].log = make([]Log,3)		
		sm[i].matchIndex = make([]int,2)
	}

	sm[0].nextIndex = []int{1,1}
	sm[0].peers = []int{2,3}
	sm[1].peers = []int{1,3}
	sm[2].peers = []int{1,2}

	sm[0].log[0].term = 1
	sm[0].log[0].command = []byte{'r','e','a','d'}
	
	sm[0].log[1].term = 1
	sm[0].log[1].command = []byte{'w','r','i','t','e'}
	
	sm[0].log[1].term = 2
	sm[0].log[1].command = []byte{'r','e','a','d'}

	sm[0].commitIndex = 1
	sm[0].lastLogIndex = 3

	sm[0].matchIndex[0] = 1
	sm[0].matchIndex[1] = 1
}

func TestLeaderAppend(t *testing.T){
		
	var sm []StateMachine
	sm = make([]StateMachine,3)
	exampleInitialise(sm[:])

	a := sm[0].ProcessEvent(Append{data:sm[0].log[0].command})
	numAppendEntryReq := 0
	numLogStore := 0
	numUnexpectedEv := 0

	for i:=0; i<len(a); i++ {
			f,ok := a[i].(Send)
			if ok{
				_,yes := f.event.(AppendEntriesReqEv)
				if yes{
					numAppendEntryReq++	
				}
			} else{
					_,ok := a[i].(LogStore)
					if ok {
						numLogStore++
					} else {
						numUnexpectedEv++
					}
			}
	}

	expect(t,strconv.Itoa(numAppendEntryReq),strconv.Itoa(len(sm[0].peers)))
	expect(t,strconv.Itoa(numLogStore),"1")
	expect(t,strconv.Itoa(numUnexpectedEv),"0")
	
}
func TestLeaderAppendEntryReq(t *testing.T){

}

func TestLeaderHeartBeat(t *testing.T){

}

func expect(t *testing.T, a string, b string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}




	/*sm.state = 1
	z:=sm.ProcessEvent(AppendEntriesReqEv{term : 10, leaderId: 1, prevLogIndex: 100, prevLogTerm: 3, 
		entries: []log{{term:1,index:2,command:"read"}},leaderCommit: 1})
	sm.state = 3
	a:=sm.ProcessEvent(VoteReqEv{candidateId:10,term:20,lastLogIndex:30,lastLogTerm:5})
	
	sm.state=2
	z :=[]byte{1,2,3,4}
	a:= sm.ProcessEvent(Append{data:z})
	
	sm.state=2
	a := sm.ProcessEvent(Timeout{})
		
	sm.state = 3
	a:=sm.ProcessEvent(AppendEntriesRespEv{senderId: 1, senderTerm: 3, response:true})

	sm.state = 3
	a:=sm.ProcessEvent(VoteRespEv{senderTerm: 3, response:true})
	*/

	/*sm.state=3
	z :=[]byte{1,2,3,4}
	a:= sm.ProcessEvent(Append{data:z})
	f,ok := a[0].(LogStore)
	//fmt.Printf("%v\n", x)
	//fmt.Println("Error is:", a[0])
	if ok {
	fmt.Printf("%v\n", f)	
	}*/

	//z :=[]byte{1,2,3,4}
	/*sm.state = 1
	sm.timer = 11
	fmt.Println(rand.Intn(2*sm.timer-sm.timer)+sm.timer)
	/*a:=sm.ProcessEvent(Timeout{})
	f,ok := a[0].(Alarm)
	if ok {
	fmt.Printf("%v\n", f)	
	}*/
	
/*
func TestAppendFollower(t *testing.T){
	var sm StateMachine

}*/