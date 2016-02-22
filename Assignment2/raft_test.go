package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
	//	"reflect"
)

func exampleInitialise(sm []StateMachine) {

	sm[0].id = 1 //leader
	sm[1].id = 2 //peer1
	sm[2].id = 3 //peer2

	sm[0].state = 3 //leader
	sm[1].state = 1 //follower
	sm[2].state = 1 //follower

	for i := 0; i < len(sm); i++ {
		sm[i].log = make([]Log, 5)
		sm[i].matchIndex = make([]int, 2)
		sm[i].votedAs = make([]int, 2)
	}

	sm[0].nextIndex = []int{1, 1}
	sm[0].peers = []int{2, 3}
	sm[1].peers = []int{1, 3}
	sm[2].peers = []int{1, 2}

	sm[0].timer = 10
	sm[1].timer = 10
	sm[2].timer = 10

	sm[0].term = 2
	sm[1].term = 1
	sm[2].term = 1
	sm[0].log[0].term = 1
	sm[0].log[0].command = []byte{'r', 'e', 'a', 'd'}
	sm[1].log[0].term = 1
	sm[1].log[0].command = []byte{'r', 'e', 'a', 'd'}
	sm[2].log[0].term = 1
	sm[2].log[0].command = []byte{'r', 'e', 'a', 'd'}

	sm[0].log[1].term = 1
	sm[0].log[1].command = []byte{'w', 'r', 'i', 't', 'e'}

	sm[0].log[2].term = 2
	sm[0].log[2].command = []byte{'r', 'e', 'a', 'd'}

	sm[0].commitIndex = 1
	sm[1].commitIndex = 0
	sm[2].commitIndex = 0
	sm[0].lastLogIndex = 2

	sm[0].matchIndex[0] = 0
	sm[0].matchIndex[1] = 0
}

func TestLeaderAppend(t *testing.T) {

	var sm []StateMachine
	sm = make([]StateMachine, 3)
	exampleInitialise(sm[:])

	a := sm[0].ProcessEvent(Append{data: sm[0].log[0].command})
	numAppendEntryReq := 0
	numLogStore := 0
	numUnexpectedEv := 0

	for i := 0; i < len(a); i++ {
		f, ok := a[i].(Send)
		if ok {
			_, yes := f.event.(AppendEntriesReqEv)
			if yes {
				numAppendEntryReq++
			}
		} else {
			_, ok := a[i].(LogStore)
			if ok {
				numLogStore++
			} else {
				numUnexpectedEv++
			}
		}
	}

	expect(t, strconv.Itoa(numAppendEntryReq), strconv.Itoa(len(sm[0].peers)))
	expect(t, strconv.Itoa(numLogStore), "1")
	expect(t, strconv.Itoa(numUnexpectedEv), "0")

}

func TestAppendEntryReqFollowerEv(t *testing.T) {
	var sm []StateMachine
	var lastIndex int
	lastIndex = -1
	sm = make([]StateMachine, 3)
	exampleInitialise(sm[:])

	a := sm[1].ProcessEvent(AppendEntriesReqEv{term: 2, senderId: 1, prevLogIndex: 0, prevLogTerm: 1,
		entries: []Log(sm[0].log[sm[0].commitIndex : sm[0].lastLogIndex+1]), senderCommitIndex: 1})

	lastIndex = giveIndexOfEvent(a, 4) //Check data need to store on follower side for LogStore
	if lastIndex >= 0 {
		expected, _ := json.Marshal(sm[0].log[sm[0].commitIndex : sm[0].lastLogIndex+1])
		actual, _ := json.Marshal(a[lastIndex].(LogStore).logEntry)
		expect(t, string(actual), string(expected))
	}

	lastIndex = giveIndexOfEvent(a, 5) //StateStore expected to update currTerm from 1 to 2
	if lastIndex >= 0 {
		expect(t, strconv.Itoa(a[lastIndex].(StateStore).currTerm), "2")
	}

	lastIndex = giveIndexOfEvent(a, 2) //Check Commit Index is updated on follower
	if lastIndex >= 0 {
		expect(t, strconv.Itoa(a[lastIndex].(Commit).index), "1")
	}
}

func TestAppendEntryRespLeaderEv(t *testing.T) {
	var sm []StateMachine
	var lastIndex int
	lastIndex = -1
	sm = make([]StateMachine, 3)
	exampleInitialise(sm[:])

	sm[1].log[1].term = 1
	sm[1].log[1].command = []byte{'w', 'r', 'i', 't', 'e'}

	sm[1].log[2].term = 2
	sm[1].log[2].command = []byte{'r', 'e', 'a', 'd'}

	sm[1].lastLogIndex = 2

	a := sm[0].ProcessEvent(AppendEntriesRespEv{senderId: 2, senderTerm: 2, response: true, lastMatchIndex: 2})
	lastIndex = giveIndexOfEvent(a, 4) //Check commit index and data is updated in Commit action
	if lastIndex >= 0 {
		expect(t, strconv.Itoa(a[lastIndex].(Commit).index), "2")
		actual, _ := json.Marshal(a[lastIndex].(Commit).data)
		expected, _ := json.Marshal(sm[0].log)
		expect(t, string(actual), string(expected))
	}

	a = sm[0].ProcessEvent(AppendEntriesRespEv{senderId: 3, senderTerm: 1, response: false, lastMatchIndex: 0})
	lastIndex = giveIndexOfEvent(a, 1) //Check Append Entry request is sent
	numAppendEntryReq := 0
	numUnexpectedEv := 0

	if lastIndex >= 0 {
		for i := 0; i < len(a); i++ {
			f, ok := a[i].(Send)
			if ok {
				_, yes := f.event.(AppendEntriesReqEv)
				if yes {
					numAppendEntryReq++
				}
			} else {
				numUnexpectedEv++
			}
		}
		expect(t, strconv.Itoa(numAppendEntryReq), "1")
		expect(t, strconv.Itoa(numUnexpectedEv), "0")

	}
}

