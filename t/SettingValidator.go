package t

type SettingValidator struct {
	EnableAuthentication *string `form:"EnableAuthentication" validate:"omitempty,oneof=on off"`
	AuthenticationType   *string `form:"AuthenticationType" validate:"omitempty,oneof=form basic"`

	EnableAutomaticScanns    *string `form:"EnableAutomaticScanns" validate:"omitempty,oneof=on off"`
	AutomaticScannsInterval  int     `form:"AutomaticScannsInterval" validate:"omitempty,number,gte=1"`
	AutomaticScannsAtStartup *string `form:"AutomaticScannsAtStartup" validate:"omitempty,oneof=on off"`

	PreCopyFileCount int `form:"PreCopyFileCount" validate:"omitempty,number,gte=0,lte=10"`

	EnableEncoding          *string `form:"EnableEncoding" validate:"omitempty,oneof=on off"`
	EncodingThreads         int     `form:"EncodingThreads" validate:"omitempty,number,gte=1"`
	EncodingCrf             int     `form:"EncodingCrf" validate:"omitempty,number,gte=1,lte=50"`
	EncodingResolution      int     `form:"EncodingResolution" validate:"omitempty,number,gte=100,lte=5000"`
	EnableHevcEncoding      *string `form:"EnableHevcEncoding" validate:"omitempty,oneof=on off"`
	EncodingMaxRetry        int     `form:"EncodingMaxRetry" validate:"omitempty,number,gte=0,lte=999"`
	EnableNvidiaGpuEncoding *string `form:"EnableNvidiaGpuEncoding" validate:"omitempty,oneof=on off"`
	EnableAmdGpuEncoding    *string `form:"EnableAmdGpuEncoding" validate:"omitempty,oneof=on off"`
	EnableImageComparison   *string `form:"EnableImageComparison" validate:"omitempty,oneof=on off"`
}
