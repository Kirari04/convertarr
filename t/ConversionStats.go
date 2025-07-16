package t

type ConversionStats struct {
	Labels     []string `json:"labels"`
	Successful []int    `json:"successful"`
	Failed     []int    `json:"failed"`
}
