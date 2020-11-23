package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

func CreateRandomNumber() string {
	return fmt.Sprintf("%015v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}


func CreateRandomDataExchangeId() string{
	current_time := time.Now().Format("20060102150405")
	ms :=time.Now().UnixNano()/1e6
	ms_str :=fmt.Sprintf("%d", ms)
	tail :=ms_str[len(ms_str)-3:]
	fmt.Println("here", tail)
	random_str := CreateRandomNumber()
	return current_time+tail+random_str
}

func CreateRandomDataExchangeId_1() string{
	current_time := time.Now().Format("20060102150405000")
	random_str := CreateRandomNumber()
	return current_time+random_str
}

func Base64Encode(input string) string{
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func Base64Decode(input string) (string, error){
	data, err :=base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(data),err
}

//func ConvertInvoiceTypeToFpdam( invoiceType string) string{
//	//to do
//
//}
func ConvertPdfTojpg(pdfFile string, jpgFile string) string{
	log.Print("ConvertPdfTojpg ",pdfFile,"\n", "ConvertPdfTojpg ",jpgFile)
	command:= exec.Command("magick","convert","-density","300",pdfFile,jpgFile)

	err := command.Run()
	if err != nil{
		log.Fatal(err)
	}
	return jpgFile
}

func CheckInputFileType(filename string) string{
	if filename[len(filename)-4:]==".jpg"{
		return filename
	}
	if filename[len(filename)-4:]==".png"{
		return filename
	}
	if filename[len(filename)-4:]=="jpeg"{
		return filename
	}

	if filename[len(filename)-4:]==".pdf"{
		jpgfile :=filename[:len(filename)-4]+".jpg"
		ConvertPdfTojpg(filename, jpgfile)
		return jpgfile
	}
	return ""
}

func ConvertFileToInvoiceJson(filename string) SingleInvoiceCheckPostData{
	InvoiceInfoJson := GetInvoiceInfoByBaiduai(filename)
	//fmt.Println("Invoice Info Json\n", string(InvoiceInfoJson))
	//
	var invoiceInfoBaiduai InvociceInfoBaiduai
	err :=json.Unmarshal(InvoiceInfoJson, &invoiceInfoBaiduai)
	if err != nil{
		log.Fatal(err)
	}

	singleInvoiceCheckPostData := SingleInvoiceCheckPostData{

	}
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
	} else{
		if invoiceInfoBaiduai.Words_result.InvoiceType=="专用发票"{
			singleInvoiceCheckPostData.Fpzl = "01"
		}
	}

	if singleInvoiceCheckPostData.Fpzl =="10"{
		singleInvoiceCheckPostData.Jym = invoiceInfoBaiduai.Words_result.CheckCode[len(invoiceInfoBaiduai.Words_result.CheckCode)-6:]
	}
	return singleInvoiceCheckPostData
}

func PrepareJsonForHttpRequest(jsonData []byte) []byte{
	jsonDataEncoded :=Base64Encode(string(jsonData))

	dataExchangeId := CreateRandomDataExchangeId_1()

	commonPostData := CommonPostData{
		ZipCode: "0",
		EncryptCode: "0",
		DataExchangeId: dataExchangeId,
		EntCode: "",
		Content: jsonDataEncoded,
	}
	commonPostDataJson,_ :=json.Marshal(commonPostData)
	return commonPostDataJson
}

func GetUrlFromFactory(RequestType string) string{
	token:= GetTokenData()
	//fmt.Println("debug 1", token,RequestType)
	var Url string
	switch RequestType{
	case "MultiInvoiceCheck":
		Url = "https://sandbox.ele-cloud.com/api/open-recipt/V1/MultilCheckInvoice"
		v, ok := token["access_token"].(string)
		if ok {
			Url += "?" + "access_token=" + v
		} else{
			log.Println("access_token is not string")
		}
	case "MultiInvoiceResultQuery":
		Url = "https://sandbox.ele-cloud.com/api/open-recipt/V1/BatchGetInvoice"
		v, ok := token["access_token"].(string)
		if ok {
			Url += "?" + "access_token=" + v
		} else{
			log.Println("access_token is not string")
		}
	case "SingleInvoiceCheck":
		Url = "https://sandbox.ele-cloud.com/api/open-recipt/V1/CheckInvoiceSingle"
		v, ok := token["access_token"].(string)
		if ok {
			Url += "?" + "access_token=" + v
		} else{
			log.Println("access_token is not string")
		}
	default:
		fmt.Println("default")
	
	}
	return Url

}

func SentHttpequestByPost(url string,commonPostDataJson [] byte) string{
	//fmt.Println("Json data", string(commonPostDataJson))
	client := &http.Client{}
	request,_ := http.NewRequest("POST", url, bytes.NewBuffer(commonPostDataJson))
	request.Header.Set("Content-Type", "application/json")
	resp, _ :=client.Do(request)
	//fmt.Println("resp", resp)

	body,_ := ioutil.ReadAll(resp.Body)
	//fmt.Println("body", string(body))

	resp_result := CommonPostData{}
	err := json.Unmarshal(body, &resp_result)
	if err != nil{
		log.Fatal(err)
	}

	result,_ := Base64Decode(resp_result.Content)
	return result
}

func GeneratePchNumber() string{
	current_time := time.Now().Format("20060102150405")
	var head0 string
	for i:=0;i<(32-len(current_time));i++{
		head0 +="0"
	}
	current_time =head0+current_time
	return current_time
}