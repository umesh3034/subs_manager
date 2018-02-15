package main

import (
	"fmt"
	"time"
	"log"
	"net/http"
	"encoding/json"
//	"net/url"
        "github.com/kirves/goradius"
	"github.com/gorilla/mux"
//	"io/ioutil"
)

type User struct{
    Username      string
    Password      string
}

var (
        server   string = "localhost"
        port     string = "1812"
        secret   string = "testing123"
        timeout  time.Duration = 5
        //user     string = "umesh"
        //password string = "redback"
        nasId    string = "localhost"
)

func Authenticate(user string, password string) int {
        fmt.Println("In Authenticate function!!!--------")
        auth := goradius.AuthenticatorWithTimeout(server, port, secret, timeout)
        //auth := goradius.Authenticator(server, port, secret)
        res, err := auth.Authenticate(user, password, nasId)
        if err != nil {
                fmt.Println(err)
                return -1
        }
        if !res {
                fmt.Println("Could not authenticate.")
                return -1
        }

        fmt.Println("Authenticate Success !!!!!!!!")
        return 0
}


/*
func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: homePage")

        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        log.Println("r.Body", string(body))

        var u User
        err = json.Unmarshal(body, &u)
        fmt.Println("User: ", u)
        values, err := url.ParseQuery(string(body))
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        fmt.Println(values)
        fmt.Println("indices from body", values["username"])
        res := Authenticate(username, password)
        if res != 0 {
	    fmt.Fprintf(w, "AUTH FAILED!")
        } else {
	    fmt.Fprintf(w, "AUTH SUCCESS!")
        }
}
*/



func homePage(w http.ResponseWriter, r *http.Request) {
    var u User
    if r.Body == nil {
        http.Error(w, "Please send a request body", 400)
        return
    }
    err := json.NewDecoder(r.Body).Decode(&u)
    if err != nil {
        fmt.Println("error ")
        http.Error(w, err.Error(), 400)
        return
    }
    fmt.Println("username: password: ", u.Username, u.Password)
    res := Authenticate(u.Username, u.Password)
    if res != 0 {
        fmt.Fprintf(w, "AUTH FAILED!")
    } else {
        fmt.Fprintf(w, "AUTH SUCCESS!")
    }
}


func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":5000", myRouter))
}

func main() {
	handleRequests()
}
