package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type TokenManager struct{}

func NewTokenService() *TokenManager {
	return &TokenManager{}
}

type TokenInterface interface {
	CreateToken(userId, userName string) (*TokenDetails, error)
	ExtractTokenMetadata(*http.Request) (*AccessDetails, error)
}

var _ TokenInterface = &TokenManager{}

func (t *TokenManager) CreateToken(userId, userName string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 30).Unix()
	td.TokenUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = td.TokenUuid + "++" + userId

	var err error

	// Creating token
	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = td.TokenUuid
	atClaims["user_id"] = userId
	atClaims["user_name"] = userName
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// Creating refresh token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userId
	rtClaims["user_name"] = userName
	rtClaims["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, nil
}

func (t *TokenManager) ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)

	if err != nil {
		return nil, err
	}

	acc, err := Extract(token)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}

	return nil
}

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

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func Extract(token *jwt.Token) (*AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		userId, userOk := claims["user_id"].(string)
		userName, userNameOk := claims["user_name"].(string)
		if !ok || !userOk || !userNameOk {
			return nil, errors.New("unauthorized")
		} else {
			return &AccessDetails{
				TokenUuid: accessUuid,
				UserId:    userId,
				UserName:  userName,
			}, nil
		}

	}

	return nil, errors.New("error trying to extract")
}
