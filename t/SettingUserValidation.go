package t

type SettingUserValidation struct {
	Username string `form:"Username" validate:"required,max=55"`
	Password string `form:"Password" validate:"omitempty,max=255"`
}
