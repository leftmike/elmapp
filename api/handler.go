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
	}
)

func validationFailed(w http.ResponseWriter, msg string, err error) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	if err != nil {
		msg = fmt.Sprintf("%s: %s", msg, err)
	}
	fmt.Fprintf(w, `{"errors":{"body":["%s"]}}`, msg)
}

func fileHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	n, ok := staticFiles[req.URL.Path]
	if !ok {
		http.NotFound(w, req)
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
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if req.URL.Path != loginPath {
		http.NotFound(w, req)
		return
	}

	var reqLogin reqLogin
	err := json.NewDecoder(req.Body).Decode(&reqLogin)
	if err != nil {
		validationFailed(w, "login json decode", err)
		return
	}
	defer req.Body.Close()

	if reqLogin.User == nil {
		validationFailed(w, "login json decode: missing user", nil)
		return
	}
	user := model.LoginUser(reqLogin.User.Email, reqLogin.User.Password)
	if user == nil {
		validationFailed(w, "login: bad email or password", nil)
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
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if req.URL.Path != registerPath {
		http.NotFound(w, req)
		return
	}

	var reqRegister reqRegister
	err := json.NewDecoder(req.Body).Decode(&reqRegister)
	if err != nil {
		validationFailed(w, "register json decode", err)
		return
	}
	defer req.Body.Close()

	if reqRegister.User == nil {
		validationFailed(w, "register json decode: missing user", nil)
		return
	}
	user, err := model.RegisterUser(reqRegister.User.Username, reqRegister.User.Email,
		reqRegister.User.Password)
	if err != nil {
		validationFailed(w, "register", err)
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
