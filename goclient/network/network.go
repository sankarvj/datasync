package network

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	DOMAIN_NAME = "localhost:8080"
)

const (
	method_get    = "GET"
	method_post   = "POST"
	method_put    = "PUT"
	method_delete = "DELETE"

	APPLICATION_JSON string = "application/json"
)

type NetworkClient struct {
	Url       string
	AuthToken string
	Method    string
}

func makeCallToServer(method string, path string, jsonbody []byte) Response {
	networkClient := formNetworkClient(method, path)
	//Send Request and Parse Json
	return requestServer(networkClient, APPLICATION_JSON, bytes.NewBuffer(jsonbody))
}

var client *http.Client

func requestServer(networkClient *NetworkClient, content_type string, bodybytes *bytes.Buffer) Response {
	if client == nil {
		client = &http.Client{}
	}
	req, _ := http.NewRequest(networkClient.Method, networkClient.Url, bodybytes)
	//Set Headers
	req.Header.Set("User-Agent", "mobile_device_9a28u6b88exx10")
	req.Header.Set("Content-Type", content_type)
	req.Header.Set("Accept", content_type)

	res, err := client.Do(req)
	if err != nil {
		log.Println("problem connecting to network ", err)
		var errorresponse = new(Response)
		errorresponse.Id = ResponseNetworkError
		return *errorresponse
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("problem reading response/problem in server ", err)
		var errorresponse = new(Response)
		errorresponse.Id = ResponseNetworkError
		return *errorresponse
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()

	var s = new(Response)
	if err := dec.Decode(&s); err != nil {
		log.Println("problem parsing/problem in server ", err)
		log.Println("problem causing response ", string(body))
		var errorresponse = new(Response)
		errorresponse.Id = ResponseNetworkError
		return *errorresponse
	}

	// -----------------UNCOMMENT IT TO SEE THE FORMATTED RESULT--------------------------
	// data, err := json.MarshalIndent(s, "", "  ")
	// if err != nil {
	// 	fmt.Printf("MarshalIndent %v", err)
	// }
	// fmt.Printf("%s\n", data)

	return *s
}

func formNetworkClient(method string, path string) *NetworkClient {
	networkClient := new(NetworkClient)
	networkClient.Url = "http://" + DOMAIN_NAME + "/" + path
	networkClient.Method = method
	return networkClient
}

func structToStr(responseStruct interface{}) string {
	out, err := json.Marshal(responseStruct)
	if err != nil {
		log.Println("can't marshal response")
	}
	return string(out)
}
