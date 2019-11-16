package main

import (
	"testing"
)

// type MasterNode struct {
// 	TrxHash        string
// 	NodeStatus     string
// 	Id             string
// 	BlockNum       string
// 	ExpiryProducer string
// 	IsselfProblock string
// }
// type Produce struct {
// 	Height    int `json:"height:"`
// 	Produceno int `json:"produceno:"`
// }

func TestHttp(t *testing.T) {
	// var masterNode MasterNode
	// var produce Produce
	// var masterNodeList []MasterNode
	// cmd := exec.Command("ulord-cli", "masternodelist", "full")
	// out, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Println(err)
	// }
	// cmd = exec.Command("ulord-cli", "masternode", "current")
	// out1, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Println(err)
	// }
	// err = json.Unmarshal(out1, &produce)
	// if err != nil {
	// 	log.Println(err)
	// }
	// str := strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(string(out)), "{"), "}")
	// linesData := strings.Split(str, ",")

	// for _, lineData := range linesData {
	// 	s := strings.SplitN(lineData, ":", 2)
	// 	txHash := strings.Split(strings.TrimSuffix(strings.TrimPrefix(strings.TrimSpace(s[0]), `"`), `"`), "-")
	// 	masterNode.TrxHash = txHash[0]

	// 	data := strings.Split(strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(s[1], " "+`"`)), `"`), " ")
	// 	masterNode.Id = data[7]
	// 	masterNode.BlockNum = data[6]
	// 	masterNode.ExpiryProducer = data[4]
	// 	masterNode.NodeStatus = data[8]
	// 	if strconv.Itoa(produce.Produceno) == data[7] {
	// 		masterNode.IsselfProblock = "true"
	// 	} else {
	// 		masterNode.IsselfProblock = "false"
	// 	}
	// 	masterNodeList = append(masterNodeList, masterNode)
	// }
	// log.Println(masterNodeList)
}
