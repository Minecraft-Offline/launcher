package main

import (
	"github.com/JoshuaDoes/go-yggdrasil"
	"github.com/JoshuaDoes/json"

	//std necessities
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Auth struct {
	Username    string `json:"username"`
	ID          string `json:"id"`
	AccessToken string `json:"accessToken"`
	DecodeToken struct {
		YGGT string `json:"yggt"`
	} `json:"-"`

	ClientToken string            `json:"-"`
	yggdrasil   *yggdrasil.Client `json:"-"`

	Email    string `json:"-"`
	Password string `json:"-"`

	hasToken bool `json:"-"`
}

func (auth *Auth) LoadToken(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	token, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(token, auth)
	if err != nil {
		return err
	}

	auth.hasToken = true
	return nil
}

func (auth *Auth) SaveToken(path string) error {
	token, err := json.Marshal(auth, false)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, token, 0644)
	return err
}

func (auth *Auth) DecodeJWT() {
	//Credit to InvoxiPlayGames for helping me understand how to do this
	if auth.AccessToken != "" {
		tokenSplit := strings.Split(auth.AccessToken, ".")
		if len(tokenSplit) != 3 {
			log.Error("Unable to split access token")
		}

		tokenJSON, err := base64.StdEncoding.DecodeString(tokenSplit[1] + "=")
		if err != nil {
			log.Error("Unable to base64 decode access token: ", err)
		}

		err = json.Unmarshal(tokenJSON, &auth.DecodeToken)
		if err != nil {
			log.Error("Unable to unmarshal decoded access token to struct: ", err)
		}
	}
}

func (auth *Auth) Login() error {
	auth.yggdrasil = &yggdrasil.Client{
		AccessToken: auth.AccessToken,
		ClientToken: clientID,
	}

	defer auth.DecodeJWT()

	var yggErr *yggdrasil.Error
	var yggErrors []*yggdrasil.Error

	if auth.hasToken {
		_, yggErr = auth.yggdrasil.Validate()
		if yggErr != nil {
			log.Error(yggErr)
			yggErrors = append(yggErrors, yggErr)
		} else {
			return nil
		}

		_, yggErr = auth.yggdrasil.Refresh()
		if yggErr != nil {
			log.Error(yggErr)
			yggErrors = append(yggErrors, yggErr)
		} else {
			auth.AccessToken = auth.yggdrasil.AccessToken

			return nil
		}
	}

	authResponse, yggErr := auth.yggdrasil.Authenticate(auth.Email, auth.Password, "Minecraft", 1)
	if yggErr != nil {
		log.Error(yggErr)
		yggErrors = append(yggErrors, yggErr)
	} else {
		auth.AccessToken = authResponse.AccessToken
		auth.Username = authResponse.SelectedProfile.Name
		auth.ID = authResponse.SelectedProfile.ID

		return nil
	}

	return fmt.Errorf("%v", yggErrors)
}
