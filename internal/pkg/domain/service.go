package domain

import (
	"context"
	"encoding/json"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/sirupsen/logrus"
)

const (
	componentService = "Service"
)

type Service struct {
	logger *logrus.Entry
	repo   Repo
	idTool IDTool
}

func NewService(logger *logrus.Logger, repo Repo, idTool IDTool) *Service {
	return &Service{
		logger: logger.WithField("component", componentService),
		repo:   repo,
		idTool: idTool,
	}
}

func (s *Service) CreateUser(ctx context.Context, userName, email, password string) (*User, *Error) {
	if userName == "" || email == "" || password == "" {
		return nil, newInvalidInputError("each of userName, email, password must not be empty", nil)
	}

	id, err := s.idTool.New()
	if err != nil {
		return nil, newSystemError("failed to generate valid ID", err)
	}

	user := NewUser(id, userName, email, password)
	if s.repo.CreateUser(ctx, *user); err != nil {
		return nil, newSystemError("failed to store new user", err)
	}

	if err := user.Validate(s.idTool); err != nil {
		return nil, newInvalidInputError("user is invalid", err)
	}

	return user, nil
}

func (s *Service) GetUser(ctx context.Context, userID string) (*User, *Error) {
	if !s.idTool.IsValid(userID) {
		return nil, newInvalidInputError("userID must be a valid UUIDv4", nil)
	}

	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, newSystemError("failed to store new user", err)
	}

	return user, nil
}

// PatchUser updates the user attributes with the given patch
func (s *Service) PatchUser(ctx context.Context, userID string, patchJSON []byte) error {
	if !s.idTool.IsValid(userID) {
		return newInvalidInputError("userID must be a valid UUIDv4", nil)
	}

	patch, err := jsonpatch.DecodePatch(patchJSON)
	if err != nil {
		return newInvalidInputError("patch could not be decoded", err)
	}

	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return newSystemError("failed to retrieve user", err)
	} else if user == nil {
		return newResourceNotFoundError("user does not exist", nil)
	}

	currentUserAttrJSON, err := json.Marshal(user.Attributes)
	if err != nil {
		return newSystemError("failed to marshal existing user", err)
	}

	patchedUserAttr, dErr := s.createPatchedUser(currentUserAttrJSON, patch)
	if dErr != nil {
		return dErr.WrapMessage("failed to patch user")
	}

	user.Attributes = *patchedUserAttr

	// need to validate user post-patch to ensure we're not left in an invalid state
	if err := user.Validate(s.idTool); err != nil {
		return newInvalidInputError("patched leaves user in invalid state", err)
	}

	if err := s.repo.UpdateUser(ctx, *user); err != nil {
		return newSystemError("failed to update user with patched attributes", err)
	}

	return nil
}

// createPatchedUser creates a new user from the current user and the patch.
// https://jsonpatch.com/
func (s *Service) createPatchedUser(userAttr []byte, patch jsonpatch.Patch) (*UserAttributes, *Error) {
	patchedUserAttr, err := patch.Apply(userAttr)
	if err != nil {
		return nil, newInvalidInputError("patch invalid", err)
	}

	// check if the patch actually changed anything
	if jsonpatch.Equal(userAttr, patchedUserAttr) {
		return nil, newInvalidInputError("patch does not effect any change", nil)
	}

	var patchedUser UserAttributes
	if err := json.Unmarshal(patchedUserAttr, &patchedUser); err != nil {
		return nil, newSystemError("could not unmarshal patched user", err)
	}

	return &patchedUser, nil
}
