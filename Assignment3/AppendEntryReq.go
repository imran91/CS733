package main

import (
//	"fmt"
)

//var actions []interface{}
func handleFollowerAppendEntryReq(sm *StateMachine, cmd *AppendEntriesReqEv) []interface{} {
	initialiseActions()
	if sm.term > cmd.Term {
		actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: false, LastMatchIndex: sm.lastMatchIndex}})
		return actions
	} else{

		actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)})
		sm.state = 1//change state to follower	
		if sm.term < cmd.Term{
				sm.term = cmd.Term
				sm.votedFor = -1
    			sm.lastMatchIndex = -1
				actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor,LastMatchIndex: sm.lastMatchIndex})

		}
		sm.leaderId = cmd.SenderId

		if cmd.PrevLogIndex > sm.lastLogIndex || (cmd.PrevLogIndex > -1 && int64(sm.log[cmd.PrevLogIndex].Term) != cmd.PrevLogTerm){
			actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: false, LastMatchIndex: sm.commitIndex}})	
		} else if int64(sm.lastMatchIndex) >= (cmd.PrevLogIndex+int64(len(cmd.Entries))){
				if sm.commitIndex < sm.lastLogIndex && sm.commitIndex < cmd.SenderCommitIndex{
					sm.commitIndex = min(cmd.SenderCommitIndex,sm.lastLogIndex)
					actions = append(actions, Commit{Index: int(sm.commitIndex), LeaderId: sm.leaderId, Data: sm.log[sm.commitIndex].Command, Err: nil})
				}
				actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: true, LastMatchIndex: sm.commitIndex}})
		}else{
				if len(cmd.Entries)>0{
					sm.lastLogIndex = int64(cmd.PrevLogIndex)
					sm.log = sm.log[:(sm.lastLogIndex+1)]
				}
				sm.lastMatchIndex = cmd.PrevLogIndex
				//append entries to the log
				for i := 0; i < len(cmd.Entries); i++{
					sm.lastLogIndex = sm.lastLogIndex + 1
					sm.log = append(sm.log, cmd.Entries[i])
					actions = append(actions, LogStore{Index: cmd.PrevLogIndex + 1, LogEntry: cmd.Entries[i]})
					sm.lastLogTerm = sm.log[sm.lastLogIndex].Term
					sm.lastMatchIndex = sm.lastLogIndex
				}
				actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor, LastMatchIndex:sm.lastMatchIndex})
				if sm.commitIndex < sm.lastLogIndex && sm.commitIndex <cmd.SenderCommitIndex {
					sm.commitIndex = min(cmd.SenderCommitIndex,sm.lastLogIndex)
					actions = append(actions, Commit{Index: int(sm.commitIndex), LeaderId: sm.leaderId, Data: sm.log[sm.commitIndex].Command, Err: nil})
				}

				actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: true, LastMatchIndex: sm.lastMatchIndex}})
		}
	}

//	fmt.Println("Inside follower AppendEntryReq")	
	return actions
}

func handleCandidateAppendEntryReq(sm *StateMachine, cmd *AppendEntriesReqEv) []interface{} {
	initialiseActions()
	if sm.term > cmd.Term {
		actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: false, LastMatchIndex: sm.lastMatchIndex}})
		return actions
	} else{

		actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)})
		sm.state = 1//change state to follower	
		if sm.term < cmd.Term{
				sm.term = cmd.Term
				sm.votedFor = -1
    			sm.lastMatchIndex = -1
				actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor,LastMatchIndex: sm.lastMatchIndex})

		}
		sm.leaderId = cmd.SenderId

		if cmd.PrevLogIndex > sm.lastLogIndex || (cmd.PrevLogIndex > -1 && int64(sm.log[cmd.PrevLogIndex].Term) != cmd.PrevLogTerm){
			actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: false, LastMatchIndex: sm.commitIndex}})	
		} else if int64(sm.lastMatchIndex) >= (cmd.PrevLogIndex+int64(len(cmd.Entries))){
				if sm.commitIndex < sm.lastLogIndex && sm.commitIndex < cmd.SenderCommitIndex{
					sm.commitIndex = min(cmd.SenderCommitIndex,sm.lastLogIndex)
					actions = append(actions, Commit{Index: int(sm.commitIndex), LeaderId: sm.leaderId, Data: sm.log[sm.commitIndex].Command, Err: nil})
				}
				actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: true, LastMatchIndex: sm.commitIndex}})
		}else{
				if len(cmd.Entries)>0{
					sm.lastLogIndex = cmd.PrevLogIndex
					sm.log = sm.log[:(sm.lastLogIndex+1)]
				}
				sm.lastMatchIndex = cmd.PrevLogIndex
				//append entries to the log
				for i := 0; i < len(cmd.Entries); i++{
					sm.lastLogIndex = sm.lastLogIndex + 1
					sm.log = append(sm.log, cmd.Entries[i])
					actions = append(actions, LogStore{Index: cmd.PrevLogIndex + 1, LogEntry: cmd.Entries[i]})
					sm.lastLogTerm = sm.log[sm.lastLogIndex].Term
					sm.lastMatchIndex = sm.lastLogIndex
				}
				actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor, LastMatchIndex:sm.lastMatchIndex})
				if sm.commitIndex < sm.lastLogIndex && sm.commitIndex <cmd.SenderCommitIndex {
					sm.commitIndex = min(cmd.SenderCommitIndex,sm.lastLogIndex)
					actions = append(actions, Commit{Index: int(sm.commitIndex), LeaderId: sm.leaderId, Data: sm.log[sm.commitIndex].Command, Err: nil})
				}

				actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: true, LastMatchIndex: sm.lastMatchIndex}})
		}
	}

