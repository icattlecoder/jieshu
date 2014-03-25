package controllers

import (
	"github.com/dchest/captcha"
	"github.com/icattlecoder/jieshu/www/models"
	"github.com/icattlecoder/tgw"
	"io"
	"labix.org/v2/mgo"
	"log"
	"net/http"
)

type Config struct {
	DoubanApiKey string
	DoubanSecret string
	ImageServer  string
}


type Server struct {
	Config *Config
	coll *mgo.Collection
	*models.UserMgr
	data map[string]interface{}
}

func NewServer(c *mgo.Collection, user_coll *mgo.Collection, cfg *Config) *Server {

	userMgr := models.NewUserMgr(user_coll)
	data := map[string]interface{}{}
	data["catalog"] = models.GetBookCatalog()
	return &Server{coll: c, UserMgr: userMgr, data: data, Config: cfg}
}

type TestArgs struct {
	Msg string
}

func (s *Server) Test(args TestArgs, env tgw.ReqEnv) {
	env.RW.Write([]byte(args.Msg))
}

func (s *Server) Verify(env tgw.ReqEnv) {

	digest := captcha.RandomDigits(4)

	img := captcha.NewImage(digest, 100, 40)

	verifyCode := 1000*int(digest[0]) + 100*int(digest[1]) + 10*int(digest[2]) + int(digest[3])

	if err := env.Session.Set("verify", verifyCode); err == nil {
		img.WriteTo(env.RW)
		return
	}
}

type UserImgArgs struct {
	Uid int
}

func (s *Server) UserImg(args UserImgArgs, env tgw.ReqEnv) {
	//TODO 性能优化，生成静态图片
	if args.Uid == 0 {
		return
	}
	user, err := s.UserMgr.Get(int64(args.Uid))
	if err != nil {
		return
	}
	resp, err := http.Get(s.Config.ImageServer + user.Email)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	io.Copy(env.RW, resp.Body)
}

func (s *Server) About()  {}
func (s *Server) Advise() {}
