package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import (
	"fmt"
	"sync"
	"labrpc"
	"time"
	"math/rand"
)

// import "bytes"
// import "labgob"



//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make(). set
// CommandValid to true to indicate that the ApplyMsg contains a newly
// committed log entry.
//
// in Lab 3 you'll want to send other kinds of messages (e.g.,
// snapshots) on the applyCh; at that point you can add fields to
// ApplyMsg, but set CommandValid to false for these other uses.
//
type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

type Log struct {
    Command     interface{}
    Term        int
}

const (
	STATE_LEADER = iota
	STATE_CANDIDATE
	STATE_FOLLOWER
)


//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.
    term      int
    voteFor   int
    curLeader int 

	voteCount int //的票数量
	log []Log //日志条目集；每一个条目包含一个用户状态机执行的指令，和收到时的任期号
    curHaveLog   int //多少日志

// 所有服务器上经常变的
	commitIndex	int //已知的最大的已经被提交的日志条目的索引值
	lastApplied	int //最后被应用到状态机的日志条目索引值（初始化为 0，持续递增）

	// 在领导人里经常改变的 （选举后重新初始化）
	nextIndex   []int
	matchIndex  []int

	state int

	receivedHeartbeat bool

	// 心跳频道
    heartBreathChan chan bool
    requestForVoteChan  chan bool
    broadCast chan bool
}

// return currentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (2A).
	return term, isleader
}


//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	// w := new(bytes.Buffer)
	// e := labgob.NewEncoder(w)
	// e.Encode(rf.xxx)
	// e.Encode(rf.yyy)
	// data := w.Bytes()
	// rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	// Your code here (2C).
	// Example:
	// r := bytes.NewBuffer(data)
	// d := labgob.NewDecoder(r)
	// var xxx
	// var yyy
	// if d.Decode(&xxx) != nil ||
	//    d.Decode(&yyy) != nil {
	//   error...
	// } else {
	//   rf.xxx = xxx
	//   rf.yyy = yyy
	// }
}




//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//

type RequestVoteArgs struct {
	Term		int//候选人的任期号
	CandidateId	int//请求选票的候选人的id
	LastLogIndex int//候选人的最后日志条目的索引值
	LastLogTerm	int//候选人最后日志条目的任期号
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Term		int //当前任期号，以便于候选人去更新自己的任期号
	VoteGranted	bool//候选人赢得了次张选票时为真
}

//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	fmt.Println("11111111111")
	rf.mu.Lock()
	defer rf.mu.Unlock()

	rf.receivedHeartbeat = true
	// 第一种情况
	if(args.Term<rf.term){
		reply.Term = rf.term
		reply.VoteGranted = false
		return
	}
	// 同步
	reply.Term = args.Term
	// 第二种情况
	// 1参数任期等于大于本身人气
	// 2没投票或者已经投票给我

	if(rf.term<args.Term||(rf.term==args.Term&&(rf.voteFor==-1||rf.voteFor==args.CandidateId))){
		if(args.LastLogTerm>rf.log[rf.curHaveLog].Term){
			reply.VoteGranted = true
		}else if(args.LastLogTerm==rf.log[rf.curHaveLog].Term&&args.LastLogIndex>=rf.curHaveLog){
			reply.VoteGranted = true
		}else{
			reply.VoteGranted = false
		}
	}else{      //same term and not vote for
		reply.VoteGranted = false
	}
	if(rf.term<args.Term){
		if(reply.VoteGranted==true){
			rf.term = args.Term
			rf.voteFor = args.CandidateId
		}else{
			rf.term = args.Term
			rf.voteFor = -1
		}
		rf.state = 2
		// 这里有bug
		// if(rf.state==4){
		// 	// 为什么关闭心跳
		// 	close(rf.heartBreathChan)
		// 	rf.requestForVoteChan = make(chan int)
		// 	go rf.startWaitingHeartbeat()
		// }
	}else{      //equal term
		if(rf.voteFor==-1){
			rf.voteFor = args.CandidateId
		}
	}      
	// Your code here (2A, 2B).
}

//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	rf.mu.Lock()
	defer rf.mu.Unlock()
	if ok{
		if(rf.term == reply.Term&&rf.state == STATE_CANDIDATE){
			if(reply.VoteGranted==true){
				rf.voteCount++
				if(rf.voteCount > len(rf.peers)/2){
					rf.curLeader = rf.me
					close(rf.requestForVoteChan)
					rf.state = STATE_LEADER
					// 这里要发送心跳
					for i,_ := range rf.peers{
						rf.matchIndex[i] = 0
						rf.nextIndex[i]=rf.curHaveLog+1
					}

					rf.heartBreathChan = make(chan bool)
					go rf.startBreath()
					rf.broadCast <- true
				}
			}
			
		}

	}
	return ok
}
// 生命周期
func (rf *Raft) startBreath(){
    for {
        select{
        case <- rf.heartBreathChan :
            return
        case <-rf.broadCast:
            go rf.HeartBreath()
        case <-time.After(time.Duration(100)*time.Millisecond):
            go rf.HeartBreath()
        }
    }
}

