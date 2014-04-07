package models

import (
	"errors"
	"fmt"
	"github.com/icattlecoder/jieshu/www/api"
	"github.com/icattlecoder/mcClient"
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
	Users mcClient.MC
	// Users map[int64]UserInfo
}

func NewUserMgr(coll *mgo.Collection, mcHosts []string) *UserMgr {

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
	return &UserMgr{coll: coll, Users: mcClient.NewGobMCClient("user", mcHosts...)}
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

func (u *UserMgr) GetInOut(uids []int64) (result []interface{}, err error) {
	result = make([]interface{}, len(uids))
	selector := bson.M{"avatar": 1, "uid": 1, "location": 1, "name": 1}
	err = u.coll.Find(bson.M{"uid": bson.M{"$in": uids}}).Select(selector).All(&result)
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
func (u *UserMgr) AddDouBan(data map[string]interface{}) (user UserInfo, err error) {

	if uid, ok := data["uid"]; ok {
		if iv, ok := uid.(string); ok {
			iuid, err2 := strconv.Atoi(iv)
			if err2 != nil {
				err = err2
				return
			}
			user.Uid = int64(iuid)

			user2, err = u.Get(user.Uid)
			if err == nil {
				user = user2
				err = errors.New("User Exsit")
				return
			}

		}
	} else {
		err = errors.New("Invalid DouBan Users")
		return
	}

	if name, ok := data["name"]; ok {
		if iv, ok := name.(string); ok {
			user.Name = iv
		}
	}

	if avatar, ok := data["avatar"]; ok {
		if iv, ok := avatar.(string); ok {
			user.Avatar = iv
		}
	}

	if loc_name, ok := data["loc_name"]; ok {
		if iv, ok := loc_name.(string); ok {
			user.Location = iv
		}
	}
	err = u.Add(user)
	return
}

func getUidStr(uid int64) string { return fmt.Sprintf("%d", uid) }

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
		user := UserInfo{}
		if err := u.Users.Get(getUidStr(uid), &user); err == nil {
			user.Email = email
			u.Users.Set(getUidStr(uid), user)
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
	go api.SendEmail("【私人借书网 | 用户注册】", strconv.Itoa(int(user.Uid)), "icattlecoder@gmail.com")
	if err == nil {
		u.Users.Set(getUidStr(user.Uid), user)
	}
	return
}

func (u *UserMgr) Get(uid int64) (user UserInfo, err error) {
	err = u.Users.Get(getUidStr(uid), &user)
	if err != nil {
		//以下是数据库操作:
		query := u.coll.Find(bson.M{"uid": uid})
		err = query.One(&user)
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
	return
}
