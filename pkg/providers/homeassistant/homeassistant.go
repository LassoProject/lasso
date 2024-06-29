/*

Copyright 2020 The Vouch Proxy Authors.
Use of this source code is governed by The MIT License (MIT) that
can be found in the LICENSE file. Software distributed under The
MIT License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES
OR CONDITIONS OF ANY KIND, either express or implied.

*/

package homeassistant

import (
	"encoding/json"
	"golang.org/x/oauth2"
	"io/ioutil"
	"net/http"

	"github.com/vouch/vouch-proxy/pkg/cfg"
	"github.com/vouch/vouch-proxy/pkg/providers/common"
	"github.com/vouch/vouch-proxy/pkg/structs"
	"go.uber.org/zap"
)

// Provider provider specific functions
type Provider struct{}

var log *zap.SugaredLogger

// Configure see main.go configure()
func (Provider) Configure() {
	log = cfg.Logging.Logger
}

// GetUserInfo provider specific call to get userinfomation
// More info: https://developers.home-assistant.io/docs/en/auth_api.html
func (Provider) GetUserInfo(r *http.Request, user *structs.User, customClaims *structs.CustomClaims, ptokens *structs.PTokens, opts ...oauth2.AuthCodeOption) (rerr error) {
	client, providerToken, err := common.PrepareTokensAndClient(r, ptokens, false, opts...)
	if err != nil {
		return err
	}
	ptokens.PAccessToken = providerToken.Extra("access_token").(string)
	userinfo, err := client.Get(cfg.GenOAuth.UserInfoURL)
	if err != nil {
		return err
	}
	defer func() {
		if err := userinfo.Body.Close(); err != nil {
			rerr = err
		}
	}()
	data, _ := ioutil.ReadAll(userinfo.Body)
	log.Infof("HA userinfo body: %s", string(data))
	haUser := structs.User{}
	if err = json.Unmarshal(data, &haUser); err != nil {
		log.Error(err)
		return err
	}
	user.Username = haUser.Username
	return nil
}
