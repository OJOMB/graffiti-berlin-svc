package domain

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

////////////////
//  GetUser  //
//////////////

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

	service := NewService(nullLogger(), mr, mIDt)
	user, err := service.GetUser(context.Background(), uID)
	assert.Nil(t, err)
	assert.EqualValues(t, expectedUser, *user)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
}

////////////////
// PatchUser //
//////////////

func TestPatchUser_successPath(t *testing.T) {
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

	service := NewService(nullLogger(), mr, mIDt)
	err := service.PatchUser(context.Background(), uID, []byte(patchJSON))
	assert.NoError(t, err)

	mr.AssertExpectations(t)
	mIDt.AssertExpectations(t)
}

func TestPatchUser_invalidJSONInPatchJSON_successPath(t *testing.T) {
	uID := "9abc46be-3bcd-42b1-aeb2-ac6ff557a580"
	patchJSON := `[{"this": "ain't valid json}]`

	mIDt := &mockIDTool{}
	mIDt.On("IsValid", uID).Return(true).Once()

	service := NewService(nullLogger(), nil, mIDt)
	err := service.PatchUser(context.Background(), uID, []byte(patchJSON))

	expectedErr := newInvalidInputError("patch could not be decoded", fmt.Errorf("unexpected end of JSON input"))
	assert.Equal(t, expectedErr.Error(), err.Error())

	mIDt.AssertExpectations(t)
}

func TestPatchUser_invalidPatchJSON_successPath(t *testing.T) {
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

	service := NewService(nullLogger(), mr, mIDt)
	err := service.PatchUser(context.Background(), uID, []byte(patchJSON))

	expectedErr := newInvalidInputError("failed to patch user, patch invalid", fmt.Errorf("Unexpected kind: unknown"))
	assert.Equal(t, expectedErr.Error(), err.Error())

	mIDt.AssertExpectations(t)
}
