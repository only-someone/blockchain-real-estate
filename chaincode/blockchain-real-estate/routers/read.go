package routers

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/togettoyou/blockchain-real-estate/chaincode/blockchain-real-estate/lib"
	"github.com/togettoyou/blockchain-real-estate/chaincode/blockchain-real-estate/utils"
)

//查询某用户上传资源
func QueryUpload(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var resourcelist []lib.Resource
	var account lib.User
	err:=json.Unmarshal(QueryAccount(stub,args).Payload,&account)
	if err != nil  {
		return shim.Error(fmt.Sprintf("用户信息验证失败%s", err))
	}

	for _,r :=range account.Upload{
		//unMarshal  []byte 字符串到结构体json等
		//marshal  结构体变json byte
		resourcelist = append(resourcelist, r)
	}

	resourcelistByte, err := json.Marshal(resourcelist)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryAccountList-序列化出错: %s", err))
	}
	return shim.Success(resourcelistByte)
}

func QueryBuyResources(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var resourcelist []lib.Resource
	var account lib.User
	err:=json.Unmarshal(QueryAccount(stub,args).Payload,&account)
	if err != nil  {
		return shim.Error(fmt.Sprintf("用户信息验证失败%s", err))
	}

	for _,r :=range account.Buy{
		//unMarshal  []byte 字符串到结构体json等
		//marshal  结构体变json byte
		resourcelist = append(resourcelist, r)
	}

	resourcelistByte, err := json.Marshal(resourcelist)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryAccountList-序列化出错: %s", err))
	}
	return shim.Success(resourcelistByte)
}

func QueryAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var account lib.User //其实就一个
	result, err := utils.GetStateByPartialCompositeKeys(stub, lib.UserKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}

	err = json.Unmarshal(result[0], &account)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryAccountList-反序列化出错: %s", err))
	}


	accountByte, err := json.Marshal(account)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryAccountList-序列化出错: %s", err))
	}


	return shim.Success(accountByte)
}

func QueryResource(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var resource lib.Resource //其实就一个
	result, err := utils.GetStateByPartialCompositeKeys(stub, lib.ResourceKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	err = json.Unmarshal(result[0], &resource)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryResource-反序列化出错: %s", err))
	}


	resourceByte, err := json.Marshal(resource)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryResource-序列化出错: %s", err))
	}
	return shim.Success(resourceByte)
}

func QueryAllAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var accountList []lib.User //其实就一个
	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.UserKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}

	for _, v := range results {
		if v != nil {
			var account lib.User
			err := json.Unmarshal(v, &account)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryAccountList-反序列化出错: %s", err))
			}
			accountList = append(accountList, account)
		}
	}

	accountListByte, err := json.Marshal(accountList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryAccountList-序列化出错: %s", err))
	}
	return shim.Success(accountListByte)
}

func QueryDeal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	//需要测试三个键是否都能查
	var dealList []lib.Deal
	results, err := utils.GetStateByPartialCompositeKeys(stub, lib.DealKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}

	for _, v := range results {
		if v != nil {
			var deal lib.Deal
			err := json.Unmarshal(v, &deal)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryAccountList-反序列化出错: %s", err))
			}
			dealList = append(dealList, deal)
		}
	}

	dealListByte, err := json.Marshal(dealList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryAccountList-序列化出错: %s", err))
	}
	return shim.Success(dealListByte)
}