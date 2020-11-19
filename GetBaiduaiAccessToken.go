package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type GetBaiduaiTokenPostData struct {
	Client_id string
	Client_secret string
}

func GetBaiduaiTokenData() map[string]interface{}{
	getBaiduaiTokenData := GetBaiduaiTokenPostData{
		Client_id: "TewQZ8lG1ypGHCwjIi4gKUdk",
		Client_secret: "8BKQscHaYiM2LxQ5IyGcr9Pe3Fknz8up",
	}


	url := "https://aip.baidubce.com/oauth/2.0/token"
	url +="?grant_type=client_credentials&client_id="+getBaiduaiTokenData.Client_id+"&client_secret="+
		getBaiduaiTokenData.Client_secret
	//fmt.Println("url", url)
	resp, err := http.Get(url)
	if err != nil{
		fmt.Println("err", err)
	}

	body, _ :=ioutil.ReadAll(resp.Body)
	var BaiduaiTokenData map[string] interface{}
	_ = json.Unmarshal(body, &BaiduaiTokenData)
	return BaiduaiTokenData

}
//
//func main() {
//	mapp := GetBaiduaiTokenData()
//	fmt.Println(mapp)
//}