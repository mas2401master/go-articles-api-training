package repository

import (
	"time"

	"github.com/mas2401master/go-articles-api-training/internal/pkg/db"
	"github.com/mas2401master/go-articles-api-training/internal/pkg/dto"
	models "github.com/mas2401master/go-articles-api-training/internal/pkg/models"
	"github.com/mas2401master/go-articles-api-training/pkg/encryption"
)

func FindAllUser(filter dto.UserFilter) ([]models.User, error) {
	var users []models.User
	rows, err := db.Connec.Query(getUserQuery(filter))
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var id, role_id uint64
		var username, firstname, lastname, password, email, name_rol string
		var status bool
		var created_at, updated_at time.Time
		err := rows.Scan(&id, &username, &firstname, &lastname, &password, &email, &role_id, &status, &created_at, &updated_at, &name_rol)
		if err != nil {
			return users, err
		}
		user := models.User{
			ID:        id,
			Username:  username,
			Firstname: firstname,
			Lastname:  lastname,
			Password:  password,
			Email:     email,
			Status:    status,
			RoleID:    role_id,
			CreatedAt: created_at.String(),
			UpdatedAt: updated_at.String(),
			Role:      models.Role{ID: role_id, Name: name_rol},
		}

		users = append(users, user)
	}
	return users, nil
}

func AddUser(user models.User) error {
	bb, _ := encryption.Encrypt([]byte(user.Password))
	pass := encryption.ToBase64(bb)
	_, err := db.Connec.Exec(queryCreateUser, user.Username, user.Firstname, user.Lastname, pass, user.Email, user.RoleID, true, time.Now(), time.Now())
	if err != nil {
		return err
	}

	return nil
}

func FindByUsername(username string) (*models.User, error) {
	var UserResponse *models.User
	var RoleResponse *models.Role
	row, err := db.Connec.Query(queryUserUserName, username)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	if row.Next() {
		var id, role_id uint64
		var firstname, lastname, password, email string
		var status bool
		var created_at, updated_at time.Time
		err := row.Scan(&id, &username, &firstname, &lastname, &password, &email, &role_id, &status, &created_at, &updated_at)
		if err != nil {
			return nil, err
		}

		RoleResponse, _ = FindByIdRol(role_id)
		UserResponse = &models.User{
			ID:        id,
			Username:  username,
			Firstname: firstname,
			Lastname:  lastname,
			Password:  password,
			Email:     email,
			RoleID:    role_id,
			Status:    status,
			CreatedAt: created_at.String(),
			UpdatedAt: updated_at.String(),
			Role:      models.Role{ID: RoleResponse.ID, Name: RoleResponse.Name},
		}
		return UserResponse, nil
	}
	return nil, nil
}

func FindById(id uint64) (*models.User, error) {
	var UserResponse *models.User
	var RoleResponse *models.Role
	row, err := db.Connec.Query(queryUseId, id)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	if row.Next() {
		var role_id uint64
		var username, firstname, lastname, password, email string
		var status bool
		var created_at, updated_at time.Time
		err := row.Scan(&id, &username, &firstname, &lastname, &password, &email, &role_id, &status, &created_at, &updated_at)
		if err != nil {
			return nil, err
		}
		RoleResponse, _ = FindByIdRol(role_id)
		UserResponse = &models.User{
			ID:        id,
			Username:  username,
			Firstname: firstname,
			Lastname:  lastname,
			Password:  password,
			Email:     email,
			RoleID:    role_id,
			Status:    status,
			CreatedAt: created_at.String(),
			UpdatedAt: updated_at.String(),
			Role:      models.Role{ID: RoleResponse.ID, Name: RoleResponse.Name},
		}
		return UserResponse, nil
	}
	return nil, nil
}

func UpdateUser(id uint64, userNew dto.UserUpdateDTO) (*models.User, error) {
	var username, firstname, lastname, password, email string
	var role_id uint64
	var status bool
	var userUpdate *models.User
	var RoleResponse *models.Role

	userOld, err := FindById(id)
	if err != nil {
		return nil, err
	}
	username = userOld.Username
	if userNew.Username != "" {
		username = userNew.Username
	}

	firstname = userOld.Firstname
	if userNew.Firstname != "" {
		username = userNew.Firstname
	}

	lastname = userOld.Lastname
	if userNew.Firstname != "" {
		lastname = userNew.Lastname
	}

	password = userOld.Password
	if userNew.Password != "" {
		bb, _ := encryption.Encrypt([]byte(userNew.Password))
		password = encryption.ToBase64(bb)
	}

	email = userOld.Email
	if userNew.Email != "" {
		email = userNew.Email
	}

	role_id = userOld.RoleID
	if userNew.RoleID != 0 {
		role_id = userNew.RoleID
	}

	status = userOld.Status
	if userNew.Status != userOld.Status {
		status = userNew.Status
	}
	_, err = db.Connec.Exec(queryUpdateUser, username, firstname, lastname, password, email, role_id, status, time.Now(), userOld.ID)
	if err != nil {
		return nil, err
	}
	RoleResponse, _ = FindByIdRol(role_id)
	userUpdate = &models.User{
		ID:        id,
		Username:  username,
		Firstname: firstname,
		Lastname:  lastname,
		Password:  password,
		Email:     email,
		RoleID:    role_id,
		Status:    status,
		CreatedAt: userOld.CreatedAt,
		UpdatedAt: time.Now().String(),
		Role:      models.Role{ID: RoleResponse.ID, Name: RoleResponse.Name},
	}
	return userUpdate, nil
}

func DeleteUser(id uint64) error {
	_, err := db.Connec.Exec(queryDeleteUser, id)
	if err != nil {
		return err
	}
	return nil
}

func FindByIdRol(id uint64) (*models.Role, error) {
	var role *models.Role
	row, err := db.Connec.Query(queryRolId, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	if row.Next() {
		var name string
		err := row.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		role = &models.Role{
			ID:   id,
			Name: name,
		}
		return role, nil
	}
	return nil, nil
}
