package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"pocketvue/helpers"
	"pocketvue/services"

	"github.com/pocketbase/pocketbase/core"
)

// CreateCheckoutRequest represents the request body for creating a checkout session
type CreateCheckoutRequest struct {
	Products      []string `json:"products"`
	WorkspaceSlug string   `json:"workspace_slug"` // Optional: for workspace-specific success URL
	ReturnPath    string   `json:"return_path"`    // Optional: custom return path (defaults to /dashboard)
}

// CreateCheckoutResponse represents the response for a successful checkout creation
type CreateCheckoutResponse struct {
	URL string `json:"url"`
}

// CreateCustomerPortalRequest represents the request body for creating a customer portal session
type CreateCustomerPortalRequest struct {
	WorkspaceSlug string `json:"workspace_slug"` // Optional: for workspace-specific return URL
	ReturnPath    string `json:"return_path"`    // Optional: custom return path (defaults to /dashboard/settings/billing)
}

// CreateCustomerPortalResponse represents the response for a successful customer portal creation
type CreateCustomerPortalResponse struct {
	URL string `json:"url"`
}

// CreateCheckoutSession creates a Polar checkout session for the authenticated user
func CreateCheckoutSession(e *core.RequestEvent) error {
	// Get authenticated user
	user, err := helpers.GetAuthenticatedUser(e)
	if err != nil {
		return err
	}

	// Parse request body
	var req CreateCheckoutRequest
	if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
		log.Printf("Error parsing checkout request: %v", err)
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// Validate required fields
	if len(req.Products) == 0 {
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "products field is required and must contain at least one product ID",
		})
	}

	// Get URLs from environment variables
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Printf("Warning: FRONTEND_URL not set, using default")
		frontendURL = "http://localhost:3000"
	}

	// Build workspace-specific URLs if workspace_slug is provided
	var successURL, returnURL string
	if req.WorkspaceSlug != "" {
		// Determine the return path (defaults to /dashboard)
		returnPath := "/dashboard"
		if req.ReturnPath != "" {
			returnPath = req.ReturnPath
		}

		successURL = frontendURL + "/" + req.WorkspaceSlug + returnPath + "?checkout=success"
		returnURL = frontendURL + "/" + req.WorkspaceSlug + returnPath
	} else {
		successURL = frontendURL + "/checkout/success"
		returnURL = frontendURL + "/dashboard"
	}

	// Extract user information
	userID := user.GetString("id")
	userEmail := user.GetString("email")
	userName := user.GetString("name")

	log.Printf("CreateCheckoutSession called by user: ID=%s, Email=%s, Products=%v",
		userID, userEmail, req.Products)

	// Create Polar service and checkout session
	polarService := services.NewPolarService()
	checkoutURL, err := polarService.CreateCheckoutSession(
		req.Products,
		successURL,
		returnURL,
		userID,
		userEmail,
		userName,
	)

	if err != nil {
		log.Printf("Error creating checkout session for user %s: %v", userID, err)
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create checkout session",
		})
	}

	// Return checkout URL
	return e.JSON(http.StatusOK, CreateCheckoutResponse{
		URL: checkoutURL,
	})
}

// CreateCustomerPortalSession creates a Polar customer portal session for the authenticated user
func CreateCustomerPortalSession(e *core.RequestEvent) error {
	// Get authenticated user
	user, err := helpers.GetAuthenticatedUser(e)
	if err != nil {
		return err
	}

	// Parse request body
	var req CreateCustomerPortalRequest
	if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
		log.Printf("Error parsing customer portal request: %v", err)
		return e.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	// Get frontend URL from environment
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		log.Printf("Warning: FRONTEND_URL not set, using default")
		frontendURL = "http://localhost:3000"
	}

	// Build workspace-specific return URL if workspace_slug is provided
	var returnURL string
	if req.WorkspaceSlug != "" {
		// Determine the return path (defaults to /dashboard/settings/billing)
		returnPath := "/dashboard/settings/billing"
		if req.ReturnPath != "" {
			returnPath = req.ReturnPath
		}

		returnURL = frontendURL + "/" + req.WorkspaceSlug + returnPath
	} else {
		returnURL = frontendURL + "/dashboard/settings/billing"
	}

	// Extract user ID
	userID := user.GetString("id")

	log.Printf("CreateCustomerPortalSession called by user: ID=%s", userID)

	// Create Polar service and customer portal session
	polarService := services.NewPolarService()
	portalURL, err := polarService.CreateCustomerSession(userID, returnURL)

	if err != nil {
		log.Printf("Error creating customer portal session for user %s: %v", userID, err)
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to create customer portal session",
		})
	}

	// Return portal URL
	return e.JSON(http.StatusOK, CreateCustomerPortalResponse{
		URL: portalURL,
	})
}
