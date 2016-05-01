package main

import (
//	"fmt"
)

func handleFollowerVoteReq(sm *StateMachine, cmd *VoteReqEv) []interface{} {
	initialiseActions()
	if sm.term < cmd.Term{
		sm.term = cmd.Term
		sm.lastMatchIndex = -1
		sm.votedFor = -1
		actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor, LastMatchIndex: sm.lastMatchIndex})	

	}

	if sm.term > cmd.Term {
		actions = append(actions, Send{PeerId: cmd.SenderId, Event: VoteRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: false}})

	} else if (sm.lastLogTerm > cmd.LastLogTerm) || ((sm.lastLogTerm == cmd.LastLogTerm) && int(sm.lastLogIndex) > cmd.LastLogIndex) {
		actions = append(actions, Send{PeerId: cmd.SenderId, Event: VoteRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: false}})

	}else if sm.votedFor != -1 && sm.votedFor != cmd.SenderId {
		actions = append(actions, Send{PeerId: cmd.SenderId, Event: VoteRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: false}})

	} else {
		actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)}) //equivalent to random(timer,2*timer)
		sm.votedFor = cmd.SenderId
		actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor, LastMatchIndex: sm.lastMatchIndex})	
		actions = append(actions, Send{PeerId: cmd.SenderId, Event: VoteRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: true}})
	}

//	fmt.Println(sm.id, "Inside follwer VoteReq",actions)
	return actions
}

func handleCandidateVoteReq(sm *StateMachine, cmd *VoteReqEv) []interface{} {
	initialiseActions()
	if sm.term < cmd.Term {
		sm.state = 1
		sm.term = cmd.Term
		sm.votedFor = -1
		sm.lastMatchIndex = -1
		actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)}) //equivalent to random(timer,2*timer)
		actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor,LastMatchIndex: sm.lastMatchIndex})

		if (sm.lastLogTerm > cmd.LastLogTerm) || ((sm.lastLogTerm == cmd.LastLogTerm) && int(sm.lastLogIndex) > cmd.LastLogIndex) {
			actions = append(actions, Send{PeerId: cmd.SenderId, Event: VoteRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: false}})
			return actions
		} else {
			sm.votedFor = cmd.SenderId
			actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor,LastMatchIndex: sm.lastMatchIndex})
			actions = append(actions, Send{PeerId: cmd.SenderId, Event: VoteRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: true}})
		}
	} else {
			actions = append(actions, Send{PeerId: cmd.SenderId, Event: VoteRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: false}})	
	}

//	fmt.Println(sm.id, "Inside candidate VoteReq")
	return actions
}

func handleLeaderVoteReq(sm *StateMachine, cmd *VoteReqEv) []interface{} {
	initialiseActions()
	if sm.term < cmd.Term {
		sm.state = 1
		sm.term = cmd.Term
		sm.votedFor = -1
		sm.lastMatchIndex = -1
		actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)}) //equivalent to random(timer,2*timer)
		actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor,LastMatchIndex: sm.lastMatchIndex})

		if (sm.lastLogTerm > cmd.LastLogTerm) || ((sm.lastLogTerm == cmd.LastLogTerm) && int(sm.lastLogIndex) > cmd.LastLogIndex) {
			actions = append(actions, Send{PeerId: cmd.SenderId, Event: VoteRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: false}})
			return actions
		} else {
			sm.votedFor = cmd.SenderId
			actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor,LastMatchIndex: sm.lastMatchIndex})
			actions = append(actions, Send{PeerId: cmd.SenderId, Event: VoteRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: true}})
		}
	} else {
			actions = append(actions, Send{PeerId: cmd.SenderId, Event: VoteRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: false}})	
	}

//	fmt.Println(sm.id, "Inside Leader VoteReq")
	return actions

}
