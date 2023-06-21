package client

import (
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
		// uncomment this line if you want to use the production server
		Host:   "final-production-2bfa.up.railway.app",
		// Host: "localhost:8080",
	}, cookies)

	c := &http.Client{
		Jar: jar,
	}

	return c, nil
}
