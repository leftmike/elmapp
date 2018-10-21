package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/leftmike/elmapp/model"
)

const (
	loginPath    = "/api/users/login"
	registerPath = "/api/users"
	userPath     = "/api/user"
)

var (
	staticFiles = map[string]string{
		"/":           "index.html",
		"/index.html": "index.html",
		"/elm.js":     "elm.js",
	}
)

func validationFailed(w http.ResponseWriter, msgs ...string) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	fmt.Fprint(w, `{"errors":{"body":[`)
	for idx, msg := range msgs {
		if idx > 0 {
			fmt.Fprint(w, ", ")
		}
		fmt.Fprintf(w, "%q", msg)
	}
	fmt.Fprint(w, "]}}")
}

func fileHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("file: %s\n", req.URL.Path)

	n, ok := staticFiles[req.URL.Path]
	if !ok {
		http.NotFound(w, req)
		return
	}

	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	f, err := os.Open(n)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: os.Open(%q): %s", n, err)
		return
	}
	defer f.Close()

	_, err = io.Copy(w, f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: io.Copy(%q): %s", n, err)
		return
	}
}

type reqLoginUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type reqLogin struct {
	User *reqLoginUser `json:"user"`
}

type respUser struct {
	User *model.User `json:"user"`
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != loginPath {
		http.NotFound(w, req)
		return
	}
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var reqLogin reqLogin
	err := json.NewDecoder(req.Body).Decode(&reqLogin)
	if err != nil {
		validationFailed(w, err.Error())
		return
	}
	defer req.Body.Close()

	if reqLogin.User == nil {
		validationFailed(w, "json is missing user field")
		return
	}
	user := model.LoginUser(reqLogin.User.Email, reqLogin.User.Password)
	if user == nil {
		validationFailed(w, "email or password is invalid")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respUser{User: user})
}

type reqRegisterUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type reqRegister struct {
	User *reqRegisterUser `json:"user"`
}

func registerHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != registerPath {
		http.NotFound(w, req)
		return
	}
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var reqRegister reqRegister
	err := json.NewDecoder(req.Body).Decode(&reqRegister)
	if err != nil {
		validationFailed(w, err.Error())
		return
	}
	defer req.Body.Close()

	if reqRegister.User == nil {
		validationFailed(w, "json is missing user field")
		return
	}
	user, msgs := model.RegisterUser(reqRegister.User.Username, reqRegister.User.Email,
		reqRegister.User.Password)
	if msgs != nil {
		validationFailed(w, msgs...)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respUser{User: user})
}

type checkToken struct {
	handlerFunc func(w http.ResponseWriter, req *http.Request, user *model.User)
}

func (ct checkToken) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
	if len(s) != 2 || s[0] != "Token" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	user, err := model.ValidateToken(s[1])
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	ct.handlerFunc(w, req, user)
}

func userHandler(w http.ResponseWriter, req *http.Request, user *model.User) {
	if req.URL.Path != userPath {
		http.NotFound(w, req)
		return
	}
	if req.Method == http.MethodPost {
		w.WriteHeader(http.StatusNotImplemented) // XXX
	} else if req.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(respUser{User: user})
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func init() {
	http.HandleFunc("/", fileHandler)
	http.HandleFunc(loginPath, loginHandler)
	http.HandleFunc(registerPath, registerHandler)
	http.Handle(userPath, checkToken{userHandler})
}
