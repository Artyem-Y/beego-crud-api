package services

import (
	"beego-crud-api/conf"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
	"time"
)

type AccessDetails struct {
	UserId int64
	Role   string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func CreateToken(userID int64, role string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Hour * 24 * 30).Unix()
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	var err error

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires
	atClaims["role"] = role
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	td.AccessToken, err = at.SignedString([]byte(conf.GetEnvConst("ACCESS_SECRET")))

	if err != nil {
		return nil, err
	}
	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["user_id"] = userID
	rtClaims["exp"] = td.RtExpires
	rtClaims["role"] = role
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(conf.GetEnvConst("REFRESH_SECRET")))

	if err != nil {
		return nil, err
	}
	return td, nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")

	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// Parse, validate, and return a token.
// keyFunc will receive the parsed token and should return the key for validating.
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
