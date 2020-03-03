package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gitlab.com/horo-go/auth0"
	"gitlab.com/horo-go/horo"
	"gitlab.com/lattetalk/lattetalk/auth"
	"gitlab.com/lattetalk/lattetalk/config"
	"gitlab.com/lattetalk/lattetalk/model"
	"gitlab.com/lattetalk/lattetalk/service"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

func Authorize(c *horo.Context) {
	var data struct {
		AccessToken string `json:"accessToken"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.BadRequest()
		return
	}

	claims, err := auth.Auth0.GetClaims(data.AccessToken, auth0.RS256)
	if err != nil {
		c.BadRequest()
		return
	}

	id := service.GetUserIDByAuthID(claims.Subject)
	if id == -1 {
		user := getAuth0UserInfo(data.AccessToken)
		id = service.CreateUser(user)
		if id == -1 {
			id = service.GetUserIDByAuthID(claims.Subject)
		}
	}

	newClaims := jwt.Claims{
		Subject: strconv.FormatInt(id, 10),
		Expiry:  claims.Expiry,
	}

	p, err := json.Marshal(newClaims)
	if err != nil {
		panic(err)
	}

	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.HS256,
		Key:       []byte(config.Config.Auth0.Secret),
	}, nil)
	if err != nil {
		panic(err)
	}

	signature, err := signer.Sign(p)
	if err != nil {
		panic(err)
	}

	token, err := signature.CompactSerialize()
	if err != nil {
		panic(err)
	}

	cookie := http.Cookie{
		Name:     auth.TokenCookie,
		Value:    token,
		Expires:  time.Unix(int64(claims.Expiry), 0),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(c.Writer, &cookie)
	c.WriteHeader(http.StatusNoContent)
}

func LogOut(c *horo.Context) {
	cookie := http.Cookie{
		Name:     auth.TokenCookie,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(c.Writer, &cookie)
	c.WriteHeader(http.StatusNoContent)
}

func getAuth0UserInfo(accessToken string) *model.User {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	req, err := http.NewRequest("GET", config.Config.Auth0.UserInfo, nil)
	req.Header.Add("Authorization", "Bearer "+accessToken)
	req = req.WithContext(ctx)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	var data struct {
		Sub            string
		GivenName      string `json:"given_name"`
		FamilyName     string `json:"family_name"`
		ProfilePicture string `json:"picture"`
	}

	if err = json.NewDecoder(res.Body).Decode(&data); err != nil {
		panic(err)
	}

	user := new(model.User)
	user.AuthID = data.Sub
	user.GivenName = data.GivenName
	user.FamilyName = data.FamilyName
	user.ProfilePicture = data.ProfilePicture

	return user
}
