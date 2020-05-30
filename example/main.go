package main

import (
	"chain_sdk/client"
	"fmt"
	"github.com/sirupsen/logrus"
)

func main() {
	client := &client.Client{
		ConfigPath:     "config.yaml",
		Username:       "Admin",
		Organization:   "org1",
		ChannelID:      "mychannel",
		Chaincode:      "mycc",
	}
	client.Init()
	// vehicle, err := client.CreateVehicle([][]byte{[]byte("uuid-1"), []byte("BWM"), []byte("1000000"), []byte(strconv.FormatInt(time.Now().Unix(), 10))})
	// if err != nil {
	// 	logrus.Error(err.Error())
	// }
	// fmt.Printf("payload:%v\n", vehicle)

	data, err := client.FindVehicle([][]byte{[]byte("uuid-1")})
	if err != nil {
		logrus.Error(err.Error())
	}
	fmt.Printf("payload:%v\n", data)
}
