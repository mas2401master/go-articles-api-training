package dto

type UserFilter struct {
	Status string `json:"status"`
	RoleID uint64 `json:"role_id"`
}

type UserDTO struct {
	Username  string `json:"username" binding:"required"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email"`
	RoleID    uint64 `json:"role_id" binding:"required"`
	Status    bool   `json:"status"`
}

type UserUpdateDTO struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	RoleID    uint64 `json:"role_id"`
	Status    bool   `json:"status"`
}
type UserDTOOrder struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

type MessageInfo struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type LoginDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Logedin struct {
	UserID    string `json:"X-USERID"`
	FirstName string `json:"firsname"`
	Lastname  string `json:"lastname"`
	Token     string `json:"token"`
}

type ItemFilter struct {
	NameItem  string `json:"name"`
	Available string `json:"available"`
}

type ItemDTOUpdate struct {
	NameItem    string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Available   bool    `json:"available"`
	Status      bool    `json:"status"`
}

type ItemDTOCreate struct {
	NameItem    string  `json:"name" binding:"required,min=3,max=50"`
	Description string  `json:"description" binding:"min=3,max=300"`
	Price       float64 `json:"price" binding:"required,gte=1"`
}

type PromotionFilter struct {
	Name string `json:"name"`
	Code string `json:"code"`
	Used string `json:"used"`
}

type PromotionDTOCreate struct {
	Code     string  `json:"code" binding:"required,min=3,max=20"`
	Name     string  `json:"name" binding:"required,min=3,max=50"`
	Discount float64 `json:"discount" binding:"required,gte=1"`
}

type PromotionDTOUpdate struct {
	Code     string  `json:"code"`
	Name     string  `json:"name"`
	Used     bool    `json:"used"`
	Discount float64 `json:"discount"`
}

type OrderFilter struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Code     string `json:"code"`
	Status   string `json:"status"`
}

type OrderDTOUpdate struct {
	Code   string `json:"code"`
	Status string `json:"status"`
}

type OrderDTOCreate struct {
	ItemID   uint64 `json:"item_id" binding:"required"`
	Quantity uint64 `json:"quantity" binding:"required,gte=1"`
	Code     string `json:"code"`
	Price    float64
}

type OrderItemDTO struct {
	ItemID   uint64 `json:"item_id" binding:"required"`
	Quantity uint64 `json:"quantity" binding:"required,gte=1"`
}
