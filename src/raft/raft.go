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

const (
	STATE_LEADER = iota
	STATE_CANDIDATE
	STATE_FOLLOWER
	// HEARTBEAT = 50 * time.Millisecond
)

type LogEntry struct {
	Command interface{}
	Term    int
}

//
// A Go object implementing a single Raft peer.
//
// 首先是根据论文补充raft结构
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]

	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

	// Persist state on all servers
	// 所有服务器上持久存在的
	currentTerm 	int //服务器最后一次知道的任期号（初始化为 0，持续递增）
	votedFor 	int //在当前获得选票的候选人的 Id
	log []LogEntry //日志条目集；每一个条目包含一个用户状态机执行的指令，和收到时的任期号

	// 所有服务器上经常变的
	commitIndex	int //已知的最大的已经被提交的日志条目的索引值
	lastApplied	int //最后被应用到状态机的日志条目索引值（初始化为 0，持续递增）

	// 在领导人里经常改变的 （选举后重新初始化）
	nextIndex   []int
	matchIndex  []int

	grantedVotesCount int
	// 对于状态是三个1:follower,2:candidate,3:leader
	state int//当前节点状态 有几个状态
	applyCh           chan ApplyMsg

	// 选举是通过定时器实现的
	// 这里定时器要求是随机时间
	// time *time.Timer

}

// return currentTerm and whether this server
// believes it is the leader.
// 在初始化时候用的了
func (rf *Raft) GetState() (int, bool) {

	var term int
	var isleader bool
	// Your code here (2A).
	// 对着一点的
	rf.mu.Lock()
	defer rf.mu.Unlock()
	term = rf.currentTerm
	isleader = (rf.state == STATE_LEADER)
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
// |term| 领导人的任期号|
// |leaderId| 领导人的 Id，以便于跟随者重定向请求|
// |prevLogIndex|新的日志条目紧随之前的索引值|
// |prevLogTerm|prevLogIndex 条目的任期号|
// |entries[]|准备存储的日志条目（表示心跳时为空；一次性发送多个是为了提高效率）
// |leaderCommit|领导人已经提交的日志的索引值|

// | 返回值| 解释|
// |---|---|
// |term|当前的任期号，用于领导人去更新自己|
// |success|跟随者包含了匹配上 prevLogIndex 和 prevLogTerm 的日志时为真|
// 这里怎么和论文不太一样
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	term		int//候选人的任期号
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
	// Your code here (2A, 2B).
	// ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	// return ok
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
	return ok
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
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {

	
	// fmt.Println(me)
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	// 每个节点的编号
	rf.me = me

	// 初始状态是跟随，所有都是跟随
	rf.state = STATE_FOLLOWER
	// 最后任期的序号
	rf.currentTerm = 0
	rf.votedFor = -1// 投票leader序号
	rf.log = make([]LogEntry, 0)// 初始化日志数组，关于日志的信息，要研究下

	//这里用的是0到时候要改变整个logcommit部分
	rf.commitIndex = 0
	rf.lastApplied = 0
	rf.nextIndex = make([]int, len(peers))
	rf.matchIndex = make([]int, len(peers))

	rf.grantedVotesCount = 0
	rf.state = STATE_FOLLOWER
	rf.applyCh = applyCh
	// rf.time = time.NewTimer(1 * time.Millisecond)


	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())

	// 开始准备选举
	rf.persist()
	// rf.resetTimer()


	// raft选举机制是透过定时器推动的，所以首先实现定时器相关。
	go rf.startHeartbeat()

	return rf
}


// 定时器推动选举
func (rf *Raft) startHeartbeat() {
	fmt.Println("选举开始")


}