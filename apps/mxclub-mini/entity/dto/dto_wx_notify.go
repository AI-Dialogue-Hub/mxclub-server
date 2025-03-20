package dto

// MiniProgram 嵌套结构体用于表示小程序信息
type MiniProgram struct {
	Appid    string `json:"appid"`    // 小程序的AppID
	Pagepath string `json:"pagepath"` // 小程序页面路径
}

// DataValue 嵌套结构体用于表示"data"字段中的每个值
type DataValue struct {
	Value any `json:"value"` // 数据的具体值
}

// WxNotifyMessageSendDTO 主结构体
type WxNotifyMessageSendDTO struct {
	Touser     string `json:"touser"`         // 接收者（用户）的 openid
	TemplateID string `json:"template_id"`    // 模板ID
	Page       string `json:"page,omitempty"` // 点击模板卡片后的跳转页面（可选）
	/*
		 {
			 "appid":"APPID",
			 "pagepath":"index?foo=bar"
		  }
	*/
	Miniprogram *MiniProgram         `json:"miniprogram,omitempty"` // 小程序跳转链接
	Lang        string               `json:"lang"`
	Data        map[string]DataValue `json:"data"` // 模板内容
}
