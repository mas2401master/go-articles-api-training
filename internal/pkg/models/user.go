package models

type User struct {
	ID        uint64 `json:"id"`
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	RoleID    uint64 `json:"role_id"`
	Status    bool   `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Role      Role   `json:"role"`
}
