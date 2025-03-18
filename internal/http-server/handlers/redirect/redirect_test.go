package redirect_test

import (
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/redirect/mocks"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://google.com",
		},
		{
			name:  "Empty alias",
			alias: "",
			url:   "https://amazon.com",
			respError: "invalid request",
		},
		{
			name:      "Empty URL",
			url:       "",
			respError: "empty url",
		},
		{
			name:      "Invalid format of URL",
			url:       "example.com",
			respError: "invalid format of URL",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.On("GetURL", tc.alias).
					Return(tc.url, tc.mockError).Once()
			}
			
			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToUrl, err := api.GetRedirect(ts.URL + "/" + tc.alias)


			if tc.respError != "" {
				assert.Error(t, err)
				// FIX: fix test
				assert.Contains(t, err.Error(), tc.respError)
				return
			}
			
			require.NoError(t, err)

			assert.Equal(t, tc.url, redirectedToUrl)
			
		})
	}

}
