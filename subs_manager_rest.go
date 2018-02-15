package main

import (
	"time"
        "github.com/kirves/goradius"
        "fmt"
        "bytes"
	"flag"
	//"os"
	"log"
	"runtime"
	"encoding/json"
	"github.com/nats-io/go-nats"
        "net/http"
        "strings"
        "io/ioutil"
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

type User struct{
    Username      string `json:"username"`
    Password      string `json:"password"`
}

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

/*func makeHttpPostReq(url string) int {

    client := http.Client{}

    var jsonStr = {"username=umesh&password=redback"}

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    resp, err := client.Do(req)
    if err != nil {
         fmt.Println("Unable to reach the server.")
         return -1
    } else {
         body, _ := ioutil.ReadAll(resp.Body)
         fmt.Println("body=", string(body))
         return 0
    }

}
*/

func makeHttpJSONReq(url string) int {
	u := User{Username: "umesh", Password: "redback"}
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(u)
        fmt.Println("Username send: ", u.Username)
        fmt.Println("String: ", b.String())
	resp, err := http.Post(url, "application/json; charset=utf-8", b)
	//resp, err := http.Get(url)
        if err != nil {
            fmt.Println("Unable to reach the server.")
            return -1
        } else {
            body, _ := ioutil.ReadAll(resp.Body)
            fmt.Println("body=", string(body))
            check := strings.Contains(string(body), "SUCCESS")
            if (check != true) {
                return -1
            } else {
                return 0
            }
        }
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
                err := makeHttpJSONReq("http://localhost:5000")
/*
                err := makeHttpJSONReq("http://10.183.255.241:4000/users")
                err := makeHttpJSONReq("http://10.183.255.241:5000")
                resp, err := http.Get("http://10.183.255.241:5000")
                defer resp.Body.Close()
                body, err := ioutil.ReadAll(resp.Body)
                str := string(body)
                fmt.Println("Response:", str)
*/
                if err != 0 {
	            reply := "NOK"
                    fmt.Println("Endpoint Error")
	            nc.Publish("RESP", []byte(reply))
                } else {
                    reply := "OK"
                    fmt.Println("Endpoint Success")
		    nc.Publish("RESP", []byte(reply))
                }
                fmt.Println("Send Response")
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
