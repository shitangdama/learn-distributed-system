package mapreduce

import "fmt"
import "sync"

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
	// fmt.Println(<-registerChan)
	var wg sync.WaitGroup
	wg.Add(ntasks)
	for taskNum:=0;taskNum<ntasks;taskNum++ {
		workerName:= <-registerChan
		go func (c int){
			deriveWork(workerName,jobName,mapFiles,phase,n_other,c,registerChan,&wg)
		}(taskNum)
	}

	wg.Wait();

	fmt.Printf("Schedule: %v phase done\n", phase)
}

func deriveWork(workerName string,jobName string,mapFiles []string,phase jobPhase,n_other int,c int,registerChan chan string,wg *sync.WaitGroup){
	
	// workername是地址
	// dotask是任务
	if(call(workerName,"Worker.DoTask",DoTaskArgs{jobName,mapFiles[c],phase,c,n_other},nil)){
		wg.Done();
		registerChan<-workerName
	}else{
		workerName := <-registerChan
		// 重新注册的worker获取
		deriveWork(workerName,jobName,mapFiles,phase,n_other,c,registerChan,wg)
	}
}