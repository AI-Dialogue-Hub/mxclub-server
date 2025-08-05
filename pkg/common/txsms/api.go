package txsms

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"mxclub/pkg/utils"
)

var (
	txSmsService *TxSmsService
)

func SendDefaultDispatchMsg(phone string) error {
	defer utils.RecoverByPrefixNoCtx("SendDefaultDispatchMsg")
	if !txSmsService.config.IsOk {
		return fmt.Errorf("tx sms service is not ok")
	}
	if txSmsService != nil {
		return txSmsService.SendDispatchMsg(phone)
	}
	return fmt.Errorf("tx sms service is nil")
}

type TxSmsService struct {
	config *TxSmsConfig
	client *sms.Client
}

func NewTxSmsService(config *TxSmsConfig) *TxSmsService {
	credential := common.NewCredential(
		config.Ak,
		config.SK,
	)
	// 使用短信服务的SDK，地域通常使用ap-guangzhou
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := sms.NewClient(credential, "ap-guangzhou", cpf)

	txSmsService = &TxSmsService{
		config: config,
		client: client,
	}

	return txSmsService
}

// SendDispatchMsg 发送默认短信，为指定打手派单短信
func (s *TxSmsService) SendDispatchMsg(phone string) error {
	request := sms.NewSendSmsRequest()

	// 设置短信应用ID
	request.SmsSdkAppId = common.StringPtr(s.config.SmsSdkAppId)
	// 设置短信签名
	request.SignName = common.StringPtr(s.config.SignName)
	// 设置短信模板ID
	request.TemplateId = common.StringPtr(s.config.TemplateId)
	// 设置模板参数，根据实际模板参数设置
	//request.TemplateParamSet = common.StringPtrs([]string{consignee, role})
	// 设置手机号，需要带国际区号，如"+86"表示中国
	phoneNumber := "+86" + phone
	request.PhoneNumberSet = common.StringPtrs([]string{phoneNumber})

	response, err := s.client.SendSms(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return fmt.Errorf("腾讯云API错误: %v", err)
	}
	if err != nil {
		return fmt.Errorf("发送短信失败: %v", err)
	}

	// 检查响应状态
	for _, status := range response.Response.SendStatusSet {
		if *status.Code != "Ok" {
			return fmt.Errorf("短信发送失败: %s", *status.Message)
		}
	}

	return nil
}
