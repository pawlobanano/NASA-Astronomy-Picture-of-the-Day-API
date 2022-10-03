package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListPicturesURL(t *testing.T) {
	type Query struct {
		From string
		To   string
	}

	testCases := []struct {
		name          string
		URLQuery      Query
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			URLQuery: Query{
				From: "2019-12-05",
				To:   "2019-12-07",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Contains(t, []int{200, 429}, recorder.Code)
			},
		},
		{
			name: "InvalidFrom",
			URLQuery: Query{
				From: "2019-11-31",
				To:   "2019-12-07",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidTo",
			URLQuery: Query{
				From: "2019-12-05",
				To:   "2019-13-01",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingFrom",
			URLQuery: Query{
				From: "",
				To:   "2019-12-01",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingTo",
			URLQuery: Query{
				From: "",
				To:   "2019-12-07",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "FromLaterThanTo",
			URLQuery: Query{
				From: "2019-12-08",
				To:   "2019-12-07",
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			server := newTestServer(t)
			recorder := httptest.NewRecorder()

			URL := "/pictures"
			request, err := http.NewRequest(http.MethodGet, URL, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("from", tc.URLQuery.From)
			q.Add("to", tc.URLQuery.To)
			request.URL.RawQuery = q.Encode()

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}
