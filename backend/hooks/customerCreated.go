package hooks

import (
	"log"

	"pocketvue/services"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// RegisterUserCreatedHook registers a hook that triggers after successful user creation
func RegisterUserCreatedHook(app *pocketbase.PocketBase) {
	polarService := services.NewPolarService()

	app.OnRecordAfterCreateSuccess("users").BindFunc(func(e *core.RecordEvent) error {
		log.Printf("New user created: %s", e.Record.Id)

		// Create Polar customer asynchronously to avoid blocking user registration
		polarService.CreateCustomerAsync(app, e.Record.Id, e.Record.GetString("email"), e.Record.GetString("name"))

		return e.Next()
	})
}
