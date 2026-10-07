package user

import (
	"context"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*User, error) {
	var user User

	err := r.db.
		WithContext(ctx).
		First(&user, id).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByUsername(
	ctx context.Context,
	username string,
) (*User, error) {
	var user User

	err := r.db.
		WithContext(ctx).
		Where("username = ?", username).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	err := r.db.
		WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetByPhone(ctx context.Context, phone string) (*User, error) {
	var user User

	err := r.db.WithContext(ctx).
		Where("phone = ?", phone).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
