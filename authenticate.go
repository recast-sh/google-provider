package google

import (
	"os"

	"recast.sh/v0/cli"
	"recast.sh/v0/provider/google/impl"
)

var (
	credsPath string = ""
)

func init() {
	cli.StringFlag(&credsPath, "google-auth", "", "", "Google compute creds (JWT json format)")
}

var auth *impl.AuthConfig

func Authenticate() {
	if credsPath == "" {
		credsPath = os.Getenv("RECAST_GOOGLE_AUTH")
	}

	var err error
	auth, err = impl.JWTAuth(credsPath)
	if err != nil {
		panic(err)
	}
}
