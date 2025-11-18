package helpers

import (
	"errors"
	"pocketvue/config"
	"strings"
)

// BuildFrontendURL constructs a frontend URL with the given path
func BuildFrontendURL(path string) string {
	baseURL := config.FrontendURL
	// Ensure baseURL doesn't end with /
	if len(baseURL) > 0 && baseURL[len(baseURL)-1] == '/' {
		baseURL = baseURL[:len(baseURL)-1]
	}
	// Ensure path starts with /
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return baseURL + path
}

// BuildWorkspaceURL constructs a workspace-specific frontend URL
func BuildWorkspaceURL(workspaceSlug, path string) string {
	// Ensure path starts with /
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return BuildFrontendURL("/" + workspaceSlug + path)
}

// BuildCheckoutSuccessURL constructs a checkout success URL
func BuildCheckoutSuccessURL(workspaceSlug, returnPath string) string {
	if workspaceSlug != "" {
		if returnPath == "" {
			returnPath = "/dashboard"
		}
		return BuildWorkspaceURL(workspaceSlug, returnPath+"?checkout=success")
	}
	return BuildFrontendURL("/checkout/success")
}

// BuildCheckoutReturnURL constructs a checkout return URL
func BuildCheckoutReturnURL(workspaceSlug, returnPath string) string {
	if workspaceSlug != "" {
		if returnPath == "" {
			returnPath = "/dashboard"
		}
		return BuildWorkspaceURL(workspaceSlug, returnPath)
	}
	return BuildFrontendURL("/dashboard")
}

// BuildCustomerPortalReturnURL constructs a customer portal return URL
func BuildCustomerPortalReturnURL(workspaceSlug, returnPath string) string {
	if workspaceSlug != "" {
		if returnPath == "" {
			returnPath = "/dashboard/settings/billing"
		}
		return BuildWorkspaceURL(workspaceSlug, returnPath)
	}
	return BuildFrontendURL("/dashboard/settings/billing")
}

// ValidateFrontendURL checks if FrontendURL is properly configured
// Returns an error if FrontendURL is empty or is the default localhost value in production
func ValidateFrontendURL() error {
	if config.FrontendURL == "" {
		return errors.New("FRONTEND_URL environment variable is not set")
	}

	// In production, don't allow the default localhost URL
	if config.AppEnv == "production" || config.PolarEnvironment == "production" {
		defaultURL := "http://localhost:3000"
		if strings.TrimSpace(config.FrontendURL) == defaultURL {
			return errors.New("FRONTEND_URL must be set to a production URL (cannot use default localhost:3000)")
		}
	}

	return nil
}
