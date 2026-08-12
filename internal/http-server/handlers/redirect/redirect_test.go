package redirect

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/redirect/mocks"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	urlservice "url-shortener/internal/service/url"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestRedirect_Success(t *testing.T) {
	urlGetterMock := mocks.NewURLGetter(t)

	urlGetterMock.
		On("GetURL", "test_alias").
		Return("https://www.google.com/", nil).
		Once()

	r := chi.NewRouter()
	r.Get("/{alias}", New(
		slogdiscard.NewDiscardLogger(),
		urlGetterMock,
	))

	ts := httptest.NewServer(r)
	defer ts.Close()

	redirectedToURL, err := api.GetRedirect(ts.URL + "/test_alias")
	require.NoError(t, err)

	require.Equal(
		t,
		"https://www.google.com/",
		redirectedToURL,
	)

	urlGetterMock.AssertExpectations(t)
}

func TestRedirect_NotFound(t *testing.T) {
	urlGetterMock := mocks.NewURLGetter(t)

	urlGetterMock.
		On("GetURL", "test_alias").
		Return("", urlservice.ErrURLNotFound).
		Once()
	router := chi.NewRouter()
	router.Get("/{alias}", New(
		slogdiscard.NewDiscardLogger(),
		urlGetterMock,
	))
	req := httptest.NewRequest(
		http.MethodGet,
		"/test_alias",
		nil,
	)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNotFound, rr.Code)
	urlGetterMock.AssertExpectations(t)
}