func TestTimeoutFollower(t *testing.T) {
	var sm []StateMachine
	sm = make([]StateMachine, 3)
	exampleInitialise(sm[:])

	a := sm[1].ProcessEvent(Timeout{})
	expect(t, strconv.Itoa(sm[1].state), "2") //state is changed to candidate
	numVoteReq := 0
	numAlarm := 0
	numStateStore := 0
	numUnexpectedEv := 0

	for i := 0; i < len(a); i++ {
		f, ok := a[i].(Send)
		if ok {
			_, yes := f.event.(VoteReqEv)
			if yes {
				numVoteReq++
			}
		} else {
			_, ok := a[i].(StateStore)
			if ok {
				numStateStore++
			} else {
				_, ok := a[i].(Alarm)
				if ok {
					numAlarm++
				} else {
					numUnexpectedEv++
				}
			}
		}
	}

	expect(t, strconv.Itoa(numVoteReq), strconv.Itoa(len(sm[0].peers)))
	expect(t, strconv.Itoa(numStateStore), "2")
	expect(t, strconv.Itoa(numAlarm), "1")
	expect(t, strconv.Itoa(numUnexpectedEv), "0")

}

func TestTimeoutCandidate(t *testing.T) {
	var sm []StateMachine
	sm = make([]StateMachine, 3)
	exampleInitialise(sm[:])

	sm[1].state = 2 //Candidate
	a := sm[1].ProcessEvent(Timeout{})
	expect(t, strconv.Itoa(sm[1].state), "2") //Candidate should not change state after timeout
	numVoteReq := 0
	numAlarm := 0
	numStateStore := 0
	numUnexpectedEv := 0

	for i := 0; i < len(a); i++ {
		f, ok := a[i].(Send)
		if ok {
			_, yes := f.event.(VoteReqEv)
			if yes {
				numVoteReq++
			}
		} else {
			_, ok := a[i].(StateStore)
			if ok {
				numStateStore++
			} else {
				_, ok := a[i].(Alarm)
				if ok {
					numAlarm++
				} else {
					numUnexpectedEv++
				}
			}
		}
	}

	expect(t, strconv.Itoa(numVoteReq), strconv.Itoa(len(sm[0].peers)))
	expect(t, strconv.Itoa(numStateStore), "2")
	expect(t, strconv.Itoa(numAlarm), "1")
	expect(t, strconv.Itoa(numUnexpectedEv), "0")

}

func TestTimeoutLeader(t *testing.T) {
	var sm []StateMachine
	sm = make([]StateMachine, 3)
	exampleInitialise(sm[:])

	a := sm[0].ProcessEvent(Timeout{})
	expect(t, strconv.Itoa(sm[0].state), "3") //Leader should not change state after timeout

	numAppendEntryReq := 0
	numAlarm := 0
	numUnexpectedEv := 0

	for i := 0; i < len(a); i++ {
		f, ok := a[i].(Send)
		if ok {
			_, yes := f.event.(AppendEntriesReqEv)
			if yes {
				numAppendEntryReq++
			}
		} else {
			_, ok := a[i].(Alarm)
			if ok {
				numAlarm++
			} else {
				numUnexpectedEv++
			}
		}
	}

	expect(t, strconv.Itoa(numAppendEntryReq), strconv.Itoa(len(sm[0].peers))) //no of heart beat msgs
	expect(t, strconv.Itoa(numAlarm), "1")
	expect(t, strconv.Itoa(numUnexpectedEv), "0")

}

func giveIndexOfEvent(a []interface{}, event int) int {
	var ind int
	ind = -1
	for i := 0; i < len(a); i++ {
		switch a[i].(type) {
		case Send:
			if event == 1 {
				ind = i
			}
		case Commit:
			if event == 2 {
				ind = i
			}
		case Alarm:
			if event == 3 {
				ind = i
			}
		case LogStore:
			if event == 4 {
				ind = i
			}
		case StateStore:
			if event == 5 {
				ind = i
			}
		}
		//fmt.Println(event.(StateStore))
	}
	//fmt.Println()
	return ind
}

func expect(t *testing.T, a string, b string) {
	if a != b {
		t.Error(fmt.Sprintf("Expected %v, found %v", b, a)) // t.Error is visible when running `go test -verbose`
	}
}
