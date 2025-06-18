package service

import (
	"errors"
	"fmt"
	"github.com/fengyuan-liang/jet-web-fasthttp/jet"
	"github.com/xuri/excelize/v2"
	"mxclub/pkg/utils"
	"strconv"
	"time"
)

func (svc *OrderService) ExportAllDasherWithDrawAmount(ctx jet.Ctx) error {
	defer utils.RecoverAndLogError(ctx)
	defer utils.TraceElapsed(ctx, "ExportAllDasherWithDrawAmount")()
	var (
		f         = excelize.NewFile()
		logger    = ctx.Logger()
		sheetName = "明星电竞-打手账单总览"
		// 列标题
		headers = []struct {
			cell string
			text string
		}{
			{"A1", "打手编号"},
			{"B1", "打手名称"},
			{"C1", "历史提现金额"},
			{"D1", "还能提现金额"},
		}
	)
	withDrawVOList, err := svc.AllDasherHistoryWithDrawAmount(ctx)
	if err != nil {
		ctx.Logger().Errorf("OrderService#ExportAllDasherWithDrawAmount err:%v", err)
		return errors.New("导出失败")
	}
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
	for rowIndex, withDrawVO := range withDrawVOList {
		if withDrawVO == nil {
			logger.Errorf("withDrawVO is nil")
			continue
		}
		row := rowIndex + 2 // 数据从第二行开始，第一行为标题
		cells := []string{
			strconv.Itoa(int(withDrawVO.DasherID)),
			withDrawVO.DasherName,
			utils.ParseString(withDrawVO.HistoryWithDrawAmount),
			utils.ParseString(withDrawVO.WithdrawAbleAmount),
		}

		for col, cellValue := range cells {
			cell := fmt.Sprintf("%c%d", 'A'+col, row)
			ifErrThrowPanic(f.SetCellValue(sheetName, cell, cellValue))
		}
	}
	// 导出到输出流中
	resp := ctx.Response()
	formatDate := time.Now().Format(time.DateTime)
	resp.Header.Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=打手账单明细-%v.xlsx", formatDate))
	resp.Header.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// 将Excel文件写入HTTP响应体
	writer := resp.BodyWriter()
	err = f.Write(writer)
	if err != nil {
		return err
	}

	return nil
}
