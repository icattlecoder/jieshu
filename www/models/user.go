package models

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"github.com/icattlecoder/tgw"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"reflect"
)

var (
	session_userid = "userInfo"
)

type UserInfo struct {
	Uid      int64
	Password string
	Email    string
	In []string
	out	 []string
}

type UserMgr struct {
	coll *mgo.Collection
}

func NewUserMgr(coll *mgo.Collection) *UserMgr {
	return &UserMgr{coll: coll}
}

func crypto(password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (u *UserMgr) InOut(user UserInfo,book_id string,typ string) error {
	log.Println(".................................")
	log.Println(user.Email,book_id,typ)
	return u.coll.Update(bson.M{"email":user.Email},bson.M{"$push":bson.M{typ:book_id}})
}

func (u *UserMgr) Add(user UserInfo) (err error) {
	n, err := u.coll.Find(bson.M{"email": user.Email}).Count()
	if err != nil {
		return err
	}
	if n > 0 {
		err = errors.New("Email :" + user.Email + " 已被注册!")
		return
	}
	crypedPasswd := crypto(user.Password)
	user.Password = crypedPasswd
	err = u.coll.Insert(user)
	return
}

func (u *UserMgr) Valid(email, password string) (err error) {

	user := UserInfo{}
	err = u.coll.Find(bson.M{"email": email}).One(&user)
	if err != nil || user.Email == "" {
		log.Println(err)
		err = errors.New(email + " 不存在!")
		return
	}
	if user.Password != crypto(password) {
		err = errors.New("密码不正确!")
	}
	return
}

func (u *UserMgr) Parse(env *tgw.ReqEnv, typ reflect.Type) (val reflect.Value, parsed bool) {

	if typ.Elem().Name() != "UserInfo" {
		return
	}
	parsed = true
	user := UserInfo{}
	err := env.Session.Get(session_userid, &user)
	if err != nil {
		val = reflect.ValueOf((*UserInfo)(nil))
		return
	}
	val = reflect.ValueOf(&user)

	log.Println("++++++++", user)

	// val = reflect.ValueOf(&v)
	return
}
