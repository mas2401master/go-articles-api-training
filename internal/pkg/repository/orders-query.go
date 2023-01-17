package repository

import (
	"strings"

	"github.com/mas2401master/go-articles-api-training/internal/pkg/dto"
)

const (
	queryCreateOrder = `insert into orders
	(user_id, promotion_id, subtotal, total_discount, total, quantity, status, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6,$7,$8,$9);`
	queryOrderId = `select a.id id, order_number, user_id, coalesce(promotion_id,0) promotion_id, subtotal, coalesce(total_discount,0) total_discount, total, quantity, 
					a.status status, a.created_at created_at, a.updated_at updated_at, coalesce(code,'') code, 
					coalesce(name,''), coalesce(discount,0) discount,username,
					coalesce(firstname,'') firstname,coalesce(lastname,'') lastname
					from orders a  left join promotion b on promotion_id = b.id 
					left join users c on user_id =c.id where a.id=$1;`
	queryOrderUserStatus = `select id, order_number, user_id, coalesce(promotion_id,0) promotion_id, subtotal, coalesce(total_discount,0), total, quantity, status, 
					created_at, updated_at from orders where user_id=$1 and status=$2;`
	queryOrderUser = `select id, order_number, user_id, coalesce(promotion_id,0) promotion_id, subtotal, coalesce(total_discount,0), total, quantity, status, 
					created_at, updated_at from orders where user_id=$1;`
	queryUpdatePromoOrder = `update orders set 
								promotion_id=$1, 
								total_discount= $2, 
								total=$3, 
								updated_at=$4
							where id=$5;`
	queryUpdateStatus         = `update orders set status=$1, updated_at=$2 where id=$3;`
	queryUpdatePromoOrderZero = `update orders o set 
								promotion_id=$1, 
								total_discount= 0, 
								total=subtotal, 
								updated_at=$2
							from  promotion p
							where o.id=$3 and o.promotion_id=p.id;`

	queryUpdateOrderAmount = `update orders set quantity=$1, subtotal=$2, total_discount=$3,total=$4, updated_at=$5 where id=$6;`

	queryDeleteOrder = `delete from orders where id=$1;`
)

func getOrderQuery(filter dto.OrderFilter) string {
	var (
		code, status, username, userid string
		clause                         = false
		query                          = "select a.id id, order_number, user_id, coalesce(promotion_id,0) promotion_id, subtotal,coalesce(total_discount,0) total_discount, total, quantity, a.status status, a.created_at created_at, a.updated_at updated_at, coalesce(code,'*') code, coalesce(name,'*') as name, coalesce(discount,0) discount, username,coalesce(firstname,'*') firstname,coalesce(lastname,'*') lastname from orders a  left join promotion b on promotion_id = b.id left join users c on user_id =c.id"
		b                              = strings.Builder{}
	)

	b.WriteString(query)

	if filter.Username != "" {
		username = "username LIKE '%" + filter.Username + "%'"
		clause = true
	}

	if filter.Code != "" {
		code = "code like '%" + filter.Code + "%'"
		clause = true
	}

	if filter.Status != "" {
		status = "a.status LIKE '%" + filter.Status + "%'"
		clause = true
	}

	if filter.UserId != "" {
		userid = "a.user_id =" + filter.UserId
		clause = true
	}

	if clause {
		b.WriteString(" where ")
		if username != "" {
			b.WriteString(username)
			if code != "" {
				b.WriteString(" and ")
				b.WriteString(code)
			}
			if status != "" {
				b.WriteString(" and ")
				b.WriteString(status)
			}
		} else {
			if code != "" {
				b.WriteString(code)
				if status != "" {
					b.WriteString(" and ")
					b.WriteString(status)
				}
				if userid != "" {
					b.WriteString(" and ")
					b.WriteString(userid)
				}
			} else {
				if status != "" {
					b.WriteString(status)
					if userid != "" {
						b.WriteString(" and ")
						b.WriteString(userid)
					}
				} else {
					b.WriteString(userid)
				}
			}
		}
	}
	return b.String()
}
