package update

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/update/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"
	urlservice "url-shortener/internal/service/url"

	"github.com/stretchr/testify/require"
)

func TestUpdate_Success(t *testing.T) {
	urlUpdaterMock := mocks.NewURLUpdater(t)

	urlUpdaterMock.
		On("UpdateURL", "test_alias", "https://yandex.com").
		Return(int64(1), nil).
		Once()

	handler := New(
		slogdiscard.NewDiscardLogger(),
		urlUpdaterMock,
	)

	input := `{
		"url": "https://yandex.com",
		"alias": "test_alias"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/url",
		bytes.NewBufferString(input),
	)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusOK, rr.Code)

	var resp Response
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, int64(1), resp.CountUpdated)
	urlUpdaterMock.AssertExpectations(t)
}

func TestUpdate_BadRequest(t *testing.T) {
	urlUpdaterMock := mocks.NewURLUpdater(t)

	handler := New(
		slogdiscard.NewDiscardLogger(),
		urlUpdaterMock,
	)

	input := `{
		"url": "https://yandex.com",
		"alias": "test_alias"
	`

	req := httptest.NewRequest(
		http.MethodPut,
		"/url",
		bytes.NewBufferString(input),
	)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusBadRequest, rr.Code)

	var resp Response
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, "failed to decode request body", resp.Error)

	urlUpdaterMock.AssertExpectations(t)
}

func TestUpdate_NotFound(t *testing.T) {
	urlUpdaterMock := mocks.NewURLUpdater(t)
	urlUpdaterMock.
		On("UpdateURL", "test_alias", "https://yandex.com").
		Return(int64(0), urlservice.ErrURLNotFound).
		Once()

	handler := New(
		slogdiscard.NewDiscardLogger(),
		urlUpdaterMock,
	)

	input := `{
		"url": "https://yandex.com",
		"alias": "test_alias"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/url",
		bytes.NewBufferString(input),
	)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	require.Equal(t, http.StatusNotFound, rr.Code)

	var resp Response
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
	require.Equal(t, "url not found", resp.Error)
	urlUpdaterMock.AssertExpectations(t)

}
