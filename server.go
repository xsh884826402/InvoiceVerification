package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

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
	if len(InvoiceFiles) == 1{
		result := FlowSingleInvoiceCheck(InvoiceFiles[0])
		fmt.Fprintf(w, result)
		return
	}else{
		log.Println("SingleInvoiceCheck must upload 1 file")
		fmt.Fprintf(w, "failed. please check the log ")
	}

}

func ServerFuncMultiInvoiceCheck(w http.ResponseWriter, r *http.Request){
	r.ParseMultipartForm(32<<20)
	files := r.MultipartForm.File["file"]
	InvoiceFiles := CopyHttpfilesToLocalFiles(files)

	if len(InvoiceFiles)>1{
		result := FlowMultiInvoiceCheck(InvoiceFiles)
		fmt.Fprintf(w, result)
		return
	}else{
		log.Println("MultiInvoiceCheck must upload more than 1 file")
		fmt.Fprintf(w, "failed. please check the log ")
	}

}

func ServerFuncMultiInvoiceResultQuery(w http.ResponseWriter, r *http.Request){
	r.ParseMultipartForm(32<<20)
	PchNumbers:=r.MultipartForm.Value["PchNumber"]
	if len(PchNumbers)==1{
		res := flowMultiInvoiceResultQuery(string(PchNumbers[0]))
		fmt.Fprintf(w, res)
		return
	}else{
		log.Println("Can not found the PchNumber ")
		fmt.Fprintf(w,"error please check the log")
	}
}
func main() {
	http.HandleFunc("/hello", sayHello)
	http.HandleFunc("/SingleInvoiceCheck", ServerFuncSingleInvoiceCheck)
	http.HandleFunc("/MultiInvoiceCheck", ServerFuncMultiInvoiceCheck)
	http.HandleFunc("/MultiInvoiceRequestQuery", ServerFuncMultiInvoiceResultQuery)
	err :=http.ListenAndServe(":9090", nil)
	if err != nil{
		log.Fatal("listenAnd Serve, err")
	}
}