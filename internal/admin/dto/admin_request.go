package admin

type AdminRequestBody struct {
	Email     string  `json:"email" binding:"email"`
	Name      string  `json:"name" binding:"required"`
	Password  *string `json:"password"`
	CreatedBy *string `json:"created_by"`
	UpdatedBy *string `json:"updated_by"`
}
