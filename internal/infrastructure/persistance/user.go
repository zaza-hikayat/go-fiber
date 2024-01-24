package persistance

import (
	"context"

	"github.com/zaza-hikayat/go-fiber/domain"
	"github.com/zaza-hikayat/go-fiber/dto/models"
	"gorm.io/gorm"
)

type userRepo struct {
	conn *gorm.DB
}

func NewUserRepository(conn *gorm.DB) domain.UserRepository {
	return &userRepo{
		conn: conn,
	}
}

func (r *userRepo) FindByUserId(ctx context.Context, userId string) (models.User, error) {
	var user models.User
	err := r.conn.Model(&models.User{}).First(&user).Error
	return user, err
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
	err := r.conn.Model(&models.User{}).Where("email = ?", email).First(&user).Error
	return user, err
}

func (r *userRepo) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	err := r.conn.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
