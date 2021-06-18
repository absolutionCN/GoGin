package util

import (
	"GoGin/config"
	"GoGin/config/logging"
	"encoding/json"
	"errors"
	"github.com/unknwon/com"
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiUrl = "/api/interface/list"
const svcUrl = "/api/project/get"

type ApiResponse struct {
	ErrCode int         `json:"errcode"`
	ErrMsg  string      `json:"errmsg"`
	Data    PageContent `json:"data"`
}

type PageContent struct {
	Count int          `json:"count"`
	Total int          `json:"total"`
	List  []ApiContent `json:"list"`
}

type ApiContent struct {
	ApiId     int           `json:"_id"`
	EditUid   int           `json:"edit_uid"`
	Status    string        `json:"status"`
	ApiOpen   bool          `json:"api_opened"`
	Tag       []interface{} `json:"tagModel"`
	Method    string        `json:"method"`
	Title     string        `json:"title"`
	Path      string        `json:"path"`
	ProjectId int           `json:"project_id"`
	CatId     int           `json:"catid"`
	UId       int           `json:"uid"`
}

func GetSvcApis(srvId int, token string) ([]ApiContent, error) {
	var content = make([]ApiContent, 0)
	params := url.Values{}
	params.Set("project_id", com.ToStr(srvId))
	params.Set("token", token)
	params.Set("limit", "20")
	// new http client
	client := &http.Client{}
	// init page number
	var init = 1
	for {
		reqUrl, err := url.Parse(config.YapiHost + apiUrl)
		if err != nil {
			logging.Error("url path parse error ", err)
			return content, err
		}
		params.Set("page", com.ToStr(init))
		reqUrl.RawQuery = params.Encode()
		req, _ := http.NewRequest("GET", reqUrl.String(), nil)
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			logging.Error("Get Yapi list error", err)
			return content, err
		}
		body, _ := ioutil.ReadAll(resp.Body)
		// 忽略异常，继续请求
		_ = resp.Body.Close()
		var apiResponse ApiResponse
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			logging.Error("Deserialization yapi response error", err)
		}
		// token过期返回
		if apiResponse.ErrCode > 0 {
			return content, errors.New(apiResponse.ErrMsg)
		}
		if apiResponse.ErrCode == 0 && len(apiResponse.Data.List) > 0 {
			content = append(content, apiResponse.Data.List...)
		}
		if apiResponse.Data.Total <= init {
			break
		}
		init++
	}
	return content, nil
}

func GetSvcBasicInfo(token string) (string, error) {
	params := url.Values{}
	params.Set("token", token)
	reqUrl, err := url.Parse(config.YapiHost + svcUrl)
	if err != nil {
		logging.Error("url path parse error ", err)
		return "", err
	}
	reqUrl.RawQuery = params.Encode()
	// new http client
	client := &http.Client{}
	req, _ := http.NewRequest("GET", reqUrl.String(), nil)
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logging.Error("Get Yapi basic info api error", err)
		return "", err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	var s map[string]interface{}
	err = json.Unmarshal(body, &s)
	if err != nil {
		return "", err
	}
	code := int(s["errcode"].(float64))
	if code != 0 {
		return "", errors.New("basic info error code" + com.ToStr(code))
	}
	dataMap := s["data"].(map[string]interface{})
	return com.ToStr(dataMap["basepath"]), nil
}