func (rf *Raft) HeartBreath(){
    for i,_ := range rf.peers{
        if(i!=rf.me){
            go rf.SendEntriesTo(i)
        }
    }
}

func (rf *Raft) SendEntriesTo(i int){
    
//     rf.mu.Lock()
//     curterm := rf.term
//     if(rf.state!=4){
//         rf.mu.Unlock()
//         return
//     }
//         req := RequestArg{Term:rf.term,LeaderID:rf.me,PrevLogIndex:rf.nextIndex[i]-1,PrevLogTerm:rf.log[rf.nextIndex[i]-1].Term,LeaderCommit:rf.commitIndex}
//         if(rf.nextIndex[i]<=rf.curHaveLog){
//             req.Log = make([]Log,rf.curHaveLog-rf.nextIndex[i]+1)
//             copy(req.Log,rf.log[rf.nextIndex[i]:rf.curHaveLog+1])
//         }else{
//             req.Log = make([]Log,0)
//         }
    
//     rf.mu.Unlock()
    
//     res := ResponseArg{}
//     ok := rf.peers[i].Call("Raft.AppendEntries",&req,&res)
//     for{
//         for !ok {
//             if(rf.term!=curterm){
//                 return
//             }
//             ok = rf.peers[i].Call("Raft.AppendEntries",&req,&res)
//         }
    
//         if(res.Success){    //test commit
//             rf.mu.Lock()
//             if(rf.term!=curterm){
//                 rf.mu.Unlock()
//                 return
//             }
//             if(req.LeaderCommit>len(req.Log)+req.PrevLogIndex){
//                 rf.nextIndex[i] = req.LeaderCommit+1
//             }else{
//                 rf.nextIndex[i] = len(req.Log)+req.PrevLogIndex+1
//             }
//             rf.matchIndex[i] = len(req.Log)+req.PrevLogIndex
            
//             if(len(req.Log)!=0){
//                 rf.testCommit()
//             }
//             rf.mu.Unlock()
//             break
//         }else{
//             rf.mu.Lock()
//             if(curterm<res.Term||rf.term!=curterm){   //older term, turn to follower
//                 rf.mu.Unlock()
//                 return
//             }else{
//                 req.PrevLogIndex--
//                 req.PrevLogTerm = rf.log[req.PrevLogIndex].Term
//                 len := len(req.Log)
//                 req.Log = make([]Log,len+1)
//                 copy(req.Log,rf.log[req.PrevLogIndex+1:req.PrevLogIndex+1+len+1])
//                 rf.mu.Unlock()
//                 ok = rf.peers[i].Call("Raft.AppendEntries",&req,&res)
//             }
//         }
//     }
}


//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	index := -1
	term := -1
	isLeader := true
	// Your code here (2B).
	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
}

func (rf *Raft) startHeartbeat() {
    for {
        m:=rand.Intn(150)
        m+=150
        
        select {
        case <-rf.requestForVoteChan:
			return
			// 超时行为
		case <-time.After(time.Duration(m)*time.Millisecond):
			// fmt.Println("心跳")
            if(rf.receivedHeartbeat){
                rf.receivedHeartbeat = false
            }else{
				// fmt.Println("心跳")
                rf.startElection()
            }
        }
    }
}

func (rf *Raft) startElection(){
    rf.mu.Lock()
        rf.state = STATE_CANDIDATE
        rf.term++
		rf.voteFor = rf.me
		// 的票数量
        rf.voteCount = 0
    rf.mu.Unlock()
    
    for i, _ := range rf.peers{
        if(i!=rf.me){
			// 这里要求投票
			// fmt.Println(i)
            go rf.broadcastVoteReq(i)
        }
    }
}

func (rf *Raft) broadcastVoteReq(i int){
	rf.mu.Lock()
    // curterm := rf.term
    if(rf.state!=1){
        rf.mu.Unlock()
        return
	}
	// 构建请求
	req := RequestVoteArgs{
		Term:rf.term,
		CandidateId:rf.me,
		LastLogIndex:rf.curHaveLog,
		LastLogTerm:rf.log[rf.curHaveLog].Term,
	}
	
	res := RequestVoteReply{}
	rf.sendRequestVote(i, &req, &res)
	
	// ok := rf.peers[i].Call("Raft.RequestVote", &req, &res)
	// if ok{


	// }

}

//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int, persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me

	// Your initialization code here (2A, 2B, 2C).

	// 初始化第一个日志
	rf.log = make([]Log, 10000)
	rf.curHaveLog = 0
	rf.log[0].Term = -1
	rf.term = -1

	rf.state = STATE_FOLLOWER

	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())
	// 广播频道
	rf.broadCast = make(chan bool)

	rf.requestForVoteChan = make(chan bool)
	// 没有心跳接收
	rf.receivedHeartbeat = false

	go rf.startHeartbeat()

	return rf
}