package models

import (
	"errors"
	"github.com/icattlecoder/tgw"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"reflect"
	"strconv"
)

var (
	session_userid = "userInfo"
)

type UserInfo struct {
	Uid      int64
	Password string
	Email    string
	Name     string
	In       []string
	Out      []string
	Location string
	Avatar   string
}

type UserMgr struct {
	coll  *mgo.Collection
	Users map[int64]UserInfo
}

func NewUserMgr(coll *mgo.Collection) *UserMgr {

	cache := map[int64]UserInfo{}
	iter := coll.Find(nil).Iter()
	for {
		user := UserInfo{}
		if iter.Next(&user) {
			cache[user.Uid] = user
		} else {
			break
		}
	}
	return &UserMgr{coll: coll, Users: cache}
}

func (u *UserMgr) InOut(user UserInfo, book_id string, typ string) (err error) {
	n, err := u.coll.Find(bson.M{"uid": user.Uid, typ: book_id}).Count()
	if err != nil {
		return
	}
	if n > 0 {
		err = errors.New("Invalid op")
		return
	}
	return u.coll.Update(bson.M{"email": user.Email}, bson.M{"$push": bson.M{typ: book_id}})
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
func (u *UserMgr) AddDouBan(data map[string]string) (user UserInfo, err error) {

	if name, ok := data["name"]; ok {
		user.Name = name
	}

	if uid, ok := data["uid"]; ok {
		iuid, err2 := strconv.Atoi(uid)
		if err2 != nil {
			err = err2
			return
		}
		user.Uid = int64(iuid)
		if user2, ok := u.Users[user.Uid]; ok {
			user = user2
			err = errors.New("User Exsit")
			return
		}
	}

	if avatar, ok := data["avatar"]; ok {
		user.Avatar = avatar
	}

	if loc_name, ok := data["loc_name"]; ok {
		user.Location = loc_name
	}
	err = u.Add(user)
	return
}

func (u *UserMgr) UpdateEmail(uid int64, email string) (err error) {
	n, err := u.coll.Find(bson.M{"uid": uid}).Count()
	if err != nil {
		return err
	}
	if n <= 0 {
		err = errors.New("尚未注册")
	}
	err = u.coll.Update(bson.M{"uid": uid}, bson.M{"$set": bson.M{"email": email}})
	//更新缓存
	if err == nil {
		if user, ok := u.Users[uid]; ok {
			user.Email = email
			u.Users[uid] = user
		}
	}
	return
}

func (u *UserMgr) Add(user UserInfo) (err error) {
	n, err := u.coll.Find(bson.M{"uid": user.Uid}).Count()
	if err != nil {
		return err
	}
	if n > 0 {
		err = errors.New("Email :" + user.Email + " 已被注册!")
		return
	}
	err = u.coll.Insert(user)
	if err == nil {
		u.Users[user.Uid] = user
	}
	return
}

func (u *UserMgr) Get(uid int64) (user UserInfo, err error) {
	user, ok := u.Users[uid]
	if !ok {
		err = errors.New("User Not Exsit!")
	}
	//以下是数据库操作:
	// query := u.coll.Find(bson.M{"uid": uid})
	// err = query.One(&user)
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
	return
}
