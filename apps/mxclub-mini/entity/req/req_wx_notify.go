package req

type WxNotifySendReq struct {
	Message string `json:"message"`
	UserId  uint   `json:"user_id"` // user db id
}

type NotifyTemplateReq struct {
	TemplateId string `json:"template_id" form:"templateId" validate:"required" reg_error_info:"TemplateId不能为空"`
}
