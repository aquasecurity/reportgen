package pdfrender

import (
	"github.com/signintech/gopdf"
	"regexp"
)

func splitString( pdf *gopdf.GoPdf, input *string, width float64) []string {
	re := regexp.MustCompile(`[[:cntrl:]]|[\x{FFFD}]`)
	pureString := re.ReplaceAllString(*input, "")
	lines, err := pdf.SplitText(pureString, width)
	if err != nil {
		return []string { pureString }
	}
	var emptyString int
	for i,v := range lines {
		if v == "" {
			copy(lines[i:], lines[i+1:])
			lines = lines[:len(lines)-1]
			emptyString ++
		}
	}
	if emptyString > 0 {
		lines = lines[:len(lines)-emptyString]
	}
	return lines
}

func addMultiLines(pdf *gopdf.GoPdf, x, deltaY float64, lines []string)  {
	for _,line := range lines {
		pdf.SetX(x)
		pdf.Cell(nil, line)
		pdf.SetY(pdf.GetY()+deltaY)
	}
}

