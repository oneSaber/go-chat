package main

import (
	"DataServer/data"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// 问好测试
func hello(writer http.ResponseWriter, request *http.Request) {
	requestLog("hello")
	fmt.Fprintf(writer, "hello world")
}

// 注册
func registrer(writer http.ResponseWriter, request *http.Request) {
	requestLog("register")
	var user map[string]interface{}
	body, _ := ioutil.ReadAll(request.Body)
	json.Unmarshal(body, &user)
	fmt.Println(user)
	newUser := data.User{
		Account:   user["account"].(string),
		Password:  user["password"].(string),
		Name:      user["name"].(string),
		Email:     user["email"].(string),
		Signature: user["signature"].(string),
		Avatar:    "wait",
	}
	fmt.Println(newUser)
	ok := newUser.Register()
	if ok {
		res := map[string]interface{}{
			"message": "register ok",
			"redict":  "login",
		}
		result, _ := json.Marshal(res)
		writer.WriteHeader(200)
		writer.Write(result)
	} else {
		res := map[string]interface{}{
			"message": "register failure",
		}
		result, _ := json.Marshal(res)
		writer.WriteHeader(400)
		writer.Write(result)
	}
}

// 登陆
func login(writer http.ResponseWriter, request *http.Request) {
	requestLog("login")
	loginInfo := GetJsonFormRequest(request)
	var LoginUser data.User
	LoginUser = data.QueryUser(loginInfo["loginInfo"].(string))
	if LoginUser.Id == 0 {
		res := map[string]interface{}{
			"message": "no user",
			"redict":  "register",
		}
		result, _ := json.Marshal(res)
		writer.WriteHeader(404)
		writer.Write(result)
		return
	}
	if LoginUser.Password != loginInfo["password"].(string) {
		res := map[string]interface{}{
			"message": "password wrong",
			"redict":  "login",
		}
		result, _ := json.Marshal(res)
		writer.WriteHeader(401)
		writer.Write(result)
		return
	}
	sess, err := LoginUser.CreateSession()
	if err != nil {
		panic(err)
		res := map[string]interface{}{
			"message": "login failure",
			"redict":  "register",
		}
		result, _ := json.Marshal(res)
		writer.WriteHeader(400)
		writer.Write(result)
		return
	}
	cookie1 := http.Cookie{
		Name:     "UserId",
		Value:    string(LoginUser.Id),
		HttpOnly: true,
	}
	cookie2 := http.Cookie{
		Name:     "session",
		Value:    sess,
		HttpOnly: true,
	}
	http.SetCookie(writer, &cookie1)
	http.SetCookie(writer, &cookie2)
	res := map[string]interface{}{
		"message": "login successful",
	}
	result, _ := json.Marshal(res)
	writer.WriteHeader(200)
	writer.Write(result)
}

// 退出登陆
func logout(writer http.ResponseWriter, request *http.Request) {
	userId, _ := request.Cookie("userId")
	session, _ := request.Cookie("session")
	if ok, _ := data.CheckLogin(userId.Value, session.Value); ok {
		err := data.Logout(userId.Value)
		if err != nil {
			res := map[string]interface{}{
				"message": "logout failure",
				"redict":  "index",
			}
			result, _ := json.Marshal(res)
			writer.WriteHeader(400)
			writer.Write(result)
			panic(err)
			return
		} else {
			res := map[string]interface{}{
				"message": "login successful",
				"redict":  "index",
			}
			result, _ := json.Marshal(res)
			writer.WriteHeader(200)
			writer.Write(result)
			return
		}
	} else {
		res := map[string]interface{}{
			"message": "no login",
			"redict":  "login",
		}
		result, _ := json.Marshal(res)
		writer.WriteHeader(404)
		writer.Write(result)
		return
	}
}

// 上传头像
