package main

import (
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"testing"
	"time"
)

func TestCreateVehicle(T *testing.T) {
	vs := new(VehicleSharing)                      // 创建Chaincode对象
	stub := shim.NewMockStub("VehicleSharing", vs) // 创建MockStub对象
	// 先创建车型
	res := stub.MockInvoke("1", [][]byte{[]byte("createVehicle"), []byte("vehicle-1"), []byte("BMW"), []byte("1000000"), []byte(strconv.Itoa(int(time.Now().Unix())))})
	res = stub.MockInvoke("1", [][]byte{[]byte("createVehicle"), []byte("vehicle-2"), []byte("ROLLS"), []byte("2000000"), []byte(strconv.Itoa(int(time.Now().Unix())))})
	res = stub.MockInvoke("1", [][]byte{[]byte("createVehicle"), []byte("vehicle-3"), []byte("MASTERATI"), []byte("3000000"), []byte(strconv.Itoa(int(time.Now().Unix())))})

	// 查询车型
	res = stub.MockInvoke("2", [][]byte{[]byte("findVehicle"), []byte("vehicle-2")})
	fmt.Println("The value of vehicle is", string(res.Payload))

	// 根据车类型查询
	res = stub.MockInvoke("2", [][]byte{[]byte("queryVehiclesByBrand"), []byte("MASTERATI")})
	fmt.Println("The value of vehicle is", string(res.Payload))

	// 根据车类型查询
	res = stub.MockInvoke("2", [][]byte{[]byte("createLease"), []byte("lease-1"), []byte("vehicle-3"), []byte("user-1")})
	fmt.Println("The value of lease is", string(res.Payload))

	// 创建lease
	res = stub.MockInvoke("2", [][]byte{[]byte("createLease"), []byte("lease-1"), []byte("vehicle-3"), []byte("user-1")})
	fmt.Println("The value of lease is", string(res.Payload))

	res = stub.MockInvoke("2", [][]byte{[]byte("findLease"), []byte("lease-1")})
	fmt.Println("The value of lease is", string(res.Payload))
}
