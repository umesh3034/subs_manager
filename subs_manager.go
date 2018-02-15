package main

import (
	"time"
        "github.com/kirves/goradius"
        "fmt"
	"flag"
	"log"
	"runtime"
	"github.com/nats-io/go-nats"
        "strings"
)

var (
	server   string = "localhost"
	port     string = "1812"
	secret   string = "testing123"
	timeout  time.Duration = 5
	user     string = "test"
	password string = "test"
	nasId    string = "localhost"
)

func Authenticate(username string, pass string) int {
	auth := goradius.AuthenticatorWithTimeout(server, port, secret, timeout)
	res, err := auth.Authenticate(username, pass, nasId)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	if err == nil {
		fmt.Println("AUTH sucess for ", "username:", username, " password:", pass, "!!!--------")
		return 0
	}
	if !res {
		fmt.Println("Could not authenticate.")
		return -1
	}

	return -1
}

// NOTE: Use tls scheme for TLS, e.g. nats-rply -s tls://demo.nats.io:4443 foo hello
func usage() {
	log.Fatalf("Usage: nats-rply [-s server][-t] <subject> <reponse>\n")
}

func printMsg(m *nats.Msg, i int) {
	log.Printf("[#%d] Received on [%s]: '%s'\n", i, m.Subject, string(m.Data))
}

func main() {
	var urls = flag.String("s", nats.DefaultURL, "The nats server URLs (separated by comma)")
	var showTime = flag.Bool("t", false, "Display timestamps")

	log.SetFlags(0)
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	nc, err := nats.Connect(*urls)
	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
	}

	subj, i := args[0], 0

	nc.Subscribe(subj, func(msg *nats.Msg) {
		i++
		printMsg(msg, i)
                input := strings.Split(string(msg.Data), ":")
                fmt.Println("username:", input[0], " password:", input[1])
                res := Authenticate(input[0], input[1])
		if res != 0 {
			reply := "NOK"
			nc.Publish("RESP", []byte(reply))
		} else {
			reply := "OK"
			nc.Publish("RESP", []byte(reply))
		}
		
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Listening on [%s]\n", subj)
	if *showTime {
		log.SetFlags(log.LstdFlags)
	}

	runtime.Goexit()
}


// func TestPacketCreation(t *testing.T) {
// 	auth := Authenticator(server, port, secret)
// 	v := auth.generateAuthenticator()
// 	encpass, _ := auth.radcrypt(v, []byte(password))
// 	pkg := auth.createRequest(v, []byte(user), encpass)
// }
