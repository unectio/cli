package main

import (
	"os"
	"log"
	"time"
	"github.com/unectio/util"
	"github.com/unectio/api/uauth"
	"encoding/base64"
	"github.com/dgrijalva/jwt-go"
)

type Config struct {
	Apilet	string		`yaml:"apilet"`
	KeyId	string		`yaml:"keyid"`
	Key	string		`yaml:"key"`
}

type Login struct {
	address	string
	token	string
}

func config() string {
	cfg := os.Getenv("LETS_CONFIG")
	if cfg == "" {
		cfg = defaultConfig
	}

	return cfg
}

func getLogin() (*Login, error) {
	if *dryrun {
		return &Login{}, nil
	}

	var conf Config

	err := util.LoadYAML(config(), &conf)
	if err != nil {
		return nil, err
	}

	return getSelfSignedLogin(&conf)
}

func getSelfSignedLogin(conf *Config) (*Login, error) {
	key, err :=  base64.StdEncoding.DecodeString(conf.Key)
	if err != nil {
		return nil, err
	}

	claims := &jwt.StandardClaims {
		ExpiresAt: time.Now().Add(uauth.SelfKeyLifetime).Unix(),
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok.Header["kid"] = conf.KeyId

	if *debug {
		log.Printf("===[ new request, key: (%02x%02x%02x%02x...%02x/%d) ]===\n",
				key[0], key[1], key[2], key[3], key[len(key)-1], len(key))
	}

	toks, err := tok.SignedString(key)
	if err != nil {
		return nil, err
	}

	return &Login{ address: conf.Apilet, token: toks }, nil
}
