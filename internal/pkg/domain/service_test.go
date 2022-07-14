package domain

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

///////////////////
//  CreateUser  //
/////////////////

func TestCreateUser_successPath(t *testing.T) {
	mr := &mockRepo{}
	mIDt := &mockIDTool{}
	mpt := &mockPasswordTool{}

	var (
		uID        = "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
		userName   = "JohnDoe"
		email      = "test@example.com"
		password   = "password"
		saltedHash = "$2a$10$zKDq1KOCqy430Fa1oyZs5eqSvyk7U6e8.wlgXTGEUDy7nX/a7lnWK"
	)

	expectedUser := User{
		ID: uID,
		Attributes: UserAttributes{
			UserName: userName,
			Email:    email,
		},
		Password: saltedHash,
	}

	mIDt.On("New").Return(uID, nil).Once()
	mIDt.On("IsValid", uID).Return(true).Once()

	mpt.On("New", password).Return(saltedHash, nil).Once()
	mpt.On("IsValid", saltedHash).Return(true).Once()

	mr.On("CreateUser", mock.Anything, expectedUser).Return(nil).Once()

	service := NewService(nullLogger(), mr, mIDt, mpt)
	user, err := service.CreateUser(context.Background(), userName, email, password)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedUser, *user)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
	mpt.AssertExpectations(t)
}

func TestCreateUser_emptyInput_failurePath(t *testing.T) {
	const (
		userName = "JohnDoe"
		email    = "test@example.com"
		password = "password"
	)

	testCases := []struct {
		name     string
		userName string
		email    string
		password string
	}{
		{
			name:     "empty user name",
			userName: "",
			email:    email,
			password: password,
		},
		{
			name:     "empty email",
			userName: userName,
			email:    "",
			password: password,
		},
		{
			name:     "empty password",
			userName: userName,
			email:    email,
			password: "",
		},
	}

	service := NewService(nullLogger(), nil, nil, nil)
	expectedErr := newInvalidInputError("each of userName, email, password must not be empty", nil)

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx, tc.name), func(t *testing.T) {
			user, err := service.CreateUser(context.Background(), tc.userName, tc.email, tc.password)
			assert.Nil(t, user)
			assert.Equal(t, expectedErr, err)
		})
	}
}

func TestCreateUser_errorGeneratingID_successPath(t *testing.T) {
	const (
		userName = "JohnDoe"
		email    = "test@example.com"
		password = "password"
	)

	expectedIDErr := fmt.Errorf("no ID for you")
	mIDt := &mockIDTool{}
	mIDt.On("New").Return("", fmt.Errorf("no ID for you")).Once()

	service := NewService(nullLogger(), nil, mIDt, nil)

	user, err := service.CreateUser(context.Background(), userName, email, password)
	assert.Nil(t, user)
	assert.Equal(t, newSystemError("failed to generate valid ID", expectedIDErr), err)

	mIDt.AssertExpectations(t)
}

func TestCreateUser_repoError_failurePath(t *testing.T) {
	mr := &mockRepo{}
	mIDt := &mockIDTool{}
	mpt := &mockPasswordTool{}

	const (
		uID        = "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
		userName   = "JohnDoe"
		email      = "test@example.com"
		password   = "password"
		saltedHash = "$2a$10$zKDq1KOCqy430Fa1oyZs5eqSvyk7U6e8.wlgXTGEUDy7nX/a7lnWK"
	)

	expectedUser := User{
		ID: uID,
		Attributes: UserAttributes{
			UserName: userName,
			Email:    email,
		},
		Password: saltedHash,
	}

	mIDt.On("New").Return(uID, nil).Once()
	mIDt.On("IsValid", uID).Return(true).Once()

	mpt.On("New", password).Return(saltedHash, nil).Once()
	mpt.On("IsValid", saltedHash).Return(true).Once()

	repoErr := fmt.Errorf("repo error")
	mr.On("CreateUser", mock.Anything, expectedUser).Return(repoErr).Once()

	service := NewService(nullLogger(), mr, mIDt, mpt)
	user, err := service.CreateUser(context.Background(), userName, email, password)
	assert.Nil(t, user)
	assert.Equal(t, newSystemError("failed to store new user", repoErr), err)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
	mpt.AssertExpectations(t)
}

