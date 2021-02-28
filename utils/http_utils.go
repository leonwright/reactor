package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/golang/gddo/httputil/header"
	"github.com/leonwright/reactor/logger"

	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
)

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type CustomClaims struct {
	Audience string `json:"aud,omitempty"`
	// Subject  string `jsonL"sub,omitempty"`
	// ExpiresAt int64  `json:"exp,omitempty"`
	// Id        string `json:"jti,omitempty"`
	// IssuedAt  int64  `json:"iat,omitempty"`
	// Issuer    string `json:"iss,omitempty"`
	// NotBefore int64  `json:"nbf,omitempty"`
	// Subject   string `json:"sub,omitempty"`
	Scope string `json:"scope,omitempty`
	jwt.StandardClaims
}

func ValidateRequestIsJSON(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			http.Error(w, msg, http.StatusUnsupportedMediaType)
			return errors.New("request is not json")
		}
		return nil
	}
	return errors.New("missing header")
}

func SendJSONResponse(resp interface{}, w http.ResponseWriter) {
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// SendResponseWithoutData
func SendResponseWithoutData(message string, status int, w http.ResponseWriter) {
	var resp Response = Response{Message: message}
	jsonResp, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonResp)
}

func CreateMiddleWare() *jwtmiddleware.JWTMiddleware {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// Verify 'aud' claim
			aud := "reactorcore.nerderbur.tech"
			checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAud {
				return token, errors.New("Invalid audience.")
			}
			// Verify 'iss' claim
			iss := "https://reactorapp.us.auth0.com/"
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("Invalid issuer.")
			}

			cert, err := getPemCert(token)
			if err != nil {
				panic(err.Error())
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})
	return jwtMiddleware
}

// main.go

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get("https://reactorapp.us.auth0.com/.well-known/jwks.json")

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k, _ := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("Unable to find appropriate key.")
		return cert, err
	}

	return cert, nil
}

func CheckScope(scope string, tokenString string) (bool, string) {
	deb.Infof("Enter method CheckScope with params %s %s", scope, logger.TruncateString(tokenString, 40))
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})
	if err != nil {
		fmt.Println(err)
	}

	claims, ok := token.Claims.(*CustomClaims)

	hasScope := false
	if ok && token.Valid {
		result := strings.Split(claims.Scope, " ")
		for i := range result {
			if result[i] == scope {
				hasScope = true
			}
		}
	}

	deb.Infof("Exit method CheckScope with result Has Scope (%s) %t, Username: %s", scope, hasScope, claims.Subject)
	return hasScope, claims.Subject
}

func ExtractToken(r *http.Request) string {
	deb.Info("Extracting token from header...")
	return strings.Split(r.Header.Get("Authorization"), " ")[1]
}

func GetUserNameFromRequest(tokenString string) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		cert, err := getPemCert(token)
		if err != nil {
			return nil, err
		}
		result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
		return result, nil
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(claims["sub"])
}
