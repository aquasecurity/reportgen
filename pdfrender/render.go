package pdfrender

import (
	"../data"
	"github.com/signintech/gopdf"
	"strconv"
	"time"
)

var pdf *gopdf.GoPdf

func Render(output string, report *data.Report) error {
	currentDate := time.Now().Format(dateFormat)
	var title string
	switch report.RequestType {
	case data.ImageRequest: title = "Image"
	case data.HostRequest:  title = "Host"
	}

	pdf = &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{ PageSize: *gopdf.PageSizeA4 }) //595.28, 841.89 = A4
	pdf.SetLeftMargin(leftMargin)
	pdf.AddPage()

	err := pdf.AddTTFFont(fontType, ttfPathRegular)
	if err != nil {
		return err
	}
	err = pdf.AddTTFFont(fontTypeBold, ttfPathBold )
	if err != nil {
		return err
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
	pdf.SetFont(fontType, "", 14)
	pdf.Cell(nil, title + " Vulnerability Report")

	pdf.SetX(leftMargin)
	pdf.SetY(yTitleBase + 3*padding+15)
	pdf.SetTextColor(0,0,0)
	pdf.SetFont(fontType, "", 10)
	pdf.Cell(nil, "Aqua Server – ")
	linkXBegin := pdf.GetX()
	pdf.Cell(nil, report.ServerUrl)
	linkXEnd := pdf.GetX()
	pdf.AddExternalLink(report.ServerUrl, linkXBegin, pdf.GetY(), linkXEnd-linkXBegin, 15)

	// line after 1
	yLine1 := pdf.GetY()+ 2* padding
	pdf.SetLineWidth(1.0)
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
	if report.RequestType == data.HostRequest {
		summaryBlochHeight -= rowSize
	}

	pdf.SetX(leftMargin)
	pdf.SetY(yLine2)
	addHr(yLine2)

	pdf.SetFillColor(246,249,250)
	pdf.RectFromUpperLeftWithStyle(leftMargin, yLine2, width, summaryBlochHeight, "F")
	pdf.SetY(yLine2+padding)
	pdf.SetX(leftMargin+padding)
	pdf.SetFont(fontTypeBold, "", 12)
	pdf.Cell(nil, title + " name \"" + report.General.ImageName+"\"")
	pdf.Br(brSize)
	pdf.SetFont(fontType, "", 10)

	opt := gopdf.CellOption{
		Align:  gopdf.Right,
		Border: 0,
		Float:  gopdf.Right,
	}
	rect := gopdf.Rect{
		W: width-padding,
		H: 10,
	}
	pdf.CellWithOption(&rect, "Report generated on "+currentDate, opt)

	pdf.SetFont(fontTypeBold, "", 10)
	if report.RequestType == data.ImageRequest {
		pdf.Br(brSize)
		pdf.SetX(leftMargin+padding)
		pdf.Cell(nil, "Registry: " + report.General.Registry)

		pdf.Br(brSize)
		pdf.SetX(leftMargin+padding)
		timeCreated,_ := time.Parse("2006-01-02T15:04:05.999999999Z07:00", report.General.Created)
		pdf.Cell(nil, "Image Creation Date: " + timeCreated.Format(dateFormat))
	}

	pdf.Br(brSize)
	pdf.SetX(leftMargin+padding)

	os := "OS: " + report.General.Os
	if report.General.OsVersion != "" {
		os += " (" + report.General.OsVersion+ ")"
	}
	pdf.Cell(nil, os)

	if report.RequestType == data.HostRequest {
		pdf.Br(brSize)
		pdf.SetX(leftMargin+padding)
		pdf.Cell(nil, "Address: " + report.General.Address )
	}

	pdf.Br(brSize)
	pdf.SetX(leftMargin+padding)
	scanDate, _ := time.Parse("2006-01-02T15:04:05.999999999Z07:00", report.General.ScanDate)
	pdf.Cell(nil, "Scan date: " + scanDate.Format(dateFormat) )

	// after Image Name block
	pdf.SetX(leftMargin)
	pdf.SetY(yLine2 + summaryBlochHeight + 2*padding)
	pdf.SetFont(fontTypeBold, "", 12)

	pdf.Cell(nil, title + " is ")

	addCompliantText(report.General.AssuranceResults.Disallowed)

	// Block Number of Vulnerabilities
	pdf.Br(brSize*1.5)
	pdf.SetX(leftMargin)
	pdf.Cell(nil, title + " Vulnerabilities")
	pdf.Br(brSize)
	yTable1 := pdf.GetY()
	pdf.SetStrokeColor(0, 0, 0)
	pdf.SetLineWidth(0.5)

	pdf.SetFont(fontTypeBold, "", 10)
	showColorfulTable(yTable1)
	showTextIntoFiveColumnsTable(leftMargin, yTable1, &[2][5]string{
		{"CRITICAL","HIGH","MEDIUM","LOW","NEGLIGIBLE",},
		{strconv.Itoa( report.General.Critical), strconv.Itoa(report.General.High),strconv.Itoa(report.General.Medium),strconv.Itoa(report.General.Low),strconv.Itoa(report.General.Negligible),},
	})

	// Image Assurance Policies
	pdf.SetY(yTable1 + cellHeight*2+ brSize)
	pdf.SetX(leftMargin)
	pdf.SetFont(fontTypeBold, "", 12)
	pdf.Cell(nil, title + " Assurance Policies")
	pdf.SetFont(fontTypeBold, "", 10)

	policiesTotal, policiesChecks := report.GetImageAssurancePolicies()

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
		pdf.Cell(nil, title + " Assurance Checks for " + name)
		showImageAssuranceChecks(policiesChecks[name])
		pdf.Br(brSize)
	}

	// Sensitive Data
	checkEndOfPageWithBr( rowSize*4)

	pdf.SetFont(fontTypeBold, "", 12)
	pdf.Cell(nil, "Sensitive Data")
	pdf.Br(brSize)
	pdf.SetX(leftMargin)

	pdf.SetFont(fontType, "", 9)
	if report.Sensitive != nil && report.Sensitive.Count > 0 {
		for _, result := range report.Sensitive.Results {
			addCellText( leftMargin, pdf.GetY(), width, 15, "Type:")
			addCellText( leftMargin, pdf.GetY()+15, width, 15, result.Type)
			addCellText( leftMargin, pdf.GetY()+15, width, 15, "Path:")
			addCellText( leftMargin, pdf.GetY()+15, width, 15, result.Path)
			checkEndOfPageWithBr( brSize+60)
		}
	} else {
		pdf.SetFont(fontType, "", 10)
		pdf.Cell(nil, "None found.")
	}
	//Malware
	checkEndOfPageWithBr( 100)
	pdf.SetFont(fontTypeBold, "", 12)
	pdf.Cell(nil, "Malware")
	pdf.Br(brSize)
	pdf.SetX(leftMargin)

	pdf.SetFont(fontType, "", 9)

	if report.Malware != nil && report.Malware.Count > 0 {
		malwareTitleWidth := 40.0+2*padding
		for _, result := range report.Malware.Results {
			addCellText( leftMargin, pdf.GetY(), malwareTitleWidth, 15, "Malware")
			addCellText( leftMargin+malwareTitleWidth, pdf.GetY(), width-malwareTitleWidth, rowSize, result.Malware)

			addCellText( leftMargin, pdf.GetY()+15, malwareTitleWidth, 15, "Path")
			addCellText( leftMargin+malwareTitleWidth, pdf.GetY(), width-malwareTitleWidth, rowSize, result.Path)

			addCellText( leftMargin, pdf.GetY()+15, malwareTitleWidth, 15, "Hash")
			addCellText( leftMargin+malwareTitleWidth, pdf.GetY(), width-malwareTitleWidth, rowSize, result.Hash)
			checkEndOfPageWithBr( brSize+45)
		}
	} else {
		pdf.SetFont(fontType, "", 10)
		pdf.Cell(nil, "None found")
	}

	// Scan History
	if report.ScanHistory != nil {
		checkEndOfPageWithBr( rowSize*3+padding*3+brSize)

		pdf.SetFont(fontTypeBold, "", 12)
		pdf.Cell(nil, "Scan History")
		pdf.Br(brSize)
		showScanHistory( report.ScanHistory)
	}
	// end of ScanHistory

	// The bench results for a host
	if report.BenchResults != nil {
		addBenchResults( report.BenchResults)
	}

	// end of BenchResults

	checkEndOfPageWithBr( heightPage/2)
	addHr( pdf.GetY())
	// line
	pdf.Br(brSize)
	pdf.SetX(leftMargin)

	pdf.SetFont(fontTypeBold, "", 19)
	pdf.Cell(nil, "Detailed Finding Descriptions")
	pdf.Br(brSize)
	pdf.SetX(leftMargin)
	pdf.SetFont(fontType, "", 10)
	pdf.Cell(nil, "This section contains the findings in more detail, ordered by severity")

	for _, vuln := range report.Vulnerabilities.Results {
		addVulnBlock( vuln)
	}
	return pdf.WritePdf(output)
}

