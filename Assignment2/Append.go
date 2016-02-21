package main
import (
	"errors"
	//"fmt"
)
var ERR_NOT_LEADER = errors.New("Contact leader")
var ERR_NO_LEADER_ELECTED = errors.New("No leader is elected yet...Try later")
func handleFollowerAppend(sm *StateMachine,cmd *Append) []interface{}{
	
	actions = append(actions,Commit{index:-1,leaderId:sm.leaderId,data:cmd.data,err:ERR_NOT_LEADER})
	return actions
}

func handleCandidateAppend(sm *StateMachine,cmd *Append) []interface{}{
	actions = append(actions,Commit{index:-1,leaderId:sm.leaderId,data:cmd.data,err:ERR_NO_LEADER_ELECTED})
	return actions
}

func handleLeaderAppend(sm *StateMachine,cmd *Append) []interface{}{
	sm.lastLogIndex = sm.lastLogIndex+1
	sm.lastLogTerm = sm.term
	var temp []Log
	temp[0].term = sm.term
	temp[0].command = cmd.data
	actions = append(actions,LogStore{index:sm.lastLogIndex,logEntry:temp})

	for i:=0; i<(len(sm.peers)); i++ {
			prevLogIndex := sm.nextIndex[sm.peers[i]]-1
			prevLogTerm := sm.log[prevLogIndex].term
			entries := sm.log[sm.nextIndex[sm.peers[i]]:]
			//sm.sentIndex[sm.peers[i]] = sm.nextIndex[]

			actions = append(actions,Send{peerId:sm.peers[i],event:AppendEntriesReqEv{term : sm.term, senderId: sm.id, prevLogIndex: prevLogIndex, 
				prevLogTerm: prevLogTerm, entries: entries,senderCommitIndex: sm.commitIndex}})
	}
	
	return actions
}