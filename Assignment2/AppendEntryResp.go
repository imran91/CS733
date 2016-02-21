package main
/*import (
	"fmt"
)*/

func handleFollowerAppendEntryResp(sm *StateMachine,cmd *AppendEntriesRespEv) []interface{}{
	
	return actions
}

func handleCandidateAppendEntryResp(sm *StateMachine,cmd *AppendEntriesRespEv) []interface{}{
	
	return actions
}


func handleLeaderAppendEntryResp(sm *StateMachine,cmd *AppendEntriesRespEv) []interface{}{
		
	if sm.term < cmd.senderTerm {
		sm.state = 1
		sm.term = cmd.senderTerm
		return actions
	}

	if cmd.response == true {
		sm.nextIndex[cmd.senderId] = sm.sentIndex[cmd.senderId]+1
		sm.matchIndex[cmd.senderId] = sm.sentIndex[cmd.senderId]

		checkCommit(sm.matchIndex,cmd.senderId,&sm.commitIndex,sm.peers,sm.log,sm.term)
	}else {
		sm.nextIndex[cmd.senderId] = sm.nextIndex[cmd.senderId]-1
		prevLogIndex := sm.nextIndex[cmd.senderId]-1
		prevLogTerm := sm.log[prevLogIndex].term
		entries := sm.log[sm.nextIndex[cmd.senderId]:sm.lastLogIndex]
		sm.sentIndex[cmd.senderId] = sm.lastLogIndex

		actions = append(actions,Send{peerId:cmd.senderId,event:AppendEntriesReqEv{term : sm.term, senderId: sm.id, prevLogIndex: prevLogIndex, 
				prevLogTerm: prevLogTerm, entries: entries,senderCommitIndex: sm.commitIndex}})		
	}

	if sm.lastLogIndex > sm.matchIndex[cmd.senderId]{
		prevLogIndex := sm.nextIndex[cmd.senderId]-1
		prevLogTerm := sm.log[prevLogIndex].term
		entries := sm.log[sm.nextIndex[cmd.senderId]:sm.lastLogIndex]
		sm.sentIndex[cmd.senderId] = sm.lastLogIndex

		actions = append(actions,Send{peerId:cmd.senderId,event:AppendEntriesReqEv{term : sm.term, senderId: sm.id, prevLogIndex: prevLogIndex, 
				prevLogTerm: prevLogTerm, entries: entries,senderCommitIndex: sm.commitIndex}})		
	}

	return actions
}

func checkCommit(matchIndex []int,senderId int,commitIndex *int,peers []int,log []Log,term int){
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
			break
		}
		cnt =0
	}

}