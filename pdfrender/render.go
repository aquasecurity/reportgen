package pdfrender

import (
	"../data"
	"github.com/signintech/gopdf"
	"log"
	"time"
)

const (
	leftMargin = 30
	rightMargin = 30
	topMargin = 30

	padding = 10
	width = 595-leftMargin-rightMargin
	brSize = 20

	cellHeight = 25
	cellWidth = width/5.0

	ttfPathRegular = "./pdfrender/calibri.ttf"
	ttfPathBold = "./pdfrender/calibri-bold.ttf"
	longPath = "./pdfrender/logo.png"

	dateFormat = "2006-01-02 15:04"
)

var cellOption = gopdf.CellOption{
	Align:  gopdf.Middle | gopdf.Center,
	Border: 0,
	Float:  gopdf.Right,
}

func Render(output string, data *data.Report)  {
	currentDate := time.Now().Format(dateFormat)

	fontType := "calibri"
	fontTypeBold := "calibri-bold"

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{ PageSize: *gopdf.PageSizeA4 }) //595.28, 841.89 = A4
	pdf.AddPage()

	err := pdf.AddTTFFont(fontType, ttfPathRegular)
	if err != nil {
		log.Print(err.Error())
		return
	}
	err = pdf.AddTTFFont(fontTypeBold, ttfPathBold )
	if err != nil {
		log.Print(err.Error())
		return
	}
	err = pdf.SetFont(fontType, "", 14)
	if err != nil {
		log.Print(err.Error())
		return
	}
	// Logo
	pdf.Image(longPath, leftMargin, topMargin, nil)

	// Image Vulnerability Report
	yTitleBase := 110.0
	pdf.SetFillColor(67,77,98)
	pdf.RectFromUpperLeftWithStyle(leftMargin, yTitleBase, width, 15+2*padding, "F")

	pdf.SetX(leftMargin+padding)
	pdf.SetY(yTitleBase + padding)
	pdf.SetTextColor(255,255,255)
	pdf.Cell(nil, "Image Vulnerability Report")

	pdf.SetX(leftMargin)
	pdf.SetY(yTitleBase + 3*padding+15)
	pdf.SetTextColor(0,0,0)
	pdf.SetFont(fontType, "", 10)
	pdf.Cell(nil, "Aqua Server – " + data.Server)

	// line after 1
	yLine1 := pdf.GetY()+ 2* padding
	pdf.SetLineWidth(1)
	pdf.Line(leftMargin, yLine1, leftMargin+width, yLine1)

	pdf.SetX(leftMargin)
	pdf.SetY(yLine1+padding)
	pdf.SetFont(fontType, "", 9)
	pdf.Cell(nil, "Report generated on "+ currentDate)

	pdf.Br(brSize*1.5)

	pdf.SetX(leftMargin)
	pdf.SetFont(fontTypeBold, "", 19)
	pdf.SetTextColor(67,77,98)
	pdf.Cell(nil, "Summary")
	pdf.Br(brSize)

	pdf.SetY(pdf.GetY()+padding)
	pdf.SetX(leftMargin)
	pdf.SetFont(fontType, "", 10)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, "This section contains the image summary")
	pdf.SetY(pdf.GetY()+15)

	// Block after Summary
	yLine2 := pdf.GetY()+padding
	summaryBlochHeight := 110.0
	pdf.SetX(leftMargin)
	pdf.SetY(yLine2)
	addHr(&pdf, yLine2)

	pdf.SetFillColor(246,249,250)
	pdf.RectFromUpperLeftWithStyle(leftMargin, yLine2, width, summaryBlochHeight, "F")
	pdf.SetY(yLine2+padding)
	pdf.SetX(leftMargin+padding)
	pdf.SetFont(fontTypeBold, "", 10)
	pdf.Cell(nil, "Image name \"" + data.ImageName+"\"")
	pdf.Br(brSize)
	pdf.SetFont(fontType, "", 8)

	opt := gopdf.CellOption{
		Align:  gopdf.Right,
		Border: 0,
		Float:  gopdf.Right,
	}
	rect := gopdf.Rect{
		W: width,
		H: 10,
	}
	pdf.CellWithOption(&rect, "Report generated on "+currentDate, opt)

	pdf.Br(brSize)
	pdf.SetX(leftMargin+padding)
	pdf.SetFont(fontTypeBold, "", 8)
	pdf.Cell(nil, "Registry: " + data.Registry)

	pdf.Br(brSize)
	pdf.SetX(leftMargin+padding)
	pdf.Cell(nil, "Image Creation Date: " + data.Created.Format(dateFormat))

	pdf.Br(brSize)
	pdf.SetX(leftMargin+padding)
	pdf.Cell(nil, "OS: " + data.Os + "(" + data.OsVersion+ ")")

	// after Image Name block
	pdf.SetX(leftMargin)
	pdf.SetY(yLine2 + summaryBlochHeight + 2*padding)
	pdf.SetFont(fontTypeBold, "", 10)
	var imageAllowed string
	if data.ImageAllowed {imageAllowed = "Allowed"} else {imageAllowed="Disallowed"}

	pdf.Cell(nil, "Image is " + imageAllowed)

	// Block Number of Vulnerabilities
	pdf.Br(brSize*1.5)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, "Number of Vulnerabilities")
	pdf.Br(brSize)
	yTable1 := pdf.GetY()
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(0.5)
	showColorfulTable(&pdf, yTable1)

	showTextIntoTable(&pdf, leftMargin, yTable1, &[2][5]string{
		{"CRITICAL","HIGH","MEDIUM","LOW","NEGLIGIBLE",},
		{countCritical,countHigh,countMedium,countLow,countNegligible,},
	})

	// Image Assurance Checks
	pdf.SetY(yTable1 + cellHeight*2+ brSize)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, "Image Assurance Checks")

	showTable(&pdf, leftMargin, pdf.GetY()+brSize)
	showTextIntoTable(&pdf, leftMargin, pdf.GetY()+brSize, &[2][5]string{
		{"Approved Base Image","CVE Blacklist","Mallware","MicroEnforcer","OS Package Manager",},
		{passFail,passFail,passFail,passFail,passFail,},
	})

	pdf.SetY( pdf.GetY()+ 2*brSize)
	pdf.SetX(leftMargin)
	showTable(&pdf, leftMargin, pdf.GetY())
	showTextIntoTable(&pdf, leftMargin, pdf.GetY(), &[2][5]string{
		{"OSS License Blacklist","OSS License Whitelist","Package Blacklist","Required Packages","SCAP",},
		{passFail,passFail,passFail,passFail,passFail,},
	})

	pdf.SetY( pdf.GetY()+ 2*brSize)
	pdf.SetX(leftMargin)
	showTable(&pdf, leftMargin, pdf.GetY())
	showTextIntoTable(&pdf, leftMargin, pdf.GetY(), &[2][5]string{
		{"Sensitive Data","Superuser","Ultrabox","Vulnerability Score","Vulnerability Severity",},
		{passFail,passFail,passFail,passFail,passFail,},
	})

	//--- Second page
	pdf.AddPage()
	pdf.SetX(leftMargin)
	pdf.SetY(topMargin)

	// Sensitive Data
	// [List of sensitive data]
	pdf.SetX(leftMargin)
	pdf.SetY(topMargin)
	pdf.SetFont(fontTypeBold, "", 10)
	pdf.Cell(nil, "Sensitive Data")
	pdf.Br(brSize)
	listData,_ := pdf.SplitText(malwareListData, width)
	addMultiLines( &pdf, leftMargin, 15, listData)

	//Malware
	pdf.Br(brSize)
	pdf.SetX(leftMargin)
	pdf.SetFont(fontTypeBold, "", 10)
	pdf.Cell(nil, "Malware")
	pdf.Br(brSize)

	listDataMalware,_ := pdf.SplitText(malwareListData, width)
	addMultiLines( &pdf, leftMargin, 15, listDataMalware)

	pdf.Br(brSize)
	addHr(&pdf, pdf.GetY())

	// line
	pdf.Br(brSize)
	pdf.SetX(leftMargin)

	pdf.SetFont(fontTypeBold, "", 19)
	pdf.Cell(nil, "Detailed Finding Descriptions")
	pdf.Br(brSize)
	pdf.SetX(leftMargin)
	pdf.SetFont(fontType, "", 10)
	pdf.Cell(nil, "This section contains the findings in more detail, ordered by severity")
	pdf.Br(brSize)

	greyBlockH := padding*2.0+13.0
	pdf.SetFont(fontType, "", 11)
	pdf.SetFillColor(236, 239, 241)
	addBlock( &pdf, leftMargin, pdf.GetY(), 200.0, greyBlockH, cveNumber)
	pdf.SetY( pdf.GetY()+greyBlockH)

	pdf.Br(brSize)
	negligbleBlockW := 100.0
	pdf.SetFillColor(200, 236, 252)
	pdf.SetTextColor(0,117, 191)
	addBlock( &pdf, leftMargin, pdf.GetY(), negligbleBlockW, greyBlockH, negligible)

	pdf.SetTextColor(0,0,0)
	pdf.SetFillColor(236, 239, 241)
	addBlock( &pdf, leftMargin+negligbleBlockW+padding*2, pdf.GetY(), 150.0, greyBlockH, cvssScrore)

	pdf.SetY( pdf.GetY()+greyBlockH)

	pdf.Br(brSize)
	addHrGrey( &pdf, pdf.GetY())

	//-- table
	pdf.Br(brSize)
	pdf.SetX(leftMargin)
	tableCellH := 16.0
	tableCellW := 155.0

	pdf.SetLineWidth(0.5)
	pdf.SetStrokeColor(0,0,0)
	pdf.SetTextColor(124, 151, 182)
	addCell( &pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, "Resource")
	addCell( &pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, "Full Resource Name")
	addCell( &pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, "Fixed Version")

	pdf.SetY(pdf.GetY()+tableCellH)
	pdf.SetX(leftMargin)
	pdf.SetTextColor(0,0,0)
	addCell( &pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, resource)
	addCell( &pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, resourceFullName)
	addCell( &pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, fixedVersion)

	pdf.Br(brSize*2)
	pdf.SetX(leftMargin)
	pdf.SetTextColor(124, 151, 182)
	pdf.Cell(nil, "Solution:")
	pdf.Br(brSize)
	pdf.SetX(leftMargin)
	pdf.SetTextColor(0,0,0)
	pdf.Cell(nil, solution)

	pdf.Br(brSize)
	addHrGrey( &pdf, pdf.GetY())

	//--- Third page
	pdf.AddPage()
	pdf.SetX(leftMargin)
	pdf.SetY(topMargin)
	pdf.SetFont(fontTypeBold, "", 10)
	pdf.Cell(nil, "VULNERABILITY DESCRIPTION")
	pdf.Br(brSize)
	pdf.SetFont(fontType, "", 10)

	multilinesVulnDescription,_ := pdf.SplitText(vulnDescription, width-2*padding)
	hBlockVulnDescription := len(multilinesVulnDescription)*14+padding*2
	pdf.SetFillColor(246,249,250)
	pdf.RectFromUpperLeftWithStyle(leftMargin, pdf.GetY(), width, float64(hBlockVulnDescription), "F")
	pdf.Br(padding)
	addMultiLines( &pdf,leftMargin+padding, 15, multilinesVulnDescription )

	pdf.WritePdf(output)
}

