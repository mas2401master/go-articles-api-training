package repository

import (
	"strconv"
	"strings"

	"github.com/mas2401master/go-articles-api-training/internal/pkg/dto"
)

const (
	queryCreateUser   = `insert into users(username, firstname, lastname, password, email, role_id, status, created_at,updated_at) values($1, $2, $3, $4, $5, $6,$7,$8,$9);`
	queryUpdateUser   = `update users set username=$1, firstname=$2, lastname=$3, password=$4, email=$5, role_id=$6, status=$7, updated_at=$8 where id=$9;`
	queryDeleteUser   = `delete from users where id=$1`
	queryUserUserName = `select id, username, firstname, lastname, password, email, role_id, status,created_at,updated_at from users where username=$1`
	queryUseId        = `select id, username, firstname, lastname, password, email, role_id, status,created_at,updated_at from users where id=$1`
	queryRolId        = `select id, name from role where id=$1`
)

func getUserQuery(filter dto.UserFilter) string {
	var (
		role, status string
		clause       = false
		query        = "select a.id id,username, firstname, lastname, password, email, role_id, a.status status, a.created_at created_at, a.updated_at, name  from users a left join role b on a.role_id = b.id"
		b            = strings.Builder{}
	)

	b.WriteString(query)

	if filter.RoleID != 0 {
		role = "a.role_id = " + strconv.FormatUint(filter.RoleID, 10)
		clause = true
	}

	if filter.Status != "" {
		status = "a.status = " + filter.Status
		clause = true
	}

	if clause {
		b.WriteString(" where ")

		if role != "" {
			b.WriteString(role)
			if status != "" {
				b.WriteString(" and ")
				b.WriteString(status)
			}
		} else {
			if status != "" {
				b.WriteString(status)
			}
		}
	}
	return b.String()
}
