package pdfrender

import "github.com/signintech/gopdf"

func splitString( pdf *gopdf.GoPdf, data *string, width float64) ([]string) {
	lines, err := pdf.SplitText(*data, width)
	if err != nil {
		return []string{}
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

