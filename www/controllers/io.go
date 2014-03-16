package controllers

import (
	"github.com/icattlecoder/jieshu/www/models"
	"github.com/icattlecoder/tgw"
	"log"
)

// =============================================
type InArgs struct {
	B string
}

type IoArgs struct {
	Id string
	Io string
}

func (s *Server) IoDo(args IoArgs, env tgw.ReqEnv) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}

	if args.Id == "" {
		data["success"] = false
		return
	}
	if args.Io != "in" && args.Io != "out" {
		data["success"] = false
		return
	}

	user := models.UserInfo{}
	err = env.Session.Get("userInfo", &user)
	if err != nil {
		data["success"] = false
		data["info"] = err.Error()
		data["needLogin"] = true
		data["directUrl"] = "/douban/login"
		return
	}

	email := user.Email
	if email == "" {
		data["success"] = false
		data["needLogin"] = true
		data["directUrl"] = "/usercomplete"
		return
	}

	//判断是否已经添加过
	n, err := s.coll.Find(models.D{"id": args.Id, args.Io: user.Uid}).Count()
	if err != nil {
		data["success"] = false
		data["info"] = err.Error()
		return
	}
	if n > 0 {
		data["success"] = false
		data["info"] = "无效的重复操作"
		return
	}

	err = s.coll.Update(models.D{"id": args.Id}, models.D{"$push": models.D{args.Io: user.Uid}})
	if err != nil {
		data["success"] = false
		data["info"] = err.Error()
		return
	}

	err = s.UserMgr.InOut(user, args.Id, args.Io)
	if err != nil {
		log.Println(err)
	}

	data["success"] = true
	return

}

func (s *Server) Io(args InArgs, user *models.UserInfo) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}
	if user != nil {
		data["user"] = user
	}
	data["catalog"] = s.data["catalog"]
	book := models.Book{}
	err = s.coll.Find(models.D{"id": args.B}).One(&book)
	if err != nil {
		return
	}
	InUsers := []models.UserInfo{}
	for _, v := range book.In {
		InUsers = append(InUsers, s.UserMgr.Users[v])
	}

	OutUsers := []models.UserInfo{}
	for _, v := range book.Out {
		OutUsers = append(OutUsers, s.UserMgr.Users[v])
	}

	data["book"] = book
	data["inUsers"] = InUsers
	data["outUsers"] = OutUsers
	return
}

// =============================================
type InAddArgs struct {
	Email string
	Id    string
}
