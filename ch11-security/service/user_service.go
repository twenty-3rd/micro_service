package service

import (
	"context"
	"errors"
	"micro_server/ch11-security/model"
)

var (
	ErrUserNotExist = errors.New("username is not exist")
	ErrPassword     = errors.New("invalid password")
)

// UserDetailsService Service Define a service interface
type UserDetailsService interface {
	// GetUserDetailByUsername Get UserDetails By username
	GetUserDetailByUsername(ctx context.Context, username, password string) (*model.UserDetails, error)
}

// InMemoryUserDetailsService UserService implement Service interface
type InMemoryUserDetailsService struct {
	userDetailsDict map[string]*model.UserDetails
}

func (service *InMemoryUserDetailsService) GetUserDetailByUsername(ctx context.Context, username, password string) (*model.UserDetails, error) {
	userDetails, ok := service.userDetailsDict[username]
	if ok {
		if userDetails.Password == password {
			return userDetails, nil
		} else {
			return nil, ErrPassword
		}
	} else {
		return nil, ErrUserNotExist
	}
}

func NewInMemoryUserDetailsService(userDetailsList []*model.UserDetails) *InMemoryUserDetailsService {
	userDetailsDict := make(map[string]*model.UserDetails)

	if userDetailsList != nil {
		for _, value := range userDetailsList {
			userDetailsDict[value.Username] = value
		}
	}

	return &InMemoryUserDetailsService{
		userDetailsDict: userDetailsDict,
	}
}