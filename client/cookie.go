package client

import (
	"a21hc3NpZ25tZW50/config"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

func GetClientWithCookie(token string, cookies ...*http.Cookie) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	cookies = append(cookies, &http.Cookie{
		Name:  "session_token",
		Value: token,
	})

	jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   config.BaseURL,
	}, cookies)

	c := &http.Client{
		Jar: jar,
	}

	return c, nil
}
