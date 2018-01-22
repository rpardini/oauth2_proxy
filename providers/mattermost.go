package providers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/rpardini/oauth2_proxy/api"
)

type MattermostProvider struct {
	*ProviderData
}

func NewMattermostProvider(p *ProviderData) *MattermostProvider {
	p.ProviderName = "Mattermost"
	if p.LoginURL.String() == "" {
		p.LoginURL = &url.URL{Scheme: "https",
			Host: "mm.misy.me",
			Path: "/oauth/authorize",
			// ?granted_scopes=true
		}
	}
	if p.RedeemURL.String() == "" {
		p.RedeemURL = &url.URL{Scheme: "https",
			Host: "mm.misy.me",
			Path: "/oauth/access_token",
		}
	}
	if p.ProfileURL.String() == "" {
		p.ProfileURL = &url.URL{Scheme: "https",
			Host: "mm.misy.me",
			Path: "/api/v4/users/me",
		}
	}
	if p.ValidateURL.String() == "" {
		p.ValidateURL = p.ProfileURL
	}
	if p.Scope == "" {
		p.Scope = "public_profile email"
	}
	return &MattermostProvider{ProviderData: p}
}

func getMattermostHeader(access_token string) http.Header {
	header := make(http.Header)
	header.Set("Accept", "application/json")
	//header.Set("x-li-format", "json")
	header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))
	return header
}

func (p *MattermostProvider) GetEmailAddress(s *SessionState) (string, error) {
	if s.AccessToken == "" {
		return "", errors.New("missing access token")
	}
	req, err := http.NewRequest("GET", p.ProfileURL.String()+"", nil)
	if err != nil {
		return "", err
	}
	req.Header = getMattermostHeader(s.AccessToken)

	type result struct {
		Email string
	}
	var r result
	err = api.RequestJson(req, &r)
	if err != nil {
		return "", err
	}
	if r.Email == "" {
		return "", errors.New("no email")
	}
	return r.Email, nil
}
func (p *MattermostProvider) GetUserName(s *SessionState) (string, error) {
	if s.AccessToken == "" {
		return "", errors.New("missing access token")
	}
	req, err := http.NewRequest("GET", p.ProfileURL.String()+"", nil)
	if err != nil {
		return "", err
	}
	req.Header = getMattermostHeader(s.AccessToken)

	type result struct {
		Username string
	}
	var r result
	err = api.RequestJson(req, &r)
	if err != nil {
		return "", err
	}
	if r.Username == "" {
		return "", errors.New("no username")
	}
	return r.Username, nil
}

func (p *MattermostProvider) ValidateSessionState(s *SessionState) bool {
	return validateToken(p, s.AccessToken, getMattermostHeader(s.AccessToken))
}
