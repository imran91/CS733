package main
/*import (
	"fmt"
)*/

func handleFollowerAppendEntryResp(sm *StateMachine,cmd *AppendEntriesRespEv) []interface{}{
	initialiseActions()
	return actions
}

func handleCandidateAppendEntryResp(sm *StateMachine,cmd *AppendEntriesRespEv) []interface{}{
	initialiseActions()
	return actions
}


func handleLeaderAppendEntryResp(sm *StateMachine,cmd *AppendEntriesRespEv) []interface{}{
	var ind int
	initialiseActions()	
	if sm.term < cmd.senderTerm {
		sm.state = 1
		sm.votedFor = -1
		sm.votedAs = make([]int,len(sm.peers))
		sm.term = cmd.senderTerm
		actions = append(actions,StateStore{currTerm:sm.term,votedFor:sm.votedFor})
		return actions
	}
	ind = findSenderIndex(sm.peers,cmd.senderId)
	sm.nextIndex[ind] = cmd.lastMatchIndex + 1
	
	if cmd.response == true {
		sm.matchIndex[ind] = cmd.lastMatchIndex

		isCommited := checkCommit(sm.matchIndex,cmd.senderId,&sm.commitIndex,sm.peers,sm.log,sm.term)
		if isCommited{
				actions = append(actions,Commit{index:sm.commitIndex,leaderId:sm.leaderId,data:sm.log[sm.commitIndex].command,err:nil})
			}
	}else {

		prevLogIndex := sm.nextIndex[ind]-1
		prevLogTerm := sm.log[prevLogIndex].term
		entries := sm.log[sm.nextIndex[ind]:sm.lastLogIndex]

		actions = append(actions,Send{peerId:cmd.senderId,event:AppendEntriesReqEv{term : sm.term, senderId: sm.id, prevLogIndex: prevLogIndex, 
				prevLogTerm: prevLogTerm, entries: entries,senderCommitIndex: sm.commitIndex}})		
	
		if sm.lastLogIndex > sm.matchIndex[ind]{
			prevLogIndex := sm.nextIndex[ind]-1
			prevLogTerm := sm.log[prevLogIndex].term
			entries := sm.log[sm.nextIndex[ind]:sm.lastLogIndex]
			//sm.sentIndex[cmd.senderId] = sm.lastLogIndex
	
			actions = append(actions,Send{peerId:cmd.senderId,event:AppendEntriesReqEv{term : sm.term, senderId: sm.id, prevLogIndex: prevLogIndex, 
					prevLogTerm: prevLogTerm, entries: entries,senderCommitIndex: sm.commitIndex}})		
		}
	}
	
	return actions
}

func checkCommit(matchIndex []int,senderId int,commitIndex *int,peers []int,log []Log,term int) bool{
	cnt :=0
	var k int
	for j := matchIndex[senderId]; j > *commitIndex; j-- {
		for k =0;k<=len(peers)+1;k++{
			if matchIndex[k] >= j{
				cnt++
			}
		}
		if (cnt > (len(peers)+1)/2) && log[j].term == term {
			*commitIndex = j
			return true
			break
		}
		cnt =0
	}
	return false
}

func findSenderIndex(peers []int,senderId int) int {
	var ind int
	ind = -1
	for i:=0; i<len(peers); i++ {
			if peers[i] == senderId {
				ind = i
				return ind
			}
		}
	return ind
}