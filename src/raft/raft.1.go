package raft

// //
// // this is an outline of the API that raft must expose to
// // the service (or tester). see comments below for
// // each of these functions for more details.
// //
// // rf = Make(...)
// //   create a new Raft server.
// // rf.Start(command interface{}) (index, term, isleader)
// //   start agreement on a new log entry
// // rf.GetState() (term, isLeader)
// //   ask a Raft for its current term, and whether it thinks it is leader
// // ApplyMsg
// //   each time a new entry is committed to the log, each Raft peer
// //   should send an ApplyMsg to the service (or tester)
// //   in the same server.
// //

// import (
// 	"fmt"
// 	"sync"
// 	"labrpc"
// 	"time"
// 	"math/rand"
// )

// // import "bytes"
// // import "labgob"



// //
// // as each Raft peer becomes aware that successive log entries are
// // committed, the peer should send an ApplyMsg to the service (or
// // tester) on the same server, via the applyCh passed to Make(). set
// // CommandValid to true to indicate that the ApplyMsg contains a newly
// // committed log entry.
// //
// // in Lab 3 you'll want to send other kinds of messages (e.g.,
// // snapshots) on the applyCh; at that point you can add fields to
// // ApplyMsg, but set CommandValid to false for these other uses.
// //
// type ApplyMsg struct {
// 	CommandValid bool
// 	Command      interface{}
// 	CommandIndex int
// }

// const (
// 	STATE_LEADER = iota
// 	STATE_CANDIDATE
// 	STATE_FOLLOWER
// 	// HEARTBEAT = 50 * time.Millisecond
// 	// 定时器时间，是random的
// 	HEARTBEAT_INTERVAL    = 50
// 	MIN_ELECTION_INTERVAL = 150
// 	MAX_ELECTION_INTERVAL = 300
// )

// type LogEntry struct {
// 	Command interface{} //命令
// 	Term    int //谁的任期干的
// 	Index   int //命令索引
// }

// //
// // A Go object implementing a single Raft peer.
// //
// // 首先是根据论文补充raft结构
// type Raft struct {
// 	mu        sync.Mutex          // Lock to protect shared access to this peer's state
// 	peers     []*labrpc.ClientEnd // RPC end points of all peers
// 	persister *Persister          // Object to hold this peer's persisted state
// 	me        int                 // this peer's index into peers[]

// 	// Your data here (2A, 2B, 2C).
// 	// Look at the paper's Figure 2 for a description of what
// 	// state a Raft server must maintain.

// 	// Persist state on all servers
// 	// 所有服务器上持久存在的
// 	currentTerm 	int //服务器最后一次知道的任期号（初始化为 0，持续递增）
// 	votedFor 	int //在当前获得选票的候选人的 Id

// 	voteCount int //的票数量
// 	log []LogEntry //日志条目集；每一个条目包含一个用户状态机执行的指令，和收到时的任期号

// 	// 所有服务器上经常变的
// 	commitIndex	int //已知的最大的已经被提交的日志条目的索引值
// 	lastApplied	int //最后被应用到状态机的日志条目索引值（初始化为 0，持续递增）

// 	// 在领导人里经常改变的 （选举后重新初始化）
// 	nextIndex   []int
// 	matchIndex  []int

// 	grantedVotesCount int
// 	// 对于状态是三个1:follower,2:candidate,3:leader
// 	state int//当前节点状态 有几个状态
// 	applyCh           chan ApplyMsg

// 	// 选举是通过定时器实现的
// 	// 这里定时器要求是随机时间
// 	time *time.Timer

// 	// 心跳压制频道
// 	chanHeartbeat chan bool
// 	chanLeader    chan bool

// }

// // return currentTerm and whether this server
// // believes it is the leader.
// // 在初始化时候用的了
// func (rf *Raft) GetState() (int, bool) {

// 	var term int
// 	var isleader bool
// 	// Your code here (2A).
// 	// 对着一点的
// 	rf.mu.Lock()
// 	defer rf.mu.Unlock()
// 	term = rf.currentTerm
// 	isleader = (rf.state == STATE_LEADER)
// 	return term, isleader
// }


