package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CommonPostData struct{
	ZipCode string `json:"zipCode"`
	EncryptCode string	`json:"encryptCode"`
	DataExchangeId string	`json:"dataExchangeId"`
	EntCode string	`json:"entCode"`
	Content string	`json:"content"`
}

type SingleInvoiceCheckPostData struct{
	Jym string `json:"jym"`
	Fpje string `json:"fpje"`
	Fpdm string `json:"fpdm"`
	Kprq string `json:"kprq"`
	Fphm string `json:"fphm"`
	Fpzl string `json:"fpzl"`
}

func SingleInvoiceCheck(url string){
	// build postdata
	singleInvoiceCheckPostData := SingleInvoiceCheckPostData{
		Jym : "",
		Fpje: "12092.26",
		Fpdm: "4100191130",
		Kprq: "20190906",
		Fphm: "07537241",
		Fpzl: "01",
	}

	singleInvoiceCheckPostDataJson,_ :=json.Marshal(singleInvoiceCheckPostData)
	fmt.Println("Single Json", string(singleInvoiceCheckPostDataJson))
	singleInvoiceCheckPostDataEncoded :=Base64Encode(string(singleInvoiceCheckPostDataJson))

	dataExchangeId := CreateRandomDataExchangeId_1()

	commonPostData := CommonPostData{
		ZipCode: "0",
		EncryptCode: "0",
		DataExchangeId: dataExchangeId,
		EntCode: "",
		Content: singleInvoiceCheckPostDataEncoded,
	}
	commonPostDataJson,_ :=json.Marshal(commonPostData)


	fmt.Println("PostDataJson",string(commonPostDataJson))
	client := &http.Client{}
	request,_ := http.NewRequest("POST", url, bytes.NewBuffer(commonPostDataJson))
	request.Header.Set("Content-Type", "application/json")
	resp, _ :=client.Do(request)
	fmt.Println("resp", resp)

	body,_ := ioutil.ReadAll(resp.Body)
	fmt.Println("body", string(body))

	resp_result := CommonPostData{}
	err := json.Unmarshal(body, &resp_result)
	if err != nil{
		log.Fatal(err)
	}

	result,_ := Base64Decode(resp_result.Content)
	fmt.Println("result", result)

}

func main() {
	//tokendata := GetTokenData()
	//fmt.Println(tokendata)
	//fmt.Println("haha")
	token:= GetTokenData()
	UrlSinggleInvoiceCheck := "https://sandbox.ele-cloud.com/api/open-recipt/V1/CheckInvoiceSingle"
	v, ok := token["access_token"].(string)
	if ok {
		UrlSinggleInvoiceCheck += "?" + "access_token=" + v
	} else{
		log.Println("access_token is not string")
	}
	fmt.Println("url",UrlSinggleInvoiceCheck)
	SingleInvoiceCheck(UrlSinggleInvoiceCheck)
}