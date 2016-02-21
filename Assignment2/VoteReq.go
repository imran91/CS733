package main
/*import (
	"fmt"
)*/

func handleFollowerVoteReq(sm *StateMachine,cmd *VoteReqEv) []interface{}{
	
	if sm.term > cmd.term {
		actions = append(actions,Send{peerId:cmd.senderId,event:VoteRespEv{senderTerm:sm.term,response:false}})
		return actions		
	}

	if (sm.lastLogTerm > cmd.lastLogTerm) || ((sm.lastLogTerm == cmd.lastLogTerm) && sm.lastLogIndex > cmd.lastLogIndex){
		actions = append(actions,Send{peerId:cmd.senderId,event:VoteRespEv{senderTerm:sm.term,response:false}})
		return actions	
	}

	if (sm.votedFor != 0 && sm.votedFor!= cmd.senderId){
		actions = append(actions,Send{peerId:cmd.senderId,event:VoteRespEv{senderTerm:sm.term,response:false}})
		return actions
	}

	sm.votedFor = cmd.senderId
	actions = append(actions,Send{peerId:cmd.senderId,event:VoteRespEv{senderTerm:sm.term,response:true}})
	return actions
}


func handleCandidateVoteReq(sm *StateMachine,cmd *VoteReqEv) []interface{}{

	if sm.term < cmd.term {
		sm.state = 1
		sm.state = cmd.term
		if (sm.lastLogTerm > cmd.lastLogTerm) || ((sm.lastLogTerm == cmd.lastLogTerm) && sm.lastLogIndex > cmd.lastLogIndex){
			actions = append(actions,Send{peerId:cmd.senderId,event:VoteRespEv{senderTerm:sm.term,response:false}})
		return actions	
		}
		sm.votedFor = cmd.senderId
		actions = append(actions,Send{peerId:cmd.senderId,event:VoteRespEv{senderTerm:sm.term,response:true}})
		return actions
	}
	
	actions = append(actions,Send{peerId:cmd.senderId,event:VoteRespEv{senderTerm:sm.term,response:false}})
	return actions
}


func handleLeaderVoteReq(sm *StateMachine,cmd *VoteReqEv) []interface{}{
	
	if sm.term < cmd.term {
		sm.state = 1
		sm.state = cmd.term
		if (sm.lastLogTerm > cmd.lastLogTerm) || ((sm.lastLogTerm == cmd.lastLogTerm) && sm.lastLogIndex > cmd.lastLogIndex){
			actions = append(actions,Send{peerId:cmd.senderId,event:VoteRespEv{senderTerm:sm.term,response:false}})
		return actions	
		}
		sm.votedFor = cmd.senderId
		actions = append(actions,Send{peerId:cmd.senderId,event:VoteRespEv{senderTerm:sm.term,response:true}})
		return actions
	}
	
	actions = append(actions,Send{peerId:cmd.senderId,event:VoteRespEv{senderTerm:sm.term,response:false}})
	return actions
	
}