// //
// // save Raft's persistent state to stable storage,
// // where it can later be retrieved after a crash and restart.
// // see paper's Figure 2 for a description of what should be persistent.
// //
// func (rf *Raft) persist() {
// 	// Your code here (2C).
// 	// Example:
// 	// w := new(bytes.Buffer)
// 	// e := labgob.NewEncoder(w)
// 	// e.Encode(rf.xxx)
// 	// e.Encode(rf.yyy)
// 	// data := w.Bytes()
// 	// rf.persister.SaveRaftState(data)
// }

// //
// // restore previously persisted state.
// //
// func (rf *Raft) readPersist(data []byte) {
// 	if data == nil || len(data) < 1 { // bootstrap without any state?
// 		return
// 	}
// 	// Your code here (2C).
// 	// Example:
// 	// r := bytes.NewBuffer(data)
// 	// d := labgob.NewDecoder(r)
// 	// var xxx
// 	// var yyy
// 	// if d.Decode(&xxx) != nil ||
// 	//    d.Decode(&yyy) != nil {
// 	//   error...
// 	// } else {
// 	//   rf.xxx = xxx
// 	//   rf.yyy = yyy
// 	// }
// }




// //
// // example RequestVote RPC arguments structure.
// // field names must start with capital letters!
// //
// type RequestVoteArgs struct {
// 	// Your data here (2A, 2B).
// 	Term		int//候选人的任期号
// 	CandidateId	int//请求选票的候选人的id
// 	LastLogIndex int//候选人的最后日志条目的索引值
// 	LastLogTerm	int//候选人最后日志条目的任期号
// }

// //
// // example RequestVote RPC reply structure.
// // field names must start with capital letters!
// //
// type RequestVoteReply struct {
// 	// Your data here (2A).
// 	Term		int //当前任期号，以便于候选人去更新自己的任期号
// 	VoteGranted	bool//候选人赢得了次张选票时为真
// }

// // 基本心跳格式
// type AppendEntriesArgs struct {
// 	// 当前领导任期
// 	Term     int
// 	// 领导者
// 	LeaderID int
// }

// type AppendEntriesReply struct {
// 	Term    int
// 	Success bool
// }
// //
// // example RequestVote RPC handler.
// //
// // rpc远程调用这个投票函数
// // 该函数主要工作是进行投票
// func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
// 	// Your code here (2A, 2B).
// 	rf.mu.Lock()
// 	defer rf.mu.Unlock()
// 	// 在这里进行状态的判断
// 	// 1. 如果 `term < currentTerm` 就返回 false （5.1 节）
// 	// 2. 如果日志在 prevLogIndex 位置处的日志条目的任期号和 prevLogTerm 不匹配，则返回 false （5.3 节）
// 	// 3. 如果已经已经存在的日志条目和新的产生冲突（相同偏移量但是任期号不同），删除这一条和之后所有的 （5.3 节）
// 	// 4. 附加任何在已有的日志中不存在的条目
// 	// 5. 如果 `leaderCommit > commitIndex`，令 commitIndex 等于 leaderCommit 和 新日志条目索引值中较小的一个

// 	// 如果要求请求的帐号直接返回不参与投票
// 	if args.Term < rf.currentTerm {
// 		reply.Term = rf.currentTerm
// 		reply.VoteGranted = false
// 	}
// 	// 第一次进入，会进入这个判断
// 	// 若term比自己大，变跟随者，重启定时器。
// 	if args.Term > rf.currentTerm {
// 		// 这里行为就是先把自己变成跟随者，同时更新term版本
// 		// rf.currentTerm = args.Term
// 		rf.updateStateTo(STATE_FOLLOWER)

// 		index := rf.getLastIndex()
// 		//1 当当前任期大于参数任期
// 		//2 当当前任期等于参数任期，并且，日志编号大于等于我的日志编号
// 		// 是否有复制日志的行为
// 		if args.LastLogTerm > args.Term {
// 			rf.currentTerm = args.Term
// 			rf.votedFor = args.CandidateId
// 			reply.VoteGranted = true
// 		}else if args.LastLogTerm == rf.currentTerm && args.LastLogIndex >= index {
// 			fmt.Println("22222222222222222")
// 			fmt.Println(args.LastLogTerm, rf.currentTerm, args.Term)
// 			fmt.Println("22222222222222222")
// 			rf.currentTerm = args.Term
// 			rf.votedFor = args.CandidateId
// 			reply.VoteGranted = true
// 		}else{
// 			reply.VoteGranted = false
// 		}
// 	}


