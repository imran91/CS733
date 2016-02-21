package main

type Append struct{
	data []byte
}

type Timeout struct {
}

type AppendEntriesReqEv struct {
	term int //leaders term
	senderId int //leaderId
	prevLogIndex int
	prevLogTerm int
	entries []Log
	senderCommitIndex int
}

type AppendEntriesRespEv struct {
	senderId int
	senderTerm int
	response bool
}

type VoteReqEv struct {
	senderId int
	term int
	lastLogIndex int
	lastLogTerm int
}

type VoteRespEv struct{
	senderTerm int
	response bool

}

type Log struct{
	term int
	command []byte
}