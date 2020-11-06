package lotusMiner

import (
	"context"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/lotus/api/apistruct"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"go-rpc/etcd"
	"go-rpc/minerList"
	"log"
	"modernc.org/mathutil"
	"net/http"
	"time"
)

type State int
const (
	Pledge State = iota
	PledgeDone
	PledgeError
	PreCommit1
	PreCommit1Done
	PreCommit1Error
	PreCommit2
	PreCommit2Done
	PreCommit2Error
	Commit1
	Commit1Done
	Commit1Error
	Commit2
	Commit2Done
	Commit2Error
	Finalize
	FinalizeDone
	FinalizeError
	WinningPost
	WinningPostDone
	WinningPostError
	WindowPost
	WindowPostDone
	WindowPostError
	Unknown
)
const TASKNUMS int = 3
const TASKKEY string = "/task_states/"
var TaskMap = make(map[string]int)

func Run (conn *clientv3.Client , root *minerList.MinerNode) {
	for {
		kvs, err := etcd.Get(conn, TASKKEY)
		if err != nil {
			log.Fatalln("error.EtcdGet : ", err.Error())
		}
		tasks := getTaskNum(kvs)

		if tasks > 0 {
			log.Println("tasks ready to execute is : ", tasks)
			for i:=0; i<tasks; i++ {
				if root.Data == 0 {
					root = root.Next
				}
				apiMiner := root.Data.(apistruct.StorageMinerStruct)
				root = root.Next
				PledgeSector(apiMiner)
			}

			// ensure tasks being inserted into etcd
			for {
				etcdkvs,_ := etcd.Get(conn, TASKKEY)
				nums := getTaskNum(etcdkvs)
				if nums == 0 {
					log.Println("break")
					break
				}
			}
		}
	}
}

func getTaskNum (kvs []*mvccpb.KeyValue) int {
	pr1s := 0
	for _, kv := range kvs {
		if string(kv.Value) == "PreCommit1" || string(kv.Value) == "PreCommit1Error" {
			pr1s++
		}
	}
	return TASKNUMS-pr1s
}

func PledgeSector (api apistruct.StorageMinerStruct ) {//, conn *clientv3.Client, key string) {
	err := api.Internal.PledgeSector(context.Background())
	if err != nil {
		log.Fatalln("error.pledgeSector : ", err)
	}
	minerAddr, err := api.ActorAddress(context.Background())
	if err != nil {
		log.Fatalln("error.ActorAddress : ", err)
	}
	log.Println("miner to pledge sector is : ", minerAddr)
	/*for {
		kvs, err := etcd.Get(conn, key)
		if err != nil {
			log.Fatalln("error.etcdGet : ", err.Error())
		}
		isFinalize := false
		for _, kv := range kvs {
			switch string(kv.Value) {
			case "PreCommit1Error","PreCommit2Error","Commit1Error","Commit2Error"  : {
				log.Println("Error")
				time.Sleep(time.Second*10)
			}
			case "Finalize" : {
				log.Println("Finalize")
				//etcd.Put(conn, string(kv.Key), "FinalizeDone")
				isFinalize = true
				break
			}
			}
		}
		if isFinalize {
			break
		}
	}
	channel <- key+" done"*/
}

func Connect (authToken string, ipAddr string) (apistruct.StorageMinerStruct) {
	headers := http.Header{"Authorization": []string{"Bearer "+authToken}}
	var api apistruct.StorageMinerStruct
	closer, err := jsonrpc.NewMergeClient(
		context.Background(),
		"http://"+ipAddr+"/rpc/v0", "Filecoin",
		[]interface{}{
			&api.CommonStruct.Internal,
			&api.Internal,
		},
		headers)
	if err != nil {
		log.Fatalln("error.connectToMiner: ", err.Error())
	}
	defer closer()
	return api
}

func GetAddress (api apistruct.StorageMinerStruct) (string) {
	// get minerAddress like 'f01000'
	minerAddr, err := api.ActorAddress(context.Background())
	if err != nil {
		log.Fatalln("error.ActorAddress : ", err.Error())
	}
	log.Println("miner address is : ", minerAddr)

	// generate taskid like '/task_states/seal_1000_2/state'
	key := "/task_states/seal_"+minerAddr.String()[2:]
	log.Println("task_state is : ", key)
	return key
}

func initTask () {
	for i:=0; i<25; i++ {
		var state = State(i)
		TaskMap[state.String()] = 0
	}
}

func (this State) String() string {
	res := ""
	switch this {
	case Pledge : res = "Pledge"
	case PledgeDone : res = "PledgeDone"
	case PledgeError : res = "PledgeError"
	case PreCommit1 : res = "PreCommit1"
	case PreCommit1Done : res = "PreCommit1Done"
	case PreCommit1Error : res = "PreCommit1Error"
	case PreCommit2 : res = "PreCommit2"
	case PreCommit2Done : res = "PreCommit2Done"
	case PreCommit2Error : res = "PreCommit2Error"
	case Commit1 : res = "Commit1"
	case Commit1Done : res = "Commit1Done"
	case Commit1Error : res = "Commit1Error"
	case Commit2 : res = "Commit2"
	case Commit2Done : res = "Commit2Done"
	case Commit2Error : res = "Commit2Error"
	case Finalize : res = "Finalize"
	case FinalizeDone : res = "FinalizeDone"
	case FinalizeError : res = "FinalizeError"
	case WinningPost : res = "WinningPost"
	case WinningPostDone : res = "WinningPostDone"
	case WinningPostError : res = "WinningPostError"
	case WindowPost : res = "WindowPost"
	case WindowPostDone : res = "WindowPostDone"
	case WindowPostError : res = "WindowPostError"
	case Unknown : res = "Unknown"
	}
	return res
}

func updateTask (kvs []*mvccpb.KeyValue) {
	for _, kv := range kvs {
		TaskMap[string(kv.Value)]++
	}
}

func test (conn *clientv3.Client, root *minerList.MinerNode) {
	initTask()
	for {
		kvs, err := etcd.Get(conn, TASKKEY)
		if err != nil {
			log.Fatalln("error.EtcdGet : ", err.Error())
		}
		tasks := getTaskNum(kvs)

		if tasks <= TASKNUMS {
			log.Println("tasks ready to execute is : ", TASKNUMS-tasks)
			for i:=0; i<mathutil.Min(TASKNUMS-tasks, TASKNUMS); i++ {
				apiMiner := root.Data.(apistruct.StorageMinerStruct)
				root = root.Next
				PledgeSector(apiMiner)
			}
			// time.Sleep(time.Second*20)
		} else {
			time.Sleep(time.Second*100)
		}
	}
}