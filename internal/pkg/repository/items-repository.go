package repository

import (
	"time"

	"github.com/mas2401master/go-articles-api-training/internal/pkg/db"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/dto"
	models "github.com/mas2401master/go-articles-api-training/internal/pkg/models"
)

func FindAllItem(filter dto.ItemFilter) ([]models.Item, error) {
	var items []models.Item
	rows, err := db.Connec.Query(getItemQuery(filter))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id uint64
		var price float64
		var name, description string
		var status bool

		var available bool
		var created_at, updated_at time.Time
		//id, name, description, price, available, created_at, updated_at,status
		err := rows.Scan(&id, &name, &description, &price, &available, &created_at, &updated_at, &status)
		if err != nil {
			return nil, err
		}
		newItem := models.Item{
			ID:          id,
			NameItem:    name,
			Description: description,
			Price:       price,
			Available:   available,
			Status:      status,
			CreatedAt:   created_at,
			UpdatedAt:   updated_at,
		}

		items = append(items, newItem)
	}
	return items, nil
}

func AddItem(newItem models.Item) error {
	_, err := db.Connec.Exec(queryCreateItem, newItem.NameItem, newItem.Description, newItem.Price, true, time.Now(), time.Now(), true)
	if err != nil {
		return err
	}
	return nil
}

func FindByItemId(id uint64) (*models.Item, error) {
	var ItemResponse *models.Item

	row, err := db.Connec.Query(queryItemId, id)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	if row.Next() {
		var id uint64
		var price float64
		var name, description string
		var available, status bool
		var created_at, updated_at time.Time
		err := row.Scan(&id, &name, &description, &price, &available, &created_at, &updated_at, &status)
		if err != nil {
			return nil, err
		}

		ItemResponse = &models.Item{
			ID:          id,
			NameItem:    name,
			Description: description,
			Price:       price,
			Available:   available,
			Status:      status,
			CreatedAt:   created_at,
			UpdatedAt:   updated_at,
		}

		return ItemResponse, nil
	}
	return nil, nil
}

func UpdateItem(id uint64, itemNew dto.ItemDTOUpdate) (*models.Item, error) {
	var price float64
	var name, description string
	var available, status bool
	var itemUpdate *models.Item

	itemOld, err := FindByItemId(id)
	if err != nil {
		return nil, err
	}

	price = itemOld.Price
	if itemNew.Price != 0 && itemOld.Price != itemNew.Price {
		price = itemNew.Price
	}

	name = itemOld.NameItem
	if itemNew.NameItem != "" && itemOld.NameItem != itemNew.NameItem {
		name = itemNew.NameItem
	}

	description = itemOld.NameItem
	if itemNew.Description != "" && itemOld.Description != itemNew.Description {
		description = itemNew.Description
	}

	available = itemOld.Available
	if itemNew.Available != itemOld.Available {
		available = itemNew.Available
	}

	status = itemOld.Status
	if itemNew.Status != itemOld.Status {
		status = itemNew.Status

	}
	_, err = db.Connec.Exec(queryUpdateItem, name, description, price, available, time.Now(), status, itemOld.ID)
	if err != nil {
		return nil, err
	}

	itemUpdate = &models.Item{
		ID:          itemOld.ID,
		NameItem:    name,
		Description: description,
		Price:       price,
		Available:   available,
		Status:      status,
		CreatedAt:   itemOld.CreatedAt,
		UpdatedAt:   time.Now(),
	}
	return itemUpdate, nil
}

func DeleteItem(id uint64) error {
	_, err := db.Connec.Exec(queryDeleteItem, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateAvailableItem(id uint64, available bool) error {
	_, err := db.Connec.Exec(queryUpdateAvailableItem, available, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}
