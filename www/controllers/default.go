package controllers

import (
	"github.com/icattlecoder/jieshu/www/models"
	"github.com/icattlecoder/tgw"
	"labix.org/v2/mgo"
)

import "github.com/dchest/captcha"

type Server struct {
	coll *mgo.Collection
	*models.UserMgr
	data map[string]interface{}
}

func NewServer(c *mgo.Collection, user_coll *mgo.Collection) *Server {
	userMgr := models.NewUserMgr(user_coll)

	data := map[string]interface{}{}
	data["catalog"] = models.GetBookCatalog()
	return &Server{coll: c, UserMgr: userMgr, data: data}
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