// 	if args.Term == rf.currentTerm {


// 	}

// 	// 是不是要开始定时任务
// 	// fmt.Println(rf.me)
// 	// fmt.Println(args.term)
// }
// // 心跳压制
// func (rf *Raft) AppendEntries(args AppendEntriesArgs, reply *AppendEntriesReply) {
// 	// 三种情况
// 	rf.mu.Lock()
// 	defer rf.mu.Unlock()
// 	//deal with the msg and reply

// 	// 三种情况
// 	// 主要判断任期问题

// 	if args.Term < rf.currentTerm {
// 		reply.Success = false
// 		reply.Term = rf.currentTerm
// 	} else if args.Term > rf.currentTerm {
// 		rf.currentTerm = args.Term
// 		rf.updateStateTo(STATE_FOLLOWER)
// 		reply.Success = true
// 	} else {
// 		reply.Success = true
// 	}
// 	reply.Term = args.Term
// 	// 这时候要求进入自己的心跳循环
// 	go func() {
// 		rf.chanHeartbeat <- true
// 	}()
// }

// //
// // example code to send a RequestVote RPC to a server.
// // server is the index of the target server in rf.peers[].
// // expects RPC arguments in args.
// // fills in *reply with RPC reply, so caller should
// // pass &reply.
// // the types of the args and reply passed to Call() must be
// // the same as the types of the arguments declared in the
// // handler function (including whether they are pointers).
// //
// // The labrpc package simulates a lossy network, in which servers
// // may be unreachable, and in which requests and replies may be lost.
// // Call() sends a request and waits for a reply. If a reply arrives
// // within a timeout interval, Call() returns true; otherwise
// // Call() returns false. Thus Call() may not return for a while.
// // A false return can be caused by a dead server, a live server that
// // can't be reached, a lost request, or a lost reply.
// //
// // Call() is guaranteed to return (perhaps after a delay) *except* if the
// // handler function on the server side does not return.  Thus there
// // is no need to implement your own timeouts around Call().
// //
// // look at the comments in ../labrpc/labrpc.go for more details.
// //
// // if you're having trouble getting RPC to work, check that you've
// // capitalized all field names in structs passed over RPC, and
// // that the caller passes the address of the reply struct with &, not
// // the struct itself.
// //
// func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
// 	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
// 	// 在这里进行一系列的数据的判断，
// 	// 对当前rf节点进行判断
// 	if ok{
// 		if reply.VoteGranted {
// 			rf.voteCount++
// 			if rf.state == STATE_CANDIDATE && rf.voteCount > len(rf.peers)/2 {
// 				fmt.Println("我成为领导者")
// 				fmt.Println("我是", rf.me)
// 				rf.updateStateTo(STATE_LEADER)
// 			}
// 		}
// 	}

// 	return ok
// }

// func (rf *Raft) sendAppendEntries(server int, args AppendEntriesArgs, reply *AppendEntriesReply) bool {
// 	ok := rf.peers[server].Call("Raft.AppendEntries", args, reply)

// 	return ok
// }


// // the service using Raft (e.g. a k/v server) wants to start
// // agreement on the next command to be appended to Raft's log. if this
// // server isn't the leader, returns false. otherwise start the
// // agreement and return immediately. there is no guarantee that this
// // command will ever be committed to the Raft log, since the leader
// // may fail or lose an election.
// //
// // the first return value is the index that the command will appear at
// // if it's ever committed. the second return value is the current
// // term. the third return value is true if this server believes it is
// // the leader.
// //
// func (rf *Raft) Start(command interface{}) (int, int, bool) {
// 	index := -1
// 	term := -1
// 	isLeader := true

// 	// Your code here (2B).

// 	return index, term, isLeader
// }

// //
// // the tester calls Kill() when a Raft instance won't
// // be needed again. you are not required to do anything
// // in Kill(), but it might be convenient to (for example)
// // turn off debug output from this instance.
// //
// func (rf *Raft) Kill() {
// 	// Your code here, if desired.
// }

// func randElectionDuration() time.Duration {
// 	r := rand.New(rand.NewSource(time.Now().UnixNano()))
// 	return time.Millisecond * time.Duration(r.Int63n(MAX_ELECTION_INTERVAL-MIN_ELECTION_INTERVAL)+MIN_ELECTION_INTERVAL)
// }

