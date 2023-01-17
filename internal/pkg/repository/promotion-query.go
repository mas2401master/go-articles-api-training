package repository

import (
	"strings"

	"github.com/mas2401master/go-articles-api-training/internal/pkg/dto"
)

const (
	queryCreatePromotion = `insert into promotion(code, name, used, discount, created_at, updated_at)
							values ( $1,$2, $3, $4, $5, $6);`
	queryUpdatePromotion = `update promotion
							set code=$1, name=$2, used=$3, discount=$4, updated_at=$5
							where id=$6;`

	queryPromtionId      = `select id, code, name, used, discount, created_at, updated_at from promotion where id = $1`
	queryPromtionCode    = `select id, code, name, used, discount, created_at, updated_at from promotion where code = $1`
	queryDeletePromotion = `delete from promotion where  id = $1`
	querySearchPromoCode = `select count(*) from promotion where code=$1`
	queryUpdateStsProm   = `update promotion
							set used=$1, updated_at=$2
							where id=$3;`
	querySearchPromoOrder = `select id from orders where promotion_id=$1;`
)

func getPromotionQuery(filter dto.PromotionFilter) string {
	var (
		code, name, used string
		clause           = false
		query            = "select id, code, name, used, discount, created_at, updated_at from promotion"
		b                = strings.Builder{}
	)

	b.WriteString(query)
	if filter.Name != "" {
		name = "name like '%" + filter.Name + "%'"
		clause = true
	}

	if filter.Code != "" {
		code = "code like '%" + filter.Code + "%'"
		clause = true
	}

	if filter.Used != "" {
		used = "used = " + filter.Used
		clause = true
	}

	if clause {
		b.WriteString(" where ")
		if name != "" {
			b.WriteString(name)
			if code != "" {
				b.WriteString(" and ")
				b.WriteString(code)
			}
			if used != "" {
				b.WriteString(" and ")
				b.WriteString(used)
			}
		} else {
			if code != "" {
				b.WriteString(code)
				if used != "" {
					b.WriteString(" and ")
					b.WriteString(used)
				}
			} else {
				b.WriteString(used)
			}
		}
	}
	return b.String()
}
