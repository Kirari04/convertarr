package t

type LoginValidator struct {
	Username string `form:"Username" validate:"required,max=55"`
	Password string `form:"Password" validate:"required,max=255"`
}
