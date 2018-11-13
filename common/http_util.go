package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type responseData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// post访问url
func PostUrl(_url string, _data interface{}) error {
	bytesData, err := json.Marshal(_data)
	if err != nil {
		return fmt.Errorf("序列化json出错. %v", err)
	}

	resp, err := http.Post(_url, "application/json;charset=utf-8",
		bytes.NewBuffer(bytesData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("获取返回数据失败. %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("code: %d. 访问pili失败. %s", resp.StatusCode, string(body))
	}

	result := new(responseData)
	err = json.Unmarshal(body, result)
	if err != nil {
		return fmt.Errorf("解析返回数据失败. %v", err)
	}

	if result.Code != 20000 {
		return fmt.Errorf("%v", result.Message)
	}

	return nil
}

// get请求
func GetUrl(_url string, _query string) error {
	_url += _query
	resp, err := http.Get(_url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("获取返回数据失败. %v", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("code: %d. 访问pili失败. %s", resp.StatusCode, string(body))
	}

	result := new(responseData)
	err = json.Unmarshal(body, result)
	if err != nil {
		return fmt.Errorf("解析返回数据失败. %v", err)
	}

	if result.Code != 20000 {
		return fmt.Errorf("%v", result.Message)
	}

	return nil
}

// 下载文件
func DownloadFile(_url string, _std string) error {
	resp, err := http.Get(_url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, err := ioutil.ReadAll(resp.Body)
		errMSG := fmt.Sprintf("http状态码: %d. %s", resp.StatusCode, string(b))
		if err != nil {
			errMSG = fmt.Sprintf("%s. %v", errMSG, err)
		}
		return fmt.Errorf("%v", errMSG)
	}

	out, err := os.Create(_std)
	if err != nil {
		return fmt.Errorf("创建命令文件失败. %s", _std)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("拷贝数据出错. %v", err)
	}

	return nil
}

func PutUrl(_url string, _data interface{}) error {
	d, err := json.Marshal(_data)
	if err != nil {
		return err
	}

	client := &http.Client{}
	body := bytes.NewBuffer(d)
	req, err := http.NewRequest(http.MethodPut, _url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d, err = ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("code: %d. 访问pili失败. %s", resp.StatusCode)
	}

	result := new(responseData)
	err = json.Unmarshal(d, result)
	if err != nil {
		return err
	}

	if result.Code != 20000 {
		return fmt.Errorf("%v", result.Message)
	}

	return nil
}

func GetURLQuery(data map[string]interface{}) string {
	params := make([]string, 0, 1)
	for key, value := range data {
		d := fmt.Sprintf("%v=%v", key, value)
		params = append(params, d)
	}

	if len(params) == 0 {
		return ""
	}

	return fmt.Sprintf("?%s", strings.Join(params, "&"))
}

// get请求
func GetUrlRaw(_url string, _query string) ([]byte, error) {
	_url += _query
	resp, err := http.Get(_url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("获取返回数据失败. %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("code: %d. 访问pili失败. %s", resp.StatusCode, string(body))
	}

	return body, nil
}
