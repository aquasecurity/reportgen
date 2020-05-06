package pdfrender

import (
	"../data"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func addBenchResults( benchResults *data.BenchResultsType ) {
	checkEndOfPageWithBr( rowSize*10 )
	addHr(pdf.GetY())
	pdf.Br(padding)

	pdf.SetFont(fontTypeBold, "", 13)
	pdf.Cell(nil, "Compliance Results")
	pdf.Br(brSize)
	showTestBlock(benchResults.Cis.Result.Tests, "CIS")
	showTestBlock(benchResults.KubeBench.Result.Tests, "kube-bench")
	showTestBlock(benchResults.Linux.Result.Tests, "Linux")
	showTestBlock(benchResults.Openshift.Result.Tests, "Openshift")
}

func showTestBlock( tests []data.TestBenchType, title string )  {
	if len(tests) == 0 {
		return
	}
	checkEndOfPageWithoutBr(4*rowSize+4*padding)
	pdf.SetFont(fontTypeBold, "", 10)
	pdf.Cell(nil, title)

	pdf.SetFont(fontType, "", 10)
	pdf.Br(brSize)
	for i, test := range tests  {
		checkEndOfPageWithoutBr(2*rowSize+2*padding)
		title := fmt.Sprintf("%d. %s", i+1, test.Desc)
		pdf.Cell(nil, title)
		pdf.Br(brSize)
		showColorfulData( test.Fail, test.Warn, test.Info, test.Pass)
		pdf.Br(padding)

		columnNumW := 40.0
		columnStatusW := 50.0
		columnDescrW := 210.0
		columnInfoW := 240.0
		for _, testResults := range test.Results {
			description := splitString( &testResults.TestDesc, columnDescrW-padding)
			var info []string
			for _,infoStr := range testResults.TestInfo {
				part := splitString( &infoStr, columnInfoW-padding)
				info = append(info, part...)
			}
			h := math.Max( float64(len(description)*rowSize), float64(len(info)*rowSize ))

			checkEndOfPageWithoutBr(h)
			y := pdf.GetY()

			addCellTextLeftTop(columnNumW, h, testResults.TestNumber )

			pdf.SetY(y)
			switch strings.ToLower(testResults.Status) {
			case "fail": failTextBlock( pdf.GetX(), y, columnStatusW-padding, testResults.Status, true )
			case "warn": warnTextBlock( pdf.GetX(), y, columnStatusW-padding, testResults.Status, true )
			case "pass": passTextBlock( pdf.GetX(), y, columnStatusW-padding, testResults.Status, true )
			default:	 infoTextBlock( pdf.GetX(), y, columnStatusW-2*padding, testResults.Status, true )
			}

			pdf.SetY(y)
			addMultiLines( pdf.GetX()+padding, rowSize, description )

			pdf.SetY(y)
			pdf.SetX( columnNumW + columnStatusW + columnDescrW+3*padding)
			addMultiLines( pdf.GetX(), rowSize, info )
			pdf.Br(brSize)
			pdf.SetY(y+h+padding)
		}
		pdf.Br(padding)
	}
}

func showColorfulData(fail, warn, pass, info int)  {
	baseY := pdf.GetY()
	baseX := pdf.GetX()+padding*2
	w := 20.0+padding*2

	failTextBlock(baseX, baseY, w, strconv.Itoa(fail), false)
	baseX += w+padding
	warnTextBlock(baseX, baseY, w, strconv.Itoa(warn), false)
	baseX += w+padding
	passTextBlock(baseX, baseY, w, strconv.Itoa(pass), false)
	baseX += w+padding
	infoTextBlock(baseX, baseY, w, strconv.Itoa(info), false)
	pdf.SetY( baseY+rowSize+padding*0.5)
	setDefaultBackgroundColor()
	pdf.SetTextColor(0,0,0)
}

func showColorfulBlock(x,y, w float64, text string, recoverColors bool)  {
	pdf.SetTextColor(255,255,255)
	addBlock( x, y, w, rowSize, text)
	if recoverColors {
		setDefaultBackgroundColor()
		pdf.SetTextColor(0,0,0)
	}
}

func failTextBlock(x,y, w float64, text string, recoverColors bool) {
	setHighBackgroundColor()
	showColorfulBlock(x,y,w,text,recoverColors)
}

func warnTextBlock(x,y, w float64, text string, recoverColors bool) {
	setOrangeBackgroundColor()
	showColorfulBlock(x,y,w,text,recoverColors)
}

func passTextBlock(x,y, w float64, text string, recoverColors bool) {
	setDarkGreenBackgroundColor()
	showColorfulBlock(x,y,w,text,recoverColors)
}
func infoTextBlock(x,y, w float64, text string, recoverColors bool) {
	setGreyBackgroundColor()
	showColorfulBlock(x,y,w,text,recoverColors)
}