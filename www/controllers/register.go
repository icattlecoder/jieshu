//用户注册
package controllers

import (
	"encoding/json"

	"github.com/icattlecoder/jieshu/www/models"
	"github.com/icattlecoder/tgw"
	"io/ioutil"
	"log"
	"net/url"
	// "log"
	"net/http"
)

func (s *Server) Logout(env tgw.ReqEnv) {
	env.Session.Clear("userInfo")
	http.Redirect(env.RW, env.Req, "http://"+env.Req.Host, 302)
}

func (s *Server) DoubanLogin(env tgw.ReqEnv) (data map[string]interface{}, err error) {
	http.Redirect(env.RW, env.Req, "https://www.douban.com/service/auth2/auth?client_id=05fe71c588b205e811fb55509a1611b8&redirect_uri=http://www.4jieshu.com/douban/callback&response_type=code", 302)
	return
}

type UserCompleteArgs struct {
	Email string
}

//UserComplete
func (s *Server) Usercomplete(args UserCompleteArgs, env tgw.ReqEnv) (data map[string]interface{}, err error) {
	if env.Req.Method == "GET" {
		return
	}
	data = map[string]interface{}{}
	userInfo := models.UserInfo{}
	err = env.Session.Get("userInfo", &userInfo)
	log.Println(userInfo)
	if userInfo.Uid == 0 {
		data["tips"] = `DouBan未授权!去<a href="/douban/login">授权</a>`
		return
	}
	if args.Email == "" {
		data["tips"] = "Email不能为空!"
		return
	}
	err = s.UserMgr.UpdateEmail(userInfo.Uid, args.Email)
	if err != nil {
		data["tips"] = err.Error()
		return
	}
	userInfo.Email = args.Email
	env.Session.Set("userInfo", userInfo)
	http.Redirect(env.RW, env.Req, "http://"+env.Req.Host+"/index", 302)
	return
}

type DoubanCallbackArgs struct {
	Code string
}

func (s *Server) DoubanCallback(args DoubanCallbackArgs, env tgw.ReqEnv) {
	access := getAccessToken(args.Code)
	data := map[string]interface{}{}
	err := json.Unmarshal([]byte(access), &data)
	if err != nil {
		log.Println(err)
		return
	}
	if token, ok := data["access_token"]; ok {
		if strToken, ok := token.(string); ok {
			usrData := getDoubanUserInfo(strToken)
			user, err := s.UserMgr.AddDouBan(usrData)
			if err != nil {
				log.Println(err)
			}
			err = env.Session.Set("userInfo", user)
			if err != nil {
				log.Println(err)
			}
			log.Println(user)
			if user.Email == "" {
				http.Redirect(env.RW, env.Req, "http://"+env.Req.Host+"/usercomplete", 302)
				return
			}
		}
	}
	http.Redirect(env.RW, env.Req, "http://"+env.Req.Host+"/index", 302)
}

func getAccessToken(code string) (token string) {
	resp, err := http.PostForm("https://www.douban.com/service/auth2/token",
		url.Values{"client_id": {"05fe71c588b205e811fb55509a1611b8"},
			"client_secret": {"44fc1984a367a30b"},
			"redirect_uri":  {"http://www.4jieshu.com/douban/callback"},
			"grant_type":    {"authorization_code"},
			"code":          {code},
		})
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	str, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	token = string(str)
	return
}

func getDoubanUserInfo(accessToken string) (data map[string]string) {
	data = map[string]string{}
	req, err := http.NewRequest("GET", "https://api.douban.com/v2/user/~me", nil)
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	json.Unmarshal(content, &data)
	return
}

func (s *Server) User(args UserImgArgs) (data map[string]interface{}) {
	data = map[string]interface{}{}
	user, err := s.UserMgr.Get(int64(args.Uid))
	if err != nil {
		log.Println(err)
	}
	data["user"] = user
	return
}