func addMultiLines(pdf *gopdf.GoPdf, x, deltaY float64, lines []string)  {
	for _,line := range lines {
		pdf.SetX(x)
		pdf.Cell(nil, line)
		pdf.SetY(pdf.GetY()+deltaY)
	}
}

func addCell(pdf *gopdf.GoPdf, x, y, w, h float64, text string) {
	pdf.RectFromUpperLeftWithStyle(x,y, w+padding, h, "D")
	opt := gopdf.CellOption{
		Align:  gopdf.Middle | gopdf.Left,
		Border: 0,
		Float:  gopdf.Right,
	}
	rect := gopdf.Rect{
		W: w,
		H: h,
	}
	pdf.SetX(x+padding)
	pdf.SetY(y)
	pdf.CellWithOption(&rect, text, opt)
}

func addBlock(pdf *gopdf.GoPdf, x, y, w,h float64, text string)  {
	pdf.RectFromUpperLeftWithStyle(x,y, w, h, "F")
	rect := gopdf.Rect{
		W: w,
		H: h,
	}
	pdf.SetX(x)
	pdf.SetY(y)
	pdf.CellWithOption(&rect, text, cellOption)
}

func addHrGrey(pdf *gopdf.GoPdf, yLeft float64) {
	pdf.SetStrokeColor(236, 239, 241)
	pdf.SetLineWidth(2)
	pdf.Line(leftMargin, yLeft, leftMargin+width, yLeft)
}

