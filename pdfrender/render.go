package pdfrender

import (
	"../data"
	"github.com/signintech/gopdf"
	"log"
	"strconv"
	"time"
)

const (
	leftMargin = 30
	rightMargin = 30
	topMargin = 30
	bottomMargin = 20

	padding = 10
	width = 595-leftMargin-rightMargin
	heightPage = 842- bottomMargin
	brSize = 20

	cellHeight = 25
	cellWidth = width/5.0

	ttfPathRegular = "./pdfrender/calibri.ttf"
	ttfPathBold = "./pdfrender/calibri-bold.ttf"
	longPath = "./pdfrender/logo.png"

	dateFormat = "2006-01-02 15:04"

	fontType = "calibri"
	fontTypeBold = "calibri-bold"
)

var cellOption = gopdf.CellOption{
	Align:  gopdf.Middle | gopdf.Center,
	Border: 0,
	Float:  gopdf.Right,
}

func Render(output string, data *data.Report)  {
	currentDate := time.Now().Format(dateFormat)

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
	pdf.Cell(nil, "Aqua Server – ")
	linkXBegin := pdf.GetX()
	pdf.Cell(nil, data.ServerUrl)
	linkXEnd := pdf.GetX()
	pdf.AddExternalLink(data.ServerUrl, linkXBegin, pdf.GetY(), linkXEnd-linkXBegin, 15)

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
	pdf.Cell(nil, "Image name \"" + data.General.ImageName+"\"")
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
	pdf.Cell(nil, "Registry: " + data.General.Registry)

	pdf.Br(brSize)
	pdf.SetX(leftMargin+padding)
	timeCreated,_ := time.Parse("2006-01-02T15:04:05.999999999Z07:00",data.General.Created)
	pdf.Cell(nil, "Image Creation Date: " + timeCreated.Format(dateFormat))

	pdf.Br(brSize)
	pdf.SetX(leftMargin+padding)
	pdf.Cell(nil, "OS: " + data.General.Os + "(" + data.General.OsVersion+ ")")

	// after Image Name block
	pdf.SetX(leftMargin)
	pdf.SetY(yLine2 + summaryBlochHeight + 2*padding)
	pdf.SetFont(fontTypeBold, "", 10)
	var imageAllowed string
	if data.General.AssuranceResults.Disallowed {imageAllowed = "Disallowed"} else {imageAllowed="Allowed"}

	pdf.Cell(nil, "Image is " + imageAllowed)
	// Block Number of Vulnerabilities
	pdf.Br(brSize*1.5)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, "Image Vulnerabilities")
	pdf.Br(brSize)
	yTable1 := pdf.GetY()
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(0.5)
	showColorfulTable(&pdf, yTable1)

	showTextIntoTable(&pdf, leftMargin, yTable1, &[2][5]string{
		{"CRITICAL","HIGH","MEDIUM","LOW","NEGLIGIBLE",},
		{strconv.Itoa( data.General.Critical), strconv.Itoa(data.General.High),strconv.Itoa(data.General.Medium),strconv.Itoa(data.General.Low),strconv.Itoa(data.General.Negligible),},
	})

	// Image Assurance Policies
	pdf.SetY(yTable1 + cellHeight*2+ brSize)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, "Image Assurance Policies")

	for _,v := range data.General.AssuranceResults.ChecksPerformed {
		pdf.Br(brSize)
		pdf.SetX(leftMargin)
		pdf.Cell(nil, "Policy \"" + v.PolicyName + "\": ")
		if v.Failed {
			pdf.SetTextColor(255, 0, 0)
			pdf.Cell(nil, "FAILED")
		} else {
			pdf.SetTextColor(0, 255, 0)
			pdf.Cell(nil, "PASS")
		}
		pdf.SetTextColor(0, 0, 0)
		checkEndOfPage( &pdf, 2*brSize)
	}
	pdf.Br(brSize)

	// Image Assurance Checks
	checks := data.MappingImageAssuranceChecks()
	checkEndOfPage( &pdf, 2*brSize+2*cellHeight)
	pdf.SetLineWidth(0.5)

	pdf.SetX(leftMargin)
	pdf.Cell(nil, "Image Assurance Checks")
	showTable(&pdf, leftMargin, pdf.GetY()+brSize)
	showTextIntoTable(&pdf, leftMargin, pdf.GetY()+brSize, &[2][5]string{
		{"Approved Base Image","CVE Blacklist","Mallware","MicroEnforcer","OS Package Manager",},
		{GetPassOrFailCheck(checks, "trusted_base_images"),
			GetPassOrFailCheck(checks, "cve_blacklist"),
			GetPassOrFailCheck(checks, "malware"),
			GetPassOrFailCheck(checks, "force_microenforcer"),
			GetPassOrFailCheck(checks, ""),},
	})

	///GetPassOrFailCheck(checks, "")
	pdf.SetY( pdf.GetY()+ 2*brSize)
	pdf.SetX(leftMargin)
	checkEndOfPage( &pdf, brSize+2*cellHeight)

	showTable(&pdf, leftMargin, pdf.GetY())
	showTextIntoTable(&pdf, leftMargin, pdf.GetY(), &[2][5]string{
		{"OSS License Blacklist","OSS License Whitelist","Package Blacklist","Required Packages","SCAP",},
		{
			GetPassOrFailCheck(checks, "license"),
			GetPassOrFailCheck(checks, "whitelisted_licenses"),
			GetPassOrFailCheck(checks, "blacklisted_packages"),
			GetPassOrFailCheck(checks, "required_packages"),
			GetPassOrFailCheck(checks, ""),},
	})

	pdf.SetY( pdf.GetY()+ 2*brSize)
	pdf.SetX(leftMargin)
	checkEndOfPage( &pdf, brSize+2*cellHeight)

	showTable(&pdf, leftMargin, pdf.GetY())
	showTextIntoTable(&pdf, leftMargin, pdf.GetY(), &[2][5]string{
		{"Sensitive Data","Superuser","Ultrabox","Vulnerability Score","Vulnerability Severity",},
		{
			GetPassOrFailCheck(checks, "sensitive_data"),
			GetPassOrFailCheck(checks, "root_user"),
			GetPassOrFailCheck(checks, ""),
			GetPassOrFailCheck(checks, "max_score"),
			GetPassOrFailCheck(checks, "max_severity"),},
	})

	// Sensitive Data
	// [List of sensitive data]

	pdf.Br(brSize*2)
	pdf.SetX(leftMargin)

	pdf.SetFont(fontTypeBold, "", 10)
	pdf.Cell(nil, "Sensitive Data")
	pdf.Br(brSize)
	pdf.SetX(leftMargin)

	pdf.SetFont(fontType, "", 9)
	for _, result := range data.Sensitive.Results {
		addCell( &pdf, leftMargin, pdf.GetY(), width, 15, "Type:")
		addCell( &pdf, leftMargin, pdf.GetY()+15, width, 15, result.Type)
		addCell( &pdf, leftMargin, pdf.GetY()+15, width, 15, "Path:")
		addCell( &pdf, leftMargin, pdf.GetY()+15, width, 15, result.Path)
		pdf.Br(brSize)
		checkEndOfPage( &pdf, brSize+60)
	}
	pdf.Br(brSize)

	//Malware
	checkEndOfPage( &pdf, 100)
	pdf.SetX(leftMargin)
	pdf.SetFont(fontTypeBold, "", 10)
	pdf.Cell(nil, "Malware")
	pdf.Br(brSize)

	pdf.SetFont(fontType, "", 9)
	malwareTitleWidth := 40.0
	for _, result := range data.Malware.Results {
		addCell( &pdf, leftMargin, pdf.GetY(), malwareTitleWidth, 15, "Malware")
		addCell( &pdf, leftMargin+malwareTitleWidth+padding, pdf.GetY(), width-malwareTitleWidth-2*padding, 15, result.Malware)

		addCell( &pdf, leftMargin, pdf.GetY()+15, malwareTitleWidth, 15, "Path")
		addCell( &pdf, leftMargin+malwareTitleWidth+padding, pdf.GetY(), width-malwareTitleWidth-2*padding, 15, result.Path)

		addCell( &pdf, leftMargin, pdf.GetY()+15, malwareTitleWidth, 15, "Hash")
		addCell( &pdf, leftMargin+malwareTitleWidth+padding, pdf.GetY(), width-malwareTitleWidth-2*padding, 15, result.Hash)
		pdf.Br(brSize)
		checkEndOfPage( &pdf, brSize+45)
	}

	pdf.Br(brSize)

	checkEndOfPage( &pdf, heightPage/2)
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

	for _, vuln := range data.Vulnerabilities.Results {
		addVulnBlock( &pdf, vuln)
	}

	pdf.WritePdf(output)
}

