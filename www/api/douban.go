package api

import (
	"qbox.us/encoding/json"
	"io/ioutil"
	"net/url"
	"net/http"
	"errors"
)

type DoubanClient struct{
	ApiKey string
	Secret string
}

func NewDoubanClient(apiKey,secret string) *DoubanClient {
	return &DoubanClient{ApiKey: apiKey, Secret: secret }
}

func (s *DoubanClient)getAccessToken(code string) (token string,err error) {
	resp, err := http.PostForm("https://www.douban.com/service/auth2/token",
		url.Values{"client_id": {s.ApiKey},
			"client_secret": {s.Secret},
			"redirect_uri":  {"http://www.4jieshu.com/douban/callback"},
			"grant_type":    {"authorization_code"},
			"code":          {code},
		})
	if err != nil {
		return
	}
	defer resp.Body.Close()


	decoder := json.NewDecoder(resp.Body)

	data2 := map[string]interface{}{}
	if err = decoder.Decode(&data2);err == nil {
		if objtoken,ok := data2["access_token"];ok{
			token,_ = objtoken.(string)
			return
		}
		err =  errors.New("No access_token ")
	}
	return
}

func (s *DoubanClient)GetDoubanUserInfo(code string) (data map[string]string,err error) {

	access, err := s.getAccessToken(code)
	if err != nil{
		return
	}
	data = map[string]string{}
	req, err2 := http.NewRequest("GET", "https://api.douban.com/v2/user/~me", nil)
	if err2 != nil {
		err = err2
		return
	}
	req.Header.Add("Authorization", "Bearer " + access)
	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		err = err2
		return
	}
	defer resp.Body.Close()
	content, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		err = err2
		return
	}
	err = json.Unmarshal(content, &data)
	return
}
