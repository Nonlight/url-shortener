package delete

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	urlservice "url-shortener/internal/service/url"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestDelete_Success(t *testing.T) {
	urlDeleterMock := mocks.NewURLDeleter(t)

	urlDeleterMock.
		On("DeleteURL", "test_alias").
		Return(int64(1), nil).
		Once()

	handler := New(
		slogdiscard.NewDiscardLogger(),
		urlDeleterMock,
	)
	router := chi.NewRouter()
	router.Delete("/{alias}", handler)

	req := httptest.NewRequest("DELETE", "/test_alias", nil)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var resp Response
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, int64(1), resp.CountDeleted)

	urlDeleterMock.AssertExpectations(t)
}

func TestDelete_NotFound(t *testing.T) {
	urlDeleterMock := mocks.NewURLDeleter(t)
	urlDeleterMock.
		On("DeleteURL", "test_alias").
		Return(int64(0), urlservice.ErrURLNotFound).
		Once()

	handler := New(
		slogdiscard.NewDiscardLogger(),
		urlDeleterMock,
	)
	router := chi.NewRouter()
	router.Delete("/{alias}", handler)

	req := httptest.NewRequest(http.MethodDelete, "/test_alias", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNotFound, rr.Code)

	var resp Response
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, "url not found", resp.Error)
	urlDeleterMock.AssertExpectations(t)

}
