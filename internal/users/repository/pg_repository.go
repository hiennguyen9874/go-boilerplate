package repository

import (
	"context"

	"github.com/hiennguyen9874/stockk-go/internal/models"
	"github.com/hiennguyen9874/stockk-go/internal/repository"
	"github.com/hiennguyen9874/stockk-go/internal/users"
	"gorm.io/gorm"
)

type UserRepo struct {
	repository.PgRepo[models.User]
}

func CreateUserRepository(db *gorm.DB) users.UserRepository {
	return &UserRepo{
		PgRepo: repository.CreatePgRepo[models.User](db),
	}
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (res *models.User, err error) {
	var obj *models.User
	if result := r.DB.WithContext(ctx).First(&obj, "email = ?", email); result.Error != nil {
		return nil, result.Error
	}
	return obj, nil
}

func (r *UserRepo) UpdatePassword(ctx context.Context, exp *models.User, newPassword string) (res *models.User, err error) {
	if result := r.DB.WithContext(ctx).Model(&exp).Select("password").Updates(map[string]interface{}{"password": newPassword}); result.Error != nil {
		return nil, result.Error
	}
	return exp, nil
}