//	fmt.Println("Inside follower AppendEntryReq")	
	return actions
}

func handleLeaderAppendEntryReq(sm *StateMachine, cmd *AppendEntriesReqEv) []interface{} {
	initialiseActions()
	if sm.term > cmd.Term {
		actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: false, LastMatchIndex: sm.lastMatchIndex}})
		return actions
	} else{

		actions = append(actions, Alarm{T:randomNoInRange(2 * sm.electionAlarm, 3 * sm.electionAlarm)})
		sm.state = 1//change state to follower	
		if sm.term < cmd.Term{
				sm.term = cmd.Term
				sm.votedFor = -1
    			sm.lastMatchIndex = -1
				actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor,LastMatchIndex: sm.lastMatchIndex})

		}
		sm.leaderId = cmd.SenderId

		if cmd.PrevLogIndex > sm.lastLogIndex || (cmd.PrevLogIndex > -1 && int64(sm.log[cmd.PrevLogIndex].Term) != cmd.PrevLogTerm){
			actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: false, LastMatchIndex: sm.commitIndex}})	
		} else if int64(sm.lastMatchIndex) >= (cmd.PrevLogIndex+int64(len(cmd.Entries))){
				if sm.commitIndex < sm.lastLogIndex && sm.commitIndex < cmd.SenderCommitIndex{
					sm.commitIndex = min(cmd.SenderCommitIndex,sm.lastLogIndex)
					actions = append(actions, Commit{Index: int(sm.commitIndex), LeaderId: sm.leaderId, Data: sm.log[sm.commitIndex].Command, Err: nil})
				}
				actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: true, LastMatchIndex: sm.commitIndex}})
		}else{
				if len(cmd.Entries)>0{
					sm.lastLogIndex = cmd.PrevLogIndex
					sm.log = sm.log[:(sm.lastLogIndex+1)]
				}
				sm.lastMatchIndex = cmd.PrevLogIndex
				//append entries to the log
				for i := 0; i < len(cmd.Entries); i++{
					sm.lastLogIndex = sm.lastLogIndex + 1
					sm.log = append(sm.log, cmd.Entries[i])
					actions = append(actions, LogStore{Index: cmd.PrevLogIndex + 1, LogEntry: cmd.Entries[i]})
					sm.lastLogTerm = sm.log[sm.lastLogIndex].Term
					sm.lastMatchIndex = sm.lastLogIndex
				}
				actions = append(actions, StateStore{CurrTerm: sm.term, VotedFor: sm.votedFor, LastMatchIndex:sm.lastMatchIndex})
				if sm.commitIndex < sm.lastLogIndex && sm.commitIndex <cmd.SenderCommitIndex {
					sm.commitIndex = min(cmd.SenderCommitIndex,sm.lastLogIndex)
					actions = append(actions, Commit{Index: int(sm.commitIndex), LeaderId: sm.leaderId, Data: sm.log[sm.commitIndex].Command, Err: nil})
				}

				actions = append(actions, Send{PeerId: cmd.SenderId, Event: AppendEntriesRespEv{SenderId: sm.id, SenderTerm: sm.term, Response: true, LastMatchIndex: sm.lastMatchIndex}})
		}
	}

//	fmt.Println("Inside follower AppendEntryReq")	
	return actions
}
