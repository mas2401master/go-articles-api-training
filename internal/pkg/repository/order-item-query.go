package repository

const (
	queryCreateOrderItem = `insert into order_items(
		order_id, item_id, price, quantity, total, status, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6,$7,$8);`

	queryTotalQuantityOrderItem = `select coalesce(sum(quantity),0) quantity from order_items where order_id=$1`
	queryTotalOrderItem         = `select sum(total) total from order_items where order_id=$1`
	queryOrderItemId            = `select o.id id, order_id, item_id, o.price price, quantity, total, o.status status, o.created_at created_at, 
							o.updated_at updated_at,name
							from order_items o left join items i on item_id= i.id
							where o.id=$1;`
	queryDeleteOrderItem = `delete from order_items where id=$1`
	queryUpdateOrderItem = `update order_items set quantity=$1, total=$2, updated_at=$3 where id=$4;`

	queryFindAllOrderItem = `select o.id id, order_id, item_id, o.price price, quantity, total, o.status status, o.created_at created_at, 
							o.updated_at updated_at,name
							from order_items o left join items i on item_id= i.id
							where o.order_id=$1;`
	queryUpdateOrderItemSts = `update order_items set status=$1 where order_id=$2;`

	queryQuantityOrderItem = `select coalesce(sum(quantity),0) quantity from order_items where order_id=$1 and item_id=$2`
	queryRealeaseItemOrder = `update items i set available=true from order_items io where i.id=io.item_id and order_id=$1`
)
