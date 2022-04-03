package pdfrender

import (
	"github.com/signintech/gopdf"
	"strings"
)

const (
	leftMargin   = 30
	rightMargin  = 30
	topMargin    = 30
	bottomMargin = 20

	padding    = 10
	width      = 595 - leftMargin - rightMargin
	heightPage = 842 - topMargin
	brSize     = 20

	rowSize = 15

	numberCellInRow = 3
	cellHeight      = 25
	cellWidth       = width / numberCellInRow

	ttfPathRegular = "./assets/calibri.ttf"
	ttfPathBold    = "./assets/calibri-bold.ttf"
	longPath       = "./assets/logo.png"

	dateFormat = "2006-01-02 15:04"

	fontType     = "calibri"
	fontTypeBold = "calibri-bold"
)

var cellOption = gopdf.CellOption{
	Align:  gopdf.Middle | gopdf.Center,
	Border: 0,
	Float:  gopdf.Right,
}

var imageAssurance = map[string]string{
	"blacklisted_packages": "Denied Packages",
	"custom_checks":        "Custom Checks",
	"cve_blacklist":        "Denied CVE",
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
	"whitelisted_licenses": "Allowed Licenses",
}

func checkEndOfPageWithoutBr(deltaY float64) {
	checkEndOfPage(deltaY, false)
}

func checkEndOfPageWithBr(deltaY float64) {
	checkEndOfPage(deltaY, true)
}

func checkEndOfPage(deltaY float64, needBr bool) {
	if (pdf.GetY() + deltaY) > heightPage {
		pdf.AddPage()
		pdf.SetY(topMargin)
		pdf.SetLineWidth(0.5)
	} else {
		if needBr {
			pdf.Br(brSize)
		}
	}
}

func addHrGreyH(yLeft, h float64) {
	pdf.SetStrokeColor(236, 239, 241)
	pdf.SetLineWidth(h)
	pdf.Line(leftMargin, yLeft, leftMargin+width, yLeft)
}

func addHrGrey(pdf *gopdf.GoPdf, yLeft float64) {
	addHrGreyH(yLeft, 2)
}

func addHr(yLeft float64) {
	pdf.SetStrokeColor(0, 172, 195)
	pdf.SetLineWidth(1)
	pdf.Line(leftMargin, yLeft, leftMargin+width, yLeft)
}

func addCompliantText(disallowed bool) {
	if disallowed {
		pdf.SetTextColor(255, 151, 47)
		pdf.Cell(nil, "Non-Compliant")
	} else {
		pdf.SetTextColor(0, 255, 0)
		pdf.Cell(nil, "Compliant")
	}
	pdf.SetTextColor(0, 0, 0)
}

func setLightGrayBackgroundColor() {
	pdf.SetFillColor(247, 248, 250)
}
func setDarkGrayBackgroundColor() {
	pdf.SetFillColor(223, 223, 223)
}

func setDefaultBackgroundColor() {
	pdf.SetFillColor(255, 255, 255)
}

func setCriticalBackgroundColor() {
	pdf.SetFillColor(192, 0, 0)
}

func setHighBackgroundColor() {
	pdf.SetFillColor(255, 0, 0)
}

func setMediumBackgroundColor() {
	pdf.SetFillColor(255, 192, 0)
}
func setLowBackgroundColor() {
	pdf.SetFillColor(255, 255, 0)
}
func setNegligibleBackgroundColor() {
	pdf.SetFillColor(0, 112, 192)
}

func setOrangeBackgroundColor() {
	pdf.SetFillColor(255, 165, 0)
}

func setDarkGreenBackgroundColor() {
	pdf.SetFillColor(76, 175, 81)
}

func setGreyBackgroundColor() {
	pdf.SetFillColor(169, 169, 169)
}

func SetSeverityColor(severity string) {
	switch strings.ToLower(severity) {
	case "critical":
		setCriticalBackgroundColor()
	case "high":
		setHighBackgroundColor()
	case "medium":
		setMediumBackgroundColor()
	case "low":
		setLowBackgroundColor()
	case "negligible":
		setNegligibleBackgroundColor()
	default:
		pdf.SetFillColor(200, 236, 252)
	}
}
