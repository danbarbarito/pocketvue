package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3395098727")
		if err != nil {
			return err
		}

		// add field
		if err := collection.Fields.AddMarshaledJSONAt(4, []byte(`{
			"convertURLs": false,
			"hidden": false,
			"id": "editor4274335913",
			"maxSize": 0,
			"name": "content",
			"presentable": false,
			"required": true,
			"system": false,
			"type": "editor"
		}`)); err != nil {
			return err
		}

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("pbc_3395098727")
		if err != nil {
			return err
		}

		// remove field
		collection.Fields.RemoveById("editor4274335913")

		return app.Save(collection)
	})
}
