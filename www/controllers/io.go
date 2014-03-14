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
		log.Println("IO:=", args.Io)
		data["success"] = false
		return
	}

	user := models.UserInfo{}

	err = env.Session.Get("userInfo", &user)
	if err != nil {
		data["success"] = false
		data["info"] = err.Error()
		data["needLogin"] = true
		return
	}

	email := user.Email

	// if user, ok := v.(models.UserInfo); !ok {
	// 	data["success"] = false
	// 	data["info"] = ok
	// 	return
	// } else {
	// 	email = user.Email
	// }

	err = s.coll.Update(models.D{"id": args.Id}, models.D{"$push": models.D{args.Io: email}})
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

func (s *Server) Io(args InArgs) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}
	book := models.Book{}
	err = s.coll.Find(models.D{"id": args.B}).One(&book)
	if err != nil {
		return
	}
	data["book"] = book
	return
}

// =============================================
type InAddArgs struct {
	Email string
	Id    string
}
