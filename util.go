package main

import (
	"encoding/base64"
	"fmt"
	"math/rand"
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