func TestCreateUser_badInput_failurePath(t *testing.T) {
	const (
		uID              = "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
		userName         = "JohnDoe"
		email            = "notAnEmail"
		password         = "password"
		saltedHash       = "$2a$10$zKDq1KOCqy430Fa1oyZs5eqSvyk7U6e8.wlgXTGEUDy7nX/a7lnWK"
		randomCharsLen21 = "WIeIluERPJhLEDXq5yIhO"
	)

	testCases := []struct {
		name     string
		userName string
		email    string
	}{
		{
			name:     "bad email",
			userName: userName,
			email:    "notAnEmail",
		},
		{
			name:     "bad userName",
			userName: randomCharsLen21,
			email:    email,
		},
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx, tc.name), func(t *testing.T) {
			mr := &mockRepo{}
			mIDt := &mockIDTool{}
			mpt := &mockPasswordTool{}

			mIDt.On("New").Return(uID, nil).Once()
			mIDt.On("IsValid", uID).Return(true).Once()

			mpt.On("New", password).Return(saltedHash, nil).Once()

			service := NewService(nullLogger(), mr, mIDt, mpt)
			user, err := service.CreateUser(context.Background(), userName, email, password)
			assert.Nil(t, user)
			assert.Equal(t, "user is invalid", err.Msg)
			assert.Equal(t, InvalidInput, err.Code)

			mr.AssertExpectations(t)
			mIDt.AssertExpectations(t)
			mpt.AssertExpectations(t)
		})
	}
}

// ////////////////
// //  GetUser  //
// //////////////

func TestGetUser_successPath(t *testing.T) {
	mr := &mockRepo{}
	mIDt := &mockIDTool{}

	uID := "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
	expectedUser := User{
		ID: uID,
		Attributes: UserAttributes{
			UserName: "foo",
			Email:    "bar@example.com",
		},
		Password: "password",
	}

	mr.On("GetUser", mock.Anything, uID).Return(&expectedUser, nil).Once()
	mIDt.On("IsValid", uID).Return(true).Once()

	service := NewService(nullLogger(), mr, mIDt, nil)
	user, err := service.GetUser(context.Background(), uID)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedUser, *user)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
}

func TestGetUser_invalidID_failurePath(t *testing.T) {
	mIDt := &mockIDTool{}

	uID := "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"

	mIDt.On("IsValid", uID).Return(false).Once()

	service := NewService(nullLogger(), nil, mIDt, nil)
	user, err := service.GetUser(context.Background(), uID)
	assert.Nil(t, user)

	expectedErr := newInvalidInputError("format of userID is invalid", nil)
	assert.EqualValues(t, expectedErr, err)

	mIDt.AssertExpectations(t)
}

func TestGetUser_repoError_failurePath(t *testing.T) {
	mr := &mockRepo{}
	mIDt := &mockIDTool{}

	uID := "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
	repoErr := fmt.Errorf("db failed")

	mr.On("GetUser", mock.Anything, uID).Return(nil, repoErr).Once()
	mIDt.On("IsValid", uID).Return(true).Once()

	service := NewService(nullLogger(), mr, mIDt, nil)
	user, err := service.GetUser(context.Background(), uID)
	assert.Nil(t, user)

	expectedErr := newSystemError("failed to retrieve user", repoErr)
	assert.EqualValues(t, expectedErr, err)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
}

func TestGetUser_userNotFound_failurePath(t *testing.T) {
	mr := &mockRepo{}
	mIDt := &mockIDTool{}

	uID := "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"

	mr.On("GetUser", mock.Anything, uID).Return(nil, nil).Once()
	mIDt.On("IsValid", uID).Return(true).Once()

	service := NewService(nullLogger(), mr, mIDt, nil)
	user, err := service.GetUser(context.Background(), uID)
	assert.Nil(t, user)

	expectedErr := newResourceNotFoundError("user does not exist", nil)
	assert.EqualValues(t, expectedErr, err)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
}

// ////////////////
// // PatchUser //
// //////////////

