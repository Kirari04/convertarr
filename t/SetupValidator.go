package t

type SetupValidator struct {
	EnableAuthentication *string `form:"EnableAuthentication" validate:"omitempty,oneof=on off"`
	AuthenticationType   *string `form:"AuthenticationType" validate:"omitempty,oneof=form basic"`
	Username             *string `form:"Username" validate:"required_with=EnableAuthentication,max=55"`
	Password             *string `form:"Password" validate:"required_with=EnableAuthentication,max=255"`
}
