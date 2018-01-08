package gtok

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	assertionExpirationDelay       = 10 * time.Minute
	gcpAccessRequestUrl            = "https://www.googleapis.com/oauth2/v4/token"
	gcpAccessRequestGrantTypeKey   = "grant_type"
	gcpAccessRequestGrantTypeValue = "urn:ietf:params:oauth:grant-type:jwt-bearer"
	gcpAccessRequestAssertionKey   = "assertion"
)

type accessToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	IssuedAt    time.Time
	ExpiresAt   time.Time
}
type gcpServiceAccountCredential struct {
	ProjectId               string `json:"project_id"`
	PrivateKeyId            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientId                string `json:"client_id"`
	AuthUri                 string `json:"auth_uri"`
	TokenUri                string `json:"token_uri"`
	AuthProviderX509CertUrl string `json:"auth_provider_x509_cert_url"`
	ClientX509CertUrl       string `json:"client_x509_cert_url"`

	key *rsa.PrivateKey
}
type GcpServiceAccount struct {
	credential        gcpServiceAccountCredential
	assertion         string
	assertedAt        time.Time
	assertionExprieAt time.Time
	token             accessToken
}

func NewGcpServiceAccountFromJson(j string) (GcpServiceAccount, error) {
	var err error
	credential := gcpServiceAccountCredential{}
	json.Unmarshal([]byte(j), &credential)
	credential.key, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(credential.PrivateKey))
	if err != nil {
		return GcpServiceAccount{}, err
	}
	gsa := GcpServiceAccount{}
	gsa.credential = credential
	return gsa, nil
}
func NewGcpServiceAccountFromFile(filename string) (GcpServiceAccount, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return GcpServiceAccount{}, err
	}
	return NewGcpServiceAccountFromJson(string(bytes))
}
func (gsa GcpServiceAccount) AccessToken() (string, error) {
	if (accessToken{}) == gsa.token || gsa.TokenExpired() {
		// request token
		assertion, err := gsa.Assertion()
		if err != nil {
			return "", err
		}
		resp, err := http.PostForm(gcpAccessRequestUrl,
			url.Values{
				gcpAccessRequestGrantTypeKey: {gcpAccessRequestGrantTypeValue},
				gcpAccessRequestAssertionKey: {assertion},
			})
		if err != nil {
			return "", err
		}

		// parse token from the response
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		respTime, err := time.Parse(time.RFC1123, resp.Header.Get("date"))
		if err != nil {
			return "", err
		}
		token := accessToken{}
		json.Unmarshal(bytes, &token)
		token.IssuedAt = respTime
		token.ExpiresAt = respTime.Add(time.Duration(token.ExpiresIn) * time.Second)
		gsa.token = token
	}
	return fmt.Sprintf("%s %s", gsa.token.TokenType, gsa.token.AccessToken), nil
}
func (gsa GcpServiceAccount) Assertion() (string, error) {
	if "" == gsa.assertion || time.Now().After(gsa.assertionExprieAt) {
		iat := time.Now()
		exp := iat.Add(assertionExpirationDelay)
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"iss":   "spropd@ccinpp-190005.iam.gserviceaccount.com",
			"scope": "https://www.googleapis.com/auth/devstorage.read_write",
			"aud":   "https://www.googleapis.com/oauth2/v4/token",
			"iat":   iat.Unix(),
			"exp":   exp.Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		assertionString, err := token.SignedString(gsa.credential.key)
		if err != nil {
			return assertionString, err
		}
		gsa.assertion = assertionString
		gsa.assertedAt = iat
		gsa.assertionExprieAt = exp
	}
	return gsa.assertion, nil
}
func (gsa GcpServiceAccount) TokenExpired() bool {
  return time.Now().After(gsa.token.ExpiresAt)
}
