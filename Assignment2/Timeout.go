package main
import (
	//"fmt"
	"math/rand"
)


func handleFollowerTimeout(sm *StateMachine,cmd *Timeout) []interface{}{
	var totalServers int
	totalServers = len(sm.peers)+1
	sm.state = 2
	sm.term = sm.term + 1
	sm.votedFor = -1
	sm.votedAs = make(map[int]int)
	actions = append(actions,StateStore{currTerm:sm.term,votedFor:sm.votedFor})

	for i:=0; i<totalServers;i++{
		sm.votedAs[i] = 0
	}
	sm.votedAs[sm.id] = 1//self vote
	sm.votedFor = sm.id
	actions = append(actions,StateStore{currTerm:sm.term,votedFor:sm.votedFor})
	actions = append(actions,Alarm{t:rand.Intn(2*sm.timer-sm.timer)+sm.timer})  //equivalent to random(timer,2*timer)

	for i:=0; i<(len(sm.peers)); i++ {
		actions= append(actions,Send{peerId:sm.peers[i],event:VoteReqEv{senderId:sm.id,
			term:sm.term,lastLogIndex:sm.lastLogIndex,lastLogTerm:sm.lastLogTerm}})
	}

	return actions
}

func handleCandidateTimeout(sm *StateMachine,cmd *Timeout) []interface{}{
	var totalServers int
	totalServers = len(sm.peers)+1
	sm.term = sm.term + 1
	sm.votedFor = -1
	sm.votedAs = make(map[int]int)
	actions = append(actions,StateStore{currTerm:sm.term,votedFor:sm.votedFor})
	//numVotes=1   new logic need to be added
	for i:=0; i<totalServers;i++{
		sm.votedAs[i] = 0
	}
	sm.votedAs[sm.id] = 1//self vote
	sm.votedFor = sm.id
	actions = append(actions,StateStore{currTerm:sm.term,votedFor:sm.votedFor})
	actions = append(actions,Alarm{t:rand.Intn(2*sm.timer-sm.timer)+sm.timer}) 

	for i:=0; i<(len(sm.peers)); i++ {
		actions= append(actions,Send{peerId:sm.peers[i],event:VoteReqEv{senderId:sm.id,
			term:sm.term,lastLogIndex:sm.lastLogIndex,lastLogTerm:sm.lastLogTerm}})
	}
	
	return actions
}


func handleLeaderTimeout(sm *StateMachine,cmd *Timeout) []interface{}{

	actions = append(actions,Alarm{t:rand.Intn(sm.timer)}) 	//change this
	var temp []Log

	for i:=0; i<(len(sm.peers)); i++ {
			prevLogIndex := sm.nextIndex[i]-1
			prevLogTerm := sm.log[prevLogIndex].term
			entries := temp

			actions = append(actions,Send{peerId:sm.peers[i],event:AppendEntriesReqEv{term : sm.term, senderId: sm.id, 
				prevLogIndex: prevLogIndex, prevLogTerm: prevLogTerm, entries: entries,senderCommitIndex: sm.commitIndex}})
	}
	
	return actions
}