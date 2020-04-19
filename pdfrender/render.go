package pdfrender

import (
	"../data"
	"github.com/signintech/gopdf"
	"log"
	"strconv"
	"strings"
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

	numberCellInRow = 3
	cellHeight = 25
	cellWidth = width/numberCellInRow

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

var imageAssurance = map[string]string {
	"blacklisted_packages": "Blacklisted Packages",
	"custom_checks":        "Custom Checks",
	"cve_blacklist":        "Blacklisted CVE",
	"force_microenforcer":  "MicroEnforcer Deployed",
	"license":              "Approved Licenses",
	"malware":              "Malware",
	"max_score":            "Maximum Vulnerability Score",
	"max_severity":         "Maximum Vulnerability Severity",
	"partial_results":      "Partial Scan Results",
	"required_packages":    "Required Packages",
	"root_user":            "Run as Superuser",
	"sensitive_data":       "Sensitive Data",
	"trusted_base_images":  "Trusted Base Image",
	"whitelisted_licenses": "Whitelisted Licenses",
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
	showTextInto5thColumnsTable(&pdf, leftMargin, yTable1, &[2][5]string{
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
	pdf.Br(brSize)
	checkEndOfPage( &pdf, brSize+80)
	pdf.SetX(leftMargin)

	pdf.SetFont(fontTypeBold, "", 12)
	pdf.Cell(nil, "Sensitive Data")
	pdf.Br(brSize)
	pdf.SetX(leftMargin)

	pdf.SetFont(fontType, "", 9)
	if data.Sensitive.Count > 0 {
		for _, result := range data.Sensitive.Results {
			addCell( &pdf, leftMargin, pdf.GetY(), width, 15, "Type:")
			addCell( &pdf, leftMargin, pdf.GetY()+15, width, 15, result.Type)
			addCell( &pdf, leftMargin, pdf.GetY()+15, width, 15, "Path:")
			addCell( &pdf, leftMargin, pdf.GetY()+15, width, 15, result.Path)
			pdf.Br(brSize)
			checkEndOfPage( &pdf, brSize+60)
		}
	} else {
		pdf.SetFont(fontType, "", 10)
		pdf.Cell(nil, "None found.")
	}

	pdf.Br(brSize)

	//Malware
	checkEndOfPage( &pdf, 100)
	pdf.SetX(leftMargin)
	pdf.SetFont(fontTypeBold, "", 12)
	pdf.Cell(nil, "Malware")
	pdf.Br(brSize)
	pdf.SetX(leftMargin)

	pdf.SetFont(fontType, "", 9)

	if data.Malware.Count > 0 {
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
	} else {
		pdf.SetFont(fontType, "", 10)
		pdf.Cell(nil, "None found")
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
	addBlock( pdf, leftMargin, pdf.GetY(), 200.0, greyBlockH, "Vulnerability: " +vuln.Name)
	pdf.SetY( pdf.GetY()+greyBlockH)

	pdf.Br(brSize)
	negligbleBlockW := 100.0
	pdf.SetFillColor(200, 236, 252)
	pdf.SetTextColor(0,117, 191)
	addBlock( pdf, leftMargin, pdf.GetY(), negligbleBlockW, greyBlockH, "Severity: " + strings.Title(vuln.AquaSeverity))

	pdf.SetTextColor(0,0,0)
	pdf.SetFillColor(236, 239, 241)
	addBlock( pdf, leftMargin+negligbleBlockW+padding*2, pdf.GetY(), 150.0, greyBlockH,
		"Score: " + strconv.FormatFloat(vuln.AquaScore, 'f', 2, 64))

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
	addCell( pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, "Fix Version")

	pdf.SetY(pdf.GetY()+tableCellH)
	pdf.SetX(leftMargin)
	pdf.SetTextColor(0,0,0)
	var name, version, fixVersion string
	if vuln.Resource.Name != "" {
		name = vuln.Resource.Name
	} else {
		name = "No name"
	}

	if vuln.Resource.Version != "" {
		version = vuln.Resource.Version
	} else {
		version = "No version"
	}

	if vuln.FixVersion != "" {
		fixVersion = vuln.FixVersion
	} else {
		fixVersion = "no fix"
	}
	addCell( pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, name)
	addCell( pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, version)
	addCell( pdf, pdf.GetX(), pdf.GetY(), tableCellW, tableCellH, fixVersion )

	pdf.Br(brSize*1.3)
	pdf.SetX(leftMargin)
	pdf.SetTextColor(124, 151, 182)
	pdf.Cell(nil, "Solution:")
	pdf.Br(brSize*0.8)
	pdf.SetX(leftMargin)
	pdf.SetTextColor(0,0,0)

	if vuln.Solution != "" {
		multilineSolution,_ := pdf.SplitText(vuln.Solution, width)
		addMultiLines( pdf, leftMargin, 12, multilineSolution )
		pdf.Br(brSize*0.5)
	} else {
		pdf.Cell(nil, "none")
		pdf.Br(brSize)
	}
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
	localCellWidth := width/5.0
	pdf.SetFillColor(192,0,0)
	pdf.RectFromUpperLeftWithStyle(leftMargin, yLeft, localCellWidth, cellHeight, "FD")
	pdf.RectFromUpperLeftWithStyle(leftMargin, yLeft+cellHeight, localCellWidth, cellHeight, "FD")

	pdf.SetFillColor(255,0,0)
	pdf.RectFromUpperLeftWithStyle(leftMargin+localCellWidth, yLeft, localCellWidth, cellHeight, "FD")
	pdf.RectFromUpperLeftWithStyle(leftMargin+localCellWidth, yLeft+cellHeight, localCellWidth, cellHeight, "FD")

	pdf.SetFillColor(255,192,0)
	pdf.RectFromUpperLeftWithStyle(leftMargin+2*localCellWidth, yLeft, localCellWidth, cellHeight, "FD")
	pdf.RectFromUpperLeftWithStyle(leftMargin+2*localCellWidth, yLeft+cellHeight, localCellWidth, cellHeight, "FD")

	pdf.SetFillColor(255,255,0)
	pdf.RectFromUpperLeftWithStyle(leftMargin+3*localCellWidth, yLeft, localCellWidth, cellHeight, "FD")
	pdf.RectFromUpperLeftWithStyle(leftMargin+3*localCellWidth, yLeft+cellHeight, localCellWidth, cellHeight, "FD")

	pdf.SetFillColor(0,112,192)
	pdf.RectFromUpperLeftWithStyle(leftMargin+4*localCellWidth, yLeft, localCellWidth, cellHeight, "FD")
	pdf.RectFromUpperLeftWithStyle(leftMargin+4*localCellWidth, yLeft+cellHeight, localCellWidth, cellHeight, "FD")
}

func addCellOfTable(pdf *gopdf.GoPdf, xTop, yLeft float64) {
	pdf.RectFromUpperLeftWithStyle(xTop, yLeft, cellWidth, cellHeight, "D")
	pdf.RectFromUpperLeftWithStyle(xTop, yLeft+cellHeight, cellWidth, cellHeight, "D")
}

func showTable(pdf *gopdf.GoPdf, xTop, yLeft float64, rows int) {
	for i := 0; i < rows; i++ {
		delta := float64(i) * cellWidth
		addCellOfTable( pdf, xTop+delta, yLeft)
	}
}

func addTextToCellOfTable(pdf *gopdf.GoPdf, xTop, yLeft, width float64, text string) {
	rect := gopdf.Rect{
		W: width,
		H: cellHeight,
	}
	pdf.SetX(xTop)
	pdf.SetY(yLeft)

	if text == "FAIL" {
		pdf.SetTextColor(255,0,0)
	}
	if text == "PASS" {
		pdf.SetTextColor(0,255,0)
	}

	pdf.CellWithOption(&rect, text, cellOption)
	pdf.SetTextColor(0,0,0)

}

func showTextInto5thColumnsTable(pdf *gopdf.GoPdf, xTop, yLeft float64, text *[2][5]string)  {
	localCellWidth := width/5.0
	for i := 0; i <2; i++ {
		deltaH := float64(i)*cellHeight
		for j :=0; j < 5; j++ {
			addTextToCellOfTable( pdf, xTop + float64(j)*localCellWidth, yLeft + deltaH, localCellWidth, text[i][j] )
		}
	}
}

func showTextIntoTable(pdf *gopdf.GoPdf, xTop, yLeft float64, text *[2][5]string, row int) {
	for i := 0; i <2; i++ {
		deltaH := float64(i)*cellHeight
		for j :=0; j < row; j++ {
			addTextToCellOfTable( pdf, xTop + float64(j)*cellWidth, yLeft + deltaH, cellWidth, text[i][j] )
		}
	}
}

func showImageAssuranceChecks( pdf *gopdf.GoPdf, checksDa []data.CheckPerformedType)  {
	var checks [2][5]string
	var i int
	for i=0;i < len(checksDa); {
		description, ok := imageAssurance[ checksDa[i].Control ]
		if ok {
			checks[0][i%numberCellInRow] = description
		} else {
			checks[0][i%numberCellInRow] = checksDa[i].Control
		}
		if checksDa[i].Failed {
			checks[1][i%numberCellInRow] = "FAIL"
		} else {
			checks[1][i%numberCellInRow] = "PASS"
		}
		i++
		if i%numberCellInRow == 0 {
			pdf.SetY( pdf.GetY()+ padding)
			pdf.SetX(leftMargin)
			checkEndOfPage( pdf, brSize+2*cellHeight)
			showTable(pdf, leftMargin, pdf.GetY()+brSize, numberCellInRow)
			showTextIntoTable(pdf, leftMargin, pdf.GetY()+brSize, &checks, numberCellInRow)
		}
	}
	if i%numberCellInRow != 0 {
		pdf.SetY( pdf.GetY()+ padding)
		pdf.SetX(leftMargin)
		checkEndOfPage( pdf, brSize+2*cellHeight)
		showTable(pdf, leftMargin, pdf.GetY()+brSize, i%numberCellInRow)
		showTextIntoTable(pdf, leftMargin, pdf.GetY()+brSize, &checks, i%numberCellInRow)
	}

}
