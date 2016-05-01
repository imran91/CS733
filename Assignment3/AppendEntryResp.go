package main

import (
	"fmt"
)

func handleFollowerAppendEntryResp(sm *StateMachine, cmd *AppendEntriesRespEv) []interface{} {
	initialiseActions()
	return actions
}

func handleCandidateAppendEntryResp(sm *StateMachine, cmd *AppendEntriesRespEv) []interface{} {
	initialiseActions()
	return actions
}

func handleLeaderAppendEntryResp(sm *StateMachine, cmd *AppendEntriesRespEv) []interface{} {
	var prevLogIndex,prevLogTerm int
	initialiseActions()
	if sm.term < cmd.SenderTerm {
		sm.state = 1
		sm.votedFor = -1
		sm.term = cmd.SenderTerm
		sm.lastMatchIndex = -1
		actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)}) //equivalent to random(timer,2*timer)
		actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor, LastMatchIndex: sm.lastMatchIndex})
	} else if cmd.Response == true {
		sm.matchIndex[cmd.SenderId] = int(max(int64(sm.matchIndex[cmd.SenderId]), cmd.LastMatchIndex)) //handling reordering
		sm.nextIndex[cmd.SenderId] = sm.matchIndex[cmd.SenderId] + 1
		
		for j,cnt := sm.matchIndex[cmd.SenderId],1; j > int(sm.commitIndex); j-- {
			for k := 0; k < len(sm.peers); k++ {
				if sm.matchIndex[sm.peers[k]] >= j {
					cnt++
					}
				}
				if (cnt > sm.clusterSize/2) && (sm.log[j].Term == sm.term) {
					sm.commitIndex = int64(j)
					actions = append(actions, Commit{Index: int(sm.commitIndex), LeaderId: sm.id, Data: sm.log[sm.commitIndex].Command, Err: nil})
					break
				}
				cnt = 1
			}

		/*isCommited := checkCommit(sm.matchIndex, cmd.SenderId, &sm.commitIndex, sm.peers, sm.log, sm.term)
		if isCommited {
			actions = append(actions, Commit{Index: sm.id, LeaderId: sm.leaderId, Data: sm.log[sm.commitIndex], Err: nil})
		}*/

		if sm.lastLogIndex > int64(sm.matchIndex[cmd.SenderId]) {
			for i := 0; i < 1; i++ {
					prevLogIndex = int(sm.nextIndex[sm.peers[i]]) - 1
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

					actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesReqEv{Term: sm.term, SenderId: sm.id, PrevLogIndex: int64(prevLogIndex),
					PrevLogTerm: int64(prevLogTerm), Entries: entries, SenderCommitIndex: sm.commitIndex}})
			}
		}

	} else if cmd.Response == false {
			sm.nextIndex[cmd.SenderId]--
			for i := 0; i < 1; i++ {
					prevLogIndex = int(sm.nextIndex[sm.peers[i]]) - 1
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

					actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesReqEv{Term: sm.term, SenderId: sm.id, PrevLogIndex: int64(prevLogIndex),
					PrevLogTerm: int64(prevLogTerm), Entries: entries, SenderCommitIndex: sm.commitIndex}})
			}
	}

	//fmt.Println("Inside leader AppendEntryResp")	
	return actions
}

func checkCommit(matchIndex map[int]int, senderId int, commitIndex *int64, peers []int, log []Log, term int) bool {
	cnt := 1
	var k int
	var j int

	for j = matchIndex[senderId]; j > int(*commitIndex); j-- {
		for k = 0; k < len(peers); k++ {
			fmt.Println("inside k loop",)
			if matchIndex[peers[k]] >= j {
				fmt.Println("inside if cond",)
				cnt++
			}
		}

		if (cnt > (len(peers)+1)/2) && log[j].Term == term {
			*commitIndex = int64(j)
			return true
		}
		cnt = 1
	}
	return false
}

