package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/domain"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleGetUser_successPath(t *testing.T) {
	ms := &mockService{}
	app := New(nil, nullLogger(), nil, "", "", nil, ms)

	u := domain.NewUser("9abc46be-3bcd-42b1-aeb2-ac6ff557a580", "test", "test@example.com", "ac6ff557a5804eff")
	ms.On("GetUser", mock.Anything, "9abc46be-3bcd-42b1-aeb2-ac6ff557a580").Return(u, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/users/9abc46be-3bcd-42b1-aeb2-ac6ff557a580", nil)
	r = mux.SetURLVars(r, map[string]string{"userID": "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"})

	app.handleGetUser()(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	expectedRespBody, err := json.Marshal(u)
	assert.NoError(t, err)
	assert.Equal(t, expectedRespBody, w.Body.Bytes())

	ms.AssertExpectations(t)
}

func TestHandleGetUser_serviceErr_failurePath(t *testing.T) {
	testCases := []struct {
		name           string
		domainErrType  domain.ErrorType
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "service returns InvalidInput error",
			domainErrType:  domain.InvalidInput,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "invalid input data - test error"}`,
		},
		{
			name:           "service returns system error",
			domainErrType:  domain.SystemError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error": "unexpected system error - test error"}`,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx, tc.name), func(t *testing.T) {
			ms := &mockService{}
			ms.On("GetUser", mock.Anything, "9abc46be-3bcd-42b1-aeb2-ac6ff557a580").
				Return(nil, &domain.Error{Code: tc.domainErrType, Msg: "test error"})

			app := New(nil, nullLogger(), nil, "", "", nil, ms)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/api/v1/users/9abc46be-3bcd-42b1-aeb2-ac6ff557a580", nil)
			r = mux.SetURLVars(r, map[string]string{"userID": "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"})

			app.handleGetUser()(w, r)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, strings.TrimRight(w.Body.String(), "\n"))

			ms.AssertExpectations(t)
		})
	}
}

func TestHandleGetUser_userNotFound_failurePath(t *testing.T) {
	ms := &mockService{}
	ms.On("GetUser", mock.Anything, "9abc46be-3bcd-42b1-aeb2-ac6ff557a580").
		Return(nil, nil)

	app := New(nil, nullLogger(), nil, "", "", nil, ms)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/users/9abc46be-3bcd-42b1-aeb2-ac6ff557a580", nil)
	r = mux.SetURLVars(r, map[string]string{"userID": "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"})

	app.handleGetUser()(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, `{"error": "user not found"}`, strings.TrimRight(w.Body.String(), "\n"))

	ms.AssertExpectations(t)
}
