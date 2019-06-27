package main

import (
	"bytes"
	"client/bson"
	"client/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

type getveryReq struct {
	Cc      string `json:"cc"`
	Phone   string `json:"phone"`
	Scode   string `json:"scode"`
	Lg      string `json:"lg"`
	Imei    string `json:"imei"`
	Imsi    string `json:"imsi"`
	Company string `json:"company"`
	Os      string `json:"os"`
	Model   string `json:"model"`
	Screenh string `json:"screen_h"`
	Screenw string `json:"screen_w"`
	Source  string `json:"source"`
	Sign    string `json:"sign"`
}
type registerReq struct {
	Cc         string `json:"cc"`
	Phone      string `json:"phone"`
	Scode      string `json:"scode"`
	Lg         string `json:"lg"`
	VeryCode   string `json:"verycode"`
	Imei       string `json:"imei"`
	Imsi       string `json:"imsi"`
	Company    string `json:"company"`
	Os         string `json:"os"`
	Model      string `json:"model"`
	Screenh    string `json:"screen_h"`
	Screenw    string `json:"screen_w"`
	Source     string `json:"source"`
	Sign       string `json:"sign"`
	Sourceuuid string `json:"sourceuuid"`
}
type routerReq struct {
	Cc      string `json:"cc"`
	Phone   string `json:"phone"`
	Scode   string `json:"scode"`
	Imei    string `json:"imei"`
	Imsi    string `json:"imsi"`
	Company string `json:"company"`
	Os      string `json:"os"`
	Model   string `json:"model"`
	Screenh string `json:"screen_h"`
	Screenw string `json:"screen_w"`
	Source  string `json:"source"`
	Sign    string `json:"sign"`
}

var bs = bson.NewBson()
var ut = utils.NewUtils()

func main() {
	socket()
}
func socket() {
	server := "127.0.0.1:8000"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	defer conn.Close()
	go auth(conn, 0)
	// go syncContacts(conn)
	for {
		buffer := make([]byte, 2048)
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("read error: ", err.Error())
		}
		m, _, err := bs.Get(buffer[11:n])
		fmt.Println(m)

	}
}
func getVerCode() {

	data := getveryReq{
		Cc:      "+86",
		Phone:   "+8617600113331",
		Scode:   "CN",
		Lg:      "0",
		Imei:    "123456789999999",
		Imsi:    "460100515667308",
		Company: "ZM001",
		Os:      "MTK60D",
		Model:   "ZM_TEST",
		Screenh: "160",
		Screenw: "128",
		Source:  "0",
		Sign:    "d18cb5b6cb54a5b9fdc2ff5227905b9d",
	}
	bytesData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	url := "http://0.0.0.0:1101/ztalk/getvery"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(respBytes))
}

func register() {

	data := registerReq{
		Cc:       "+86",
		Phone:    "+8617600113331",
		Scode:    "CN",
		Lg:       "0",
		VeryCode: "1234",
		Imei:     "123456789999999",
		Imsi:     "460100515667308",
		Company:  "ZM001",
		Os:       "MTK60D",
		Model:    "ZM_TEST",
		Screenh:  "160",
		Screenw:  "128",
		Source:   "0",
		Sign:     "7bc4994e0806e271baa45b2a9a887333",
	}
	bytesData, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	url := "http://0.0.0.0:1101/ztalk/register"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(string(respBytes))
}
func auth(conn net.Conn, reqType int) {
	data := make(map[string]interface{})
	data["phone"] = "+8617600113331"
	data["source"] = "0"
	data["type"] = reqType
	if reqType == 1 {
		sign := fmt.Sprintf("%s%s%02x%s", data["phone"].(string), data["source"].(string), []byte("kjhi|X`jR]0]bsAG"), "006fef6cce9e2900d49f906bef179bf1")
		fmt.Println(sign)
		data["sign"] = ut.Md5(sign)
	} else {
		if contents, err := ioutil.ReadFile("./nonce.dat"); err == nil {
			result := strings.Replace(string(contents), "\n", "", 1)
			fmt.Println(result)
			sign := fmt.Sprintf("%s%s%02x%s%s", data["phone"].(string), data["source"].(string), []byte("kjhi|X`jR]0]bsAG"), result, "006fef6cce9e2900d49f906bef179bf1")
			data["sign"] = ut.Md5(sign)
		}
	}
	b := bs.Set(data, 1)
	conn.Write(b)
}
func syncContacts(conn net.Conn) {
	data := make(map[string]interface{})
	data["phone"] = "+8617600113331"
	data["add"] = []interface{}{
		"+8618768144506",
		"123123123",
	}
	data["del"] = []interface{}{}
	b := bs.Set(data, 5)
	conn.Write(b)
}
