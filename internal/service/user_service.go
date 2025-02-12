package service

import (
	"context"
	"go-scheduler/config"
	"go-scheduler/internal/repository"
	"go-scheduler/internal/repository/entity"
	"go-scheduler/logger"
	"go-scheduler/model"
	"sync"
)

type UserServiceInterface interface {
	GetUsers(ctx context.Context) ([]*entity.GetUser, error)
	BatchUpdateUser(ctx context.Context, users []*model.UpdateUserRequest, batchSize int, workerCount int, maxRetries int)
}

type userService struct {
	config         *config.Config
	userRepository repository.UserRepositoryInterface
}

func NewUserService(config *config.Config, userRepo repository.UserRepositoryInterface) UserServiceInterface {
	return &userService{
		config:         config,
		userRepository: userRepo,
	}
}

func (u *userService) GetUsers(ctx context.Context) ([]*entity.GetUser, error) {
	users, err := u.userRepository.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	logger.InfofCtx(ctx, "Total users: %d", len(users))
	return users, nil
}

func (u *userService) workerBatchUpdateUser(ctx context.Context, jobs <-chan []*model.UpdateUserRequest, wg *sync.WaitGroup, maxRetries int) {
	defer wg.Done()
	totalUpdated := 0
	for batch := range jobs {
		err := u.userRepository.BatchUpdateUser(ctx, batch, maxRetries)
		if err != nil {
			logger.ErrorfCtx(ctx, "Batch update failed: %v", err)
		} else {
			countBatch := len(batch)
			totalUpdated += countBatch
			logger.InfofCtx(ctx, "Updated %d users", countBatch)
		}
	}

	logger.InfofCtx(ctx, "Total updated: %d", totalUpdated)
}

func (u *userService) BatchUpdateUser(ctx context.Context, users []*model.UpdateUserRequest, batchSize int, workerCount int, maxRetries int) {
	jobs := make(chan []*model.UpdateUserRequest, len(users)/batchSize+1)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go u.workerBatchUpdateUser(ctx, jobs, &wg, maxRetries)
	}

	// Enqueue batches
	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}
		jobs <- users[i:end]
	}
	close(jobs)

	// Wait for all workers to complete
	wg.Wait()
}