// //
// // the service or tester wants to create a Raft server. the ports
// // of all the Raft servers (including this one) are in peers[]. this
// // server's port is peers[me]. all the servers' peers[] arrays
// // have the same order. persister is a place for this server to
// // save its persistent state, and also initially holds the most
// // recent saved state, if any. applyCh is a channel on which the
// // tester or service expects Raft to send ApplyMsg messages.
// // Make() must return quickly, so it should start goroutines
// // for any long-running work.
// //
// func Make(peers []*labrpc.ClientEnd, me int,
// 	persister *Persister, applyCh chan ApplyMsg) *Raft {

	
// 	// fmt.Println(me)
// 	rf := &Raft{}
// 	rf.peers = peers
// 	rf.persister = persister
// 	// 每个节点的编号
// 	rf.me = me

// 	// 初始状态是跟随，所有都是跟随
// 	rf.state = STATE_FOLLOWER
// 	// 这时候需要一些投票的chan


// 	// 最后任期的序号
// 	rf.currentTerm = 0
// 	rf.votedFor = -1// 投票leader序号
// 	rf.log = make([]LogEntry, 0)// 初始化日志数组，关于日志的信息，要研究下

// 	//这里用的是0到时候要改变整个logcommit部分
// 	rf.commitIndex = 0
// 	rf.lastApplied = 0
// 	rf.nextIndex = make([]int, len(peers))
// 	rf.matchIndex = make([]int, len(peers))

// 	rf.grantedVotesCount = 0
// 	rf.state = STATE_FOLLOWER
// 	rf.applyCh = applyCh
// 	// rf.time = time.NewTimer(1 * time.Millisecond)

// 	rf.chanHeartbeat = make(chan bool)

// 	// initialize from state persisted before a crash
// 	rf.readPersist(persister.ReadRaftState())

	
	
// 	rf.persist()
// 	// rf.resetTimer()

// 	// 选举用chan


// 	// raft选举机制是透过定时器推动的，所以首先实现定时器相关。
// 	// 领导者周期性的向所有跟随者发送心跳包（不包含日志项内容的附加日志项 RPCs）来维持自己的权威。
// 	// 这样一开始应该是一个超时状态
// 	go rf.startHeartbeat()

// 	return rf
// }


// // 定时器推动选举
// func (rf *Raft) startHeartbeat() {
// 	// fmt.Println("选举开始")
// 	// 这里还是不是很明白选举人每个节点
// 	// 定时器赋值
// 	rf.time = time.NewTimer(randElectionDuration())
// 	// 
// 	for{
// 		switch rf.state {
// 		case STATE_FOLLOWER:
// 			// 有三个情况
// 			// 1.开始，2超时导致选举转换，3有其他领导人高于我，我从leader转化
// 			select {
// 				// 心跳的的压制
// 				case <-rf.chanHeartbeat:
// 				// case <-rf.voteCh:
// 				// 	rf.electionTimer.Reset(randElectionDuration())
// 				// case <-rf.appendCh:
// 				// 	rf.electionTimer.Reset(randElectionDuration())
// 				// 超时时候将转换成选举者
// 				case <-rf.time.C:
// 					rf.mu.Lock()
// 					rf.updateStateTo(STATE_CANDIDATE)
// 					rf.mu.Unlock()
// 			}
// 			// 成为leader后对其他节点进行压制
// 		case STATE_LEADER:
// 			// fmt.Printf("Leader:%v %v\n",rf.me,"broadcastAppendEnrties	")
// 			// 这里需要心跳压制
// 			rf.broadcastAppendEnrties()
// 			time.Sleep(HEARTBEAT_INTERVAL) 

// 		case STATE_CANDIDATE:
// 			rf.startElection()
// 			// select {
// 			// 	// case <-rf.chanHeartbeat:
// 			// 	rf.startElection()

// 			// }
// 		}
// 	}
// }

// // 心跳广播
// func (rf *Raft) broadcastAppendEnrties() {
// 	fmt.Println("心跳压制")
// 	var args AppendEntriesArgs
// 	// rf.mu.Lock()
// 	for i := range rf.peers {
// 		// 对所有人进行广播
// 		if i != rf.me && rf.state == STATE_LEADER {
// 			go func(id int) {
// 				var reply AppendEntriesReply
// 				rf.sendAppendEntries(id, args, &reply)
// 			}(i)
// 		}
// 	}
// }

