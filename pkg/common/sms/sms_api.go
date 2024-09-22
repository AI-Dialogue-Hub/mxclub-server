package sms

import (
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/fengyuan-liang/jet-web-fasthttp/pkg/xlog"
	"mxclub/pkg/utils"
)

var logger = xlog.NewWith("sms_svc")

type ISmsService interface {
	// SendDispatchMsg consignee=用户昵称；role=用户昵称；
	SendDispatchMsg(phone, consignee, role string) error
}

type AliyunSmsService struct {
	config *Config
	client *dysmsapi20170525.Client
}

func (svc AliyunSmsService) SendDispatchMsg(phone, consignee, role string) error {
	req := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  utils.Ptr(phone),
		SignName:      utils.Ptr(svc.config.SignName),
		TemplateCode:  utils.Ptr(svc.config.SmsCode),
		TemplateParam: utils.Ptr(utils.ObjToJsonStr(NewDispatchReq(consignee, role))),
	}
	result, err := svc.client.SendSms(req)
	if err != nil {
		logger.Errorf("SendSms ERROR, %v", err)
		return err
	}
	logger.Infof("send sms success, result is: %v", utils.ObjToJsonStr(result))
	return nil
}

func NewAliSmsService(config *Config) ISmsService {
	// init client
	c, err := CreateClient(config)
	if err != nil {
		logger.Infof("NewAliSmsService init error: %v", err)
		panic(err)
	}
	return &AliyunSmsService{config: config, client: c}
}
