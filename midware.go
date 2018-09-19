package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func jsonfy(result map[string]interface{}) interface{} {
	b, err := json.Marshal(result)
	if err != nil {
		return nil
	}
	fmt.Println("b:" + string(b))
	return b
}

// 控制台日志
func requestLog(url string) {
	requestTime := time.Now()
	fmt.Println("request time : " + requestTime.String() + "\nrequest url: " + url)
}

// 从request中提取json
func GetJsonFormRequest(request *http.Request) map[string]interface{} {
	var jsonData map[string]interface{}
	body, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(body, &jsonData)
	return jsonData
}
