package mapreduce

import "fmt"

//
// schedule() starts and waits for all tasks in the given phase (Map
// or Reduce). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
		case mapPhase:
			ntasks = len(mapFiles)
			n_other = nReduce
		case reducePhase:
			ntasks = nReduce
			n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)
	// Schedule: 20 Map tasks (10 I/Os)
	// Schedule: 10 Reduce tasks (20 I/Os)
	// var wg sync.WaitGroup
	// 需要声明一个WaitGroup组等待所
	fmt.Println(ntasks)

	fmt.Printf("Schedule: %v phase done\n", phase)
}
