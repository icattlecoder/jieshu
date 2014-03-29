//用户注册
package controllers

import (
	"github.com/icattlecoder/jieshu/www/models"
	"github.com/icattlecoder/tgw"
	"log"

	// "log"
	"net/http"
)

// Get /logout
func (s *Server) Logout(env tgw.ReqEnv) {
	env.Session.Clear("userInfo")
	http.Redirect(env.RW, env.Req, "http://"+env.Req.Host, 302)
}

// Get /douban/login
func (s *Server) DoubanLogin(env tgw.ReqEnv) (data map[string]interface{}, err error) {
	http.Redirect(env.RW, env.Req, "https://www.douban.com/service/auth2/auth?client_id="+s.Config.DoubanApiKey+ "&redirect_uri=http://www.4jieshu.com/douban/callback&response_type=code", 302)
	return
}
// Get /mock/login
func (s *Server) MockLogin(env tgw.ReqEnv){
	user,err := s.UserMgr.Get(84779859)
	if err !=nil{
		return
	}
	if err = env.Session.Set("userInfo", user);err == nil{
		env.RW.Write([]byte("模拟登录成功"))
	}
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
	userData,err := s.douban.GetDoubanUserInfo(args.Code)
	if err == nil{
		user, err := s.UserMgr.AddDouBan(userData)
		if err != nil {
			log.Println(err)
		}
		err = env.Session.Set("userInfo", user)
		if err != nil {
			log.Println(err)
		}
		//如果没有Email信息，跳转至usercomplete页面
		if user.Email == "" {
			http.Redirect(env.RW, env.Req, "http://"+env.Req.Host+"/usercomplete", 302)
			return
		}
	}
	// env.RW.Write([]byte(err.Error()))
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
