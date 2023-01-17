package repository

import (
	"strings"
	"time"

	"github.com/mas2401master/go-articles-api-training/internal/pkg/db"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/dto"
	models "github.com/mas2401master/go-articles-api-training/internal/pkg/models"
)

func FindAllIPromotion(filter dto.PromotionFilter) ([]models.Promotion, error) {
	var promotions []models.Promotion
	rows, err := db.Connec.Query(getPromotionQuery(filter))
	if err != nil {
		return promotions, err
	}
	defer rows.Close()

	for rows.Next() {
		var id uint64
		var discount float64
		var name, code string

		var used bool
		var created_at, updated_at time.Time
		err := rows.Scan(&id, &code, &name, &used, &discount, &created_at, &updated_at)
		if err != nil {
			return promotions, err
		}
		newPromotion := models.Promotion{
			ID:        id,
			Name:      name,
			Code:      code,
			Used:      used,
			Discount:  discount,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
		}

		promotions = append(promotions, newPromotion)
	}
	return promotions, nil
}

func AddPromotion(newPromotion models.Promotion) error {
	_, err := db.Connec.Exec(queryCreatePromotion, newPromotion.Code, newPromotion.Name, false, newPromotion.Discount, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func FindByIdPromotion(id uint64) (*models.Promotion, error) {
	row, err := db.Connec.Query(queryPromtionId, id)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	if row.Next() {
		var discount float64
		var name, code string

		var used bool
		var created_at, updated_at time.Time
		err := row.Scan(&id, &code, &name, &used, &discount, &created_at, &updated_at)
		if err != nil {
			return nil, err
		}

		PromotionResponse := &models.Promotion{
			ID:        id,
			Name:      name,
			Code:      code,
			Used:      used,
			Discount:  discount,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
		}

		return PromotionResponse, nil
	}
	return nil, nil
}

func FindByCodePromotion(code string) (*models.Promotion, error) {
	row, err := db.Connec.Query(queryPromtionCode, code)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	if row.Next() {
		var id uint64
		var discount float64
		var name, code string

		var used bool
		var created_at, updated_at time.Time
		err := row.Scan(&id, &code, &name, &used, &discount, &created_at, &updated_at)
		if err != nil {
			return nil, err
		}

		PromotionResponse := &models.Promotion{
			ID:        id,
			Name:      name,
			Code:      code,
			Used:      used,
			Discount:  discount,
			CreatedAt: created_at,
			UpdatedAt: updated_at,
		}

		return PromotionResponse, nil
	}
	return nil, nil
}

func UpdatePromotion(id uint64, PromoNew dto.PromotionDTOUpdate) (*models.Promotion, error) {
	var discount float64
	var name, code string
	var used bool

	promoOld, err := FindByIdPromotion(id)
	if err != nil {
		return nil, err
	}

	discount = promoOld.Discount
	if PromoNew.Discount != 0 && promoOld.Discount != PromoNew.Discount {
		discount = PromoNew.Discount
	}

	name = promoOld.Name
	if PromoNew.Name != "" && promoOld.Name != PromoNew.Name {
		name = PromoNew.Name
	}

	code = promoOld.Code
	if strings.ToUpper(PromoNew.Code) != "" && promoOld.Code != PromoNew.Code {
		code = (PromoNew.Code)
	}

	used = promoOld.Used
	if PromoNew.Used != promoOld.Used {
		used = PromoNew.Used
	}
	_, err = db.Connec.Exec(queryUpdatePromotion, code, name, used, discount, time.Now(), promoOld.ID)
	if err != nil {
		return nil, err
	}

	PromoUpdate := &models.Promotion{
		ID:        promoOld.ID,
		Name:      name,
		Code:      code,
		Used:      used,
		Discount:  discount,
		CreatedAt: promoOld.CreatedAt,
		UpdatedAt: time.Now(),
	}
	return PromoUpdate, nil
}

func DeletePromotion(id uint64) error {
	_, err := db.Connec.Exec(queryDeletePromotion, id)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUsedPromotion(id uint64, used bool) error {
	_, err := db.Connec.Exec(queryUpdateStsProm, used, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func OrderPromotionUsed(idpromotion uint64) (uint64, error) {
	var orderid uint64
	row, err := db.Connec.Query(querySearchPromoOrder, idpromotion)
	if err != nil {
		return orderid, err
	}
	if row.Next() {
		err = row.Scan(&orderid)
		if err != nil {
			return orderid, err
		}
	}
	return orderid, nil
}