func addHr(pdf *gopdf.GoPdf, yLeft float64) {
	pdf.SetStrokeColor(0, 172, 195)
	pdf.SetLineWidth(2)
	pdf.Line(leftMargin, yLeft, leftMargin+width, yLeft)
}

func showColorfulTable(pdf *gopdf.GoPdf, yLeft float64)  {
	pdf.SetFillColor(192,0,0)
	pdf.RectFromUpperLeftWithStyle(leftMargin, yLeft, cellWidth, cellHeight, "FD")
	pdf.RectFromUpperLeftWithStyle(leftMargin, yLeft+cellHeight, cellWidth, cellHeight, "FD")

	pdf.SetFillColor(255,0,0)
	pdf.RectFromUpperLeftWithStyle(leftMargin+cellWidth, yLeft, cellWidth, cellHeight, "FD")
	pdf.RectFromUpperLeftWithStyle(leftMargin+cellWidth, yLeft+cellHeight, cellWidth, cellHeight, "FD")

	pdf.SetFillColor(255,192,0)
	pdf.RectFromUpperLeftWithStyle(leftMargin+2*cellWidth, yLeft, cellWidth, cellHeight, "FD")
	pdf.RectFromUpperLeftWithStyle(leftMargin+2*cellWidth, yLeft+cellHeight, cellWidth, cellHeight, "FD")

	pdf.SetFillColor(255,255,0)
	pdf.RectFromUpperLeftWithStyle(leftMargin+3*cellWidth, yLeft, cellWidth, cellHeight, "FD")
	pdf.RectFromUpperLeftWithStyle(leftMargin+3*cellWidth, yLeft+cellHeight, cellWidth, cellHeight, "FD")

	pdf.SetFillColor(0,112,192)
	pdf.RectFromUpperLeftWithStyle(leftMargin+4*cellWidth, yLeft, cellWidth, cellHeight, "FD")
	pdf.RectFromUpperLeftWithStyle(leftMargin+4*cellWidth, yLeft+cellHeight, cellWidth, cellHeight, "FD")
}

func showTable(pdf *gopdf.GoPdf, xTop, yLeft float64) {
	for i := 0; i < 5; i++ {
		delta := float64(i) * cellWidth
		pdf.RectFromUpperLeftWithStyle(xTop+delta, yLeft, cellWidth, cellHeight, "D")
		pdf.RectFromUpperLeftWithStyle(xTop+delta, yLeft+cellHeight, cellWidth, cellHeight, "D")
	}
}

func showTextIntoTable(pdf *gopdf.GoPdf, xTop, yLeft float64, text *[2][5]string) {
	rect := gopdf.Rect{
		W: cellWidth,
		H: cellHeight,
	}
	for i := 0; i <2; i++ {
		deltaH := float64(i)*cellHeight
		for j :=0; j < 5; j++ {
			pdf.SetX(xTop + float64(j)*cellWidth)
			pdf.SetY(yLeft + deltaH)
			pdf.CellWithOption(&rect, text[i][j], cellOption)
		}
	}
}
