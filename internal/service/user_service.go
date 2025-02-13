package service

import (
	"context"
	"go-scheduler/config"
	"go-scheduler/internal/repository"
	"go-scheduler/internal/repository/entity"
	"go-scheduler/logger"
	"go-scheduler/model"
	"sync"
	"sync/atomic"
	"time"
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

func (u *userService) workerBatchUpdateUser(ctx context.Context, workerID int, jobs <-chan []*model.UpdateUserRequest, maxRetries int, batchCounter *int32, resultChan chan<- int64) {
	var totalUpdated int64 = 0

	logger.InfofCtx(ctx, "[Worker %d] 🚀 Start", workerID)

	for batch := range jobs {
		batchID := atomic.AddInt32(batchCounter, 1)
		startTime := time.Now()

		// logger.InfofCtx(ctx, "[Worker %d] 🚀 Processing batch %d with %d users", workerID, batchID, len(batch))

		// Attempt batch update with retries
		var err error
		for retry := 0; retry < maxRetries; retry++ {
			// Forced error for testing retry
			// if retry > 0 {
			// 	err = u.userRepository.BatchUpdateUser(ctx, batch, workerID, int(batchID))
			// 	if err == nil {
			// 		break // Jika berhasil, keluar dari loop retry
			// 	}
			// } else {
			// 	err = fmt.Errorf("forced error on retry %d", retry+1)
			// }
			err = u.userRepository.BatchUpdateUser(ctx, batch, workerID, int(batchID))
			if err == nil {
				break // Jika berhasil, keluar dari loop retry
			}
			logger.WarnfCtx(ctx, "[Worker %d] Retry %d for batch %d failed: %v", workerID, retry+1, batchID, err)
			time.Sleep(time.Second) // Delay antara retries
		}

		duration := time.Since(startTime)

		if err != nil {
			logger.ErrorfCtx(ctx, "[Worker %d] ❌ Batch %d failed after %v", workerID, batchID, duration)
			logger.ErrorfCtx(ctx, "[Worker %d] ❌ Error: %v", workerID, err)
		} else {
			countBatch := int64(len(batch))
			totalUpdated += countBatch
			// logger.InfofCtx(ctx, "[Worker %d] ✅ Batch %d completed in %v (Updated %d users)", workerID, batchID, duration, countBatch)
		}
	}

	logger.InfofCtx(ctx, "[Worker %d] 🏁 Done (Total updated: %d)", workerID, totalUpdated)

	resultChan <- totalUpdated
}

func (u *userService) BatchUpdateUser(ctx context.Context, users []*model.UpdateUserRequest, batchSize int, workerCount int, maxRetries int) {
	jobs := make(chan []*model.UpdateUserRequest, len(users)/batchSize+1)
	resultChan := make(chan int64, workerCount)
	var batchCounter int32

	// Start worker goroutines
	var wg sync.WaitGroup
	for workerID := 0; workerID < workerCount; workerID++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			u.workerBatchUpdateUser(ctx, workerID, jobs, maxRetries, &batchCounter, resultChan)
		}(workerID + 1)
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

	// count totalUpdate from all workers
	totalUpdated := int64(0)
	for i := 0; i < workerCount; i++ {
		totalUpdated += <-resultChan
	}

	logger.InfofCtx(ctx, "Total updated: %d", totalUpdated)
}
