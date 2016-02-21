package main

type Send struct {
	peerId int
	event Event //change
}

type Event interface{}

type Commit struct{
	index int
	leaderId int
	data []byte
	err error
}

type Alarm struct{
	t int
}

type LogStore struct{
	index int
	logEntry []Log
}

type StateStore struct{
	currTerm int
	votedFor int
}
 