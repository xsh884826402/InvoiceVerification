package main

import (
	"encoding/json"
	"fmt"
)

type MultiInvoiceResultQueryPostData struct {
	Pch string `json:"pch"`
}

func flowMultiInvoiceResultQuery(PchNumber string) string{
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
	return result
}

//func main() {
//	flowMultiInvoiceResultQuery("00000000000000000020201123105043")
//}