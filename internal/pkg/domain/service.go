package domain

import (
	"context"
	"encoding/json"
	"strings"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/sirupsen/logrus"
)

const (
	componentService = "Service"
)

type Service struct {
	logger *logrus.Entry
	repo   Repo

	passWordTool PasswordTool
	idTool       IDTool
}

func NewService(logger *logrus.Logger, repo Repo, idTool IDTool, passwordTool PasswordTool) *Service {
	return &Service{
		logger:       logger.WithField("component", componentService),
		repo:         repo,
		idTool:       idTool,
		passWordTool: passwordTool,
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

	// generate salted hash from plaintext password
	saltedHash, err := s.passWordTool.New(password)
	if err != nil {
		return nil, newSystemError("failed to hash password", err)
	}

	// create and validate user in memory
	user := NewUser(id, userName, email, saltedHash)
	if err := user.Validate(s.idTool, s.passWordTool); err != nil {
		return nil, newInvalidInputError("user is invalid", err)
	}

	if err := s.repo.CreateUser(ctx, *user); err != nil {
		return nil, newSystemError("failed to store new user", err)
	}

	return user, nil
}

func (s *Service) GetUser(ctx context.Context, userID string) (*User, *Error) {
	if !s.idTool.IsValid(userID) {
		return nil, newInvalidInputError("format of userID is invalid", nil)
	}

	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, newSystemError("failed to retrieve user", err)
	} else if user == nil {
		return nil, newResourceNotFoundError("user does not exist", nil)
	}

	return user, nil
}

// PatchUser updates the user attributes with the given patch
func (s *Service) PatchUser(ctx context.Context, userID string, patchJSON []byte) *Error {
	if !s.idTool.IsValid(userID) {
		return newInvalidInputError("format of userID is invalid", nil)
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
	if err := user.Validate(s.idTool, s.passWordTool); err != nil {
		return newInvalidInputError("patch would leave user in invalid state", err)
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

// ValidateUserCredentials checks if the given credentials are valid. If they are, we return the User, if not we return an error
func (s *Service) ValidateUserCredentials(ctx context.Context, userName, email, password string) (*User, *Error) {
	if password == "" {
		return nil, newInvalidInputError("password must not be empty", nil)
	}

	var user *User
	var err error
	if userName != "" {
		user, err = s.repo.GetUserByUserName(ctx, userName)
	} else if email != "" {
		user, err = s.repo.GetUserByEmail(ctx, email)
	} else {
		return nil, newInvalidInputError("must provide either username or email", nil)
	}

	if user == nil {
		return nil, newResourceNotFoundError("user does not exist", nil)
	} else if err != nil {
		return nil, newSystemError("failed to retrieve user", err)
	}

	if err := s.passWordTool.Check(user.Password, password); err != nil && strings.Contains(err.Error(), "mismatched hash and password") {
		return nil, newUnauthorizedError("credentials are invalid", nil)
	} else if err != nil {
		return nil, newSystemError("failed to validate password", err)
	}

	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, *Error) {
	if email == "" {
		return nil, newInvalidInputError("email must not be empty", nil)
	}

	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, newSystemError("failed to retrieve user", err)
	} else if user == nil {
		return nil, newResourceNotFoundError("user does not exist", nil)
	}

	return user, nil
}

func (s *Service) GetUserByUserName(ctx context.Context, userName string) (*User, *Error) {
	if userName == "" {
		return nil, newInvalidInputError("username must not be empty", nil)
	}

	user, err := s.repo.GetUserByUserName(ctx, userName)
	if err != nil {
		return nil, newSystemError("failed to retrieve user", err)
	} else if user == nil {
		return nil, newResourceNotFoundError("user does not exist", nil)
	}

	return user, nil
}
