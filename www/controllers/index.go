package controllers

import (
	"github.com/icattlecoder/jieshu/www/models"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	// "log"
)

// =============================================
type IndexArgs struct {
	Page int
	T    string
}

//  /index
func (s *Server) Index(args IndexArgs, user *models.UserInfo) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}
	if user != nil {
		data["user"] = user
	}
	data["catalog"] = s.data["catalog"]

	return
}

type CatalogArgs struct {
	C       string
	T       string
	Start   int
	Keyword string
}

// ===
func (s *Server) Catalog(args CatalogArgs, user *models.UserInfo) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}
	if user != nil {
		data["user"] = user
	}
	data["catalog"] = s.data["catalog"]
	data["currentCtg"] = args.C
	data["currentTag"] = args.T
	data["tags"] = models.Tags[args.C]

	return
}

func (s *Server) CatalogSearch(args CatalogArgs, user *models.UserInfo) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}
	if user != nil {
		data["user"] = user
	}
	var cond = bson.M{}
	if args.T == "" {
		cond = bson.M{"catalog": args.C}
	} else if args.Keyword == "" {
		cond = bson.M{"catalog": args.C, "tag": args.T}
	} else {
		args.Start = 0
		cond = bson.M{"catalog": args.C, "title": bson.M{"$regex": "^" + args.Keyword}}
	}

	query := s.coll.Find(cond)

	if err != nil {
		return
	}
	books := make([]models.Book, 100)
	err = query.Skip(args.Start * 100).Limit(100).All(&books)
	if err != nil {
		return
	}
	data["books"] = books
	return
}

//=================================================

type MoreArgs struct {
	catalog string
	tag     string
	Start   int
}

func (s *Server) More(args MoreArgs) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}
	query := s.coll.Find(nil)

	if err != nil {
		return
	}

	s.coll.Find(query)

	books := make([]models.Book, 100)
	err = query.Skip(args.Start * 100).Limit(100).All(&books)
	if err != nil {
		return
	}
	data["books"] = books
	return
}

//=================================================

type SearchArgs struct {
	Keyword string
	Start   int
}

func (s *Server) Search(args SearchArgs) (data map[string]interface{}, err error) {
	data = map[string]interface{}{}
	query := &mgo.Query{}
	if args.Keyword == "" {
		query = s.coll.Find(nil)
	} else {
		cond := bson.M{"title": bson.M{"$regex": "^" + args.Keyword}}
		query = s.coll.Find(cond)
	}
	books := make([]models.Book, 100)
	err = query.Skip(args.Start * 100).Limit(100).All(&books)
	if err != nil {
		return
	}
	data["books"] = books
	return
}
