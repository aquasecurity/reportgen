package pdfrender

import "github.com/signintech/gopdf"

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
