package main
import (
	//"fmt"
	"github.com/cs733-iitb/cluster/mock"
	"errors"
	"os"
	//"strconv"
)

type Node interface {
	// Client's message to Raft node
	Append([]byte) (bool)

	// A channel for client to listen on. What goes into Append must come out of here at some point.
	CommitChannel() (<- chan Commit,bool)

	// Last known committed index in the log. This could be -1 until the system stabilizes.
	CommittedIndex() (int,bool)

	ConvertIntoMockServ(*mock.MockServer)//convert server to mock server

	// Returns the data at a log index, or an error.
	Get(index int) (Log,error)

	// Node's id
	Id() (int,bool)

	// Id of leader. -1 if unknown
	LeaderId() (int,bool)

	// Signal to shut down all goroutines, stop sockets, flush log and close it, cancel timers.
	Shutdown()
}

func (rn *RaftNode) Append(data []byte) bool {
	rn.MutLock.Lock()
	defer rn.MutLock.Unlock()
	if rn.isOn {
		rn.appendCh <- Append{Data:data}
		return true
	} else {
		return false
	}
}

func (rn *RaftNode) CommitChannel() (<- chan Commit,bool) {
	rn.MutLock.Lock()
	defer rn.MutLock.Unlock()
	if rn.isOn {
		return rn.CommitCh,true
	} else {
		return make(chan Commit,ChBufSize),false
	}	
}

func (rn *RaftNode) CommittedIndex() (int, bool) {
	rn.MutLock.Lock()
	defer rn.MutLock.Unlock()
	if rn.isOn {
		return int(rn.sm.commitIndex), true
	} else {
		return -1, false
	}	
}


func (rn *RaftNode) ConvertIntoMockServ(mockServer *mock.MockServer){
	rn.MutLock.Lock()
	defer rn.MutLock.Unlock()
	if rn.isOn{
		rn.serv.Close()
		rn.serv = mockServer
	}

}

func (rn *RaftNode) Get(index int) (Log,error) {
	var logentry Log
	var temp interface{}
	var err error
	rn.MutLock.Lock()
	defer rn.MutLock.Unlock()
	
	if rn.isOn {
		temp, err = rn.logFile.Get(int64(index))		//verify
		logentry = temp.(Log)
		return logentry, err		
	} else {
		return Log{}, errors.New("ERR_SM_OFF")
	}
}

func (rn *RaftNode) Id() (int,bool) {
	rn.MutLock.Lock()
	defer rn.MutLock.Unlock()
	if rn.isOn {
		return rn.sm.id,true
	} 
	return -1,false
}

func (rn *RaftNode) LeaderId() (int, bool) {
	rn.MutLock.Lock()
	defer rn.MutLock.Unlock()
	if rn.isOn {
		return rn.sm.leaderId, true
	}
	return -1, false
}

func (rn *RaftNode) Shutdown() {
	rn.MutLock.Lock()
	defer rn.MutLock.Unlock()
	if rn.isOn {
		rn.isOn = false
		rn.timer.Stop()
		rn.serv.Close()
		close(rn.appendCh)
		close(rn.CommitCh)
		rn.logFile.Close()
		//rn.stateFile.Close()
		os.RemoveAll(rn.fileDir)
	}
}