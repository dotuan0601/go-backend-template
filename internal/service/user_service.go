package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/dotuan0601/go-backend-template/internal/model"
	"github.com/dotuan0601/go-backend-template/internal/repository"
)

type UserService interface {
	GetUser(
		ctx context.Context,
		id int64,
	) (*model.User, error)

	GetUsers(
		ctx context.Context,
	) ([]*model.User, error)

	CreateUser(
		ctx context.Context,
		user *model.User,
	) error
}

type userService struct {
	userRepository repository.UserRepository
	jobRepository  repository.JobRepository
	redisClient    *redis.Client
}

func NewUserService(
	userRepository repository.UserRepository,
	jobRepository repository.JobRepository,
	redisClient *redis.Client,
) UserService {

	return &userService{
		userRepository: userRepository,
		jobRepository:  jobRepository,
		redisClient:    redisClient,
	}
}

func (s *userService) GetUser(
	ctx context.Context,
	id int64,
) (*model.User, error) {

	key := fmt.Sprintf(
		"user:%d",
		id,
	)

	// =====================
	// Redis cache
	// =====================

	if s.redisClient != nil {

		data, err := s.redisClient.Get(
			ctx,
			key,
		).Bytes()

		if err == nil {

			var user model.User

			err := json.Unmarshal(
				data,
				&user,
			)

			if err == nil {

				return &user, nil
			}
		}
	}

	// =====================
	// MySQL
	// =====================

	user, err := s.userRepository.GetByID(
		ctx,
		id,
	)

	if err != nil {

		return nil, err
	}

	// =====================
	// Save Redis cache
	// =====================

	if s.redisClient != nil {

		data, err := json.Marshal(
			user,
		)

		if err == nil {

			err = s.redisClient.Set(
				ctx,
				key,
				data,
				5*time.Minute,
			).Err()

			if err != nil {

				// cache fail không ảnh hưởng business
			}
		}
	}

	return user, nil
}

func (s *userService) GetUsers(
	ctx context.Context,
) ([]*model.User, error) {

	users, err := s.userRepository.GetAll(
		ctx,
	)

	if err != nil {

		return nil, err
	}

	return users, nil
}

func (s *userService) CreateUser(
	ctx context.Context,
	user *model.User,
) error {

	if user.Name == "" {

		return errors.New(
			"name is required",
		)
	}

	if user.Email == "" {

		return errors.New(
			"email is required",
		)
	}

	// =====================
	// Create user
	// =====================

	err := s.userRepository.Create(
		ctx,
		user,
	)

	if err != nil {

		return err
	}

	// =====================
	// Create async job
	// =====================

	payload := fmt.Sprintf(
		`{"user_id":%d}`,
		user.ID,
	)

	job := &model.Job{

		Type: model.JobTypeSendEmail,

		Payload: payload,

		Status: model.JobStatusPending,
	}

	err = s.jobRepository.Create(
		ctx,
		job,
	)

	if err != nil {

		return err
	}

	return nil
}
