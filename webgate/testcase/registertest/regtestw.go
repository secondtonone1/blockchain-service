package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func regpost() (content string) {
	//data := `{"name":"lemon", "passwd":"lemon123"}`
	v := url.Values{}
	v.Add("name", "zack")
	v.Add("passwd", "zack123")
	v.Add("email", "zack@163.com")

	request, reqerr := http.NewRequest("POST", "http://127.0.0.1:12222/register", strings.NewReader(v.Encode()))
	if reqerr != nil {
		fmt.Println("http.NewRequest error!")
		return
	}

	// 表单方式(必须)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	defer request.Body.Close()
	client := &http.Client{Timeout: 5 * time.Second}
	resp, rsperr := client.Do(request)
	if rsperr != nil {
		fmt.Println("http request rsp error!")
		return
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	content = string(result)
	fmt.Println(content)
	return
}

func main() {
	//httpPostForm()
	regpost()
}
