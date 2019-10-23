package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func httpGet() {
	resp, err := http.Get("http://www.01happy.com/demo/accept.php?id=1")
	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func httpPost() {
	resp, err := http.Post("http://www.01happy.com/demo/accept.php",
		"application/x-www-form-urlencoded",
		strings.NewReader("name=cjb"))
	if err != nil {
		fmt.Println(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func httpPostForm() {
	resp, err := http.PostForm("http://127.0.0.1:12222",
		url.Values{"name": {"lemon"}, "passwd": {"lemon123"}})

	if err != nil {
		// handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))

}

func loginpost() (content string) {
	//data := `{"name":"lemon", "passwd":"lemon123"}`
	v := url.Values{}
	v.Add("name", "lemon")
	v.Add("passwd", "lemon123")

	request, reqerr := http.NewRequest("POST", "http://127.0.0.1:12222", strings.NewReader(v.Encode()))
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
	loginpost()
}
