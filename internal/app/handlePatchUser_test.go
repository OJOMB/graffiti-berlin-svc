package app

import (
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

func TestHandlePatchUser_successPath(t *testing.T) {
	ms := &mockService{}
	app := New(nil, nullLogger(), nil, "", "", nil, ms)

	patchJSON := `[
		{ "op": "replace", "path": "/user_name", "value": "foo" },
		{ "op": "replace", "path": "/email", "value": "bar@example.com" }
	]`

	ms.On("PatchUser", mock.Anything, "9abc46be-3bcd-42b1-aeb2-ac6ff557a580", []byte(patchJSON)).Return(nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/users/9abc46be-3bcd-42b1-aeb2-ac6ff557a580", strings.NewReader(patchJSON))
	r = mux.SetURLVars(r, map[string]string{"userID": "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"})

	app.handlePatchUser()(w, r)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "", w.Body.String())

	ms.AssertExpectations(t)
}

func TestHandlePatchUser_requestBodyUnreadable_successPath(t *testing.T) {
	app := New(nil, nullLogger(), nil, "", "", nil, nil)

	var unreadableBody errReader = 0
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/users/9abc46be-3bcd-42b1-aeb2-ac6ff557a580", unreadableBody)
	r = mux.SetURLVars(r, map[string]string{"userID": "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"})

	app.handlePatchUser()(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"error": "request body unreadable"}`, strings.TrimRight(w.Body.String(), "\n"))
}

func TestHandlePatchUser_serviceErr_failurePath(t *testing.T) {
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
			app := New(nil, nullLogger(), nil, "", "", nil, ms)

			patchJSON := `[
				{ "op": "replace", "path": "/user_name", "value": "foo" },
				{ "op": "replace", "path": "/email", "value": "bar@example.com" }
			]`

			ms.On("PatchUser", mock.Anything, "9abc46be-3bcd-42b1-aeb2-ac6ff557a580", []byte(patchJSON)).
				Return(&domain.Error{Code: tc.domainErrType, Msg: "test error"})

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPatch, "/users/9abc46be-3bcd-42b1-aeb2-ac6ff557a580", strings.NewReader(patchJSON))
			r = mux.SetURLVars(r, map[string]string{"userID": "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"})

			app.handlePatchUser()(w, r)

			assert.Equal(t, tc.expectedStatus, w.Code)
			assert.Equal(t, tc.expectedBody, strings.TrimRight(w.Body.String(), "\n"))

			ms.AssertExpectations(t)
		})
	}
}

func TestHandlePatchUser_userNotFound_failurePath(t *testing.T) {
	ms := &mockService{}
	app := New(nil, nullLogger(), nil, "", "", nil, ms)

	patchJSON := `[
		{ "op": "replace", "path": "/user_name", "value": "foo" },
		{ "op": "replace", "path": "/email", "value": "bar@example.com" }
	]`
	ms.On("PatchUser", mock.Anything, "9abc46be-3bcd-42b1-aeb2-ac6ff557a580", []byte(patchJSON)).
		Return(&domain.Error{Code: domain.ResourceNotFound, Msg: "user does not exist"})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/users/9abc46be-3bcd-42b1-aeb2-ac6ff557a580", strings.NewReader(patchJSON))
	r = mux.SetURLVars(r, map[string]string{"userID": "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"})

	app.handlePatchUser()(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, `{"error": "resource not found - user does not exist"}`, strings.TrimRight(w.Body.String(), "\n"))

	ms.AssertExpectations(t)
}
