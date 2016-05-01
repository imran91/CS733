package main

import(
	"github.com/cs733-iitb/cluster"
//	"github.com/cs733-iitb/log"
	"math/rand"
	"strconv"
	//"errors"
//	"fmt"
	"time"
)

func retConfigFile(config Config)(cluster.Config){
	var configFile cluster.Config
	var peers []cluster.PeerConfig
	var inboxSize int
	var outboxSize int

	for _,temp:= range(config.cluster) {
		id := temp.Id
		address := temp.Host + ":" + strconv.Itoa(temp.Port)
		peers = append(peers, cluster.PeerConfig{id, address})
	}
	inboxSize = 10000
	outboxSize = 10000
	configFile = cluster.Config{peers,inboxSize,outboxSize}
	return configFile
}

func retPeerIds(config Config)([]int){
	var peerIds []int
	for _,temp := range config.cluster {
		if temp.Id != config.Id {
			peerIds = append(peerIds, temp.Id)
		}
	}
	return peerIds
}

func (rn *RaftNode) retLog()([]Log, error){
	var temp interface{}
    var err error
    var LogEntries []Log
    var entry Log
    var i int64

	lastIndex := rn.logFile.GetLastIndex()	
	for i = 0; i < lastIndex; i++ {
		temp, err = rn.logFile.Get(i)
		if err != nil {
			return []Log{}, err
		}
		entry = temp.(Log)
		LogEntries = append(LogEntries, entry)
	}
	return LogEntries, nil
}



func initNextIndex(peerIds []int)(map[int]int){
	var nextIndex map[int]int
	nextIndex = make(map[int]int,len(peerIds))
	for i := 0; i < len(peerIds); i++ {
		nextIndex[peerIds[i]] = -1
	}
	return nextIndex
}

func initMatchIndex(peerIds []int)(map[int]int){
	var matchIndex map[int]int
	matchIndex = make(map[int]int,len(peerIds))
	for i := 0; i < len(peerIds); i++ {
		matchIndex[peerIds[i]] = -1
	}
	return matchIndex
}

func initVotedAs(peerIds []int)(map[int]int){
	var votedAs map[int]int
	votedAs = make(map[int]int,len(peerIds))
	for i := 0; i < len(peerIds); i++ {
		votedAs[peerIds[i]] = 0
	}
	return votedAs
}


func (sm *StateMachine) initStateMachine(currTerm int, votedFor int, log []Log, myId int, peerIds []int, electionTimeout int64, 
						heartbeatTimeout int64, lastMatchIndex int64, currState int, commitIndex int64, leaderId int,
							lastLogIndex int, lastLogTerm int, votedAs map[int]int, nextIndex map[int]int, matchIndex map[int]int){
	sm.state = currState
	sm.term = currTerm
	sm.votedFor = votedFor
	sm.log = log
	sm.id = myId
	sm.peers = peerIds	
	sm.electionAlarm = electionTimeout
	sm.heartbeatAlarm = heartbeatTimeout
	sm.lastMatchIndex = lastMatchIndex
	
	sm.commitIndex = commitIndex
	sm.leaderId = leaderId
	sm.lastLogIndex = int64(lastLogIndex)
	sm.lastLogTerm = lastLogTerm
	sm.clusterSize = len(peerIds) + 1
	sm.votedAs = votedAs
	sm.nextIndex = nextIndex
	sm.matchIndex = matchIndex

	// For Generating Random Time between timer and 2 * timer
	rand.Seed(time.Now().Unix())
}