package main

import (
	"github.com/icattlecoder/jieshu/www/controllers"
	"github.com/icattlecoder/tgw"
	"labix.org/v2/mgo"
	"log"
)

func main() {

	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Fatal("mgo.Dial Err:", err)
	}
	defer session.Close()

	c := session.DB("jieshu")
	coll := c.C("book")
	coll_user := c.C("user")
	ser := controllers.NewServer(coll, coll_user)

	_tgw := tgw.NewTGW()

	store := tgw.NewMemcachedSessionStore("127.0.0.1:11211")
	log.Fatal(_tgw.SetSessionStore(store).AddParser(ser.UserMgr).Register(&ser).Run(":8080"))

}
