package main

import (
	"chain_sdk/client"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"sync"
	"time"
)

type Resource struct {
	ResourceID   uint64 `json:"resourceID"`
	ResourceName string `json:"resourceName"`
}

func main() {
	client := &client.Client{
		ConfigPath:   "config.yaml",
		Username:     "Admin",
		Organization: "org1",
		ChannelID:    "mychannel",
		Chaincode:    "mycc",
	}
	client.Init()
	// vehicle, err := client.CreateVehicle([][]byte{[]byte("uuid-1"), []byte("BWM"), []byte("1000000"), []byte(strconv.FormatInt(time.Now().Unix(), 10))})
	// if err != nil {
	// 	logrus.Error(err.Error())
	// }
	// fmt.Printf("payload:%v\n", vehicle)

	resource := &Resource{
		ResourceID:   uint64(time.Now().Unix()),
		ResourceName: "wjj",
	}
	resourceData, err := json.Marshal(resource)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(3)
	for i := 0; i < 3; i++ {
		go func(i int) {
			defer wg.Done()
			data, err := client.Client.ChannelExecute(
				channel.Request{ChaincodeID: client.Chaincode, Fcn: "CheckDBConstraint", Args: [][]byte{resourceData}},
				channel.WithTargetEndpoints(client.TargetEndPoint...),
			)
			if err != nil {
				fmt.Printf("%d==============ChannelExecute err:%s\n", i, err.Error())
			} else {
				fmt.Printf("%d==============TransactionID:%s, TxValidationCode:%d \n", i, data.TransactionID, data.TxValidationCode)
			}

		}(i)
	}

	wg.Wait()
}
