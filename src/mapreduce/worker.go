package mapreduce

//
// Please do not modify this file.
//

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"
)

// Worker holds the state for a server waiting for DoTask or Shutdown RPCs
type Worker struct {
	sync.Mutex

	name       string
	Map        func(string, string) []KeyValue
	Reduce     func(string, []string) string
	nRPC       int // quit after this many RPCs; protected by mutex
	nTasks     int // total tasks executed; protected by mutex
	concurrent int // number of parallel DoTasks in this worker; mutex
	l          net.Listener
}

// DoTask is called by the master when a new task is being scheduled on this
// worker.
// 这里是schedule调用的dotask
func (wk *Worker) DoTask(arg *DoTaskArgs, _ *struct{}) error {
	fmt.Printf("%s: given %v task #%d on file %s (nios: %d)\n",
		wk.name, arg.Phase, arg.TaskNumber, arg.File, arg.NumOtherPhase)

	wk.Lock()
	wk.nTasks += 1
	wk.concurrent += 1
	nc := wk.concurrent
	wk.Unlock()

	if nc > 1 {
		// schedule() should never issue more than one RPC at a
		// time to a given worker.
		log.Fatal("Worker.DoTask: more than one DoTask sent concurrently to a single worker\n")
	}
	// 在这里进行了任务处理，运行了doMap，和doReduce
	switch arg.Phase {
	case mapPhase:
		doMap(arg.JobName, arg.TaskNumber, arg.File, arg.NumOtherPhase, wk.Map)
	case reducePhase:
		doReduce(arg.JobName, arg.TaskNumber, mergeName(arg.JobName, arg.TaskNumber), arg.NumOtherPhase, wk.Reduce)
	}

	wk.Lock()
	wk.concurrent -= 1
	wk.Unlock()

	fmt.Printf("%s: %v task #%d done\n", wk.name, arg.Phase, arg.TaskNumber)
	return nil
}

// Shutdown is called by the master when all work has been completed.
// We should respond with the number of tasks we have processed.
func (wk *Worker) Shutdown(_ *struct{}, res *ShutdownReply) error {
	debug("Shutdown %s\n", wk.name)
	wk.Lock()
	defer wk.Unlock()
	res.Ntasks = wk.nTasks
	wk.nRPC = 1
	return nil
}

//这一步很重要
// 这个地方是worker的register
// 这个 Master.Register
// 是那个
// Tell the master we exist and ready to work
func (wk *Worker) register(master string) {
	args := new(RegisterArgs)
	args.Worker = wk.name
	ok := call(master, "Master.Register", args, new(struct{}))
	if ok == false {
		fmt.Printf("Register: RPC %s register error\n", master)
	}
}

// RunWorker sets up a connection with the master, registers its address, and
// waits for tasks to be scheduled.
// 这个是work的rpcserver吧
// 这里是两个worker
func RunWorker(MasterAddress string, me string,
	MapFunc func(string, string) []KeyValue,
	ReduceFunc func(string, []string) string,
	nRPC int,
) {
	debug("RunWorker %s\n", me)
	wk := new(Worker)
	// 使用wk作为公共对象
	// 链接生成server
	// 这里重点
	// 生成 worker的rpcserver
	wk.name = me
	wk.Map = MapFunc
	wk.Reduce = ReduceFunc
	wk.nRPC = nRPC



	rpcs := rpc.NewServer()
	rpcs.Register(wk)


	os.Remove(me) // only needed for "unix"
	l, e := net.Listen("unix", me)
	if e != nil {
		log.Fatal("RunWorker: worker ", me, " error: ", e)
	}
	wk.l = l

	// 这个是之前声明的主节点地址
	fmt.Println("11111111111111111111111")
	// fmt.Println(MasterAddress)
	fmt.Println(nRPC)
	fmt.Println("11111111111111111111111")

	wk.register(MasterAddress)

	// DON'T MODIFY CODE BELOW
	for {
		wk.Lock()
		if wk.nRPC == 0 {
			wk.Unlock()
			break
		}
		wk.Unlock()
		conn, err := wk.l.Accept()
	// 在这里就是10 然后减少
		if err == nil {
			wk.Lock()
			wk.nRPC--
			wk.Unlock()
			go rpcs.ServeConn(conn)
		} else {
			break
		}
	}
	wk.l.Close()
	debug("RunWorker %s exit\n", me)
}
