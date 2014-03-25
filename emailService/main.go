package main

import (
	"encoding/json"
	"flag"
	"github.com/icattlecoder/tgw"
	"log"
	"net/smtp"
	"os"
	"time"
)

type Config struct {
	SmtpHost     string `json:"smtp_host"`
	Port         string `json:"port"`
	SmtpUserName string `json:"smtp_username"`
	SmtpPassword string `json:"smtp_password"`
}

type Server struct {
	Config
	auth   smtp.Auth
	buffer chan *SendArgs
}

func readConfig(filename string, val interface{}) (err error) {
	r, err := os.Open(filename)
	if err != nil {
		return
	}
	defer r.Close()
	decoder := json.NewDecoder(r)
	err = decoder.Decode(val)
	return
}

func NewServer(cfg Config) *Server {
	auth := smtp.PlainAuth("", cfg.SmtpUserName, cfg.SmtpPassword, cfg.SmtpHost)
	buffer := make(chan *SendArgs, 1000)
	return &Server{Config: cfg, auth: auth, buffer: buffer}
}

func (s *Server) run() {
	for {
		args := <-s.buffer
		s.send(*args)
	}
}

func (s *Server) send(args SendArgs) {

	body := "To: " + args.To + "\r\n" +
		"From: " + s.SmtpUserName + "\r\n" +
		"Subject: " + args.Sub + "\r\n" +
		"Date: " + time.Now().String() + "\r\n\r\n" +
		args.Body
	log.Println("Send Email ...")
	log.Println("To:", args.To)
	log.Println("Subject:", args.Sub)
	now := time.Now()
	smtp.SendMail(s.SmtpHost+":25", s.auth, s.SmtpUserName, []string{args.To}, []byte(body))
	log.Println("Taken(s)", time.Now().Sub(now).Seconds())
	log.Println("Send!")
}

type SendArgs struct {
	To   string
	Sub  string
	Body string
}

func (s *Server) Send(args SendArgs) {
	if args.To != "" {
		s.buffer <- &args
	}
}

var config = flag.String("c", "value", `emailSever -c <config> 
	config:
	{
		"port":"8090",
		"smtp_host":"smtp.126.com",
		"smtp_username":"wmshfu@hotmail",
		"smtp_password":"***********"
	}
`)

func main() {
	flag.Parse()
	tg := tgw.NewTGW()
	cfg := Config{}
	err := readConfig(*config, &cfg)
	if err != nil {
		log.Fatal("readConfig error:", config, err)
	}
	svr := NewServer(cfg)
	go svr.run()
	log.Fatal(tg.Register(&svr).Run(":"+cfg.Port))
}
