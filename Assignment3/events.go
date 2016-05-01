package main

type Append struct {
	Data []byte
}

type Timeout struct {
}

type AppendEntriesReqEv struct {
	Term              int //leaders term
	SenderId          int //leaderId
	PrevLogIndex      int64
	PrevLogTerm       int64
	Entries           []Log
	SenderCommitIndex int64
}

type AppendEntriesRespEv struct {
	SenderId       int
	SenderTerm     int
	Response       bool
	LastMatchIndex int64
}

type VoteReqEv struct {
	SenderId     int
	Term         int
	LastLogIndex int
	LastLogTerm  int
}

type VoteRespEv struct {
	SenderId   int
	SenderTerm int
	Response   bool
}

type Log struct {
	Term    int
	Command []byte
}
