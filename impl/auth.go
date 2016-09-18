package impl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2/jwt"
)

// Google's OAuth 2.0 token URL to use with the JWT flow.
const jwtTokenURL = "https://accounts.google.com/o/oauth2/token"

type AuthConfig struct {
	ProjectID string
	jwt.Config
}

func JWTAuth(file string) (*AuthConfig, error) {
	if file == "" {
		return nil, fmt.Errorf("Missing creds file")
	}
	jsonKey, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var key struct {
		Email        string `json:"client_email"`
		PrivateKey   string `json:"private_key"`
		PrivateKeyID string `json:"private_key_id"`
		ProjectID    string `json:"project_id"`
	}
	if err := json.Unmarshal(jsonKey, &key); err != nil {
		return nil, err
	}
	return &AuthConfig{
		ProjectID: key.ProjectID,
		Config: jwt.Config{
			Email:        key.Email,
			PrivateKey:   []byte(key.PrivateKey),
			PrivateKeyID: key.PrivateKeyID,
			TokenURL:     jwtTokenURL,
		},
	}, nil
}

var (
	CloudPlatformScope         = "https://www.googleapis.com/auth/cloud-platform"
	CloudPlatformReadOnlyScope = "https://www.googleapis.com/auth/cloud-platform.read-only"

	NdevClouddnsReadonlyScope  = "https://www.googleapis.com/auth/ndev.clouddns.readonly"
	NdevClouddnsReadwriteScope = "https://www.googleapis.com/auth/ndev.clouddns.readwrite"

	ComputeScope         = "https://www.googleapis.com/auth/compute"
	ComputeReadonlyScope = "https://www.googleapis.com/auth/compute.readonly"

	DevstorageFullControlScope = "https://www.googleapis.com/auth/devstorage.full_control"
	DevstorageReadOnlyScope    = "https://www.googleapis.com/auth/devstorage.read_only"
	DevstorageReadWriteScope   = "https://www.googleapis.com/auth/devstorage.read_write"
)

func (c *AuthConfig) Requires(scopes ...string) {
	c.Config.Scopes = append(c.Config.Scopes, scopes...)
}
