package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"reflect"
	"time"
)

var MonitorCoreAddress string
var LinkedseeUrl string
var LinkedseeToken string

type SendAlarm struct {
	Receiver string `json:"receiver"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Content  string `json:"content"`
}

type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    []User `json:"data"`
}

type User struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email"`
}

func httpClient(url string, method string, head map[string]string, body interface{}) ([]byte, error) {
	//1.建立客户端
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*2)
				if err != nil {
					return nil, err
				}
				_ = conn.SetDeadline(time.Now().Add(time.Second * 2))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 300,
			DisableKeepAlives:     true,
			MaxIdleConns:          2,
			MaxIdleConnsPerHost:   2,
		},
	}
	var buf []byte

	if reflect.TypeOf(body).String() == "string" {
		buf = []byte(body.(string))
	} else {
		buf, _ = json.Marshal(body)
	}

	if req, err := http.NewRequest(method, url, bytes.NewBuffer(buf)); err != nil {
		log.Println(err.Error())
		return nil, err
	} else {
		req.Close = true
		for key, value := range head {
			req.Header.Set(key, value)
		}

		resp, err := client.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		} else if err != nil {
			return nil, err
		} else {
			log.Println("respone is nul")
			return nil, fmt.Errorf("respone is nul")
		}

		return ioutil.ReadAll(resp.Body)
	}

}
