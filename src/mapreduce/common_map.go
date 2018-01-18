package mapreduce

import (
	// "fmt"
	"hash/fnv"
	"io/ioutil"
	"encoding/json"
	"os"
)

// doMap manages one map task: it reads one of the input files
// (inFile), calls the user-defined map function (mapF) for that file's
// contents, and partitions the output into nReduce intermediate files.
func doMap(
	jobName string, // the name of the MapReduce job
	mapTaskNumber int, // which map task this is
	inFile string,
	nReduce int, // the number of reduce tasks that will run ("R" in the paper)
	mapF func(file string, contents string) []KeyValue,
) {
	// 
	// fmt.Println("11111111111111111111111")
	// 这里的infile就是对应的mapfunc中创建的1-10w的文件
	// 824-mrinput-0.txt
	data, e := ioutil.ReadFile(inFile)
	if e != nil {
	  panic("ReadFile")
	}
	// data是文件读取的中的数据,是数组格式，数字是1到10w
	// pairs是一个数组，里面对应的格式是一个对象，key是数字，后面是一个空字符
	pairs := mapF(inFile, string(data))
	encoders := make([]*json.Encoder, nReduce)
	// fmt.Println(encoders)
	for i, _ := range encoders {
		// 这里reduceName(jobName, mapTaskNumber, i)产生map的结果名字文件
		// mrtmp.test-0-0
		// fmt.Printf(reduceName(jobName, mapTaskNumber, i))
		// 不存在就创建他
		f, e := os.OpenFile(reduceName(jobName, mapTaskNumber, i), os.O_WRONLY | os.O_CREATE, 0666)
		if e != nil {
		  	panic("OpenFile")
		}
		// 对于json.NewEncoder函数
		// 这里的json.NewEncoder函数是产生的一个可写入的json的pipline，将这个pipline存入encoders数组里面

		encoders[i] = json.NewEncoder(f)
		// fmt.Println(encoders)
		defer f.Close()
	}
	// 上一布针对每一个file产生里一个独特的pipline，存储在encoders中，供调用者写入使用
	// 将很长的pairs按照reduce的数目，分到几个reduce文件中
	for _, pair := range pairs {
		i := ihash(pair.Key) % nReduce
		// 在这里是一次一次输入，所以造成了value都为空的原因
		e := encoders[i].Encode(&pair)
		if e != nil {
			panic("Encode")
		}
	}
}

func ihash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32() & 0x7fffffff)
}
