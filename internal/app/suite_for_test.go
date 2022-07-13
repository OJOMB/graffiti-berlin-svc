package app

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/OJOMB/graffiti-berlin-svc/internal/pkg/domain"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

type mockService struct {
	mock.Mock
}

func (ms *mockService) CreateUser(ctx context.Context, userName, email, password string) (*domain.User, *domain.Error) {
	args := ms.Called(ctx, userName, email, password)

	var user *domain.User
	if args.Get(0) == nil {
		user = nil
	} else {
		user = args.Get(0).(*domain.User)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return user, err
}

func (ms *mockService) GetUser(ctx context.Context, userID string) (*domain.User, *domain.Error) {
	args := ms.Called(ctx, userID)

	var user *domain.User
	if args.Get(0) == nil {
		user = nil
	} else {
		user = args.Get(0).(*domain.User)
	}

	var err *domain.Error
	if args.Get(1) == nil {
		err = nil
	} else {
		err = args.Get(1).(*domain.Error)
	}

	return user, err
}

func (ms *mockService) PatchUser(ctx context.Context, userID string, patchJSON []byte) *domain.Error {
	args := ms.Called(ctx, userID, patchJSON)
	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(*domain.Error)
}

// Creates a logger instance that discards all output
func nullLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	return logger
}

// errReader is intended to help us test - mainly in the corner case of error handling in the case of defective request/response bodies
type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed to read")
}
