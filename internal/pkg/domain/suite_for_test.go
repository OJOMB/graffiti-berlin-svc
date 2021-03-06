package domain

import (
	"context"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type mockRepo struct {
	mock.Mock
}

func (mr *mockRepo) CreateUser(ctx context.Context, user User) error {
	args := mr.Called(ctx, user)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	return args.Error(0)
}

func (mr *mockRepo) GetUser(ctx context.Context, userID string) (*User, error) {
	args := mr.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*User), args.Error(1)
}

func (mr *mockRepo) UpdateUser(ctx context.Context, user User) error {
	args := mr.Called(ctx, user)
	if args.Error(0) != nil {
		return args.Error(0)
	}

	return args.Error(0)
}

func (mr *mockRepo) GetUserByUserName(ctx context.Context, userName string) (*User, error) {
	args := mr.Called(ctx, userName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*User), args.Error(1)
}

func (mr *mockRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	args := mr.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*User), args.Error(1)
}

type mockIDTool struct {
	mock.Mock
}

func (mIDt *mockIDTool) New() (string, error) {
	args := mIDt.Called()

	return args.Get(0).(string), args.Error(1)
}

func (mIDt *mockIDTool) IsValid(ID string) bool {
	args := mIDt.Called(ID)
	return args.Get(0).(bool)
}

type mockPasswordTool struct {
	mock.Mock
}

func (mpt *mockPasswordTool) New(password string) (string, error) {
	args := mpt.Called(password)

	return args.Get(0).(string), args.Error(1)
}

func (mpt *mockPasswordTool) IsValid(password string) bool {
	args := mpt.Called(password)
	return args.Get(0).(bool)
}

func (mpt *mockPasswordTool) Check(hash, password string) error {
	args := mpt.Called(hash, password)
	return args.Error(0)
}

// Creates a silent logger instance that discards all output
func nullLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	return logger
}
