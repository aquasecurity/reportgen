package pdfrender

import (
	"../data"
	"github.com/signintech/gopdf"
	"log"
	"strconv"
	"time"
)

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

	// Block Summary
	yLine2 := pdf.GetY()+padding
	summaryBlochHeight := 125.0
	pdf.SetX(leftMargin)
	pdf.SetY(yLine2)
	addHr(&pdf, yLine2)

	pdf.SetFillColor(246,249,250)
	pdf.RectFromUpperLeftWithStyle(leftMargin, yLine2, width, summaryBlochHeight, "F")
	pdf.SetY(yLine2+padding)
	pdf.SetX(leftMargin+padding)
	pdf.SetFont(fontTypeBold, "", 12)
	pdf.Cell(nil, "Image name \"" + data.General.ImageName+"\"")
	pdf.Br(brSize)
	pdf.SetFont(fontType, "", 10)

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
	pdf.SetFont(fontTypeBold, "", 10)
	pdf.Cell(nil, "Registry: " + data.General.Registry)

	pdf.Br(brSize)
	pdf.SetX(leftMargin+padding)
	timeCreated,_ := time.Parse("2006-01-02T15:04:05.999999999Z07:00",data.General.Created)
	pdf.Cell(nil, "Image Creation Date: " + timeCreated.Format(dateFormat))

	pdf.Br(brSize)
	pdf.SetX(leftMargin+padding)
	pdf.Cell(nil, "OS: " + data.General.Os + "(" + data.General.OsVersion+ ")")

	pdf.Br(brSize)
	pdf.SetX(leftMargin+padding)
	scanDate, _ := time.Parse("2006-01-02T15:04:05.999999999Z07:00", data.General.ScanDate)
	pdf.Cell(nil, "Scan date: " + scanDate.Format(dateFormat) )

	// after Image Name block
	pdf.SetX(leftMargin)
	pdf.SetY(yLine2 + summaryBlochHeight + 2*padding)
	pdf.SetFont(fontTypeBold, "", 12)

	pdf.Cell(nil, "Image is ")
	if data.General.AssuranceResults.Disallowed {
		pdf.SetTextColor(255,151,47)
		pdf.Cell(nil, "Non-Compliant")
	} else {
		pdf.SetTextColor(0,255,0)
		pdf.Cell(nil, "Compliant")
	}
	pdf.SetTextColor(0,0,0)

	// Block Number of Vulnerabilities
	pdf.Br(brSize*1.5)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, "Image Vulnerabilities")
	pdf.Br(brSize)
	yTable1 := pdf.GetY()
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(0.5)

	pdf.SetFont(fontTypeBold, "", 10)
	showColorfulTable(&pdf, yTable1)
	showTextIntoFiveColumnsTable(&pdf, leftMargin, yTable1, &[2][5]string{
		{"CRITICAL","HIGH","MEDIUM","LOW","NEGLIGIBLE",},
		{strconv.Itoa( data.General.Critical), strconv.Itoa(data.General.High),strconv.Itoa(data.General.Medium),strconv.Itoa(data.General.Low),strconv.Itoa(data.General.Negligible),},
	})

	// Image Assurance Policies
	pdf.SetY(yTable1 + cellHeight*2+ brSize)
	pdf.SetX(leftMargin)
	pdf.SetFont(fontTypeBold, "", 12)
	pdf.Cell(nil, "Image Assurance Policies")
	pdf.SetFont(fontTypeBold, "", 10)

	policiesTotal, policiesChecks := data.GetImageAssurancePolicies()

	for name,value := range policiesTotal {
		pdf.Br(brSize)
		pdf.SetX(leftMargin)
		pdf.Cell(nil, "Policy \"" + name + "\": ")
		if value {
			pdf.SetTextColor(255, 0, 0)
			pdf.Cell(nil, "FAILED")
		} else {
			pdf.SetTextColor(0, 255, 0)
			pdf.Cell(nil, "PASS")
		}
		pdf.SetTextColor(0, 0, 0)
		// Image Assurance Checks
		pdf.Br(brSize)
		pdf.SetX(leftMargin)
		pdf.Cell(nil, "Image Assurance Checks for " + name)
		showImageAssuranceChecks(&pdf, policiesChecks[name])
		pdf.Br(brSize)
	}

	// Sensitive Data
	checkEndOfPageWithBr( &pdf, brSize+80)

	pdf.SetFont(fontTypeBold, "", 12)
	pdf.Cell(nil, "Sensitive Data")
	pdf.Br(brSize)
	pdf.SetX(leftMargin)

	pdf.SetFont(fontType, "", 9)
	if data.Sensitive.Count > 0 {
		for _, result := range data.Sensitive.Results {
			addCellText( &pdf, leftMargin, pdf.GetY(), width, 15, "Type:")
			addCellText( &pdf, leftMargin, pdf.GetY()+15, width, 15, result.Type)
			addCellText( &pdf, leftMargin, pdf.GetY()+15, width, 15, "Path:")
			addCellText( &pdf, leftMargin, pdf.GetY()+15, width, 15, result.Path)
			checkEndOfPageWithBr( &pdf, brSize+60)
		}
	} else {
		pdf.SetFont(fontType, "", 10)
		pdf.Cell(nil, "None found.")
	}
	//Malware
	checkEndOfPageWithBr( &pdf, 100)
	pdf.SetFont(fontTypeBold, "", 12)
	pdf.Cell(nil, "Malware")
	pdf.Br(brSize)
	pdf.SetX(leftMargin)

	pdf.SetFont(fontType, "", 9)

	if data.Malware.Count > 0 {
		malwareTitleWidth := 40.0+2*padding
		for _, result := range data.Malware.Results {
			addCellText( &pdf, leftMargin, pdf.GetY(), malwareTitleWidth, 15, "Malware")
			addCellText( &pdf, leftMargin+malwareTitleWidth, pdf.GetY(), width-malwareTitleWidth, rowSize, result.Malware)

			addCellText( &pdf, leftMargin, pdf.GetY()+15, malwareTitleWidth, 15, "Path")
			addCellText( &pdf, leftMargin+malwareTitleWidth, pdf.GetY(), width-malwareTitleWidth, rowSize, result.Path)

			addCellText( &pdf, leftMargin, pdf.GetY()+15, malwareTitleWidth, 15, "Hash")
			addCellText( &pdf, leftMargin+malwareTitleWidth, pdf.GetY(), width-malwareTitleWidth, rowSize, result.Hash)
			checkEndOfPageWithBr( &pdf, brSize+45)
		}
	} else {
		pdf.SetFont(fontType, "", 10)
		pdf.Cell(nil, "None found")
	}

	pdf.Br(brSize)

	checkEndOfPageWithBr( &pdf, heightPage/2)
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

