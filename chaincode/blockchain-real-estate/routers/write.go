/**
 * @Author: 夜央 Oh oh oh oh oh oh (https://github.com/togettoyou)
 * @Email: zoujh99@qq.com
 * @Date: 2020/3/5 4:18 下午
 * @Description: 账户相关合约路由
 */
package routers

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/togettoyou/blockchain-real-estate/chaincode/blockchain-real-estate/lib"
	"github.com/togettoyou/blockchain-real-estate/chaincode/blockchain-real-estate/utils"
	"strconv"
	"time"
)

const Layout = "2006-01-02 15:04:05"//时间常量

//新建用户
func CreateUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	Id :=args[0]
	Score := args[1]//初始积分
	if Id == "" || Score == ""  {
		return shim.Error("参数存在空值")
	}

	// 参数数据格式转换
	var formattedTotalScore float64
	if val, err := strconv.ParseFloat(Score, 64); err != nil {
		return shim.Error(fmt.Sprintf("积分参数格式转换出错: %s", err))
	} else {
		formattedTotalScore = val
	}

	//判断用户是否存在在关系型处理

	NewUser := &lib.User{
		Id: args[0],
		Score:    formattedTotalScore,
	}
	// 写入账本
	if err := utils.WriteLedger(NewUser, stub, lib.UserKey, []string{NewUser.Id}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	NewUserByte, err := json.Marshal(NewUser)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(NewUserByte)
}

//创建数据资源
func UploadResource(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	Id := args[0]
	Hash:=args[1]
	Uploader:= args[2]
	Time:= args[3]
	State := args[4]
	Cost :=  args[5]
	GetScore := args[6]
	if Id == "" || Hash == "" || Uploader==""||Time=="" ||State==""|| Cost ==""{
		return shim.Error("参数存在空值")
	}

	// 参数数据格式转换
	var formattedCost float64
	if val, err := strconv.ParseFloat(Cost, 64); err != nil {
		return shim.Error(fmt.Sprintf("积分参数格式转换出错: %s", err))
	} else {
		formattedCost = val
	}
	var formattedScore float64
	if val, err := strconv.ParseFloat(GetScore, 64); err != nil {
		return shim.Error(fmt.Sprintf("积分参数格式转换出错: %s", err))
	} else {
		formattedScore = val
	}

	NewResource := &lib.Resource{
		Id: args[0],
		Hash: args[1],
		Uploader: args[2],
		Time: args[3],
		State: args[4],
		Cost: formattedCost,
	}
	// 写入资源账本
	if err := utils.WriteLedger(NewResource, stub, lib.ResourceKey, []string{NewResource.Id}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//写入用户帐本
	var account lib.User
	err:=json.Unmarshal(QueryAccount(stub,[]string{args[2]}).Payload,&account)
	if err != nil  {
		return shim.Error(fmt.Sprintf("用户信息验证失败%s", err))
	}

	old_resources:=account.Upload
	new_resources:=append(old_resources,*NewResource)
	account.Upload=new_resources
	account.Score=account.Score+formattedScore

	if err := utils.WriteLedger(account, stub, lib.UserKey, []string{account.Id}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	accountByte, err := json.Marshal(account)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(accountByte)
}

func CreateDeal(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	Sell_id	 := 	args[0]
	Buy_id	 := 	args[1]
	Rescource_id  :=    args[2]
	Cost      :=   args[3]


	if Sell_id == "" || Buy_id == "" || Rescource_id==""|| Cost ==""{
		return shim.Error("参数存在空值")
	}

	// 参数数据格式转换
	var formattedCost float64
	if val, err := strconv.ParseFloat(Cost, 64); err != nil {
		return shim.Error(fmt.Sprintf("积分参数格式转换出错: %s", err))
	} else {
		formattedCost = val
	}
	// time转换。先检查time格式，再转换成time
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return shim.Error(fmt.Sprintf("时区设置失败%s", err))
	}
	time.Local = timeLocal
	formattedTime:= time.Now()

	NewDeal := &lib.Deal{
		Sell_id	 :	args[0],
		Buy_id	 :	args[1],
		Rescource_id  :    args[2],
		Cost     :   formattedCost,
		Time     :   formattedTime,
	}
	// 写入交易账本
	if err := utils.WriteLedger(NewDeal, stub, lib.DealKey, []string{NewDeal.Sell_id,NewDeal.Buy_id,NewDeal.Rescource_id}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}

	//查找卖方
	var account_seller lib.User
	err_seller:=json.Unmarshal(QueryAccount(stub,[]string{Sell_id}).Payload,&account_seller)
	if err_seller != nil  {
		return shim.Error(fmt.Sprintf("用户信息验证失败%s", err_seller))
	}

	//查找买方
	var account_buyer lib.User
	err_buyer:=json.Unmarshal(QueryAccount(stub,[]string{Buy_id}).Payload,&account_buyer)
	if err_buyer != nil  {
		return shim.Error(fmt.Sprintf("用户信息验证失败%s", err_buyer))
	}

	//查找资源
	var resource lib.Resource
	err_resource:=json.Unmarshal(QueryResource(stub,[]string{Rescource_id}).Payload,&resource)
	if err_resource != nil  {
		return shim.Error(fmt.Sprintf("资源信息验证失败%s", err_resource))
	}

	//更新卖方积分，买方buy列表和积分，更新resource的时间

	resource.Time=formattedTime.Format(Layout)
	old_buyresources:=account_buyer.Buy
	new_resources:=append(old_buyresources,resource)
	account_buyer.Buy=new_resources
	account_buyer.Score=account_buyer.Score  -  formattedCost
	account_seller.Score=account_seller.Score + formattedCost

	if err := utils.WriteLedger(account_buyer, stub, lib.UserKey, []string{account_buyer.Id}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	if err := utils.WriteLedger(account_seller, stub, lib.UserKey, []string{account_seller.Id}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	accountByte, err := json.Marshal(account_buyer)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(accountByte)
}



