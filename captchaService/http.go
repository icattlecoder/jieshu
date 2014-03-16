// CServer project main.go
package main

import (
	"flag"
	// "fmt"
	"image/png"
	"log"
	"net/http"

	"github.com/hanguofeng/gocaptcha"
)

var (
	ccaptcha *gocaptcha.Captcha
)

var configFile = flag.String("c", "gocaptcha.conf", "the config file")

func ShowImageHandler(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	if len(key) >= 0 {
		cimg := ccaptcha.GenImage(key)
		w.Header().Add("Content-Type", "image/png")
		png.Encode(w, cimg)
	}
}

func main() {
	captcha, err := gocaptcha.CreateCaptchaFromConfigFile(*configFile)

	if nil != err {
		panic(err.Error())
	} else {
		ccaptcha = captcha
	}

	http.HandleFunc("/getimage", ShowImageHandler)

	s := &http.Server{Addr: ":8001"}
	log.Fatal(s.ListenAndServe())
}
