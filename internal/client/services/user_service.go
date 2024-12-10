package service

import (
	"context"
	"errors"

	"github.com/ShvetsovYura/oykeeper/internal/logger"
)

type UserService struct {
	otpService *Totp
}

func NewUserService(otpService *Totp) *UserService {
	return &UserService{
		otpService: otpService,
	}
}

func (s *UserService) Register(ctx context.Context, login string) error {
	s.otpService.Generate(login)
	key := s.otpService.GetKey()
	if key == nil {
		logger.Log.Error("не удалось получить ключ регистрации")
		return errors.New("error on generate key")
	}
	s.otpService.GenerateImage()
	return nil
}

func (s *UserService) Login(ctx context.Context, url string, code string) error {
	s.otpService.Load(url)
	valid := s.otpService.Validate(code)
	if !valid {
		return errors.New("not valid login data")
	}
	return nil
}
