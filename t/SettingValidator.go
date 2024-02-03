package t

type SettingValidator struct {
	EnableAuthentication *string `form:"EnableAuthentication" validate:"omitempty,oneof=on off"`
	AuthenticationType   *string `form:"AuthenticationType" validate:"omitempty,oneof=form basic"`

	EnableAutomaticScanns    *string `form:"EnableAutomaticScanns" validate:"omitempty,oneof=on off"`
	AutomaticScannsInterval  int     `form:"AutomaticScannsInterval" validate:"omitempty,number,gte=1"`
	AutomaticScannsAtStartup *string `form:"AutomaticScannsAtStartup" validate:"omitempty,oneof=on off"`
}
