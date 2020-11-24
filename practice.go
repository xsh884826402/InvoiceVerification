package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func sayHello(w http.ResponseWriter, r *http.Request){
	count := 1
	r.ParseForm()
	for k,v := range r.Form{
		fmt.Println("key", k)
		fmt.Println("value", strings.Join(v, ""))
	}
	fmt.Println(time.Now().Format("20060102150405"))
	time.Sleep(10*time.Second)

	fmt.Fprintf(w, "hello", count)
	count += 1
}


func main() {
	http.HandleFunc("/hello", sayHello)

	err :=http.ListenAndServe(":9090", nil)
	if err != nil{
		log.Fatal("listenAnd Serve, err")
	}
}