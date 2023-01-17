package repository

import (
	"time"

	"github.com/mas2401master/go-articles-api-training/internal/pkg/db"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/dto"
	models "github.com/mas2401master/go-articles-api-training/internal/pkg/models"
)

func FindAllOrder(filter dto.OrderFilter) ([]models.Order, error) {

	var orders []models.Order
	rows, err := db.Connec.Query(getOrderQuery(filter))
	if err != nil {
		return orders, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, order_number, user_id, promotion_id, quantity uint64
		var subtotal, total_discount, total, discount float64
		var status, firstname, lastname, code, name, username string
		var created_at, updated_at time.Time

		err := rows.Scan(&id, &order_number, &user_id, &promotion_id, &subtotal, &total_discount, &total, &quantity, &status, &created_at, &updated_at, &code, &name, &discount, &username, &firstname, &lastname)
		if err != nil {
			return orders, err
		}
		detail, _ := FindAllOrderItem(id)
		newOrder := models.Order{
			ID:            id,
			OrderNumber:   order_number,
			UserID:        user_id,
			PromotionID:   promotion_id,
			Subtotal:      subtotal,
			TotalDiscount: total_discount,
			Total:         total,
			Quantity:      quantity,
			Status:        status,
			CreatedAt:     created_at,
			UpdatedAt:     updated_at,
			Promotion:     dto.PromotionDTOCreate{Name: name, Code: code, Discount: discount},
			UserOrder:     dto.UserDTOOrder{Username: username, Firstname: firstname, Lastname: lastname},
			Details:       detail,
		}

		orders = append(orders, newOrder)
	}
	return orders, nil
}

func AddOrder(newOrder models.Order, orderDTO dto.OrderDTOCreate) error {
	_, err := db.Connec.Exec(queryCreateOrder, newOrder.UserID, newOrder.PromotionID, newOrder.Subtotal, newOrder.TotalDiscount, newOrder.Total, newOrder.Quantity, newOrder.Status, time.Now(), time.Now())
	if err != nil {
		return err
	}
	ordencreate, _ := FindByOrderUserStatus(newOrder.UserID, newOrder.Status)

	orderItem := models.OrderItems{
		ItemID:   orderDTO.ItemID,
		OrderID:  ordencreate.ID,
		Price:    orderDTO.Price,
		Quantity: orderDTO.Quantity,
		Total:    newOrder.Subtotal,
		Status:   newOrder.Status,
	}
	err = AddOrderItem(orderItem)
	if err != nil {
		return err
	}
	if newOrder.PromotionID != 0 {
		err = UpdateUsedPromotion(newOrder.PromotionID, true)
		if err != nil {
			return err
		}
	}

	return nil
}

