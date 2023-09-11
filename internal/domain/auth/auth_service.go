package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/evermos/boilerplate-go/configs"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user *User) error
	Login(username, password string) (string, error)
}

type AuthServiceImpl struct {
	AuthRepository AuthRepository
}

func ProvideAuthServiceImpl(authRepository AuthRepository) *AuthServiceImpl {
	return &AuthServiceImpl{
		AuthRepository: authRepository,
	}
}

func (s *AuthServiceImpl) Register(user *User) error {
	existingUser, err := s.AuthRepository.IsExist(user.Username)
	if err != nil {
		log.Error().Err(err).Msg("Something went wrong")
		return err
	}
	if existingUser {
		log.Error().Msg("Username already exists")
		return ErrUserExist
	}
	return s.AuthRepository.Register(user)
}

func (s *AuthServiceImpl) UserCheck(username, password string) (*Access, error) {
	// Fetch the user by username from the repository
	user, err := s.AuthRepository.GetUserByUsername(username)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get user by username")
		return nil, ErrNotFound
	}

	// Check if the user exists and the password is correct
	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		log.Error().Err(err).Msg("User is not exist or password incorrect")
		return nil, ErrUnauthorized
	}

	return user, nil
}

func (s *AuthServiceImpl) Login(username, password string) (string, error) {
	// also return token expired date
	user, err := s.UserCheck(username, password)
	if err != nil {
		log.Error().Err(err).Msg("Failed to check user")
		return "", err
	}

	// Generate JWT token
	token, err := GenerateJWT(user)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate jwt")
		return "", err
	}

	return token, nil
}

func GenerateJWT(access *Access) (string, error) {
	// Create the claims for the JWT token
	claims := jwt.MapClaims{
		"user_id":  access.ID,
		"username": access.Username,
		"role":     access.Role,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	// Create the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	config := configs.Get()

	tokenString, err := token.SignedString([]byte(config.App.JWTAccessKey))
	if err != nil {
		log.Error().Err(err).Msg("Service: Failed to generate jwt")
		return "", err
	}

	return tokenString, nil
}
