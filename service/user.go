package service

import (
	"a21hc3NpZ25tZW50/model"
	repo "a21hc3NpZ25tZW50/repository"
	"a21hc3NpZ25tZW50/utils"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type UserService interface {
	Register(user *model.User) (model.User, error)
	Login(user *model.User) (token *string, err error)
	GetUserTaskCategory(userID int) ([]model.UserTaskCategory, error)
	GetUserProfile(userID int) (model.UserProfile, error)
}

type userService struct {
	userRepo     repo.UserRepository
	sessionsRepo repo.SessionRepository
}

func NewUserService(userRepository repo.UserRepository, sessionsRepo repo.SessionRepository) UserService {
	return &userService{userRepository, sessionsRepo}
}

func (s *userService) Register(user *model.User) (model.User, error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return *user, err
	}

	// user.Password = utils.EncryptAES(user.Password)
	// user.IDCard = utils.EncryptAES(user.IDCard)
	user.IDCard, err = utils.EncryptRC4(user.IDCard)
	if err != nil {
		return *user, errors.New("error encrypting id card")
	}
	user.Address, err = utils.EncryptRC4(user.Address)
	// fmt.Printf("user address: %+v\n", user.Address)
	if err != nil {
		return *user, errors.New("error encrypting address")
	}
	// user.NIK, err = utils.EncryptDES(user.NIK)
	user.NIK = utils.EncryptAES(user.NIK)
	// fmt.Printf("user nik: %+v\n", user.NIK)
	user.Password = utils.EncryptAES(user.Password)

	if dbUser.Email != "" || dbUser.ID != 0 {
		return *user, errors.New("email already exists")
	}
	PrivateKey, PublicKey, err := utils.GenerateKeyPair()
	user.PrivateKey = PrivateKey
	user.PublicKey = PublicKey

	user.CreatedAt = time.Now()

	newUser, err := s.userRepo.CreateUser(*user)
	if err != nil {
		return *user, err
	}

	return newUser, nil
}

func (s *userService) Login(user *model.User) (token *string, err error) {
	dbUser, err := s.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if dbUser.Email == "" || dbUser.ID == 0 {
		return nil, errors.New("user not found")
	}

	fmt.Printf("dbUser password: %+v\n", dbUser.Password)

	decryptedPassword, err := utils.DecryptAES(dbUser.Password)
	// decryptedPassword, err := utils.DecryptRC4(dbUser.Password)

	fmt.Printf("decrypt dbUser password: %+v\n", decryptedPassword)
	fmt.Printf("user password: %+v\n", user.Password)

	if user.Password != decryptedPassword {
		return nil, errors.New("wrong email or password")
	}

	expirationTime := time.Now().Add(20 * time.Minute)
	claims := &model.Claims{
		ID:    dbUser.ID,
		Email: dbUser.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := t.SignedString(model.JwtKey)
	if err != nil {
		return nil, err
	}

	session := model.Session{
		Token:  tokenString,
		Email:  user.Email,
		Expiry: expirationTime,
	}

	_, err = s.sessionsRepo.SessionAvailEmail(session.Email)
	if err != nil {
		err = s.sessionsRepo.AddSessions(session)
	} else {
		err = s.sessionsRepo.UpdateSessions(session)
	}

	return &tokenString, nil
}

func (s *userService) GetUserTaskCategory(userID int) ([]model.UserTaskCategory, error) {
	taskCategory, err := s.userRepo.GetUserTaskCategory(userID)
	for i := 0; i < len(taskCategory); i++ {
		// decryptFilePath, err := utils.DecryptAES(taskCategory[i].DocumentPath)
		decryptFilePath, err := utils.DecryptAES(taskCategory[i].DocumentPath)
		if err != nil {
			return nil, err
		}
		taskCategory[i].DocumentPath = decryptFilePath
	}

	if err != nil {
		return nil, err
	}
	return taskCategory, nil
}

func (s *userService) GetUserProfile(userID int) (model.UserProfile, error) {
	userData, err := s.userRepo.GetUserProfile(userID)

	if err != nil {
		return userData, err
	}

	userData.NIK, err = utils.DecryptAES(userData.NIK)
	if err != nil {
		return userData, err
	}
	userData.Address, err = utils.DecryptRC4(userData.Address)
	if err != nil {
		return userData, err
	}
	// userData.IDCard, err = utils.DecryptAES(userData.IDCard)
	userData.IDCard, err = utils.DecryptRC4(userData.IDCard)
	if err != nil {
		return userData, err
	}

	return userData, nil
}
