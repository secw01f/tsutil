package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"golang.org/x/oauth2/clientcredentials"
)

func GetSMSecret(smarn string) string {
	type SecretValue struct {
		Key string
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatal(err)
	}

	client := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(smarn),
	}

	output, err := client.GetSecretValue(context.TODO(), input)
	if err != nil {
		log.Fatal(err)
	}

	var keyValue SecretValue

	jerr := json.Unmarshal([]byte(*output.SecretString), &keyValue)
	if jerr != nil {
		log.Fatal(jerr)
	}

	secret := keyValue.Key

	return secret
}

func GetTsAuth(id string, secret string, tags string) string {
	type AuthToken struct {
		ID           string
		Key          string
		Created      string
		Expires      string
		Revoked      string
		Capabilities struct {
			Devices struct {
				Create struct {
					Reusable      bool
					Ephemeral     bool
					Preauthorized bool
					Tags          []string
				}
			}
		}
		Description string
	}

	cfg := &clientcredentials.Config{
		ClientID:     id,
		ClientSecret: secret,
		TokenURL:     "https://api.tailscale.com/api/v2/oauth/token",
	}

	client := cfg.Client(context.Background())

	var data = []byte(fmt.Sprintf(`{"capabilities": {"devices": {"create": {"reusable": false, "ephemeral": false, "preauthorized": true, "tags": ["%s"]}}}}`, strings.Join(strings.Split(tags, ","), `","`)))

	req, err := http.NewRequest("POST", "https://api.tailscale.com/api/v2/tailnet/-/keys", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, rerr := client.Do(req)
	if rerr != nil {
		log.Fatal(rerr)
	}

	body, berr := ioutil.ReadAll(resp.Body)
	if berr != nil {
		log.Fatal(berr)
	}

	var keyData AuthToken

	jerr := json.Unmarshal(body, &keyData)
	if jerr != nil {
		log.Fatal(jerr)
	}

	return keyData.Key
}

func GetTsApi(id string, secret string) string {
	cfg := &clientcredentials.Config{
		ClientID:     id,
		ClientSecret: secret,
		TokenURL:     "https://api.tailscale.com/api/v2/oauth/token",
	}

	client, err := cfg.Token(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return client.AccessToken
}
