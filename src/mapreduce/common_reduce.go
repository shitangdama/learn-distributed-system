package mapreduce

import (
	// "fmt"
	"os"
	"encoding/json"
	"log"
	"sort"
)

// doReduce manages one reduce task: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	outFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	debug("DEBUG[doReduce]: jobName:%s, reduceTaskNumber:%d, mapTask:%d\n", jobName, reduceTaskNumber, nMap);
	// fmt.Printf("1111111111")
	keyValues := make(map[string][]string)
	for i := 0; i < nMap; i++ {
		// 这里产生了和map中文件一样的filename是和之前map中一样的，为什么没有使用传入name方式
		fileName := reduceName(jobName, i, reduceTaskNumber)

		debug("DEBUG[doReduce]: Reduce File %s\n", fileName)
		file, err := os.Open(fileName)
		if err != nil {
			log.Fatal("Open error: ", err)
		}
		dec := json.NewDecoder(file)
		// fmt.Printf(dec)
		//iterate each pairs in the file
		//append the pair's Value for which has the same Key
		//e.g. before reduce: ["a", "1"] ["b", "2"] ["a", "3"] ["b", "5"] ["c", "6"]
		//     after reduce:  ["a", ["1", "3"] ["b", ["2", "5"]] ["c", ["6"]]
		//send this result to reduceF(user-defined)
		for {
			// 这个是kv结构
			var kv KeyValue
			// type KeyValue struct {
			// 	Key    string
			// 	Value    string
			// }
			err := dec.Decode(&kv)
			if err != nil {
				break
			}
			// json 包提供了Decoder 和 Encoder 用来支持JSON 数据流的读写。函数NewDecoder 和NewEncoder 封装了io.Reader和io.Writer 接口类型。
			_, ok := keyValues[kv.Key]
			// 将数据输出成
			// fmt.Println(kv.Key)
			// fmt.Println(kv.Value)
			// 使用for，循环输出
			if !ok {
				keyValues[kv.Key] = make([]string, 0)
			}
			keyValues[kv.Key] = append(keyValues[kv.Key], kv.Value)
		}
		file.Close()

		var keys []string
		for k, _:= range keyValues {
			keys = append(keys, k)
		}
		// keys是一个数组
	
		sort.Strings(keys)
		// 增加mergefilename文件名
		mergeFileName := mergeName(jobName, reduceTaskNumber)
		debug("DEBUG[doReduce]: mergeFileName: %v\n", mergeFileName)
		merge_file, err := os.Create(mergeFileName)
		if err != nil {
			log.Fatal("ERROR[doReduce]: Create file error: ", err)
		}
		enc := json.NewEncoder(merge_file)
		for _, k := range keys {
			res := reduceF(k, keyValues[k])
			enc.Encode(&KeyValue{k, res})
		}
		merge_file.Close()
	}
}
