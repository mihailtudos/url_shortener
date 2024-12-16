package tests

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/gavv/httpexpect/v2"
	"github.com/mihailtudos/url-shortener/internal/http-server/handlers/url/save"
	"github.com/mihailtudos/url-shortener/internal/lib/api"
	rand "github.com/mihailtudos/url-shortener/internal/lib/random"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8080"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())
	alias, err := rand.GenerateRandomString(6)
	require.NoError(t, err)

	e.POST("/url").
		WithJSON(save.Request{
			URL:   "https://google.com",
			Alias: alias,
		}).
		WithBasicAuth("myuser", "mypass").
		Expect().
		Status(200).
		JSON().
		Object().
		ContainsKey("alias")
}

func TestURLShortener_SaveRedirect(t *testing.T) {
	cases := []struct {
		name string
		url string
		alias string
		error string
	} {
		{
			name: "Valud URL",
			url: gofakeit.URL(),
			alias: gofakeit.Word(),
		},
		{
			name: "Invalid URL",
			url: "invalid_url",
			alias: gofakeit.Word(),
			error: "field URL is not a valid URL",
		},
		{
			name: "Empty alias",
			url: gofakeit.URL(),
			alias: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("myuser", "mypass").
				Expect().
				Status(http.StatusOK).
				JSON().
				Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")
				resp.Value("error").String().IsEqual(tc.error)
				return
			}

			alias := tc.alias
			if  tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				alias = resp.Value("alias").String().Raw()
			}

			testRedirect(t, alias, tc.url)

			// Remove alias from the database
			// send request using httpexpect
			// implement testRedirectNotFiund func checker
		})
	}
}

func testRedirect(t *testing.T, alias, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path: alias,
	}

	requirestToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, requirestToURL, urlToRedirect)
}