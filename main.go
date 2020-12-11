package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var (
	Header     string
	AuthScheme string
	Algorithm  string
	SigningKey string
	Addr       string
	Path       string
)

func main() {
	SigningKey = os.Getenv("JWT_SIGNING_KEY")
	if SigningKey == "" {
		log.Fatal("empty jwt signing key")
	}

	Header = os.Getenv("JWT_HEADER")
	if Header == "" {
		Header = "Authorization"
	}

	AuthScheme = os.Getenv("JWT_AUTH_SCHEME")
	if AuthScheme == "" {
		AuthScheme = "Bearer"
	}

	Algorithm = os.Getenv("JWT_ALGORITHM")
	if Algorithm == "" {
		Algorithm = "HS256"
	}

	Addr = os.Getenv("SERVER_ADDR")
	if Addr == "" {
		Addr = ":8080"
	}

	Path = os.Getenv("SERVER_PATH")
	if Path == "" {
		Path = "/auth"
	}

	http.HandleFunc(Path, verifyJWT)
	log.Fatal(http.ListenAndServe(Addr, nil))
}

func verifyJWT(w http.ResponseWriter, req *http.Request) {
	auth, err := jwtFromHeader(req)
	if err != nil {
		response(w, http.StatusBadRequest, err.Error())
		return
	}

	token := new(jwt.Token)
	token, err = jwt.Parse(auth, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != Algorithm {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
		}
		return []byte(SigningKey), nil
	})

	if err != nil {
		response(w, http.StatusBadRequest, err.Error())
		return
	}

	if !token.Valid {
		response(w, http.StatusBadRequest, "invalid or expired jwt")
		return
	}

	filtersParam := req.URL.Query().Get("filters")
	if filtersParam == "" {
		// no filters, allow all
		response(w, http.StatusOK, "ok")
		return
	}

	filters := strings.Split(filtersParam, "|")
	if len(filters) > 0 {
		if mapClaims, ok := token.Claims.(jwt.MapClaims); ok {
			for _, filter := range filters {
				pair := strings.Split(filter, ":")
				if len(pair) != 2 {
					continue
				}

				value, ok := mapClaims[pair[0]].(string)
				if !ok {
					response(w, http.StatusUnauthorized, fmt.Sprintf("'%s' claim is not available", pair[0]))
					return
				}

				if value != pair[1] {
					response(w, http.StatusUnauthorized, fmt.Sprintf("'%s' must be equal to '%s' receive '%s'", pair[0], pair[1], value))
					return
				}
			}
		}
	}

	response(w, http.StatusOK, "ok")
}

func jwtFromHeader(req *http.Request) (string, error) {
	auth := req.Header.Get(Header)
	l := len(AuthScheme)
	if len(auth) > l+1 && auth[:l] == AuthScheme {
		return auth[l+1:], nil
	}
	return "", errors.New("missing or malformed jwt")
}

func response(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	_, _ = w.Write([]byte(message))
}
