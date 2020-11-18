

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type GetTokenPostData struct {
	AppKey string `json:"appKey"`
	AppSecret string `json:"appSecret"`
}

func GetTokenData() map[string]interface{}{
	getTokenPostData := GetTokenPostData{
		AppKey: "Mlfs7n9kofqPMaNVJSFoDcwS",
		AppSecret: "awSW7gts8AS4StGV84HCKVCf",
	}
	getTokenPostDataJson, err :=json.Marshal(getTokenPostData)
	fmt.Println(string(getTokenPostDataJson))
	if err != nil{
		log.Fatal(err)
	}

	client := &http.Client{}

	request,err := http.NewRequest("POST","https://sandbox.ele-cloud.com/api/authen/token", bytes.NewBuffer(getTokenPostDataJson))
	request.Header.Set("Content-Type", "application/json")
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("request", request)
	resp, err := client.Do(request)

	if err != nil{
		log.Fatal(err)
	}
	body, _ :=ioutil.ReadAll(resp.Body)

	var TokenData map[string]interface{}
	err = json.Unmarshal(body, &TokenData)
	if err != nil{
		log.Fatal(err)
	}
	//fmt.Println(TokenData,reflect.TypeOf(TokenData))
	return TokenData
}

func Test(){
	fmt.Println("haha")
}
//func main(){
//	tokendata := GetTokenData()
//	fmt.Println(tokendata)
//}