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

func (s *Server) Register(args AddUserArgs, env tgw.ReqEnv) (data map[string]interface{}, err error) {
	if env.Req.Method == "GET" {
		return
	}
	data = map[string]interface{}{}

	if args.Verify == 0 {
		data["tips"] = "校验码不正确"
		return
	}
	i := 0
	err = env.Session.Get("verify", &i)
	if err != nil {
		data["tips"] = err.Error()
		return
	}

	if args.Email == "" {
		data["tips"] = "Email不能为空!"
		return
	}

	if args.Password == "" {
		data["tips"] = "密码不能为空!"
		return
	}

	userInfo := models.UserInfo{Password: args.Password, Email: args.Email}
	err = s.UserMgr.Add(userInfo)
	if err != nil {
		data["tips"] = err.Error()
		return
	}
	if err = s.login(args.Email, args.Password, &env); err == nil {
		http.Redirect(env.RW, env.Req, "http://"+env.Req.Host+"/index", 301)
		return
	}
	data["tips"] = err.Error()
	return
}

type AddUserArgs struct {
	Password string
	Email    string
	Verify   int
}

func (s *Server) login(username, password string, env *tgw.ReqEnv) (err error) {
	err = s.UserMgr.Valid(username, password)
	if err != nil {
		return
	}

	val := models.UserInfo{Email: username}

	err = env.Session.Set("userInfo", val)
	return
}

type UserLoginArgs struct {
	Password string
	Email    string
	Verify   int
}

func (s *Server) Login(args UserLoginArgs, env tgw.ReqEnv) (data map[string]interface{}, err error) {
	if env.Req.Method == "GET" {
		return
	}
	data = map[string]interface{}{}
	if args.Email == "" || args.Password == "" {
		data["tips"] = "Email或密码不能为空"
		return
	}

	err = s.UserMgr.Valid(args.Email, args.Password)
	if err != nil {
		data["tips"] = err.Error()
		return
	}
	val := models.UserInfo{Email: args.Email}
	err = env.Session.Set("userInfo", val)
	if err != nil {
		data["tips"] = err.Error()
		return
	}

	query := env.Req.URL.Query()

	if returnUrl := query.Get("returnUrl"); returnUrl != "" {
		if b := query.Get("b"); b != "" {
			http.Redirect(env.RW, env.Req, "http://"+env.Req.Host+returnUrl+"?b="+b, 301)
			return
		}
	}

	http.Redirect(env.RW, env.Req, "http://"+env.Req.Host+"/index", 301)
	return

	// env.RW.Write([]byte(err.Error()))
}

func (s *Server) Logout(env tgw.ReqEnv) (data map[string]interface{}, err error) {
	env.Session.Clear("userInfo")
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
	log.Println(data)
	if token, ok := data["access_token"]; ok {
		if strToken, ok := token.(string); ok {
			getDoubanUserInfo(strToken)
		}
	}
	http.Redirect(env.RW, env.Req, "http://"+env.Req.Host+"/index", 302)
}

func getAccessToken(code string) (token string) {
	resp, err := http.PostForm("https://www.douban.com/service/auth2/token",
		url.Values{"client_id": {"05fe71c588b205e811fb55509a1611b8"},
			"client_secret": {"44fc1984a367a30b"},
			"redirect_uri":  {"http://requestb.in/xk1idzxk"},
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

func getDoubanUserInfo(accessToken string) {
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
	/*
		{
			"loc_id":"108296",
			"name":"iwangming",
			"created":"2014-03-14 13:58:53",
			"is_banned":false,
			"is_suicide":false,
			"loc_name":"上海",
			"avatar":"http:\/\/img3.douban.com\/icon\/user_normal.jpg",
			"signature":"",
			"uid":"84779859",
			"alt":"http:\/\/www.douban.com\/people\/84779859\/",
			"desc":"","type":"user","id":"84779859",
			"large_avatar":"http:\/\/img3.douban.com\/icon\/user_large.jpg"
		}
	*/
	log.Println("userInfo:", string(content))
	return
}
