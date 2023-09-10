// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package auth

import (
	"context"
	"github.com/google/go-github/v55/github"
	"sync"
)

// Ensure, that GithubClientMock does implement GithubClient.
// If this is not the case, regenerate this file with moq.
var _ GithubClient = &GithubClientMock{}

// GithubClientMock is a mock implementation of GithubClient.
//
//	func TestSomethingThatUsesGithubClient(t *testing.T) {
//
//		// make and configure a mocked GithubClient
//		mockedGithubClient := &GithubClientMock{
//			ExchangeCodeToAccessKeyFunc: func(ctx context.Context, clientID string, clientSeret string, code string) (string, error) {
//				panic("mock out the ExchangeCodeToAccessKey method")
//			},
//			GetUserFunc: func(ctx context.Context, accessKey string, user string) (*github.User, error) {
//				panic("mock out the GetUser method")
//			},
//		}
//
//		// use mockedGithubClient in code that requires GithubClient
//		// and then make assertions.
//
//	}
type GithubClientMock struct {
	// ExchangeCodeToAccessKeyFunc mocks the ExchangeCodeToAccessKey method.
	ExchangeCodeToAccessKeyFunc func(ctx context.Context, clientID string, clientSeret string, code string) (string, error)

	// GetUserFunc mocks the GetUser method.
	GetUserFunc func(ctx context.Context, accessKey string, user string) (*github.User, error)

	// calls tracks calls to the methods.
	calls struct {
		// ExchangeCodeToAccessKey holds details about calls to the ExchangeCodeToAccessKey method.
		ExchangeCodeToAccessKey []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ClientID is the clientID argument value.
			ClientID string
			// ClientSeret is the clientSeret argument value.
			ClientSeret string
			// Code is the code argument value.
			Code string
		}
		// GetUser holds details about calls to the GetUser method.
		GetUser []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// AccessKey is the accessKey argument value.
			AccessKey string
			// User is the user argument value.
			User string
		}
	}
	lockExchangeCodeToAccessKey sync.RWMutex
	lockGetUser                 sync.RWMutex
}

// ExchangeCodeToAccessKey calls ExchangeCodeToAccessKeyFunc.
func (mock *GithubClientMock) ExchangeCodeToAccessKey(ctx context.Context, clientID string, clientSeret string, code string) (string, error) {
	if mock.ExchangeCodeToAccessKeyFunc == nil {
		panic("GithubClientMock.ExchangeCodeToAccessKeyFunc: method is nil but GithubClient.ExchangeCodeToAccessKey was just called")
	}
	callInfo := struct {
		Ctx         context.Context
		ClientID    string
		ClientSeret string
		Code        string
	}{
		Ctx:         ctx,
		ClientID:    clientID,
		ClientSeret: clientSeret,
		Code:        code,
	}
	mock.lockExchangeCodeToAccessKey.Lock()
	mock.calls.ExchangeCodeToAccessKey = append(mock.calls.ExchangeCodeToAccessKey, callInfo)
	mock.lockExchangeCodeToAccessKey.Unlock()
	return mock.ExchangeCodeToAccessKeyFunc(ctx, clientID, clientSeret, code)
}

// ExchangeCodeToAccessKeyCalls gets all the calls that were made to ExchangeCodeToAccessKey.
// Check the length with:
//
//	len(mockedGithubClient.ExchangeCodeToAccessKeyCalls())
func (mock *GithubClientMock) ExchangeCodeToAccessKeyCalls() []struct {
	Ctx         context.Context
	ClientID    string
	ClientSeret string
	Code        string
} {
	var calls []struct {
		Ctx         context.Context
		ClientID    string
		ClientSeret string
		Code        string
	}
	mock.lockExchangeCodeToAccessKey.RLock()
	calls = mock.calls.ExchangeCodeToAccessKey
	mock.lockExchangeCodeToAccessKey.RUnlock()
	return calls
}

// GetUser calls GetUserFunc.
func (mock *GithubClientMock) GetUser(ctx context.Context, accessKey string, user string) (*github.User, error) {
	if mock.GetUserFunc == nil {
		panic("GithubClientMock.GetUserFunc: method is nil but GithubClient.GetUser was just called")
	}
	callInfo := struct {
		Ctx       context.Context
		AccessKey string
		User      string
	}{
		Ctx:       ctx,
		AccessKey: accessKey,
		User:      user,
	}
	mock.lockGetUser.Lock()
	mock.calls.GetUser = append(mock.calls.GetUser, callInfo)
	mock.lockGetUser.Unlock()
	return mock.GetUserFunc(ctx, accessKey, user)
}

// GetUserCalls gets all the calls that were made to GetUser.
// Check the length with:
//
//	len(mockedGithubClient.GetUserCalls())
func (mock *GithubClientMock) GetUserCalls() []struct {
	Ctx       context.Context
	AccessKey string
	User      string
} {
	var calls []struct {
		Ctx       context.Context
		AccessKey string
		User      string
	}
	mock.lockGetUser.RLock()
	calls = mock.calls.GetUser
	mock.lockGetUser.RUnlock()
	return calls
}
