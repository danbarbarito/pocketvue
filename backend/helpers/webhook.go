package helpers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	svix "github.com/standard-webhooks/standard-webhooks/libraries/go"
)

// VerifyWebhookSignature verifies the webhook signature using Standard Webhooks
func VerifyWebhookSignature(payload []byte, headers http.Header) error {
	secret := os.Getenv("POLAR_WEBHOOK_SECRET")
	if secret == "" {
		return fmt.Errorf("POLAR_WEBHOOK_SECRET not configured")
	}

	// Standard Webhooks library expects secrets to be base64 encoded with whsec_ prefix
	// If the secret doesn't have the prefix, add it and encode
	if !strings.HasPrefix(secret, "whsec_") {
		// If it's not already base64, encode it
		if _, err := base64.StdEncoding.DecodeString(secret); err != nil {
			// Secret is plain text, encode it as base64
			secret = "whsec_" + base64.StdEncoding.EncodeToString([]byte(secret))
		} else {
			// Secret is already base64, just add prefix
			secret = "whsec_" + secret
		}
	}

	// Initialize webhook verifier
	wh, err := svix.NewWebhook(secret)
	if err != nil {
		return fmt.Errorf("failed to initialize webhook verifier: %w", err)
	}

	// Extract headers - try both canonical and lowercase forms
	msgID := headers.Get("Webhook-Id")
	if msgID == "" {
		msgID = headers.Get("webhook-id")
	}

	msgTimestamp := headers.Get("Webhook-Timestamp")
	if msgTimestamp == "" {
		msgTimestamp = headers.Get("webhook-timestamp")
	}

	msgSignature := headers.Get("Webhook-Signature")
	if msgSignature == "" {
		msgSignature = headers.Get("webhook-signature")
	}

	if msgID == "" || msgTimestamp == "" || msgSignature == "" {
		return fmt.Errorf("missing required webhook headers")
	}

	// Verify the signature - pass headers directly as http.Header
	err = wh.Verify(payload, headers)
	if err != nil {
		return fmt.Errorf("webhook signature verification failed: %w", err)
	}

	// Additional timestamp validation (prevent replay attacks)
	// Parse Unix timestamp (seconds since epoch)
	var timestamp time.Time
	var timestampInt int64
	if _, err := fmt.Sscanf(msgTimestamp, "%d", &timestampInt); err == nil {
		timestamp = time.Unix(timestampInt, 0)
	} else {
		// Try parsing as RFC3339 as fallback
		timestamp, err = time.Parse(time.RFC3339, msgTimestamp)
		if err != nil {
			return fmt.Errorf("invalid timestamp format: %w", err)
		}
	}

	// Reject if timestamp is more than 5 minutes old
	if time.Since(timestamp) > 5*time.Minute {
		return fmt.Errorf("webhook timestamp too old")
	}

	return nil
}
