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
func SingleInvoiceCheckRedis(singleInvoiceCheckPostData SingleInvoiceCheckPostData) string{
	singleInvoiceCheckPostDataJson,_ := json.Marshal(singleInvoiceCheckPostData)
	jsonData :=PrepareJsonForHttpRequest(singleInvoiceCheckPostDataJson)

	SingleInvoiceCheckUrl := GetUrlFromFactory("SingleInvoiceCheck")

	result := SentHttpequestByPost(SingleInvoiceCheckUrl, jsonData )
	var returnResult ReturnResult
	_ =json.Unmarshal([]byte(result), &returnResult)
	if returnResult.Code=="0000"{
		return result
	}else{
		return "失败"
	}
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

	//if err != nil{
	//	log.Fatal(err)
	//}
	//

	fmt.Println("debug")
	for {
		/*find id in redis
		if nil or "失败":
			redis[ID] = "已提交"
			提交请求
			if 成功：
				redis[ID] = result
			else:
				redis[ID] = "失败"
			reuturn redis[ID]
		elif "已提交"：
			循环等待
		else:
			return redis[ID]
		*/
		searchRedisResult0,err :=redis.Int(conn.Do("exists", id))
		if err != nil{
			log.Fatal(err)
		}

		if searchRedisResult0==0 {
			result := SingleInvoiceCheckRedis(singleInvoiceCheckPostData)
			conn.Do("set",id, result)
			return result
		}else{
			searchRedisResult1,err :=redis.String(conn.Do("get", id))
			if err!=nil{
				log.Fatal("in FlowSingInvoiceCheckThroughRedis",err)
			}
			switch searchRedisResult1 {
			case "失败":
				result :=SingleInvoiceCheckRedis(singleInvoiceCheckPostData)
				switch result {
				case "失败":
					return "失败"
				default:
					_,err:= conn.Do("set",id,result)
					if err!=nil{
						log.Fatal(err)
					}
					return result
				}
			case "已提交":
				continue
			default:
				return searchRedisResult1
			}
		}
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