package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//
type MultiInvoiceCheckPostData struct{
	Pch string `json:"pch"`
	MultiInvoiceInfo []SingleInvoiceCheckPostData `json:"invoiceList"`
}
func CheckMultiInputFileType(filenames_str [] string) []string{
	var checked_filenames_str []string
	for _,filename := range filenames_str{
		checked_filenames_str = append(checked_filenames_str, CheckInputFileType(filename))
	}
	return checked_filenames_str
}

func FlowMultiInvoiceCheck(filenames_str []string){
	filenames_str = CheckMultiInputFileType(filenames_str)
	var multiInvoiceInfo []SingleInvoiceCheckPostData
	for _,filename := range filenames_str{
		singleInvoiceCheckPostData := ConvertFileToInvoiceJson(filename)
		multiInvoiceInfo = append(multiInvoiceInfo,singleInvoiceCheckPostData)
		}
	PchNumber := GeneratePchNumber()
	fmt.Println("PchNumber", PchNumber)
	_ = ioutil.WriteFile("./PchNumber", []byte(PchNumber),0777)
	multiInvoiceCheckPostData := MultiInvoiceCheckPostData{
		Pch: PchNumber,
		MultiInvoiceInfo: multiInvoiceInfo,
	}
	multiInvoiceCheckPostDataJson,_ :=json.Marshal(multiInvoiceCheckPostData)
	jsonData := PrepareJsonForHttpRequest(multiInvoiceCheckPostDataJson)
	MultiInvoiceCheckUrl := GetUrlFromFactory("MultiInvoiceCheck")
	result := SentHttpequestByPost(MultiInvoiceCheckUrl, jsonData)
	fmt.Println("result", result)


}
//func flowInMultiVoiceCheck([]string){
//
//}


//func main() {
//	filenames := []string{"/Users/shenghu/Project/InvoiceVerification/doc/data/a.pdf", "/Users/shenghu/Project/InvoiceVerification/doc/data/b.png"}
//	fmt.Println("in multiVoice")
//	filenames = CheckMultiInputFileType(filenames)
//
//	FlowMultiInvoiceCheck(filenames)
//}
