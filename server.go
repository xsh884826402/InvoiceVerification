package main

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"net/http"
	"strings"
	"time"
)
type ReturnResult struct {
	Success bool
	Content string
}
func sayHello(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	for k,v := range r.Form{
		fmt.Println("key", k)
		fmt.Println("value", strings.Join(v, ""))
	}

	fmt.Fprintf(w, "hello")
}

func ServerFuncSingleInvoiceCheck(w http.ResponseWriter, r *http.Request){
	r.ParseMultipartForm(32<<20)
	files := r.MultipartForm.File["file"]
	InvoiceFiles := CopyHttpfilesToLocalFiles(files)
	returnResult :=ReturnResult{}
	if len(InvoiceFiles) == 1{
		result := FlowSingleInvoiceCheckThroughRedis(InvoiceFiles[0])
		returnResult.Success=true
		returnResult.Content=result

	}else{
		log.Println("SingleInvoiceCheck must upload 1 file")
		returnResult.Success=false
		returnResult.Content="输入不合法"
	}
	returnResultByte,_ :=json.Marshal(returnResult)
	fmt.Fprintf(w,string(returnResultByte))

}

func ServerFuncMultiInvoiceCheck(w http.ResponseWriter, r *http.Request){
	r.ParseMultipartForm(32<<20)
	files := r.MultipartForm.File["file"]
	InvoiceFiles := CopyHttpfilesToLocalFiles(files)

	returnResult :=ReturnResult{}

	if len(InvoiceFiles)>1{
		result := FlowMultiInvoiceCheckThroughRedis(InvoiceFiles)
		returnResult.Success=true
		returnResult.Content=result
	}else{
		log.Println("MultiInvoiceCheck must upload more than 1 file")
		returnResult.Success=false
		returnResult.Content="输入不合法"
	}
	returnResultByte,_ :=json.Marshal(returnResult)
	fmt.Fprintf(w,string(returnResultByte))

}

func ServerFuncMultiInvoiceResultQuery(w http.ResponseWriter, r *http.Request){
	r.ParseMultipartForm(32<<20)
	PchNumbers:=r.MultipartForm.Value["PchNumber"]
	returnResult :=ReturnResult{}

	if len(PchNumbers)==1{
		res := FlowMultiResultQueryThroughRedis(string(PchNumbers[0]))
		returnResult.Success=true
		returnResult.Content=res
	}else{
		log.Println("Can not found the PchNumber ")
		returnResult.Success=false
		returnResult.Content="输入不合法"
	}
	returnResultByte,_ :=json.Marshal(returnResult)
	fmt.Fprintf(w,string(returnResultByte))
}

func main() {
	pool = &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
	}

	http.HandleFunc("/hello", sayHello)
	http.HandleFunc("/SingleInvoiceCheck", ServerFuncSingleInvoiceCheck)
	http.HandleFunc("/MultiInvoiceCheck", ServerFuncMultiInvoiceCheck)
	http.HandleFunc("/MultiInvoiceResultQuery", ServerFuncMultiInvoiceResultQuery)
	err :=http.ListenAndServe(":9090", nil)
	if err != nil{
		log.Fatal("listenAnd Serve, err")
	}
}