package category

import (
	"context"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{
		db: db,
	}
}

func (r *CategoryRepository) Create(
	ctx context.Context,
	category *Category,
) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *CategoryRepository) GetByID(
	ctx context.Context,
	id uint,
) (*Category, error) {
	var category Category

	err := r.db.WithContext(ctx).
		First(&category, id).Error

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) GetByUserID(
	ctx context.Context,
	userID uint,
) ([]Category, error) {
	var categories []Category

	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("name ASC").
		Find(&categories).Error

	return categories, err
}

func (r *CategoryRepository) Update(
	ctx context.Context,
	category *Category,
) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *CategoryRepository) Delete(
	ctx context.Context,
	category *Category,
) error {
	return r.db.WithContext(ctx).Delete(category).Error
}
