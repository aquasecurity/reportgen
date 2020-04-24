package pdfrender

import (
	"../data"
	"fmt"
	"github.com/signintech/gopdf"
	"time"
)

func showScanHistory(pdf *gopdf.GoPdf, scans *data.ScanHistoryType) {
	cellScanWidth := width/5.0

	title := []string{"Scan Date", "Image ID", "Security Status", "Image Creation Date", "Scan results"}
	showScanRow( pdf, title, cellScanWidth)
	for _, scan := range scans.Results {
		var securityStatus string
		if scan.SecurityStatus {
			securityStatus = "Non-compliant"
		} else {
			securityStatus = "Non-compliant"
		}
		scanResult := fmt.Sprintf("%d / %d / %d / %d / %d", scan.CriticalCount, scan.HighCount, scan.MediumCount, scan.LowCount, scan.NegCount)

		scanDate,_ := time.Parse("2006-01-02T15:04:05.999999999Z07:00", scan.Date)
		scanDateFormatted := scanDate.Format(dateFormat)



		showScanRow( pdf, []string {
			scanDateFormatted, scan.ImageId, securityStatus, scan.ImageCreationDate, scanResult,
		}, cellScanWidth)

		pdf.Br(brSize)
		pdf.SetX(leftMargin)
	}
}

func showScanRow(pdf *gopdf.GoPdf, content []string, w float64) {
	for i:=0; i<len(content); i++ {
		pdf.SetX(leftMargin+float64(i)*w)
		pdf.Cell(nil, content[i])
	}
}