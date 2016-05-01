package main

import (
	"github.com/cs733-iitb/cluster"
	//"github.com/cs733-iitb/cluster/mock"	
	"github.com/cs733-iitb/log"
	"sync"
	"encoding/gob"
//	"math/rand"

	//"errors"
	"fmt"
	"time"
)
const (
	ChBufSize = 10000
)

type CommitInfo struct{
	Data []byte
	Index int64 
	Err error 
}

type RaftNode struct{
	MutLock sync.Mutex
	isOn bool
	sm StateMachine
	appendCh chan Append//Event channel
	CommitCh chan Commit//Action channel
	fileDir string
	logFile *log.Log
	stateFile *log.Log
	serv cluster.Server
	timer *time.Timer
}

type NetConfig struct{
Id int
Host string
Port int
}

type Config struct{
cluster []NetConfig // Information about all servers, including this.
Id int // this node's id. One of the cluster's entries should match.
LogDir string // Log file directory for this node
StateFileDir     string
ElectionTimeout int64
HeartbeatTimeout int64
}

func RegisterStructs() {
	gob.Register(Log{})
	gob.Register(Send{})
	gob.Register(StateStore{})
	gob.Register(AppendEntriesReqEv{})
	gob.Register(AppendEntriesRespEv{})
	gob.Register(VoteReqEv{})
	gob.Register(VoteRespEv{})
}

func (rn *RaftNode) actionHandler(actions []interface{}){
	rn.MutLock.Lock()
	defer rn.MutLock.Unlock()
	if !rn.isOn {
		return
	}
	for _,action := range actions {
		switch action.(type) {
			case Send:
				send := action.(Send)
				//fmt.Println("sendaction",send)
				rn.serv.Outbox() <- &cluster.Envelope{Pid: send.PeerId, MsgId: -1, Msg: send.Event }
			case Commit:
				commit := action.(Commit)
				//fmt.Println("commitaction",commit)
				rn.CommitCh <- commit
			case Alarm:
				//fmt.Println("alarmaction")
				alarm := action.(Alarm)
				status := rn.timer.Reset(time.Duration(alarm.T) * time.Millisecond)
				if !status {
					rn.timer = time.NewTimer(time.Duration(alarm.T) * time.Millisecond)
				}
			case LogStore:
				//fmt.Println("logstoreaction")
				logStore := action.(LogStore)
				err := rn.logFile.TruncateToEnd(int64(logStore.Index))
				if err!=nil{
					//fmt.Println("logstore action error")	
				}	
				//fmt.Println("Log is=",logStore.LogEntry)	
				/*err = rn.logFile.Append(logStore.LogEntry)
				if err!=nil{
					fmt.Println("logstore action error")	
				}*/	
				
			case StateStore:
				//stateStore := action.(StateStore)
				//rn.stateFile.TruncateToEnd(0)
				//rn.stateFile.Append(StateStore{CurrTerm: stateStore.CurrTerm, VotedFor: stateStore.VotedFor, LastMatchIndex: stateStore.LastMatchIndex})
				//fmt.Println("statestoreaction",stateStore)
				//err := rn.stateFile.writeState(stateStore)
				
		}
	}
}

func (rn *RaftNode) start(){
	var res bool
	var	event Event
	var	envelope *cluster.Envelope

	for {
		rn.MutLock.Lock()
		if !rn.isOn {
			rn.MutLock.Unlock()
			fmt.Println("rn is not on")
			break
		}
		rn.MutLock.Unlock()
	select {	
		case event, res = <- rn.appendCh:
			if res == false {
				break
			}
			actions := rn.sm.ProcessEvent(event)
			rn.actionHandler(actions)
			//fmt.Println("inside start appendchannel",actions)
		case <- rn.timer.C:
			event = Timeout{}
			actions := rn.sm.ProcessEvent(event)
			rn.actionHandler(actions)
		case envelope, _ = <- rn.serv.Inbox():
			event = envelope.Msg
			actions := rn.sm.ProcessEvent(event)
			rn.actionHandler(actions)		
		default:
		}
	}
}

func New(config Config) Node {
	var rn RaftNode
	var configCluster cluster.Config
	var peerIds []int
	var nextIndex map[int]int
	var matchIndex map[int]int
	var votedAs map[int]int
	var LogEntries []Log
	var err error
	rn.isOn = true


	rn.fileDir = config.LogDir
	rn.appendCh = make(chan Append, ChBufSize)
	rn.CommitCh = make(chan Commit, ChBufSize)
	RegisterStructs()
	configCluster = retConfigFile(config)
	peerIds = retPeerIds(config)
	
	nextIndex = initNextIndex(peerIds)
	matchIndex = initMatchIndex(peerIds)
	votedAs = initVotedAs(peerIds)
	
	rn.serv, err = cluster.New(config.Id, configCluster)
	if err!=nil{
		fmt.Println("cluster config failure")
	}
	rn.logFile, err = log.Open(rn.fileDir + "/log")
	if err!=nil{
		fmt.Println("Logfile failure")
	}

	rn.stateFile, err = log.Open(config.StateFileDir)
	if err!=nil{
		//fmt.Println("Statefile failure")
	}
	
	rn.sm.log = append(rn.sm.log,Log{Term:0,Command:nil})

	LogEntries,err = rn.retLog()
	if err!=nil{
		//fmt.Println("Log retrival failure")
	}
	LogEntries = append(LogEntries,Log{Term:0,Command:nil})
	
//	rn.stateFile.TruncateToEnd(0)
//	rn.stateFile.Append(StateStore{CurrTerm: 0, VotedFor: -1, LastMatchIndex: -1})


	rn.sm.initStateMachine( /* currTerm */ 0, 
						 	/* votedFor */ -1, 
						 	/* Log */ LogEntries,
						 	/* selfId */ config.Id, 
						 	/* peerIds */ peerIds, 
						 	/* electionAlarm */ config.ElectionTimeout, 
						 	/* heartbeatAlarm */ config.HeartbeatTimeout, 
						 	/* lastMatchIndex */ -1, 
						 	/* currState --Follower*/ 1, 
						 	/* commitIndex */ -1, 
						 	/* leaderId */ -1, 
						 	/* lastLogIndex */ -1, 
						 	/* lastLogTerm */ 0, 
						 	/* votedAs */ votedAs, 
						 	/* nextIndex */ nextIndex, 
							/* matchIndex */ matchIndex)
	
	/*if config.Id == 101 {
		rn.timer = time.NewTimer(time.Duration(100) * time.Millisecond)
	} else {
		rn.timer = time.NewTimer(time.Duration(randomNoInRange(rn.sm.electionAlarm, 2 * rn.sm.electionAlarm)) * time.Millisecond)
	}
*/
	rn.timer = time.NewTimer(time.Duration(randomNoInRange(rn.sm.electionAlarm, 2 * rn.sm.electionAlarm)) * time.Millisecond)
	//fmt.Println(configCluster,peerIds,nextIndex,matchIndex,LogEntries)
	go rn.start()
	
	return & rn
}
