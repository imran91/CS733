package main

import (
//		"fmt"
)

func handleFollowerTimeout(sm *StateMachine, cmd *Timeout) []interface{} {
	var totalServers int
	initialiseActions()
	actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)}) //equivalent to random(timer,2*timer)
	totalServers = len(sm.peers) + 1
	sm.state = 2
	sm.term = sm.term + 1
	sm.lastMatchIndex = -1

	sm.votedAs = make(map[int]int, len(sm.peers)+1)
	for i := 0; i < totalServers-1; i++ {
		sm.votedAs[sm.peers[i]] = 0
	}
	
	sm.votedAs[sm.id] = 1 //self vote //last index of votedAs is for self voting
	sm.votedFor = sm.id
	actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor})
	
	for i := 0; i < (len(sm.peers)); i++ {
		actions = append(actions, Send{PeerId: sm.peers[i], Event: VoteReqEv{SenderId: sm.id,
			Term: sm.term, LastLogIndex: int(sm.lastLogIndex), LastLogTerm: sm.lastLogTerm}})
	}
//	fmt.Println(sm.id, "Inside follwer Timeout")
	return actions
}

func handleCandidateTimeout(sm *StateMachine, cmd *Timeout) []interface{} {
	var totalServers int
	initialiseActions()
	actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)}) //equivalent to random(timer,2*timer)
	totalServers = len(sm.peers) + 1
	sm.state = 2
	sm.term = sm.term + 1
	sm.lastMatchIndex = -1

	sm.votedAs = make(map[int]int, len(sm.peers)+1)
	for i := 0; i < totalServers-1; i++ {
		sm.votedAs[sm.peers[i]] = 0
	}
	
	sm.votedAs[len(sm.peers)] = 1 //self vote //last index of votedAs is for self voting
	sm.votedFor = sm.id
	actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor})
	
	for i := 0; i < (len(sm.peers)); i++ {
		actions = append(actions, Send{PeerId: sm.peers[i], Event: VoteReqEv{SenderId: sm.id,
			Term: sm.term, LastLogIndex: int(sm.lastLogIndex), LastLogTerm: sm.lastLogTerm}})
	}
//	fmt.Println(sm.id, "Inside candidate Timeout")
	return actions
}

func handleLeaderTimeout(sm *StateMachine, cmd *Timeout) []interface{} {
	var prevLogTerm int
	initialiseActions()
	actions = append(actions, Alarm{T: sm.heartbeatAlarm}) 

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
	//fmt.Println(sm.id, "Inside leader Timeout")
	return actions
}
