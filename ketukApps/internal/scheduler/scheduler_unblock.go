package scheduler

import (
	"fmt"
	"ketukApps/internal/models"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func (s *Scheduler) RegisterUnblockJob() error {
	if s.Client == nil {
		return fmt.Errorf("scheduler not initialized")
	}

	_, err := s.Client.NewJob(
		gocron.CronJob(
			"*/1 * * * *",
			false,
		),
		gocron.NewTask(
			func() {
				s.unblockUsersTask()
			},
		),
	)
	if err != nil {
		return fmt.Errorf("failed to register unblock job: %w", err)
	}
	log.Println("Unblock job registered to run daily at midnight")
	return nil
}

func (s *Scheduler) unblockUsersTask() {
	log.Println("Running unblock users task...")
	var unblockings []models.Unblocking
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Printf("Failed to load location: %v\n", err)
		return
	}
	time.Local = loc
	s.db.Where("start_date <= ?", time.Now().Local()).Where("end_date >= ?", time.Now().Local()).Find(&unblockings)
	log.Printf("Found %d unblockings to process.\n", len(unblockings))
	if len(unblockings) > 1 {
		log.Println("Multiple unblockings found, skipping to avoid conflicts.")
		DisableUnblock()
	} else if len(unblockings) == 1 {
		for _, unblocking := range unblockings {
			log.Printf("Processing unblocking : %v\n", unblocking)
			EnableUnblock()
		}
	} else {
		log.Println("No unblockings to process at this time.")
		DisableUnblock()
	}
}
