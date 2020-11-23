package main

import (
	"encoding/json"
	"fmt"
)

type MultiInvoiceResultQueryPostData struct {
	Pch string `json:"pch"`
}

func flowMultiInvoiceResultQuery(PchNumber string){
	multiInvoiceResultQueryPostData := MultiInvoiceCheckPostData{}
	multiInvoiceResultQueryPostData.Pch = PchNumber

	multiInvoiceResultQueryPostDataJson,_ :=json.Marshal(multiInvoiceResultQueryPostData)
	jsonData := PrepareJsonForHttpRequest(multiInvoiceResultQueryPostDataJson)
	MultiInvoiceResultQueryUrl := GetUrlFromFactory("MultiInvoiceResultQuery")
	result := SentHttpequestByPost(MultiInvoiceResultQueryUrl, jsonData)
	fmt.Println("result", result)
}

//func main() {
//	flowMultiInvoiceResultQuery("00000000000000000020201123105043")
//}