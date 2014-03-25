package main

import (
	"github.com/icattlecoder/jieshu/www/api"
	"encoding/json"
	"os"
	"flag"
	"github.com/icattlecoder/jieshu/www/controllers"
	"github.com/icattlecoder/tgw"
	"labix.org/v2/mgo"
	"log"
)

/*
config File Demo:
{
    "port": "8080",
    "douban_apikey": "",
    "doubanSecret": "",
    "mgo_host": "localhost",
    "mc_host": ["127.0.0.1:11211"],
    "img_service": "http://127.0.0.1:8001/getimage?key=",
    "email_service": ""
}
*/
type Config struct {
	Port         string `json:"port"`
	DoubanApiKey string `json:"douban_apikey"`
	DoubanSecret string `json:"doubanSecret"`
	MgoHost      string `json:"mgo_host"`
	MCHost     []string `json:"mc_host"`
	ImageServer  string `json:"img_service"`
	EmailServer  string `json:"email_service"`
}

func loadConfig(path string) Config {
	r, err := os.Open(path)
	if err != nil {
		log.Fatal("load Config", path, err)
	}
	decoder := json.NewDecoder(r)
	conf := Config{}
	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatal("load Config", path, err)
	}
	return conf
}


func main() {

	cfgPath := flag.String("f", "conf.json", "usage: sitename -f <conf>")
	flag.Parse()

	cfg := loadConfig(*cfgPath)

	session, err := mgo.Dial(cfg.MgoHost)
	if err != nil {
		log.Fatal("mgo.Dial Err:", err)
	}
	defer session.Close()
	c := session.DB("jieshu")
	coll := c.C("book")
	coll_user := c.C("user")
	index := mgo.Index{
		Key:        []string{"uid"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	}
	coll_user.EnsureIndex(index)

	scfg := controllers.Config{
		DoubanApiKey: cfg.DoubanApiKey,
		DoubanSecret: cfg.DoubanSecret,
		ImageServer: cfg.ImageServer,
	}
	api.EmailServerHost = cfg.EmailServer

	ser := controllers.NewServer(coll, coll_user, &scfg)

	_tgw := tgw.NewTGW()

	store := tgw.NewMemcachedSessionStore(cfg.MCHost...)
	log.Fatal(_tgw.SetSessionStore(store).AddParser(ser.UserMgr).Register(&ser).Run(":"+cfg.Port))

}
