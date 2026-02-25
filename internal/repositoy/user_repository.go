package repositoy

import (
	"context"
	"errors"
	"sika/internal/database"
	dbModels "sika/internal/database/models"
	"sika/internal/domain"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	Get(ctx context.Context, id string) (*domain.User, error)
}

type userRepository struct {
	db database.Database
}

func NewUserRepositor(database database.Database) UserRepository {
	return &userRepository{
		db: database,
	}
}

func (u *userRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	var addressesModel []dbModels.Address
	for _, addr := range user.Addresses {
		addressesModel = append(addressesModel, dbModels.Address{
			Street:  addr.City,
			City:    addr.City,
			State:   addr.State,
			ZipCode: addr.ZipCode,
			Country: addr.Country,
		})
	}

	userModel := dbModels.User{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Addresses:   addressesModel,
	}

	if err := u.db.Db.WithContext(ctx).Create(&userModel).Error; err != nil {
		return nil, err
	}

	return nil, nil
}

func (u *userRepository) Get(ctx context.Context, id string) (*domain.User, error) {
	var userModel dbModels.User

	err := u.db.Db.
		WithContext(ctx).
		Preload("Addresses").
		First(&userModel, "id = ?", id).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	addresses := make([]domain.Address, 0, len(userModel.Addresses))
	for _, a := range userModel.Addresses {
		addresses = append(addresses, domain.Address{
			Street:  a.Street,
			City:    a.City,
			State:   a.State,
			ZipCode: a.ZipCode,
			Country: a.Country,
		})
	}

	user := &domain.User{
		ID:          userModel.ID,
		Name:        userModel.Name,
		Email:       userModel.Email,
		PhoneNumber: userModel.PhoneNumber,
		Addresses:   addresses,
	}

	return user, nil
}

func toUserDomain(user dbModels.User) domain.User {
	return domain.User{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Addresses:   toDomainAddresses(user.Addresses),
	}
}

func toDomainAddresse(addr dbModels.Address) domain.Address {
	return domain.Address{
		Street:  addr.Street,
		City:    addr.City,
		State:   addr.State,
		ZipCode: addr.ZipCode,
		Country: addr.Country,
	}
}

func toDomainAddresses(addres []dbModels.Address) []domain.Address {
	res := make([]domain.Address, len(addres))

	for i, addr := range addres {
		res[i] = toDomainAddresse(addr)
	}

	return res
}
