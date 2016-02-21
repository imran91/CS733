package main

var actions []interface{}

type StateMachine struct {
	id int // server id
	peers []int // other server ids
	log []Log
	term int 
	state int//1:follower 2:candidate 3:leader
	votedFor int
	commitIndex int
	lastLogIndex int
	lastLogTerm int
	leaderId int
	nextIndex []int
	matchIndex []int
	votedAs map[int] int 
	timer int
}

func (sm *StateMachine) ProcessEvent (ev interface{}) []interface{}{
	var act []interface{}
	switch ev.(type) {
	case Append:
		cmd := ev.(Append)
		switch sm.state {
		case 1://Follower
			act = handleFollowerAppend(sm,&cmd)
			return act
		case 2://Candidate
			act = handleCandidateAppend(sm,&cmd)
			return act
		case 3://Leader
			act = handleLeaderAppend(sm,&cmd)
			return act
		}
		break
	
	case Timeout:
		cmd := ev.(Timeout)
		switch sm.state {
		case 1://Follower
			act = handleFollowerTimeout(sm,&cmd)
			return act
		case 2://Candidate
			act = handleCandidateTimeout(sm,&cmd)
			return act
		case 3://Leader
			act = handleLeaderTimeout(sm,&cmd)
			return act
		}
		break
		
	case AppendEntriesReqEv:
		cmd := ev.(AppendEntriesReqEv)
		switch sm.state {
		case 1://Follower
			act = handleFollowerAppendEntryReq(sm,&cmd)
			return act
		case 2://Candidate
			act = handleCandidateAppendEntryReq(sm,&cmd)
			return act
		case 3://Leader
			act = handleLeaderAppendEntryReq(sm,&cmd)
			return act
		}
		break
	
	case AppendEntriesRespEv:
		cmd := ev.(AppendEntriesRespEv)
		switch sm.state {
		case 1://Follower
			act = handleFollowerAppendEntryResp(sm,&cmd)
			return act
		case 2://Candidate
			act = handleCandidateAppendEntryResp(sm,&cmd)
			return act
		case 3://Leader
			act = handleLeaderAppendEntryResp(sm,&cmd)
			return act
		}
		break
	
	case VoteReqEv:
		cmd := ev.(VoteReqEv)
		switch sm.state {
		case 1://Follower
			act = handleFollowerVoteReq(sm,&cmd)
			return act
		case 2://Candidate
			act = handleCandidateVoteReq(sm,&cmd)
			return act
		case 3://Leader
			act = handleLeaderVoteReq(sm,&cmd)
			return act
		}
		break

	case VoteRespEv:
		cmd := ev.(VoteRespEv)
		switch sm.state {
		case 1://Follower
			act = handleFollowerVoteResp(sm,&cmd)
			return act
		case 2://Candidate
			act = handleCandidateVoteResp(sm,&cmd)
			return act
		case 3://Leader
			act = handleLeaderVoteResp(sm,&cmd)
			return act
		}
		break

	default: 
		println ("Unrecognized")
		return act
		break
	}
return act
}