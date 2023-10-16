package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUserByEmail(email string) (model.User, error)
	CreateUser(user model.User) (model.User, error)
	GetUserTaskCategory(UserID int) ([]model.UserTaskCategory, error)
	GetUserProfile(userID int) (model.UserProfile, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *userRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByEmail(email string) (model.User, error) {
	user := model.User{}
	err := r.db.Where("email = ?", email).First(&model.User{}).Scan(&user)
	if err != nil {
		return user, nil
	}

	return user, nil
}

func (r *userRepository) CreateUser(user model.User) (model.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) GetUserTaskCategory(userID int) ([]model.UserTaskCategory, error) {
	userTaskCategories := []model.UserTaskCategory{}
	err := r.db.Raw("SELECT t.id as task_id, u.id AS id, u.fullname AS fullname, u.email AS email, t.title AS task, t.deadline AS deadline, t.priority AS priority, t.status AS status, c.name AS category FROM users u, tasks t, categories c WHERE  t.user_id = ? AND u.id = t.user_id AND t.category_id = c.id", userID).Scan(&userTaskCategories).Error
	if err != nil {
		return userTaskCategories, err
	}
	return userTaskCategories, nil
}

func (r *userRepository) GetUserProfile(userID int) (model.UserProfile, error) {
	userProfile := model.UserProfile{}
	err := r.db.Raw("SELECT u.id AS id, u.nik as nik, u.fullname AS fullname, u.email AS email, u.address AS address, u.id_card AS id_card FROM users u WHERE u.id = ?", userID).Scan(&userProfile).Error
	if err != nil {
		return userProfile, err
	}
	return userProfile, nil
}
