package scheduler

import (
	"log"

	"github.com/go-co-op/gocron/v2"
	"gorm.io/gorm"
)



type Scheduler struct {
	Client gocron.Scheduler
	db	 *gorm.DB
}

func NewScheduler(db *gorm.DB) (*Scheduler, error){
	schedulerClient, err := gocron.NewScheduler()
	if err != nil {
		log.Panicf("Failed to create scheduler: %s", err)
	}
	return &Scheduler{
		Client: schedulerClient,
		db:     db,
	}, nil
}

func (s *Scheduler) Start(){
	if s.Client != nil {
		s.Client.Start()
		log.Println("Scheduler started")
	}
}

func (s *Scheduler) Shutdown(){
	if s.Client != nil {
		s.Client.Shutdown()
		log.Println("Scheduler stopped")
	}
}