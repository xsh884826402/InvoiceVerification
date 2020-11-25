package main

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"sync"
)
var pool *redis.Pool
var rwMutex *sync.RWMutex
type ReturnResult struct {
	Success bool `json:"success"`
	Code string `json:"code"`
}
func SingleInvoiceCheckAfterRedis(id string,singleInvoiceCheckPostData SingleInvoiceCheckPostData) string{
	conn := pool.Get()
	defer conn.Close()
	_,err :=conn.Do("set",id,"已提交")
	if err != nil{
		log.Fatal(err)
	}
	singleInvoiceCheckPostDataJson,_ := json.Marshal(singleInvoiceCheckPostData)
	jsonData :=PrepareJsonForHttpRequest(singleInvoiceCheckPostDataJson)

	SingleInvoiceCheckUrl := GetUrlFromFactory("SingleInvoiceCheck")

	result := SentHttpequestByPost(SingleInvoiceCheckUrl, jsonData )
	var returnResult ReturnResult
	_ =json.Unmarshal([]byte(result), &returnResult)
	return result

}

//func SingleInvoiceCheckAfterRedis(id string, singleInvoiceCheckPostData SingleInvoiceCheckPostData) string{
//	conn :=pool.Get()
//	searchRedisResult0,_ :=redis.Int(conn.Do("exists", id))
//	if searchRedisResult0==0 {
//		result := SingleInvoiceCheckRedis(singleInvoiceCheckPostData)
//		return result
//	}else{
//		searchRedisResult1,err :=redis.String(conn.Do("get", id))
//		if err!=nil{
//			log.Fatal("in FlowSingInvoiceCheckThroughRedis",err)
//		}
//		switch searchRedisResult1 {
//		case "失败":
//			//发送请求
//			result := SingleInvoiceCheckRedis(singleInvoiceCheckPostData)
//			return result
//		case "已提交":
//			fmt.Println("已提交，无需多次提交")
//			return "已提交"
//		default:
//			return searchRedisResult1
//		}
//	}
//
//}

func FlowSingleInvoiceCheckThroughRedis(file_str string) string{
	file_str = CheckInputFileType(file_str)
	singleInvoiceCheckPostData := ConvertFileToInvoiceJson(file_str)

	id := singleInvoiceCheckPostData.Fpdm+singleInvoiceCheckPostData.Fphm
	fmt.Println("id", id)
	conn := pool.Get()
	defer conn.Close()

	fmt.Println("debug")
	for {

		searchRedisResult0,err :=redis.Int(conn.Do("exists", id))
		if err != nil{
			log.Fatal(err)
		}

		if searchRedisResult0==0 {
			result := SingleInvoiceCheckAfterRedis(id, singleInvoiceCheckPostData)
			conn.Do("set",id, result)
			return result
		}else{
			searchRedisResult1,err :=redis.String(conn.Do("get", id))
			if err!=nil{
				log.Fatal("in FlowSingInvoiceCheckThroughRedis",err)
			}
			switch searchRedisResult1 {

			case "已提交":
				continue
			default:
				return searchRedisResult1
			}
		}
	}




}


func FlowMultiInvoiceCheckThroughRedis(filenames_str []string) string{
	filenames_str = CheckMultiInputFileType(filenames_str)
	var multiInvoiceInfo []SingleInvoiceCheckPostData
	for _,filename := range filenames_str{
		singleInvoiceCheckPostData := ConvertFileToInvoiceJson(filename)
		multiInvoiceInfo = append(multiInvoiceInfo,singleInvoiceCheckPostData)
	}
	PchNumber := GeneratePchNumber()
	fmt.Println("PchNumber", PchNumber)
	AppendContentToFile("./PchNumber_record", PchNumber)

	multiInvoiceCheckPostData := MultiInvoiceCheckPostData{
		Pch: PchNumber,
		MultiInvoiceInfo: multiInvoiceInfo,
	}
	multiInvoiceCheckPostDataJson,_ :=json.Marshal(multiInvoiceCheckPostData)
	jsonData := PrepareJsonForHttpRequest(multiInvoiceCheckPostDataJson)
	MultiInvoiceCheckUrl := GetUrlFromFactory("MultiInvoiceCheck")
	result := SentHttpequestByPost(MultiInvoiceCheckUrl, jsonData)
	fmt.Println("result", result)

	//添加redis缓存

	return "PchNumber:"+PchNumber+result
}


func FlowMultiResultQueryThroughRedis(PchNumber string) string{
	conn :=pool.Get()
	defer conn.Close()

	searchRedisResult0,err :=redis.Int(conn.Do("exists", PchNumber))
	if err != nil{
		log.Fatal(err)
	}

	if searchRedisResult0==0 {
		multiInvoiceResultQueryPostData := MultiInvoiceCheckPostData{}
		multiInvoiceResultQueryPostData.Pch = PchNumber

		multiInvoiceResultQueryPostDataJson,_ :=json.Marshal(multiInvoiceResultQueryPostData)
		fmt.Println("json data",multiInvoiceResultQueryPostDataJson)
		jsonData := PrepareJsonForHttpRequest(multiInvoiceResultQueryPostDataJson)
		fmt.Println("json data", string(jsonData))
		MultiInvoiceResultQueryUrl := GetUrlFromFactory("MultiInvoiceResultQuery")
		fmt.Println("Url", string(MultiInvoiceResultQueryUrl))
		result := SentHttpequestByPost(MultiInvoiceResultQueryUrl, jsonData)
		//fmt.Println("result", result)
		_,err:=conn.Do("set",PchNumber,result)
		if err!=nil{
			log.Fatal(err)
		}
		return result
	}else{
		searchRedisResult1,err :=redis.String(conn.Do("get", PchNumber))
		if err!=nil{
			log.Fatal("in FlowSingInvoiceCheckThroughRedis",err)
		}
		return searchRedisResult1
	}
}
func testRedis(){
	conn :=pool.Get()
	defer conn.Close()
	result0,_ := redis.Int(conn.Do("EXISTS","c2"))
	result,_ := redis.String(conn.Do("get", "c2"))
	if result0==0 || result=="失败"{

	}
	fmt.Println(result)

}
//func main() {
//	pool = &redis.Pool{
//		MaxIdle:     10,
//		IdleTimeout: 240 * time.Second,
//		Dial: func() (redis.Conn, error) {
//			return redis.Dial("tcp", "localhost:6379")
//		},
//	}
//	testRedis()
//
//}