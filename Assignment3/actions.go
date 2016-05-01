package main

type Send struct {
	PeerId int
	Event  Event 
}

type Event interface{}

type Commit struct {
	Index    int
	LeaderId int
	Data     []byte
	Err      error
}

type Alarm struct {
	T int64
}

type LogStore struct {
	Index    int64
	LogEntry Log
}

type StateStore struct {
	CurrTerm int
	VotedFor int
	LastMatchIndex int64
}
