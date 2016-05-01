package main

import (
	"errors"
//	"fmt"
)

var ERR_NOT_LEADER = errors.New("Contact leader")
var ERR_NO_LEADER_ELECTED = errors.New("No leader is elected yet...Try later")

func handleFollowerAppend(sm *StateMachine, cmd *Append) []interface{} {
	initialiseActions()
	var x []byte
	actions = append(actions, Commit{Index: -1, LeaderId: sm.leaderId, Data: x, Err: ERR_NOT_LEADER})
//	fmt.Println("Inside follwer Append")
	return actions
}

func handleCandidateAppend(sm *StateMachine, cmd *Append) []interface{} {
	initialiseActions()
	var x []byte
	actions = append(actions, Commit{Index: -1, LeaderId: -1, Data: x, Err: ERR_NO_LEADER_ELECTED})
//	fmt.Println("Inside candidate Append")
	return actions
}

func handleLeaderAppend(sm *StateMachine, cmd *Append) []interface{} {
	var temp Log
	var prevLogTerm int
	initialiseActions()
	sm.lastLogIndex = sm.lastLogIndex + 1
	sm.lastLogTerm = sm.term
	temp.Term = sm.term
	temp.Command = cmd.Data
	sm.log = append(sm.log, temp)
	
	actions = append(actions, LogStore{Index: sm.lastLogIndex, LogEntry: temp})
			for i := 0; i < len(sm.peers); i++ {
					prevLogIndex := int(sm.nextIndex[sm.peers[i]]) - 1
					if prevLogIndex < 0 {
							prevLogTerm = 0
					} else if prevLogIndex >= len(sm.log) {
							continue
					} else{
							prevLogTerm = sm.log[prevLogIndex].Term
					}

					if sm.nextIndex[sm.peers[i]] > len(sm.log) {
					continue
					}
					entries := sm.log[sm.nextIndex[sm.peers[i]]:]

					actions = append(actions, Send{PeerId: sm.peers[i], Event: AppendEntriesReqEv{Term: sm.term, SenderId: sm.id, PrevLogIndex: int64(prevLogIndex),
					PrevLogTerm: int64(prevLogTerm), Entries: entries, SenderCommitIndex: sm.commitIndex}})
			}
//	fmt.Println("Inside leader Append")
	return actions
}
