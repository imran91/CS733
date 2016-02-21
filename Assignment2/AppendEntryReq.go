package main
import (
	//"fmt"
	"math/rand"
	"math"
)
//var actions []interface{}
func handleFollowerAppendEntryReq(sm *StateMachine,cmd *AppendEntriesReqEv) []interface{}{
	
	if sm.term > cmd.term {
		actions = append(actions,Send{peerId:cmd.senderId,event:AppendEntriesRespEv{senderId: sm.id , senderTerm: sm.term, response:false,lastMatchIndex:sm.matchIndex[sm.id]}})
		return actions		
	}
	actions = append(actions,Alarm{t:rand.Intn(2*sm.timer-sm.timer)+sm.timer})

	sm.term = cmd.term
	actions = append(actions,StateStore{currTerm:sm.term,votedFor:sm.votedFor})
	sm.leaderId = cmd.senderId

	if sm.log[cmd.prevLogIndex].term != cmd.prevLogTerm{

		actions = append(actions,Send{peerId:cmd.senderId,event:AppendEntriesRespEv{senderId: sm.id , senderTerm: sm.term, response:false,lastMatchIndex:sm.matchIndex[sm.id]}})
		return actions
	}

	actions = append(actions,LogStore{index:cmd.prevLogIndex+1,logEntry:cmd.entries})

	sm.lastLogIndex = cmd.prevLogIndex + len(cmd.entries)
	sm.lastLogTerm = sm.log[sm.lastLogIndex].term
	sm.commitIndex = int(math.Min(float64(sm.lastLogIndex),float64(cmd.senderCommitIndex)))
	actions = append(actions,Send{peerId:cmd.senderId,event:AppendEntriesRespEv{senderId: sm.id , senderTerm: sm.term, response:true,lastMatchIndex:sm.lastLogIndex}})
	return actions
}

func handleCandidateAppendEntryReq(sm *StateMachine,cmd *AppendEntriesReqEv) []interface{}{
	if sm.term > cmd.term {
		actions = append(actions,Send{peerId:cmd.senderId,event:AppendEntriesRespEv{senderId: sm.id , senderTerm: sm.term, response:false,lastMatchIndex:sm.matchIndex[sm.id]}})
		return actions		
	}

	actions = append(actions,Alarm{t:rand.Intn(2*sm.timer-sm.timer)+sm.timer})

	sm.state = 1
	sm.term = cmd.term
	actions = append(actions,StateStore{currTerm:sm.term,votedFor:sm.votedFor})
	sm.leaderId = cmd.senderId

	if sm.log[cmd.prevLogIndex].term != cmd.prevLogTerm{
		actions = append(actions,Send{peerId:cmd.senderId,event:AppendEntriesRespEv{senderId: sm.id , senderTerm: sm.term, response:false,lastMatchIndex:sm.matchIndex[sm.id]}})
		return actions
	}

	actions = append(actions,LogStore{index:cmd.prevLogIndex+1,logEntry:cmd.entries})

	sm.lastLogIndex = cmd.prevLogIndex + len(cmd.entries)
	sm.lastLogTerm = sm.log[sm.lastLogIndex].term
	sm.commitIndex = int(math.Min(float64(sm.lastLogIndex),float64(cmd.senderCommitIndex)))

	actions = append(actions,Send{peerId:cmd.senderId,event:AppendEntriesRespEv{senderId: sm.id , senderTerm: sm.term, response:true,lastMatchIndex:sm.lastLogIndex}})
	return actions
}


func handleLeaderAppendEntryReq(sm *StateMachine,cmd *AppendEntriesReqEv) []interface{}{
	if sm.term > cmd.term {
		actions = append(actions,Send{peerId:cmd.senderId,event:AppendEntriesRespEv{senderId: sm.id , senderTerm: sm.term, response:false,lastMatchIndex:sm.matchIndex[sm.id]}})
		return actions		
	}

	actions = append(actions,Alarm{t:rand.Intn(2*sm.timer-sm.timer)+sm.timer})

	sm.state = 1
	sm.term = cmd.term
	actions = append(actions,StateStore{currTerm:sm.term,votedFor:sm.votedFor})
	sm.leaderId = cmd.senderId

	if sm.log[cmd.prevLogIndex].term != cmd.prevLogTerm{
		actions = append(actions,Send{peerId:cmd.senderId,event:AppendEntriesRespEv{senderId: sm.id , senderTerm: sm.term, response:false,lastMatchIndex:sm.matchIndex[sm.id]}})
		return actions
	}

	actions = append(actions,LogStore{index:cmd.prevLogIndex+1,logEntry:cmd.entries})

	sm.lastLogIndex = cmd.prevLogIndex + len(cmd.entries)
	sm.lastLogTerm = sm.log[sm.lastLogIndex].term
	sm.commitIndex = int(math.Min(float64(sm.lastLogIndex),float64(cmd.senderCommitIndex)))

	actions = append(actions,Send{peerId:cmd.senderId,event:AppendEntriesRespEv{senderId: sm.id , senderTerm: sm.term, response:true,lastMatchIndex:sm.lastLogIndex}})

	return actions
}