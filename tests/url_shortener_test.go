package tests

import (
	"net/http"
	"net/url"
	"path"
	"testing"
	"url-shortener/internal/http-server/handlers/save"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/random"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	host = "localhost:8080"
)

func TestUrlShortener_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	e.POST("/url").WithJSON(save.Request{
		URL:   gofakeit.URL(),
		Alias: random.NewRandomString(10),
	}).
		WithBasicAuth("testuser", "testpass").
		Expect().
		Status(http.StatusOK).
		JSON().Object().
		ContainsKey("alias")
	// ContainsKey("url")
}

func TestUrlShortener_SaveRedirectDelete(t *testing.T) {
	cases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid test",
			url:   gofakeit.URL(),
			alias: random.NewRandomString(7),
		},
		{
			name:  "Invalid URL",
			url:   gofakeit.Word(),
			alias: gofakeit.Phone(),
			error: "field URL is not a valid URL",
		},
		{
			name:  "Invalid URL with spaces",
			url:   "    " + gofakeit.Word(),
			alias: gofakeit.Phone(),
			error: "field URL is not a valid URL",
		},
		{
			name:  "Empty alias",
			url:   gofakeit.URL(),
			alias: "",
		},
		{
			name:  "All is empty",
			error: "field URL is a required field",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Save

			// TODO: replace to another struct (they repeat many time)
			resp := e.POST("/url").WithJSON(save.Request{
				URL:   tc.url,
				Alias: tc.alias,
			}).WithBasicAuth("test_user", "test_password").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			} else {
				resp.ContainsKey("alias")
			}

			aliasFromResp := tc.alias

			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				resp.Value("alias").String().NotEmpty()

				aliasFromResp = resp.Value("alias").String().Raw()
			}

			// Redirect

			testRedirect(t, aliasFromResp, tc.url)

			// Delete

			reqDel := e.DELETE("/"+path.Join("url", aliasFromResp)).
				WithBasicAuth("testuser", "testpass").
				Expect().Status(http.StatusOK).
				JSON().Object()

			reqDel.Value("status").String().IsEqual("OK")

			// Redirect again

			testRedirectNotFound(t, aliasFromResp)

		})
	}

}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToUrl, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	assert.Equal(t, urlToRedirect, redirectedToUrl)

}

func testRedirectNotFound(t *testing.T, alias string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	_, err := api.GetRedirect(u.String())

	require.EqualError(t, err, "api.GetRedirect: invalid status code: 404")
	// assert.ErrorIs(t, err, api.ErrInvalidStatusCode)
}
