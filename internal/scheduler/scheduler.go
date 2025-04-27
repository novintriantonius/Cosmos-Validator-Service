package scheduler

import (
	"github.com/novintriantonius/cosmos-validator-service/internal/handlers"
	"github.com/novintriantonius/cosmos-validator-service/internal/services"
	"github.com/novintriantonius/cosmos-validator-service/internal/store"
)

// SetupScheduler initializes and configures all scheduled tasks
func SetupScheduler(
	validatorStore store.ValidatorStore,
	delegationStore store.DelegationStore,
	cosmosService *services.CosmosService,
) *handlers.Scheduler {
	// Initialize scheduler
	sched := handlers.NewScheduler()
	
	// Register all tasks
	RegisterDelegationTasks(sched, validatorStore, delegationStore, cosmosService)
	
	return sched
}