package repository

import (
	"time"

	"github.com/mas2401master/go-articles-api-training/internal/pkg/db"
	models "github.com/mas2401master/go-articles-api-training/internal/pkg/models"
)

func AddOrderItem(newOrder models.OrderItems) error {
	_, err := db.Connec.Exec(queryCreateOrderItem, newOrder.OrderID, newOrder.ItemID, newOrder.Price, newOrder.Quantity, newOrder.Total, newOrder.Status, time.Now(), time.Now())
	if err != nil {
		return err
	}
	_ = UpdateAvailableItem(newOrder.ItemID, false)
	subtotal, _ := TotalOrderItem(newOrder.OrderID)
	quantity, _ := TotalQuantityOrderItem(newOrder.OrderID)
	_, err = UpdateOrderAmount(newOrder.OrderID, quantity, subtotal)
	if err != nil {
		return err
	}
	return nil
}

func FindAllOrderItem(orderID uint64) ([]models.OrderItems, error) {
	var orderitems []models.OrderItems
	rows, err := db.Connec.Query(queryFindAllOrderItem, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, order_id, order_number, item_id, quantity uint64
		var price, total float64
		var status, name string
		var created_at, updated_at time.Time
		err := rows.Scan(&id, &order_id, &item_id, &price, &quantity, &total, &status, &created_at, &updated_at, &name)
		if err != nil {
			return nil, err
		}

		newOrderItem := models.OrderItems{
			ID:        id,
			OrderID:   order_number,
			ItemID:    item_id,
			NameItem:  name,
			Price:     price,
			Quantity:  quantity,
			Total:     total,
			Status:    status,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
		}
		orderitems = append(orderitems, newOrderItem)
	}
	return orderitems, nil
}

func FindByOrderItemId(id uint64) (*models.OrderItems, error) {
	row, err := db.Connec.Query(queryOrderItemId, id)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	if row.Next() {
		var id, order_id, item_id, quantity uint64
		var price, total float64
		var status, name string
		var created_at, updated_at time.Time
		/* o.id id, order_id, item_id, o.price price, quantity, total, o.status status, o.created_at created_at,
		o.updated_at updated_at,name*/
		err := row.Scan(&id, &order_id, &item_id, &price, &quantity, &total, &status, &created_at, &updated_at, &name)
		if err != nil {
			return nil, err
		}
		newOrder := &models.OrderItems{
			ID:        id,
			OrderID:   order_id,
			ItemID:    item_id,
			NameItem:  name,
			Price:     price,
			Quantity:  quantity,
			Total:     total,
			Status:    status,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
		}

		return newOrder, nil
	}
	return nil, nil
}

func TotalQuantityOrderItem(orderid uint64) (uint64, error) {
	var quantity uint64
	row, err := db.Connec.Query(queryTotalQuantityOrderItem, orderid)
	if err != nil {
		return quantity, err
	}
	defer row.Close()

	if row.Next() {
		err = row.Scan(&quantity)
		if err != nil {
			return quantity, err
		}
	}
	return quantity, nil
}

func TotalOrderItem(idorden uint64) (float64, error) {
	var total float64
	row, err := db.Connec.Query(queryTotalOrderItem, idorden)
	if err != nil {

		return total, err
	}
	if row.Next() {
		err = row.Scan(&total)
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

func GetQuantityOrderItem(idorden, iditem uint64) (uint64, error) {
	var total uint64
	row, err := db.Connec.Query(queryQuantityOrderItem, idorden, iditem)
	if err != nil {
		return total, err
	}
	err = row.Scan(&total)
	if err != nil {
		return total, err
	}
	return total, nil
}

func UpdateOrderItem(orderItem models.OrderItems) (*models.OrderItems, error) {
	_, err := db.Connec.Exec(queryUpdateOrderItem, orderItem.Quantity, orderItem.Total, time.Now(), orderItem.ID)
	if err != nil {
		return nil, err
	}
	subtotal, _ := TotalOrderItem(orderItem.OrderID)
	quantity, _ := TotalQuantityOrderItem(orderItem.OrderID)

	_, err = UpdateOrderAmount(orderItem.OrderID, quantity, subtotal)
	if err != nil {
		return nil, err
	}

	orderItemNew, _ := FindByOrderItemId(orderItem.ID)
	return orderItemNew, nil
}

func UpdateOrderItemSts(orderid uint64, status string) (bool, error) {
	_, err := db.Connec.Exec(queryUpdateOrderItemSts, status, orderid)
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteOrderItem(id, orderid uint64) error {
	_, err := db.Connec.Exec(queryDeleteOrderItem, id)
	if err != nil {
		return err
	}
	_ = UpdateAvailableItem(orderid, true)
	subtotal, _ := TotalOrderItem(orderid)
	quantity, _ := TotalQuantityOrderItem(orderid)

	_, err = UpdateOrderAmount(orderid, quantity, subtotal)
	if err != nil {
		return err
	}
	return nil
}

func ReleaseItemOrder(orderid uint64) (bool, error) {
	_, err := db.Connec.Exec(queryRealeaseItemOrder, orderid)
	if err != nil {
		return false, err
	}
	return true, nil
}
