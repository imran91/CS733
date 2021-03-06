package main

/*import (
	"fmt"
)*/

func handleFollowerVoteReq(sm *StateMachine, cmd *VoteReqEv) []interface{} {
	initialiseActions()
	if sm.term > cmd.term {
		actions = append(actions, Send{peerId: cmd.senderId, event: VoteRespEv{senderId: sm.id, senderTerm: sm.term, response: false}})
		return actions
	}

	if (sm.lastLogTerm > cmd.lastLogTerm) || ((sm.lastLogTerm == cmd.lastLogTerm) && sm.lastLogIndex > cmd.lastLogIndex) {
		actions = append(actions, Send{peerId: cmd.senderId, event: VoteRespEv{senderId: sm.id, senderTerm: sm.term, response: false}})
		return actions
	}

	if sm.votedFor != 0 && sm.votedFor != cmd.senderId {
		actions = append(actions, Send{peerId: cmd.senderId, event: VoteRespEv{senderId: sm.id, senderTerm: sm.term, response: false}})
		return actions
	}

	sm.votedFor = cmd.senderId
	sm.term = cmd.term
	actions = append(actions, StateStore{currTerm: sm.term, votedFor: sm.votedFor})
	actions = append(actions, Send{peerId: cmd.senderId, event: VoteRespEv{senderId: sm.id, senderTerm: sm.term, response: true}})
	return actions
}

func handleCandidateVoteReq(sm *StateMachine, cmd *VoteReqEv) []interface{} {
	initialiseActions()
	if sm.term < cmd.term {
		sm.state = 1
		sm.term = cmd.term
		sm.votedFor = -1
		sm.votedAs = make([]int, len(sm.peers)+1)
		actions = append(actions, StateStore{currTerm: sm.term, votedFor: sm.votedFor})

		if (sm.lastLogTerm > cmd.lastLogTerm) || ((sm.lastLogTerm == cmd.lastLogTerm) && sm.lastLogIndex > cmd.lastLogIndex) {
			actions = append(actions, Send{peerId: cmd.senderId, event: VoteRespEv{senderId: sm.id, senderTerm: sm.term, response: false}})
			return actions
		}
		sm.votedFor = cmd.senderId
		actions = append(actions, StateStore{currTerm: sm.term, votedFor: sm.votedFor})
		actions = append(actions, Send{peerId: cmd.senderId, event: VoteRespEv{senderId: sm.id, senderTerm: sm.term, response: true}})
		return actions
	}

	actions = append(actions, Send{peerId: cmd.senderId, event: VoteRespEv{senderId: sm.id, senderTerm: sm.term, response: false}})
	return actions
}

func handleLeaderVoteReq(sm *StateMachine, cmd *VoteReqEv) []interface{} {
	initialiseActions()
	if sm.term < cmd.term {
		sm.state = 1
		sm.term = cmd.term
		sm.votedFor = -1
		sm.votedAs = make([]int, len(sm.peers)+1)
		actions = append(actions, StateStore{currTerm: sm.term, votedFor: sm.votedFor})

		if (sm.lastLogTerm > cmd.lastLogTerm) || ((sm.lastLogTerm == cmd.lastLogTerm) && sm.lastLogIndex > cmd.lastLogIndex) {
			actions = append(actions, Send{peerId: cmd.senderId, event: VoteRespEv{senderId: sm.id, senderTerm: sm.term, response: false}})
			return actions
		}
		sm.votedFor = cmd.senderId
		actions = append(actions, StateStore{currTerm: sm.term, votedFor: sm.votedFor})
		actions = append(actions, Send{peerId: cmd.senderId, event: VoteRespEv{senderId: sm.id, senderTerm: sm.term, response: true}})
		return actions
	}

	actions = append(actions, Send{peerId: cmd.senderId, event: VoteRespEv{senderId: sm.id, senderTerm: sm.term, response: false}})
	return actions

}
