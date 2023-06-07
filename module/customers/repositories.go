package customers

import (
	"errors"
	"gorm.io/gorm"
)

type ActorsRepository interface {
	Create(actor *Actor) error
	FindByUsername(username string) (*Actor, error)
}

type actorsRepository struct {
	db *gorm.DB
}

func NewActorsRepository(db *gorm.DB) ActorsRepository {
	return &actorsRepository{
		db: db,
	}
}

func (r *actorsRepository) Create(actor *Actor) error {
	return r.db.Select("username", "password", "role_id", "flag_act").Create(actor).Error
}

func (r *actorsRepository) FindByUsername(token_key string) (*Actor, error) {
	var actor Actor
	err := r.db.Where("token_key = ?", token_key).First(&actor).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &actor, nil
}
