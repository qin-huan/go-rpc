package main

import (
	"flag"
	"fmt"
	"go-rpc/etcd"
	"go-rpc/lotusMiner"
	"go-rpc/minerList"
	"log"
	"strconv"
	"strings"
)

type Miner struct {
	authToken string
	ipAddr string
}

func initMinerList (miners []Miner, ratio []int) *minerList.MinerNode {
	if len(miners) != len(ratio) {
		log.Fatalln("length of miners and ratio is not same")
	}
	num := len(miners)

	root := minerList.New(0)
	current := root

	for i:=0; i<num; i++ {
		for j:=0; j<ratio[i]; j++ {
			apiMiner := lotusMiner.Connect(miners[i].authToken, miners[i].ipAddr)
			current.Insert(apiMiner, current)
			current = current.Next
		}
	}
	return root
}

func main() {
	minerNum := flag.Int("miners", 0, "numbers of miners")
	ratioStr := flag.String("ratio", "0:0:0", "ratio of miners to seal data")

	flag.Parse()
	input := flag.Args()

	miners := make([]Miner, *minerNum)
	ratio := make([]int, *minerNum)

	ratios := strings.Split(*ratioStr,":")
	for i:=0; i<*minerNum; i++ {
		num, err := strconv.Atoi(ratios[i])
		if err != nil {
			log.Fatalln("error.strconv.Atoi: ", err.Error())
		}
		ratio[i] = num

		miners[i].ipAddr = input[i*2]
		miners[i].authToken = input[i*2+1]
	}
	root := initMinerList(miners, ratio)
	conn := etcd.Connect()
	lotusMiner.Run(conn, root)
}

func test() {
	fmt.Println("this is a test function!")
}