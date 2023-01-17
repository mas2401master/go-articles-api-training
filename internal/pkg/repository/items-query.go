package repository

import (
	"strings"

	"github.com/mas2401master/go-articles-api-training/internal/pkg/dto"
)

const (
	queryCreateItem = `insert into items(
		name, description, price, available, created_at, updated_at,status)
		values ( $1, $2, $3, $4, $5, $6,$7);`
	queryItemId     = `select id, name, description, price, available, created_at, updated_at,status from items where id=$1;`
	queryUpdateItem = `update items
						set name=$1, description=$2, price=$3, available=$4, updated_at=$5, status=$6
						where id=$7;`
	queryDeleteItem          = `delete from items where id=$1;`
	queryUpdateAvailableItem = `update items set  available=$1, updated_at=$2 where id=$3;`
)

func getItemQuery(filter dto.ItemFilter) string {
	var (
		name, available string
		clause          = false
		query           = "select id, name, description, price, available, created_at, updated_at,status from items where status=true"
		b               = strings.Builder{}
	)

	b.WriteString(query)

	if filter.NameItem != "" {
		name = "name like '%" + filter.NameItem + "%'"
		clause = true
	}

	if filter.Available != "" {
		available = "available = " + filter.Available
		clause = true
	}

	if clause {
		b.WriteString(" and ")
		if name != "" {
			b.WriteString(name)
			if available != "" {
				b.WriteString(" and ")
				b.WriteString(available)
			}
		} else {
			b.WriteString(available)
		}
	}
	return b.String()
}
