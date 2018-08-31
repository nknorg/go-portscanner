package portscanner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var waitForResult = 1 * time.Second

type response struct {
	RequestID string `json:"request_id"`
}

type result map[string][]struct {
	Time    float32 `json:"time"`
	Address string  `json:"address"`
	Error   string  `json:"error"`
}

func getResponseBody(url string) ([]byte, error) {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(res.Body)
}

func checkResult(url string) (bool, error) {
	body, err := getResponseBody(url)
	if err != nil {
		return false, err
	}

	res := result{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return false, err
	}

	for _, connRes := range res {
		if connRes != nil && len(connRes) > 0 {
			if connRes[0].Error == "" && connRes[0].Time > 0 {
				return true, nil
			}
		}
	}

	return false, nil
}

func checkPort(url string) (bool, error) {
	body, err := getResponseBody(url)
	if err != nil {
		return false, err
	}

	r := response{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return false, err
	}

	if r.RequestID == "" {
		return false, errors.New("No results returned from check-host")
	}

	time.Sleep(waitForResult)

	resultURL := fmt.Sprintf("https://check-host.net/check-result/%s", r.RequestID)

	return checkResult(resultURL)
}

func CheckTCP(ip string, port uint16) (bool, error) {
	url := fmt.Sprintf("https://check-host.net/check-tcp?host=%s:%d", ip, port)
	return checkPort(url)
}