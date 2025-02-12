package scheduler

import (
	"context"
	"fmt"
	"go-scheduler/config"
	"go-scheduler/internal/repository"
	"go-scheduler/internal/service"
	"go-scheduler/model"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	config      *config.Config
	scheduler   *cron.Cron
	userService service.UserServiceInterface
}

func NewScheduler(config *config.Config, scheduler *cron.Cron) *Scheduler {
	pgsql := repository.NewPostgreConnection(config)
	userRepo := repository.NewUserRepository(pgsql)

	return &Scheduler{
		config:      config,
		scheduler:   scheduler,
		userService: service.NewUserService(config, userRepo),
	}
}

func (s *Scheduler) Start() {
	// every day at 3:14 AM
	s.scheduler.AddFunc("14 3 * * *", s.HXMSGetUsers)

	s.scheduler.Start()
}

func (s *Scheduler) HXMSGetUsers() {
	requestid := config.GenerateRequestID()
	ctx := context.WithValue(context.Background(), config.RequestIDKey, requestid)

	currentTime := time.Now()
	// get users from hxms
	// dummy
	getUsers, err := s.userService.GetUsers(ctx)
	if err != nil {
		panic(err)
	}

	users := []*model.UpdateUserRequest{}
	for _, user := range getUsers {
		index := strings.Index(user.FullName, " |")
		if index != -1 {
			user.FullName = user.FullName[:index]
		}
		users = append(users, &model.UpdateUserRequest{
			Username: user.Username,
			FullName: fmt.Sprintf(user.FullName+" | %s", currentTime.Format(time.DateTime)),
		})
	}

	if len(users) > 0 {
		batchSize := 500 // 500 row per batch
		workerCount := 3 // 3 worker paralel
		maxRetries := 3  // Retry maksimal 3 kali jika error

		startTime := time.Now()
		s.userService.BatchUpdateUser(ctx, users, batchSize, workerCount, maxRetries)

		fmt.Printf("Total execution time: %v\n", time.Since(startTime))
	}
}
