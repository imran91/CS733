package main

import (
//	"fmt"
)

func handleFollowerVoteResp(sm *StateMachine, cmd *VoteRespEv) []interface{} {
	initialiseActions()
	if sm.term < cmd.SenderTerm {
		actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)}) //equivalent to random(timer,2*timer)
		sm.term = cmd.SenderTerm
		sm.votedFor = -1
		sm.lastMatchIndex = -1
		actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor,LastMatchIndex: sm.lastMatchIndex})
	}
//	fmt.Println("Inside follwer VoteResp")
	return actions
}

func handleCandidateVoteResp(sm *StateMachine, cmd *VoteRespEv) []interface{} {
	var totalNumPosVotes int
	var totalNumNegVotes int
	var totalServers int
	var prevLogTerm int64
	initialiseActions()
	totalNumPosVotes = 1
	totalNumNegVotes = 0
	totalServers = len(sm.peers) + 1
	for i := 0; i < totalServers-1; i++ {
		if sm.votedAs[sm.peers[i]] == 1 { //counting all ACKs
			totalNumPosVotes++
		} else if sm.votedAs[sm.peers[i]] == -1 {
			totalNumNegVotes++
		}
	}

	if sm.term < cmd.SenderTerm {
		sm.state = 1
		sm.term = cmd.SenderTerm
		sm.votedFor = -1
		sm.lastMatchIndex = -1
		actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)}) //equivalent to random(timer,2*timer)
		actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor,LastMatchIndex: sm.lastMatchIndex})
	} else if sm.term == cmd.SenderTerm {
		if cmd.Response == true {
			if sm.votedAs[cmd.SenderId] == 0 {
					sm.votedAs[cmd.SenderId] = 1
					totalNumPosVotes++
				}
			if totalNumPosVotes > (sm.clusterSize / 2) {
				//elected as leader
				sm.state = 3
				sm.leaderId = sm.id
				sm.nextIndex = make(map[int]int)
				sm.matchIndex = make(map[int]int)

					for i:=0;i< totalServers-1;i++ {
						sm.nextIndex[i] = int(sm.lastLogIndex) + 1
						sm.matchIndex[i] = -1
					}
				actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)}) //equivalent to random(timer,2*timer)
				for i := 0; i < len(sm.peers); i++ {
					prevLogIndex := int(sm.nextIndex[sm.peers[i]]) - 1
					if prevLogIndex < 0 {
							prevLogTerm = 0
					} else if prevLogIndex >= len(sm.log) {
							continue
					} else{
							prevLogTerm = int64(sm.log[prevLogIndex].Term)
					}

					if sm.nextIndex[sm.peers[i]] > len(sm.log) {
					continue
					}
					entries := sm.log[sm.nextIndex[sm.peers[i]]:]

					actions = append(actions, Send{PeerId: sm.peers[i], Event: AppendEntriesReqEv{Term: sm.term, SenderId: sm.id, PrevLogIndex: int64(prevLogIndex),
					PrevLogTerm: int64(prevLogTerm), Entries: entries, SenderCommitIndex: sm.commitIndex}})
				}
			}		
		} else if cmd.Response == false {
				if sm.votedAs[cmd.SenderId] == 0 {
					sm.votedAs[cmd.SenderId] = 1
					totalNumPosVotes++
				}
				if totalNumNegVotes > (sm.clusterSize / 2) {
					sm.state = 1
					actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)}) //equivalent to random(timer,2*timer)
				}

		}

	}
	
//	fmt.Println("Inside candidate VoteResp")
	return actions
}

func handleLeaderVoteResp(sm *StateMachine, cmd *VoteRespEv) []interface{} {
	initialiseActions()
	if sm.term < cmd.SenderTerm {
		actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)}) //equivalent to random(timer,2*timer)
		sm.term = cmd.SenderTerm
		sm.votedFor = -1
		sm.lastMatchIndex = -1
		actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor,LastMatchIndex: sm.lastMatchIndex})
	}
//	fmt.Println("Inside leader VoteResp")
	return actions
}
