package service

import (
	"fmt"
	"mxclub/apps/mxclub-admin/config"
	"mxclub/domain/order/entity/enum"
	"mxclub/domain/order/po"
	"mxclub/domain/order/repo"
	"mxclub/pkg/common/wxpay"
	"mxclub/pkg/common/xmysql"
	"mxclub/pkg/common/xredis"
	"mxclub/pkg/utils"
	"time"

	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/pkg/errors"
	"github.com/xuri/excelize/v2"
)

func init() {
	jet.Provide(NewExcelService)
}

func NewExcelService(orderRepo repo.IOrderRepo,
	callbackRepo repo.IWxPayCallbackRepo,
	rewardRepo repo.IRewardRecordRepo) *ExcelService {
	return &ExcelService{
		orderRepo:    orderRepo,
		callbackRepo: callbackRepo,
		rewardRepo:   rewardRepo,
	}
}

type ExcelService struct {
	orderRepo    repo.IOrderRepo
	callbackRepo repo.IWxPayCallbackRepo
	rewardRepo   repo.IRewardRecordRepo
}

type exportExcelDTO struct {
	transactionId string // 交易单号
	outTradeNo    string
}

func (svc ExcelService) ExportExcel(ctx jet.Ctx, startDate, endDate string) (err error) {
	defer utils.RecoverAndLogError(ctx)
	defer utils.TraceElapsed(ctx, "ExportExcel")()
	// 防止同一时间多个用户同时导出
	if err = xredis.Debounce("export_excel", time.Minute*5); err != nil {
		return errors.New("正在有其他用户导出，五分钟内只允许一个用户进行导出")
	}
	var (
		f         = excelize.NewFile()
		logger    = ctx.Logger()
		sheetName = "发货单模板"
		appId     = config.GetConfig().WxPayConfig.AppId
		tradeType = "虚拟发货"
		tradeMode = "统一发货"
		// 列标题
		headers = []struct {
			cell string
			text string
		}{
			{"A1", "交易单号"},
			{"B1", "商户单号"},
			{"C1", "商户号"},
			{"D1", "发货方式"},
			{"E1", "发货模式"},
			{"F1", "快递公司"},
			{"G1", "快递单号（多个快递单使用;分隔）"},
			{"H1", "是否完成发货"},
			{"I1", "是否重新发货"},
			{"J1", "商品信息"},
		}
	)

	if endDate == "" {
		endDate = time.Now().Format("2006-01-02 15:04:05")
	}

	// 1. 查找指定时间内的订单
	wrapper := new(xmysql.MysqlQuery)
	wrapper.SetFilter("created_at >= ? and created_at <= ? and order_status = ?", startDate, endDate, enum.SUCCESS)
	wrapper.SetLimit(1000000)
	orderPOList, err := svc.orderRepo.ListNoCountByQuery(wrapper)
	if err != nil || orderPOList == nil || len(orderPOList) == 0 {
		logger.Errorf("cannot find any order, duration is: %v %v", startDate, endDate)
		return
	}

	orderIdList := utils.Map(orderPOList, func(in *po.Order) string { return utils.ParseString(in.OrderId) })

	if !config.GetConfig().WxPayConfig.IsBaoZaoClub() {
		// 2. 打赏订单
		rewardRecords, err := svc.rewardRepo.ListNoCountDuration(ctx, startDate, endDate, enum.SUCCESS)
		if err != nil {
			ctx.Logger().Errorf(err.Error())
		} else {
			outTradeNoList := utils.Map(rewardRecords, func(in *po.RewardRecord) string {
				return in.OutTradeNo
			})
			ctx.Logger().Infof("find outTradeNoList => %v", outTradeNoList)
			orderIdList = append(orderIdList, outTradeNoList...)
		}
	}

	callbackWrapper := new(xmysql.MysqlQuery)
	callbackWrapper.SetFilter("out_trade_no in ?", orderIdList)
	callbackWrapper.SetLimit(10000)
	callbackPOList, err := svc.callbackRepo.ListNoCountByQuery(callbackWrapper)
	if err != nil || callbackPOList == nil || len(callbackPOList) <= 0 {
		logger.Errorf("cannot find any callbackPOList, orderIdList is: %v", utils.ObjToJsonStr(orderIdList))
		return
	}

	exportExcelDTOS := utils.Map(callbackPOList, func(in *po.WxPayCallback) *exportExcelDTO {
		callbackDTO := utils.MustMapToObj[wxpay.WxPayCallbackDTO](in.RawData)
		if callbackDTO == nil {
			return nil
		}
		return &exportExcelDTO{
			transactionId: callbackDTO.TransactionID,
			outTradeNo:    callbackDTO.OutTradeNo,
		}
	})

	// 创建一个工作表
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return err
	}
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// 循环设置列标题
	for _, header := range headers {
		ifErrThrowPanic(f.SetCellValue(sheetName, header.cell, header.text))
	}

	// 写入订单数据
	for rowIndex, dto := range exportExcelDTOS {
		if dto == nil {
			logger.Errorf("exportExcelDTO is nil")
			continue
		}
		row := rowIndex + 2 // 数据从第二行开始，第一行为标题
		cells := []string{
			dto.transactionId,
			dto.outTradeNo,
			appId,
			tradeType,
			tradeMode,
			tradeType, // 快递公司暂时用发货方式代替
			"",        // 快递单号暂为空
			"是",
			"否",
			"用户自提商品",
		}

		for col, cellValue := range cells {
			cell := fmt.Sprintf("%c%d", 'A'+col, row)
			ifErrThrowPanic(f.SetCellValue(sheetName, cell, cellValue))
		}
	}

	// 导出到输出流中
	resp := ctx.Response()
	resp.Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=发货单-%v-%v.xlsx", startDate, endDate))
	resp.Header.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// 将Excel文件写入HTTP响应体
	writer := resp.BodyWriter()
	err = f.Write(writer)
	if err != nil {
		return err
	}

	return
}

func ifErrThrowPanic(err error) {
	if err != nil {
		panic(err)
	}
}
