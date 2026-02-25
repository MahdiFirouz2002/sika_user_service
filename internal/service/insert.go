package service

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sika/internal/domain"
	"sync"
)

type producer struct {
	userService UserService
	workerCount int
}

func NewProducer(userSrv UserService, wCount int) *producer {
	return &producer{
		userService: userSrv,
		workerCount: wCount,
	}
}

func (p *producer) RunInsert(ctx context.Context) error {
	// read file
	jsonFile, err := os.Open("users_data.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}

	// unmarshal on struct
	var users []domain.User
	if err := json.Unmarshal(bytes, &users); err != nil {
		return err
	}

	// insert users into channel
	userChan := make(chan domain.User, len(users))
	for _, user := range users {
		userChan <- user
	}

	close(userChan)

	// start workers
	var wg sync.WaitGroup
	for range p.workerCount {
		wg.Add(1)
		go p.worker(ctx, &wg, userChan)
	}

	wg.Wait()

	return nil
}

func (p *producer) worker(ctx context.Context, wg *sync.WaitGroup, usersChan chan domain.User) {
	defer wg.Done()

	for user := range usersChan {
		p.userService.Create(ctx, &user)
	}
}
