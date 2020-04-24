package pdfrender

import (
	"../data"
	"github.com/signintech/gopdf"
)

func showTextIntoFiveColumnsTable(pdf *gopdf.GoPdf, xTop, yLeft float64, text *[2][5]string)  {
	localCellWidth := width/5.0
	for i := 0; i <2; i++ {
		deltaH := float64(i)*cellHeight
		for j :=0; j < 5; j++ {
			addTextToCellOfTable( pdf, xTop + float64(j)*localCellWidth, yLeft + deltaH, localCellWidth, text[i][j] )
		}
	}
}

func showTextIntoTable(pdf *gopdf.GoPdf, xTop, yLeft float64, text *[2][5]string, column int) {
	for i := 0; i <2; i++ {
		deltaH := float64(i)*cellHeight
		for j :=0; j < column; j++ {
			addTextToCellOfTable( pdf, xTop + float64(j)*cellWidth, yLeft + deltaH, cellWidth, text[i][j] )
		}
	}
}

func showTwoLineTable(pdf *gopdf.GoPdf, xTop, yLeft float64, columns int) {
	for i := 0; i < columns; i++ {
		delta := float64(i) * cellWidth
		pdf.RectFromUpperLeftWithStyle(xTop+delta, yLeft, cellWidth, cellHeight, "D")
		pdf.RectFromUpperLeftWithStyle(xTop+delta, yLeft+cellHeight, cellWidth, cellHeight, "D")
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

func addCellBorder(pdf *gopdf.GoPdf, x, y, w, h float64) {
	pdf.RectFromUpperLeftWithStyle(x,y, w, h, "D")
}

func addCellText(pdf *gopdf.GoPdf, x, y, w, h float64, text string) {
	addCellBorder(pdf, x,y,w,h)
	opt := gopdf.CellOption{
		Align:  gopdf.Middle | gopdf.Left,
		Border: 0,
		Float:  gopdf.Right,
	}
	rect := gopdf.Rect{
		W: w-2*padding,
		H: h,
	}
	pdf.SetX(x+padding)
	pdf.SetY(y)
	pdf.CellWithOption(&rect, text, opt)
	pdf.SetX(pdf.GetX()+padding)
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
			checkEndOfPageWithoutBr( pdf, brSize+2*cellHeight)
			showTwoLineTable(pdf, leftMargin, pdf.GetY()+brSize, numberCellInRow)
			showTextIntoTable(pdf, leftMargin, pdf.GetY()+brSize, &checks, numberCellInRow)
			pdf.SetY( pdf.GetY())
			pdf.SetX(leftMargin)
		}
	}
	if i%numberCellInRow != 0 {
		checkEndOfPageWithBr( pdf, brSize+2*cellHeight)
		showTwoLineTable(pdf, leftMargin, pdf.GetY()+brSize, i%numberCellInRow)
		showTextIntoTable(pdf, leftMargin, pdf.GetY()+brSize, &checks, i%numberCellInRow)
		pdf.SetY( pdf.GetY()+ padding)
		pdf.SetX(leftMargin)
	}
}