func TestPatchUser_successPath(t *testing.T) {
	mr := &mockRepo{}
	mIDt := &mockIDTool{}
	mpt := &mockPasswordTool{}

	uID := "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
	originalUser := User{
		ID: uID,
		Attributes: UserAttributes{
			UserName: "JohnDoe",
			Email:    "test@example.com",
		},
		Password: "password",
	}

	patchJSON := `[
		{ "op": "replace", "path": "/user_name", "value": "foo" },
		{ "op": "replace", "path": "/email", "value": "bar@example.com" }
	]`

	// once patch is applied we expect to see the following user
	patchedUser := User{
		ID: uID,
		Attributes: UserAttributes{
			UserName: "foo",
			Email:    "bar@example.com",
		},
		Password: "password",
	}

	mr.On("GetUser", mock.Anything, uID).Return(&originalUser, nil)
	mr.On("UpdateUser", mock.Anything, patchedUser).Return(nil)

	mIDt.On("IsValid", uID).Return(true).Twice()

	mpt.On("IsValid", "password").Return(true).Once()

	service := NewService(nullLogger(), mr, mIDt, mpt)
	err := service.PatchUser(context.Background(), uID, []byte(patchJSON))
	assert.Nil(t, err)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
}

func TestPatchUser_invalidID_failurePath(t *testing.T) {
	mIDt := &mockIDTool{}
	mIDt.On("IsValid", "nope").Return(false).Once()

	service := NewService(nullLogger(), nil, mIDt, nil)
	err := service.PatchUser(context.Background(), "nope", []byte("[]"))

	expectedErr := newInvalidInputError("format of userID is invalid", nil)
	assert.Equal(t, expectedErr.Error(), err.Error())

	mIDt.AssertExpectations(t)
}

func TestPatchUser_invalidJSONInPatchJSON_failurePath(t *testing.T) {
	uID := "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
	patchJSON := `[{"this": "ain't valid json}]`

	mIDt := &mockIDTool{}
	mIDt.On("IsValid", uID).Return(true).Once()

	service := NewService(nullLogger(), nil, mIDt, nil)
	err := service.PatchUser(context.Background(), uID, []byte(patchJSON))

	expectedErr := newInvalidInputError("patch could not be decoded", fmt.Errorf("unexpected end of JSON input"))
	assert.Equal(t, expectedErr.Error(), err.Error())

	mIDt.AssertExpectations(t)
}

func TestPatchUser_invalidPatchJSON_failurePath(t *testing.T) {
	uID := "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
	originalUser := User{
		ID: uID,
		Attributes: UserAttributes{
			UserName: "JohnDoe",
			Email:    "test@example.com",
		},
		Password: "password",
	}

	patchJSON := `[{"this_aint": "a valid patch"}]`

	mr := &mockRepo{}
	mr.On("GetUser", mock.Anything, uID).Return(&originalUser, nil).Once()

	mIDt := &mockIDTool{}
	mIDt.On("IsValid", uID).Return(true).Once()

	service := NewService(nullLogger(), mr, mIDt, nil)
	err := service.PatchUser(context.Background(), uID, []byte(patchJSON))

	expectedErr := newInvalidInputError("failed to patch user, patch invalid", fmt.Errorf("Unexpected kind: unknown"))
	assert.Equal(t, expectedErr.Error(), err.Error())

	mIDt.AssertExpectations(t)
}

func TestPatchUser_repoErrorRetrievingExistingUser_failurePath(t *testing.T) {
	mr := &mockRepo{}
	mIDt := &mockIDTool{}

	uID := "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"

	patchJSON := `[
		{ "op": "replace", "path": "/user_name", "value": "foo" },
		{ "op": "replace", "path": "/email", "value": "bar@example.com" }
	]`

	repoErr := fmt.Errorf("repo error")
	mr.On("GetUser", mock.Anything, uID).Return(nil, repoErr).Once()

	mIDt.On("IsValid", uID).Return(true).Once()

	service := NewService(nullLogger(), mr, mIDt, nil)
	err := service.PatchUser(context.Background(), uID, []byte(patchJSON))

	expectedErr := newSystemError("failed to retrieve user", repoErr)
	assert.Equal(t, expectedErr, err)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
}

func TestPatchUser_userDoesNotExist_failurePath(t *testing.T) {
	mr := &mockRepo{}
	mIDt := &mockIDTool{}

	uID := "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"

	patchJSON := `[
		{ "op": "replace", "path": "/user_name", "value": "foo" },
		{ "op": "replace", "path": "/email", "value": "bar@example.com" }
	]`

	mr.On("GetUser", mock.Anything, uID).Return(nil, nil).Once()

	mIDt.On("IsValid", uID).Return(true).Once()

	service := NewService(nullLogger(), mr, mIDt, nil)
	err := service.PatchUser(context.Background(), uID, []byte(patchJSON))

	expectedErr := newResourceNotFoundError("user does not exist", nil)
	assert.Equal(t, expectedErr, err)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
}

