package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

type AccessDetails struct {
	TokenUuid string
	UserId    string
	UserName  string
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type AuthInterface interface {
	CreateAuth(string, *TokenDetails) error
	FetchAuth(string) (string, error)
	DeleteRefresh(string) error
	DeleteTokens(*AccessDetails) error
}

type RedisAuthService struct {
	client *redis.Client
}

var _ AuthInterface = &RedisAuthService{}

func (tk *RedisAuthService) CreateAuth(userId string, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atCreated, err := tk.client.Set(td.TokenUuid, userId, at.Sub(now)).Result()

	if err != nil {
		return err
	}

	rtCreated, err := tk.client.Set(td.RefreshUuid, userId, rt.Sub(now)).Result()

	if err != nil {
		return err
	}

	if atCreated == "0" || rtCreated == "0" {
		return errors.New("no records inserted")
	}

	return nil

}

func (tk *RedisAuthService) FetchAuth(tokenUuid string) (string, error) {
	userId, err := tk.client.Get(tokenUuid).Result()
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (tk *RedisAuthService) DeleteRefresh(refreshUuid string) error {
	deleted, err := tk.client.Del(refreshUuid).Result()
	if err != nil || deleted == 0 {
		return err
	}

	return nil
}

func (tk *RedisAuthService) DeleteTokens(ad *AccessDetails) error {
	refreshUuid := fmt.Sprintf("%s++%s", ad.TokenUuid, ad.UserId)

	deletedAt, err := tk.client.Del(ad.TokenUuid).Result()
	if err != nil || deletedAt == 0 {
		return err
	}

	deletedRt, err := tk.client.Del(refreshUuid).Result()
	if err != nil || deletedRt == 0 {
		return err
	}

	return nil
}
