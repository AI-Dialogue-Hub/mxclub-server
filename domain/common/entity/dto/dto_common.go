package dto

type SwiperDTO struct {
	ImageUrl string `json:"image_url"`
	GoodsId  uint   `json:"goods_id"`
}

type NotificationsDTO struct {
	Image   string `json:"image"`
	Message string `json:"message"`
}
