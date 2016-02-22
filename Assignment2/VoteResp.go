package main

func handleFollowerVoteResp(sm *StateMachine,cmd *VoteRespEv) []interface{}{
	initialiseActions()
	if sm.term < cmd.senderTerm {
		sm.term = cmd.senderTerm
		sm.votedFor = -1
		sm.votedAs = make([]int,len(sm.peers)+1)
		actions = append(actions,StateStore{currTerm:sm.term,votedFor:sm.votedFor})
	}

	return actions
}

func handleCandidateVoteResp(sm *StateMachine,cmd *VoteRespEv) []interface{}{
	var totalNumPosVotes int
	var totalNumNegVotes int
	var totalServers int
	var temp []Log
	initialiseActions()
	totalNumPosVotes = 0
	totalNumPosVotes = 0 
	totalServers = len(sm.peers)+1
	if sm.term < cmd.senderTerm {
		sm.term = cmd.senderTerm
		sm.votedFor = -1
		sm.votedAs = make([]int,len(sm.peers)+1)
		actions = append(actions,StateStore{currTerm:sm.term,votedFor:sm.votedFor})
		sm.state = 1
	}

	if cmd.response == true {
		sm.votedAs[cmd.senderId] = 1	
	} else {
		sm.votedAs[cmd.senderId] = -1	
	}
	
	for i:=0; i<totalServers;i++{
		if sm.votedAs[i] == 1{	//counting all ACKs
			totalNumPosVotes++
		} else if sm.votedAs[i] == -1{
			totalNumNegVotes++
		}
	}

	if totalNumPosVotes > (totalServers/2) {
		//elected as leader
		sm.state = 3
		sm.leaderId = sm.id

		for i:=0;i<totalServers-1;i++{
			prevLogIndex := sm.nextIndex[sm.peers[i]]-1
			prevLogTerm := sm.log[prevLogIndex].term
			entries := temp
			actions = append(actions,Send{peerId:sm.peers[i],event:AppendEntriesReqEv{term : sm.term, senderId: sm.id, prevLogIndex: prevLogIndex, 
			prevLogTerm: prevLogTerm, entries: entries,senderCommitIndex: sm.commitIndex}})		
		}
	}else if totalNumNegVotes > (totalServers/2){
		//step down and become follower
		sm.state = 1	

	} 

	return actions
}


func handleLeaderVoteResp(sm *StateMachine,cmd *VoteRespEv) []interface{}{
	initialiseActions()
	if sm.term < cmd.senderTerm {
		sm.term = cmd.senderTerm
		sm.votedFor = -1
		sm.votedAs = make([]int,len(sm.peers))
		actions = append(actions,StateStore{currTerm:sm.term,votedFor:sm.votedFor})
		sm.state = 1
	}

	return actions
}