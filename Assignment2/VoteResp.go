package main
/*import (
	"fmt"
)*/

func handleFollowerVoteResp(sm *StateMachine,cmd *VoteRespEv) []interface{}{

	if sm.term < cmd.senderTerm {
		sm.term = cmd.senderTerm
		actions = append(actions,StateStore{currTerm:sm.term,votedFor:sm.votedFor})
	}

	return actions
}

func handleCandidateVoteResp(sm *StateMachine,cmd *VoteRespEv) []interface{}{
	
	if sm.term < cmd.senderTerm {
		sm.term = cmd.senderTerm
		sm.state = 1
	}
	return actions
}


func handleLeaderVoteResp(sm *StateMachine,cmd *VoteRespEv) []interface{}{
	if sm.term < cmd.senderTerm {
		sm.term = cmd.senderTerm
		sm.state = 1
	}

	return actions
}