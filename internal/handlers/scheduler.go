package handlers

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

// Task represents a function that will be executed on a schedule
type Task func(ctx context.Context) error

// Scheduler manages scheduled tasks
type Scheduler struct {
	cron *cron.Cron
}

// NewScheduler creates a new scheduler instance
func NewScheduler() *Scheduler {
	// Create a new cron instance with seconds field enabled
	c := cron.New(cron.WithSeconds())
	
	return &Scheduler{
		cron: c,
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
	log.Println("Scheduler started")
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	log.Println("Scheduler stopped")
}

// AddHourlyTask adds a task to be executed every hour
func (s *Scheduler) AddHourlyTask(taskName string, task Task) {
	_, err := s.cron.AddFunc("0 0 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)
		defer cancel()
		
		log.Printf("Running task: %s", taskName)
		start := time.Now()
		
		if err := task(ctx); err != nil {
			log.Printf("Task %s failed: %v", taskName, err)
		} else {
			log.Printf("Task %s completed in %v", taskName, time.Since(start))
		}
	})
	
	if err != nil {
		log.Printf("Failed to schedule task %s: %v", taskName, err)
	} else {
		log.Printf("Task %s scheduled to run hourly", taskName)
	}
}

// AddCustomScheduleTask adds a task with a custom cron schedule
func (s *Scheduler) AddCustomScheduleTask(taskName string, cronSchedule string, task Task) {
	_, err := s.cron.AddFunc(cronSchedule, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)
		defer cancel()
		
		log.Printf("Running task: %s", taskName)
		start := time.Now()
		
		if err := task(ctx); err != nil {
			log.Printf("Task %s failed: %v", taskName, err)
		} else {
			log.Printf("Task %s completed in %v", taskName, time.Since(start))
		}
	})
	
	if err != nil {
		log.Printf("Failed to schedule task %s with schedule %s: %v", taskName, cronSchedule, err)
	} else {
		log.Printf("Task %s scheduled with custom schedule: %s", taskName, cronSchedule)
	}
} 