func addVulnBlock( pdf *gopdf.GoPdf, vuln data.VulnerabilitiesResultType) {
	greyBlockH := padding*2.0+13.0

	checkEndOfPage( pdf, 2*greyBlockH + 2*brSize )

	pdf.Br(brSize)
	pdf.SetFont(fontType, "", 11)
	pdf.SetFillColor(236, 239, 241)
	addBlock( pdf, leftMargin, pdf.GetY(), 200.0, greyBlockH, vuln.Name)
	pdf.SetY( pdf.GetY()+greyBlockH)

	pdf.Br(brSize)
	negligbleBlockW := 100.0
	pdf.SetFillColor(200, 236, 252)
	pdf.SetTextColor(0,117, 191)
	addBlock( pdf, leftMargin, pdf.GetY(), negligbleBlockW, greyBlockH, vuln.AquaSeverity)

	pdf.SetTextColor(0,0,0)
	pdf.SetFillColor(236, 239, 241)
	addBlock( pdf, leftMargin+negligbleBlockW+padding*2, pdf.GetY(), 150.0, greyBlockH,
		strconv.FormatFloat(vuln.AquaScore, 'f', 2, 64))

	pdf.SetY( pdf.GetY()+greyBlockH)

	pdf.Br(brSize)


	tableCellH := 16.0
	tableCellW := 155.0

	checkEndOfPage( pdf, tableCellH*2+3*brSize+30)
	addHrGrey( pdf, pdf.GetY())

	//-- table
	pdf.Br(brSize)
	pdf.SetX(leftMargin)

	pdf.SetLineWidth(0.5)
	pdf.SetStrokeColor(0,0,0)
	pdf.SetTextColor(124, 151, 182)
	addCell( pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, "Resource")
	addCell( pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, "Full Resource Name")
	addCell( pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, "Fixed Version")

	pdf.SetY(pdf.GetY()+tableCellH)
	pdf.SetX(leftMargin)
	pdf.SetTextColor(0,0,0)
	addCell( pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, vuln.Resource.Name)
	addCell( pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, vuln.Resource.Version)
	addCell( pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, vuln.FixVersion)

	pdf.Br(brSize)
	pdf.SetX(leftMargin)
	pdf.SetTextColor(124, 151, 182)
	pdf.Cell(nil, "Solution:")
	pdf.Br(brSize)
	pdf.SetX(leftMargin)
	pdf.SetTextColor(0,0,0)
	pdf.Cell(nil, vuln.Solution)

	pdf.Br(brSize)

	checkEndOfPage( pdf, heightPage/6)
	addHrGrey( pdf, pdf.GetY())

	pdf.Br(brSize)
	pdf.SetX(leftMargin)
	pdf.SetFont(fontTypeBold, "", 10)
	pdf.Cell(nil, "VULNERABILITY DESCRIPTION")
	pdf.Br(brSize)
	pdf.SetFont(fontType, "", 10)

	multilinesVulnDescription,_ := pdf.SplitText(vuln.Description, width-2*padding)
	hBlockVulnDescription := len(multilinesVulnDescription)*14+padding*2
	pdf.SetFillColor(246,249,250)
	pdf.RectFromUpperLeftWithStyle(leftMargin, pdf.GetY(), width, float64(hBlockVulnDescription), "F")
	pdf.Br(padding)
	addMultiLines( pdf,leftMargin+padding, 15, multilinesVulnDescription )
	pdf.Br(brSize)
	addHr(pdf, pdf.GetY())
}

func GetPassOrFailCheck(m map[string]bool, key string) string  {
	var result string
	v, ok := m[key]
	if ok {
		if v {
			result = "PASS"
		} else {
			result = "FAIL"
		}
	}
	return result
}

func checkEndOfPage(pdf *gopdf.GoPdf, deltaY float64) {
	if ( pdf.GetY() + deltaY ) > heightPage  {
		pdf.AddPage()
		pdf.SetY( topMargin)
		pdf.SetLineWidth(0.5)
	}
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
			if text[i][j] == "FAIL" {
				pdf.SetTextColor(255,0,0)
			}
			if text[i][j] == "PASS" {
				pdf.SetTextColor(0,255,0)
			}
			pdf.CellWithOption(&rect, text[i][j], cellOption)
			pdf.SetTextColor(0,0,0)
		}
	}
}