func TestPatchUser_patchDoesntChangeAnything_failurePath(t *testing.T) {
	mr := &mockRepo{}
	mIDt := &mockIDTool{}

	uID := "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
	originalUser := User{
		ID: uID,
		Attributes: UserAttributes{
			UserName: "JohnDoe",
			Email:    "test@example.com",
		},
		Password: "password",
	}

	patchJSON := `[
		{ "op": "replace", "path": "/user_name", "value": "JohnDoe" },
		{ "op": "replace", "path": "/email", "value": "test@example.com" }
	]`

	mr.On("GetUser", mock.Anything, uID).Return(&originalUser, nil).Once()

	mIDt.On("IsValid", uID).Return(true).Once()

	service := NewService(nullLogger(), mr, mIDt, nil)
	err := service.PatchUser(context.Background(), uID, []byte(patchJSON))

	expectedErr := newInvalidInputError("patch does not effect any change", nil).WrapMessage("failed to patch user")
	assert.Equal(t, expectedErr, err)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
}

func TestPatchUser_patchMakesInvalidChanges_successPath(t *testing.T) {
	testCases := []struct {
		name      string
		patchJSON string
	}{
		{
			name:      "patch changes user name to invalid value",
			patchJSON: `[{ "op": "replace", "path": "/user_name", "value": "" }]`,
		},
		{
			name:      "patch changes email to invalid value",
			patchJSON: `[{ "op": "replace", "path": "/email", "value": "notAnEmail" }]`,
		},
	}

	uID := "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
	originalUser := User{
		ID: uID,
		Attributes: UserAttributes{
			UserName: "JohnDoe",
			Email:    "test@example.com",
		},
		Password: "password",
	}

	for idx, tc := range testCases {
		t.Run(fmt.Sprintf("test case %d: %s", idx, tc.name), func(t *testing.T) {
			mr := &mockRepo{}
			mIDt := &mockIDTool{}

			mr.On("GetUser", mock.Anything, uID).Return(&originalUser, nil).Once()

			mIDt.On("IsValid", uID).Return(true).Twice()

			service := NewService(nullLogger(), mr, mIDt, nil)
			err := service.PatchUser(context.Background(), uID, []byte(tc.patchJSON))

			assert.Equal(t, InvalidInput, err.Code)
			assert.Equal(t, "patch would leave user in invalid state", err.Msg)

			mr.AssertExpectations(t)
			mIDt.AssertExpectations(t)
		})
	}
}

func TestPatchUser_repoErrorUpdatingUser_failurePath(t *testing.T) {
	var (
		uID        = "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
		userName   = "JohnDoe"
		email      = "test@example.com"
		saltedHash = "$2a$10$zKDq1KOCqy430Fa1oyZs5eqSvyk7U6e8.wlgXTGEUDy7nX/a7lnWK"
	)

	mr := &mockRepo{}
	mIDt := &mockIDTool{}
	mpt := &mockPasswordTool{}

	originalUser := User{
		ID: uID,
		Attributes: UserAttributes{
			UserName: userName,
			Email:    email,
		},
		Password: saltedHash,
	}

	patchJSON := `[
		{ "op": "replace", "path": "/user_name", "value": "foo" },
		{ "op": "replace", "path": "/email", "value": "bar@example.com" }
	]`

	// once patch is applied we expect to see the following user
	patchedUser := User{
		ID: uID,
		Attributes: UserAttributes{
			UserName: "foo",
			Email:    "bar@example.com",
		},
		Password: saltedHash,
	}

	repoErr := fmt.Errorf("repo error")
	mr.On("GetUser", mock.Anything, uID).Return(&originalUser, nil).Once()
	mr.On("UpdateUser", mock.Anything, patchedUser).Return(repoErr).Once()

	mpt.On("IsValid", saltedHash).Return(true).Once()

	mIDt.On("IsValid", uID).Return(true).Twice()

	service := NewService(nullLogger(), mr, mIDt, mpt)
	err := service.PatchUser(context.Background(), uID, []byte(patchJSON))

	expectedErr := newSystemError("failed to update user with patched attributes", repoErr)
	assert.Equal(t, expectedErr, err)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
}
