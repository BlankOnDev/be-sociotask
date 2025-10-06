package auth

import (
	"context"
	"log"

	"github.com/harundarat/be-socialtask/internal/store"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
)

func NewGoogleAuth() *oauth2.Config {
	oauthConfGl := &oauth2.Config{
		ClientID:     store.GetEnv("Google_Client_ID_Web"),
		ClientSecret: store.GetEnv("Google_Client_Secret_Web"),
		RedirectURL:  store.GetEnv("Google_Redirect_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", "https://www.googleapis.com/auth/userinfo.email", "openid"},
		Endpoint:     google.Endpoint,
	}

	return oauthConfGl
}

func GoogleVerifytokenID(tokenID string) (*idtoken.Payload, error) {
	// maaf kode ini jelek, kapan-kapan ku oplas ini.
	idinfo, err := idtoken.Validate(context.Background(), tokenID, store.GetEnv("Google_Client_ID_Web"))
	if err != nil {
		log.Printf(`error saat memverifikasi tokenID: %v`, err)
		return nil, err
	}

	return idinfo, nil
}
