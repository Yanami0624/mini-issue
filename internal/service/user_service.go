package service

import (
	"errors"
	. "mini-issue/internal/dao"
	"mini-issue/internal/model"
	. "mini-issue/internal/model"
	. "mini-issue/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	udao *UserDAO
}

func NewUserService(userDAO *UserDAO) *UserService {
	return &UserService{userDAO}
}

func (us UserService) Register(req RegisterRequest) error {
	exist, _ := us.udao.GetByUsername(req.Username)
	switch {
	case len(req.Username) == 0:
		return errors.New("username can not be empty")
	case len(req.Password) < 6:
		return errors.New("password should have 6 characters at least")
	case exist != nil:
		return errors.New("user already exist")
	}

	hashpassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	return us.udao.CreateUser(req.Username, string(hashpassword))
}

func (us UserService) Login(req LoginRequest) (*LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, errors.New("username and password are required")
	}
	user, err := us.udao.GetByUsername(req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("username dose not exist")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("incorrect password")
	}

	token, err := GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{Token: token}, nil
}

func (us *UserService) GetMe(userid int64) (*User, error) {
	user, err := us.udao.GetByUserID(userid)
	if err != nil {
		return nil, errors.New("invalid userid")
	}
	return user, nil
}
