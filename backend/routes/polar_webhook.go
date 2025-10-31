package routes

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"pocketvue/helpers"
	"pocketvue/services"
	"pocketvue/types"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// HandlePolarWebhook handles incoming Polar webhook events
func HandlePolarWebhook(e *core.RequestEvent) error {
	// Read the raw request body (needed for signature verification)
	body, err := io.ReadAll(e.Request.Body)
	if err != nil {
		log.Printf("Error reading webhook body: %v", err)
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "failed to read request body",
		})
	}

	// Extract headers for signature verification
	// Try both lowercase and canonical forms
	headers := map[string]string{
		"webhook-id":        e.Request.Header.Get("webhook-id"),
		"webhook-timestamp": e.Request.Header.Get("webhook-timestamp"),
		"webhook-signature": e.Request.Header.Get("webhook-signature"),
	}

	// If not found, try canonical form (Webhook-Id, etc.)
	if headers["webhook-id"] == "" {
		headers["webhook-id"] = e.Request.Header.Get("Webhook-Id")
	}
	if headers["webhook-timestamp"] == "" {
		headers["webhook-timestamp"] = e.Request.Header.Get("Webhook-Timestamp")
	}
	if headers["webhook-signature"] == "" {
		headers["webhook-signature"] = e.Request.Header.Get("Webhook-Signature")
	}

	// Verify webhook signature - pass the entire request headers
	if err := helpers.VerifyWebhookSignature(body, e.Request.Header); err != nil {
		log.Printf("Webhook signature verification failed: %v", err)
		return e.JSON(http.StatusUnauthorized, map[string]string{
			"error": "invalid signature",
		})
	}

	// Parse the webhook event
	var event types.WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		log.Printf("Error parsing webhook event: %v", err)
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid JSON payload",
		})
	}

	log.Printf("Received webhook event: type=%s, timestamp=%s", event.Type, event.Timestamp)

	// Re-marshal the data for individual handlers
	eventData, err := json.Marshal(event.Data)
	if err != nil {
		log.Printf("Error marshaling event data: %v", err)
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to process event data",
		})
	}

	// Create webhook service
	app, ok := e.App.(*pocketbase.PocketBase)
	if !ok {
		log.Printf("Failed to cast app to PocketBase")
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	}
	webhookService := services.NewWebhookService(app)

	// Route to appropriate handler based on event type
	var handlerErr error
	switch event.Type {
	case "subscription.created":
		handlerErr = webhookService.HandleSubscriptionCreated(eventData)

	case "subscription.updated":
		handlerErr = webhookService.HandleSubscriptionUpdated(eventData)

	case "subscription.active":
		handlerErr = webhookService.HandleSubscriptionActive(eventData)

	case "subscription.canceled":
		handlerErr = webhookService.HandleSubscriptionCanceled(eventData)

	case "subscription.revoked":
		handlerErr = webhookService.HandleSubscriptionRevoked(eventData)

	case "order.created":
		handlerErr = webhookService.HandleOrderCreated(eventData)

	case "order.paid":
		handlerErr = webhookService.HandleOrderPaid(eventData)

	case "product.created":
		handlerErr = webhookService.HandleProductCreated(eventData)

	case "product.updated":
		handlerErr = webhookService.HandleProductUpdated(eventData)

	default:
		log.Printf("Unhandled webhook event type: %s", event.Type)
		// Return 200 OK for unhandled events to prevent retries
		return e.JSON(http.StatusOK, map[string]string{
			"message": "event type not handled",
		})
	}

	// Check if handler returned an error
	if handlerErr != nil {
		log.Printf("Error handling webhook event %s: %v", event.Type, handlerErr)
		// Return 500 to trigger Polar's retry mechanism
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to process webhook",
		})
	}

	// Return success response
	return e.JSON(http.StatusOK, map[string]string{
		"message": "webhook processed successfully",
	})
}
