package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/novintriantonius/cosmos-validator-service/internal/handlers"
	"github.com/novintriantonius/cosmos-validator-service/internal/services"
	"github.com/novintriantonius/cosmos-validator-service/internal/store"
	"github.com/novintriantonius/cosmos-validator-service/internal/tasks"
	"github.com/robfig/cron/v3"
)

// RegisterDelegationTasks registers all delegation-related tasks with the scheduler
func RegisterDelegationTasks(
	sched *handlers.Scheduler,
	validatorStore store.ValidatorStore,
	delegationStore store.DelegationStore,
	cosmosService *services.CosmosService,
) {
	// Initialize the delegation sync task
	delegationSyncTask := tasks.NewDelegationSyncTask(
		validatorStore,
		delegationStore,
		cosmosService,
	)
	
	// Schedule delegation sync task to run hourly at the start of each hour (minute 0)
	// Cron format: second minute hour day month weekday
	sched.AddCustomScheduleTask(
		"hourly-validator-delegations-sync",
		"0 * * * * *", // Run at minute 0 of every hour
		func(ctx context.Context) error {
			return delegationSyncTask.SyncEnabledValidatorDelegations(ctx)
		},
	)
	
	log.Println("Scheduled hourly delegation sync task to run at the start of every hour")
	
	// Run the task once immediately on startup to populate initial data
	go func() {
		// Wait a few seconds to allow the server to start properly
		time.Sleep(5 * time.Second)
		
		log.Println("Running initial delegation sync...")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		
		if err := delegationSyncTask.SyncEnabledValidatorDelegations(ctx); err != nil {
			log.Printf("Initial delegation sync completed with errors: %v", err)
		} else {
			log.Println("Initial delegation sync completed successfully")
		}
	}()
}

// ScheduleDelegationSync schedules the delegation sync task
func ScheduleDelegationSync(delegationSyncTask *tasks.DelegationSyncTask) error {
	c := cron.New()

	// Schedule task to run at the start of every hour
	_, err := c.AddFunc("0 * * * *", func() {
		log.Printf("Running hourly delegation sync...")
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := delegationSyncTask.SyncEnabledValidatorDelegations(ctx); err != nil {
			log.Printf("Error syncing delegations: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("error scheduling delegation sync: %v", err)
	}

	// Also run immediately on startup
	log.Printf("Running initial delegation sync...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	if err := delegationSyncTask.SyncEnabledValidatorDelegations(ctx); err != nil {
		log.Printf("Error running initial delegation sync: %v", err)
	}

	// Start the scheduler in a goroutine
	go c.Start()

	log.Printf("Scheduled hourly delegation sync task to run at the start of every hour")
	return nil
} 