func FindByOrderId(id uint64) (*models.Order, error) {
	row, err := db.Connec.Query(queryOrderId, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if row.Next() {
		var id, order_number, user_id, promotion_id, quantity uint64
		var subtotal, total_discount, total, discount float64
		var status, firstname, lastname, code, name, username string
		var created_at, updated_at time.Time

		err := row.Scan(&id, &order_number, &user_id, &promotion_id, &subtotal, &total_discount, &total, &quantity, &status, &created_at, &updated_at, &code, &name, &discount, &username, &firstname, &lastname)
		if err != nil {
			return nil, err
		}
		var detail []models.OrderItems
		detail, _ = FindAllOrderItem(id)

		newOrder := &models.Order{
			ID:            id,
			OrderNumber:   order_number,
			UserID:        user_id,
			PromotionID:   promotion_id,
			Subtotal:      subtotal,
			TotalDiscount: total_discount,
			Total:         total,
			Quantity:      quantity,
			Status:        status,
			CreatedAt:     created_at,
			UpdatedAt:     updated_at,
			Promotion:     dto.PromotionDTOCreate{Name: name, Code: code, Discount: discount},
			UserOrder:     dto.UserDTOOrder{Username: username, Firstname: firstname, Lastname: lastname},
			Details:       detail,
		}
		return newOrder, nil
	}
	return nil, nil
}

func FindByOrderUserStatus(iduser uint64, status string) (*models.Order, error) {

	row, err := db.Connec.Query(queryOrderUserStatus, iduser, status)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if row.Next() {
		var id, order_number, user_id, promotion_id, quantity uint64
		var subtotal, total_discount, total float64
		var status string
		var created_at, updated_at time.Time

		err := row.Scan(&id, &order_number, &user_id, &promotion_id, &subtotal, &total_discount, &total, &quantity, &status, &created_at, &updated_at)
		if err != nil {
			return nil, err
		}
		newOrder := &models.Order{
			ID:            id,
			OrderNumber:   order_number,
			UserID:        user_id,
			PromotionID:   promotion_id,
			Subtotal:      subtotal,
			TotalDiscount: total_discount,
			Total:         total,
			Quantity:      quantity,
			Status:        status,
			CreatedAt:     created_at,
			UpdatedAt:     updated_at,
		}

		return newOrder, nil
	}
	return nil, nil
}

func FindByOrderUser(iduser uint64) (*models.Order, error) {
	row, err := db.Connec.Query(queryOrderUser, iduser)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if row.Next() {
		var id, order_number, user_id, promotion_id, quantity uint64
		var subtotal, total_discount, total float64
		var status string
		var created_at, updated_at time.Time

		err := row.Scan(&id, &order_number, &user_id, &promotion_id, &subtotal, &total_discount, &total, &quantity, &status, &created_at, &updated_at)
		if err != nil {
			return nil, err
		}
		newOrder := &models.Order{
			ID:            id,
			OrderNumber:   order_number,
			UserID:        user_id,
			PromotionID:   promotion_id,
			Subtotal:      subtotal,
			TotalDiscount: total_discount,
			Total:         total,
			Quantity:      quantity,
			Status:        status,
			CreatedAt:     created_at,
			UpdatedAt:     updated_at,
		}

		return newOrder, nil
	}
	return nil, nil
}

func UpdateOrder(id, promotionID uint64, discount float64, status string) (*models.Order, error) {
	orderOld, err := FindByOrderId(id)
	if err != nil {
		return nil, err
	}

	if status != "" && status != orderOld.Status {
		_, err = db.Connec.Exec(queryUpdateStatus, status, time.Now(), id)
		if err != nil {
			return nil, err
		}
		_, err = UpdateOrderItemSts(id, status)
		if err != nil {
			return nil, err
		}
	}

	if promotionID != orderOld.PromotionID {
		total_discount := orderOld.Subtotal * discount / 100
		total := orderOld.Subtotal - total_discount
		_, err = db.Connec.Exec(queryUpdatePromoOrder, promotionID, total_discount, total, time.Now(), id)
		if err != nil {
			return nil, err
		}
		if orderOld.PromotionID != 0 {
			_ = UpdateUsedPromotion(orderOld.PromotionID, false)
		}
		if promotionID != 0 {
			_ = UpdateUsedPromotion(promotionID, true)
		}
	}

	orderNew, _ := FindByOrderId(id)

	return orderNew, nil
}

func UpdateOrderAmount(id, quantity uint64, subtotal float64) (*models.Order, error) {
	var total_discount, total float64
	orderOld, err := FindByOrderId(id)
	if err != nil {
		return nil, err
	}
	total_discount = subtotal * orderOld.Promotion.Discount / 100
	total = subtotal - total_discount
	_, err = db.Connec.Exec(queryUpdateOrderAmount, quantity, subtotal, total_discount, total, time.Now(), id)
	if err != nil {
		return nil, err
	}

	orderNew, _ := FindByOrderId(id)

	return orderNew, nil
}

func DeleteOrder(id uint64) error {
	_, err := db.Connec.Exec(queryDeleteOrder, id)
	if err != nil {
		return err
	}
	return nil
}
