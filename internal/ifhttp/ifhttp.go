package ifhttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var client = http.Client{
	Timeout: 10 * time.Second,
}

func HttpPostJson(url string, data interface{}, result interface{}, header map[string]string) error {
	b, _ := json.Marshal(data)
	fmt.Println(string(b))
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	if err := encoder.Encode(data); err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")
	if header != nil {
		for k, v := range header {
			request.Header.Add(k, v)
		}
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	// decoder := json.NewDecoder(response.Body)
	// if err = decoder.Decode(&result); err != nil {
	// 	return err
	// }

	return nil
}