// //   Raft 使用一种心跳机制来触发领导人选举。当服务器程序启动时，他们都是跟随者身份。
// // 如果一个跟随者在一段时间里没有接收到任何消息，也就是选举超时，然后他就会认为系统中没有可用的领导者然后开始进行选举以选出新的领导者。
// // 要开始一次选举过程，跟随者先要增加自己的当前任期号并且转换到候选人状态。
// // 然后他会并行的向集群中的其他服务器节点发送请求投票的 RPCs 来给自己投票。
// // 候选人的状态维持直到发生以下任何一个条件发生的时候，
// // (a) 他自己赢得了这次的选举，(b) 其他的服务器成为领导者，(c) 一段时间之后没有任何一个获胜的人。
// // 当一个候选人从整个集群的大多数服务器节点获得了针对同一个任期号的选票，那么他就赢得了这次选举并成为领导人。
// // 每一个服务器最多会对一个任期号投出一张选票，按照先来先服务的原则。


// // 状态改变，从一个状态该变成另一个状态
// func (rf *Raft) updateStateTo(state int32) {
// 	// 判断暂时不要
// 	// if rf.isState(state) {
// 	// 	return
// 	// }
// 	// 用于显示
// 	stateDesc := []string{"STATE_LEADER", "STATE_CANDIDATE", "STATE_FOLLOWER"}
// 	// 这里反映各个状态的转换
// 	preState := rf.state
// 	switch state {
// 	case STATE_FOLLOWER:
// 		rf.state = STATE_FOLLOWER
// 		// 转变成follow，投票清空
// 		rf.votedFor = -1
// 	case STATE_CANDIDATE:
// 		rf.state = STATE_CANDIDATE
// 		// 如果是候选人开始选举
// 		// fmt.Println("开始选举")
// 		rf.startElection()
// 	case STATE_LEADER:
// 		// 组织决定是当选
// 		rf.state = STATE_LEADER
// 	default:
// 		fmt.Printf("Warning: invalid state %d, do nothing.\n", state)
// 	}
// 	fmt.Printf("In term %d: Server %d transfer from %s to %s\n",
// 		rf.currentTerm, rf.me, stateDesc[preState], stateDesc[rf.state])

// }

// // 开始选举函数
// func (rf *Raft) startElection() {
// 	// 任期号增加1
// 	// rf.incrementTerm()  //first increase current term 这个函数也是增加1
// 	rf.currentTerm++
// 	// 投票给自己
// 	rf.votedFor = rf.me
// 	// 初始化投票数量
// 	rf.voteCount = 1

// 	// rf.electionTimer.Reset(randElectionDuration())
// 	// 广播投票
// 	rf.broadcastVoteReq()
// }

// func (rf *Raft) getLastIndex() int {
// 	if 	len(rf.log) == 0{
// 		return 0 
// 	}else{
// 		return rf.log[len(rf.log) - 1].Index
// 	}
	
// }
// func (rf *Raft) getLastTerm() int {
// 	if 	len(rf.log) == 0{
// 		return 0
// 	}else{
// 		return rf.log[len(rf.log) - 1].Term
// 	}
// }

// // 
// func (rf *Raft) broadcastVoteReq() {
// 	// 构建投票信息
// 	// args := RequestVoteArgs{Term: atomic.LoadInt32(&rf.currentTerm), CandidateId: rf.me}
// 	// term		int//候选人的任期号
// 	// CandidateId	int//请求选票的候选人的id
// 	// LastLogIndex int//候选人的最后日志条目的索引值
// 	// LastLogTerm	int//候选人最后日志条目的任期号
// 	var args RequestVoteArgs
// 	args.Term = rf.currentTerm
// 	args.CandidateId = rf.me
// 	args.LastLogTerm = rf.getLastTerm()
// 	args.LastLogIndex = rf.getLastIndex()

// 	for i := range rf.peers {
// 		// 对所有人进行广播
// 		if i != rf.me {
// 			go func(i int) {
// 				// 声明返回变量
// 				var reply RequestVoteReply
// 				rf.sendRequestVote(i, &args, &reply)
// 			}(i)
// 		}
// 	}
// }
