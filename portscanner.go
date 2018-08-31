package portscanner

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"
)

const waitForResult = 1 * time.Second
const retries = 3

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

func CheckTCP(ip string, port uint16) (isOpen bool, err error) {
	url := fmt.Sprintf("https://check-host.net/check-tcp?host=%s:%d", ip, port)
	for retry := 0; retry < retries; retry++ {
		isOpen, err = checkPort(url)
		if err == nil {
			return isOpen, nil
		}
		delay := math.Pow(2, float64(retry))
		log.Printf("Retry after %f seconds", delay)
		time.Sleep(time.Duration(delay) * time.Second)
	}
	return
}
