/////////////////////////////////////////////////////////////////////////////////
//
// Copyright (C) 2019-2020, Unectio Inc, All Right Reserved.
//
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
/////////////////////////////////////////////////////////////////////////////////

package main

import (
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
	"github.com/unectio/api/uauth"
	"github.com/unectio/util"
	"log"
	"os"
	"time"
)

type Config struct {
	Apilet    string `yaml:"apilet"`
	KeyId     string `yaml:"keyid"`
	Key       string `yaml:"key"`
	ProjectId string `yaml:"projectid"`
}

type Claims struct {
	ProjectId string `json:"projectid"`
	jwt.StandardClaims
}

type Login struct {
	address string
	token   string
}

func config() string {
	cfg := Cfg
	if cfg == "" {
		cfg = os.Getenv(defaultCfgEnv)
		if cfg == "" {
			cfg = defaultConfig
		}
	}

	return cfg
}

func getLogin() (*Login, error) {
	if DryRun {
		return &Login{}, nil
	}

	var conf Config

	err := util.LoadYAML(config(), &conf)
	if err != nil {
		return nil, err
	}

	if Project != "" {
		conf.ProjectId = Project
	}

	return getSelfSignedLogin(&conf)
}

func getSelfSignedLogin(conf *Config) (*Login, error) {
	key, err := base64.StdEncoding.DecodeString(conf.Key)
	if err != nil {
		return nil, err
	}

	//	claims := &jwt.StandardClaims{
	//		ExpiresAt: time.Now().Add(uauth.SelfKeyLifetime).Unix(),
	//	}

	claims := &Claims{
		ProjectId: conf.ProjectId,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: time.Now().Add(uauth.SelfKeyLifetime).Unix(),
		},
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok.Header["kid"] = conf.KeyId

	if Verbose {
		log.Printf("===[ new request, key: (%02x%02x%02x%02x...%02x/%d) ]===\n",
			key[0], key[1], key[2], key[3], key[len(key)-1], len(key))
	}

	toks, err := tok.SignedString(key)
	if err != nil {
		return nil, err
	}

	return &Login{address: conf.Apilet, token: toks}, nil
}
