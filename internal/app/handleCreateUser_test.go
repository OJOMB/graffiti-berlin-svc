package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleCreateUser_successPath(t *testing.T) {
	ms := &mockService{}
	app := New(nil, nullLogger(), nil, "", false, ms)

	u := domain.NewUser("9abc46be-3bcd-42b1-aeb2-ac6ff557a580", "test", "test@example.com", "ac6ff557a5804eff")
	ms.On("CreateUser", mock.Anything, "test", "test@example.com", "ac6ff557a5804eff").Return(u, nil)

	w := httptest.NewRecorder()
	reqBody := `{"user_name":"test", "email":"test@example.com", "password":"ac6ff557a5804eff"}`
	r := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))

	app.handleCreateUser()(w, r)

	assert.Equal(t, http.StatusCreated, w.Code)

	expectedRespBody, err := json.Marshal(u)
	assert.NoError(t, err)
	assert.Equal(t, expectedRespBody, w.Body.Bytes())

	ms.AssertExpectations(t)
}

func TestHandleCreateUser_requestBodyContainsInvalidJSON_failurePath(t *testing.T) {
	app := New(nil, nullLogger(), nil, "", false, nil)

	w := httptest.NewRecorder()
	reqBody := `{"user_name":"test", "email":"test@example.com", "password":"ac6ff557a5804eff"`
	r := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))

	app.handleCreateUser()(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"error": "invalid json in request body"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestHandleCreateUser_requestBodyUnreadable_failurePath(t *testing.T) {
	app := New(nil, nullLogger(), nil, "", false, nil)

	w := httptest.NewRecorder()
	var reqBody errReader = 0
	r := httptest.NewRequest(http.MethodPost, "/users", reqBody)

	app.handleCreateUser()(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"error": "request body unreadable"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestHandleCreateUser_serviceErr_failurePath(t *testing.T) {
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
			app := New(nil, nullLogger(), nil, "", false, ms)

			ms.On("CreateUser", mock.Anything, "test", "test@example.com", "ac6ff557a5804eff").
				Return(nil, &domain.Error{Code: tc.domainErrType, Msg: "test error"})

			w := httptest.NewRecorder()
			reqBody := `{"user_name":"test", "email":"test@example.com", "password":"ac6ff557a5804eff"}`
			r := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(reqBody))

			app.handleCreateUser()(w, r)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, strings.TrimRight(w.Body.String(), "\n"))

			ms.AssertExpectations(t)
		})
	}
}
