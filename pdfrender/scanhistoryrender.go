package pdfrender

import (
	"../data"
	"github.com/signintech/gopdf"
	"strconv"
	"time"
)

func showScanHistory(pdf *gopdf.GoPdf, scans *data.ScanHistoryType) {
	cellScanWidth := width/5.0

	title := []string{"Scan Date", "Image ID", "Security Status", "Image Creation Date", "Scan results"}
	pdf.SetFont(fontTypeBold, "", 10)
	showScanRow( pdf, title, cellScanWidth)

	pdf.Br(brSize)
	for _, scan := range scans.Results {
		scanDate,_ := time.Parse("2006-01-02T15:04:05.999999999Z07:00", scan.Date)
		scanDateFormatted := scanDate.Format(dateFormat)

		creationDate,_ := time.Parse( "2006-01-02T15:04:05.999999999Z07:00", scan.ImageCreationDate)
		creationDateFormatted := creationDate.Format(dateFormat)

		imageIdMultyline,_ := pdf.SplitText( scan.ImageId, cellScanWidth-2*padding )
		imageId := imageIdMultyline[0]
		if len(imageIdMultyline) > 1 {
			imageId += "..."
		}

		pdf.SetFont( fontType, "", 9)
		pdf.SetX( leftMargin + padding)
		pdf.SetY(pdf.GetY() + padding/2)

		pdf.Cell(nil, scanDateFormatted)

		pdf.SetX( leftMargin + cellScanWidth)
		pdf.Cell(nil, imageId)

		pdf.SetX( leftMargin + cellScanWidth*2)
		addCompliantText( pdf, scan.SecurityStatus)

		pdf.SetX( leftMargin + cellScanWidth*3)
		pdf.Cell(nil, creationDateFormatted)

		pdf.SetX( leftMargin + cellScanWidth*4)
		showCountsResults(pdf, scan, cellScanWidth)

		pdf.SetY(pdf.GetY() + rowSize + padding/5)

		addHrGreyH(pdf, pdf.GetY(), 0.5)

		checkEndOfPageWithoutBr(pdf, rowSize+padding)
	}
}

func showCountsResults(pdf *gopdf.GoPdf, results data.ScanHistoryResult, width float64)  {
	maxWidth := (width-padding) / 6

	pdf.SetFont( fontType, "", 8)
	critical := strconv.Itoa(results.CriticalCount)
	high := strconv.Itoa(results.HighCount)
	medium := strconv.Itoa(results.MediumCount)
	low := strconv.Itoa(results.LowCount)
	neg := strconv.Itoa(results.NegCount)
	malware := strconv.Itoa(results.MalwareCount)

	pdf.SetY( pdf.GetY()-padding/4.0)
	xBegin := pdf.GetX()

	if results.CriticalCount > 0 {
		setCriticalBackgroundColor(pdf)
	} else {
		setLightGrayBackgroundColor(pdf)
	}
	pdf.RectFromUpperLeftWithStyle( pdf.GetX(), pdf.GetY(), maxWidth,rowSize, "F")
	addCellCount( pdf, critical, maxWidth )

	pdf.SetX(xBegin+maxWidth)
	if results.HighCount > 0 {
		setHighBackgroundColor(pdf)
	} else {
		setLightGrayBackgroundColor(pdf)
	}
	pdf.RectFromUpperLeftWithStyle( pdf.GetX(), pdf.GetY(), maxWidth,rowSize, "F")
	addCellCount( pdf, high, maxWidth)

	pdf.SetX(xBegin+2*maxWidth)
	if results.MediumCount > 0 {
		setMediumBackgroundColor(pdf)
	} else {
		setLightGrayBackgroundColor(pdf)
	}
	pdf.RectFromUpperLeftWithStyle( pdf.GetX(), pdf.GetY(), maxWidth,rowSize, "F")
	addCellCount( pdf, medium, maxWidth )

	pdf.SetX(xBegin+3*maxWidth)
	if results.LowCount > 0 {
		setLowBackgroundColor(pdf)
	} else {
		setLightGrayBackgroundColor(pdf)
	}
	pdf.RectFromUpperLeftWithStyle( pdf.GetX(), pdf.GetY(), maxWidth,rowSize, "F")
	addCellCount( pdf, low, maxWidth )

	pdf.SetX(xBegin+ 4*maxWidth)
	if results.NegCount > 0 {
		setNegligibleBackgroundColor(pdf)
	} else {
		setLightGrayBackgroundColor(pdf)
	}
	pdf.RectFromUpperLeftWithStyle( pdf.GetX(), pdf.GetY(), maxWidth,rowSize, "F")
	addCellCount( pdf, neg, maxWidth )

	pdf.SetX(xBegin+ (5*maxWidth)+(padding/2.0))
	if results.MalwareCount > 0 {
		setDarkGrayBackgroundColor(pdf)
	} else {
		setLightGrayBackgroundColor(pdf)
	}
	pdf.RectFromUpperLeftWithStyle( pdf.GetX(), pdf.GetY(), maxWidth,rowSize, "F")
	addCellCount( pdf, malware, maxWidth )

	setDefaultBackgroundColor(pdf)
}

func addCellCount(pdf *gopdf.GoPdf,  text string, w float64) {
	rect := gopdf.Rect{
		W: w,
		H: rowSize,
	}
	pdf.CellWithOption(&rect, text, cellOption)
	pdf.SetTextColor(0,0,0)
}

func showScanRow(pdf *gopdf.GoPdf, content []string, w float64) {
	setLightGrayBackgroundColor(pdf)
	pdf.RectFromUpperLeftWithStyle( leftMargin, pdf.GetY(), width, rowSize+padding, "F")
	setDefaultBackgroundColor(pdf)
	pdf.SetY(pdf.GetY()+padding*0.8)

	pdf.SetX(leftMargin+padding)
	for i:=0; i<len(content); i++ {
		pdf.Cell(nil, content[i])
		pdf.SetX(leftMargin+float64(i+1)*w)
	}
}