package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type InvociceInfoBaiduai struct{
	Words_result ValidInvoiceInfoBaiduai `json:"words_result"`
}

type ValidInvoiceInfoBaiduai struct{
	InvoiceCode string
	InvoiceNum string
	InvoiceType string
	CheckCode string
	TotalAmount string
	InvoiceDate string
}
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

func SingleInvoiceCheck(url string, singleInvoiceCheckPostData SingleInvoiceCheckPostData){
	// build postdata
	//singleInvoiceCheckPostData := SingleInvoiceCheckPostData{
	//	Jym : "",
	//	Fpje: "12092.26",
	//	Fpdm: "4100191130",
	//	Kprq: "20190906",
	//	Fphm: "07537241",
	//	Fpzl: "01",
	//}
	//singleInvoiceCheckPostData := SingleInvoiceCheckPostData{
	//	Jym : "",
	//	Fpje: "146.98",
	//	Fpdm: "011001900611",
	//	Kprq: "20200811",
	//	Fphm: "97672880",
	//	Fpzl: "10",
	//}

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
	//fmt.Println("body", string(body))

	resp_result := CommonPostData{}
	err := json.Unmarshal(body, &resp_result)
	if err != nil{
		log.Fatal(err)
	}

	result,_ := Base64Decode(resp_result.Content)
	fmt.Println("result", result)

}

func GetInvoiceInfoByBaiduai() []byte{
	Baiduai_url := "https://aip.baidubce.com/rest/2.0/ocr/v1/vat_invoice"
	BaiduaiTokenData := GetBaiduaiTokenData()
	v, ok := BaiduaiTokenData["access_token"].(string)
	if ok{
		Baiduai_url += "?access_token="+v
	} else{
		log.Fatal("access_token is not sting")
	}

	file_str :="/Users/shenghu/Project/InvoiceVerification/doc/data/"
	file_str += "31605490445_.pic_hd.jpg"
	fmt.Println("file_str",file_str)
	image_data, err := ioutil.ReadFile(file_str)
	if err != nil{
		log.Fatal(err)
	}
	image_data_base64 := base64.StdEncoding.EncodeToString(image_data)

	//fmt.Println(reflect.TypeOf(image_data_base64),image_data_base64)
	params := url.Values{}
	params.Add("image", image_data_base64)

	//fmt.Println("params:", params.Encode())
	request,_ :=http.NewRequest("POST", Baiduai_url, bytes.NewBuffer([] byte(params.Encode())))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp,_:=client.Do(request)
	//resp, _ := http.Post(Baiduai_url, "application/x-www-form-urlencoded", bytes.NewBuffer([] byte(params.Encode())))

	fmt.Println("InvoiceInfo By Baiduai resp", resp)
	body,_ :=ioutil.ReadAll(resp.Body)
	return body
	//var result map[string] interface{}
	//_ =json.Unmarshal(body, &result)
	//fmt.Println("result", result)
	//return result
}

func FlowSingleInvoiceCheck(){
	token:= GetTokenData()
	UrlSinggleInvoiceCheck := "https://sandbox.ele-cloud.com/api/open-recipt/V1/CheckInvoiceSingle"
	v, ok := token["access_token"].(string)
	if ok {
		UrlSinggleInvoiceCheck += "?" + "access_token=" + v
	} else{
		log.Println("access_token is not string")
	}
	fmt.Println("url",UrlSinggleInvoiceCheck)
	InvoiceInfoMap := GetInvoiceInfoByBaiduai()

	var invoiceInfoBaiduai InvociceInfoBaiduai
	err :=json.Unmarshal(InvoiceInfoMap, &invoiceInfoBaiduai)
	if err != nil{
		log.Fatal(err)
	}

	singleInvoiceCheckPostData := SingleInvoiceCheckPostData{
		//Jym : "",
		//Fpje: "12092.26",
		//Fpdm: "4100191130",
		//Kprq: "20190906",
		//Fphm: "07537241",
		//Fpzl: "01",
	}
	fmt.Println("debug ", invoiceInfoBaiduai.Words_result.InvoiceType)
	singleInvoiceCheckPostData.Fpje = invoiceInfoBaiduai.Words_result.TotalAmount
	singleInvoiceCheckPostData.Fpdm = invoiceInfoBaiduai.Words_result.InvoiceCode
	st := invoiceInfoBaiduai.Words_result.InvoiceDate
	st = strings.Replace(st, "年", "", -1)
	st = strings.Replace(st,"月","",-1)
	st = strings.Replace(st,"日","",-1)

	//singleInvoiceCheckPostData.Kprq = st+"haha"
	singleInvoiceCheckPostData.Kprq = st

	singleInvoiceCheckPostData.Fphm = invoiceInfoBaiduai.Words_result.InvoiceNum
	if invoiceInfoBaiduai.Words_result.InvoiceType=="电子普通发票" {
		singleInvoiceCheckPostData.Fpzl = "10"
	}

	if singleInvoiceCheckPostData.Fpzl =="10"{
		singleInvoiceCheckPostData.Jym = invoiceInfoBaiduai.Words_result.CheckCode[len(invoiceInfoBaiduai.Words_result.CheckCode)-6:]
	}
	SingleInvoiceCheck(UrlSinggleInvoiceCheck,singleInvoiceCheckPostData)
}

func main() {
	FlowSingleInvoiceCheck()

}
