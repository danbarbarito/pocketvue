package routes

import (
	"net/http"
	"pocketvue/constants"
	"pocketvue/helpers"

	"github.com/pocketbase/pocketbase/core"
)

// GetProducts returns all non-archived products
func GetProducts(e *core.RequestEvent) error {
	// Fetch all non-archived products, ordered by created date
	records, err := helpers.FindAllRecords(e.App, constants.CollectionPolarProducts)
	if err != nil {
		return e.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to fetch products",
		})
	}

	// Filter out archived products
	var activeProducts []map[string]interface{}
	for _, record := range records {
		if !record.GetBool("is_archived") {
			productMap := map[string]interface{}{
				"id":                       record.GetString("id"),
				"name":                     record.GetString("name"),
				"description":              record.GetString("description"),
				"price_amount":             record.GetInt("price_amount"),
				"price_currency":           record.GetString("price_currency"),
				"recurring_interval":       record.GetString("recurring_interval"),
				"recurring_interval_count": record.GetInt("recurring_interval_count"),
				"is_recurring":             record.GetBool("is_recurring"),
				"trial_interval":           record.GetString("trial_interval"),
				"trial_interval_count":     record.GetInt("trial_interval_count"),
				"polar_price_id":           record.GetString("polar_price_id"),
			}
			activeProducts = append(activeProducts, productMap)
		}
	}

	return e.JSON(http.StatusOK, activeProducts)
}
