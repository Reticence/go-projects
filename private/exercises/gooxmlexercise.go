// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-11-07 15:47
// @Version : 0.1
// @Software: GoLand

package exercises

import (
	"baliance.com/gooxml/color"
	"baliance.com/gooxml/common"
	"baliance.com/gooxml/document"
	"baliance.com/gooxml/measurement"
	"baliance.com/gooxml/schema/soo/wml"
	"image"
	"strconv"
)

const (
	fontTitle = 14
	fontText  = 12
	Zero      = measurement.Zero
	Point     = measurement.Point
	Inch      = measurement.Inch
)

type mydocx struct {
	docx *document.Document
}

type properties struct {
	//Paragraph属性
	Breaks bool
	Tndent bool

	//Run属性
	Bold     bool
	Italic   bool
	FontSize measurement.Distance

	//Drawing Anchored属性
	Zoom    measurement.Distance
	YOffset measurement.Distance
}

func getXY(size image.Point) (x, y measurement.Distance) {
	if size.X < size.Y {
		size.X, size.Y = size.Y, size.X
	}
	x, y = measurement.Distance(size.X)/measurement.Distance(size.Y), 1
	if x > 1.6 {
		y = 1.6 / x
		x = 1.6
	}
	return x, y
}

func (m *mydocx) addPicture(path string, p properties) {
	if p.Breaks {
		m.docx.AddParagraph().Properties().AddSection(wml.ST_SectionMarkNextPage)
	}
	img, _ := common.ImageFromFile(path)
	x, y := getXY(img.Size)
	iref, _ := m.docx.AddImage(img)
	para := m.docx.AddParagraph()
	//para.Properties().SetPageBreakBefore(p.Breaks)
	para.Properties().SetAlignment(wml.ST_JcBoth)
	para.Properties().SetSpacing(Zero, y*p.Zoom*Inch)
	anchored, _ := para.AddRun().AddDrawingAnchored(iref)
	anchored.SetSize(x*p.Zoom*Inch, y*p.Zoom*Inch)
	anchored.SetOrigin(wml.WdST_RelFromHPage, wml.WdST_RelFromVTopMargin)
	anchored.SetHAlignment(wml.WdST_AlignHCenter)
	anchored.SetYOffset(p.YOffset)
	anchored.SetTextWrapSquare(wml.WdST_WrapTextRight)
}

func (m *mydocx) addTital(title string, p properties) {
	para := m.docx.AddParagraph()
	para.Properties().SetAlignment(wml.ST_JcCenter)
	para.Properties().SetSpacing(Zero, 6*Point)
	run := para.AddRun()
	run.Properties().SetSize(p.FontSize * Point)
	run.Properties().SetBold(true)
	run.AddText(title)
}

func (m *mydocx) addParagraph(text string, p properties) {
	para := m.docx.AddParagraph()
	para.Properties().SetAlignment(wml.ST_JcBoth)
	para.Properties().SetSpacing(Zero, 6*Point)
	run := para.AddRun()
	run.Properties().SetSize(p.FontSize * Point)
	if p.Tndent {
		para.Properties().AddTabStop(2*p.FontSize*Point, wml.ST_TabJcLeft, wml.ST_TabTlcNone)
		run.AddTab()
	}
	run.Properties().SetBold(p.Bold)
	run.AddText(text)
}

func (m *mydocx) addTable(colnum int) {
	m.addTital("表1：xxxx", properties{FontSize: fontText})
	table := m.docx.AddTable()
	table.Properties().SetCellSpacingAuto()
	table.Properties().SetWidth(5 * Inch)
	table.Properties().SetAlignment(wml.ST_JcTableCenter)
	borders := table.Properties().Borders()
	borders.SetTop(wml.ST_BorderSingle, color.Auto, Zero)
	borders.SetBottom(wml.ST_BorderSingle, color.Auto, Zero)
	row := table.AddRow()
	for j := 0; j < colnum; j++ {
		cell := row.AddCell()
		cell.Properties().SetVerticalAlignment(wml.ST_VerticalJcCenter)
		para := cell.AddParagraph()
		para.Properties().SetAlignment(wml.ST_JcCenter)
		para.AddRun().AddText("Title-" + strconv.Itoa(j+1))
	}

	for i := 0; i <= 10; i++ {
		row := table.AddRow()
		for j := 0; j < colnum; j++ {
			cell := row.AddCell()
			cell.Properties().SetVerticalAlignment(wml.ST_VerticalJcCenter)
			para := cell.AddParagraph()
			if j == 1 {
				para.Properties().SetAlignment(wml.ST_JcCenter)
			} else {
				para.Properties().SetAlignment(wml.ST_JcRight)
			}
			if i < 10 {
				para.AddRun().AddText(strconv.Itoa(i+1) + "-" + strconv.Itoa(j+1))
			} else {
				para.AddRun().AddText("Bottom-" + strconv.Itoa(j+1))
			}
		}
	}
}

func (m *mydocx) addPageJjh(jd JsonData) {
	m.addPicture(jd.PicturePath[0], properties{Breaks: true, Zoom: 1, YOffset: Inch})
	m.addTital(jd.Title, properties{FontSize: fontTitle})

	for _, text := range jd.Paragraphs {
		m.addParagraph(text, properties{Tndent: true, FontSize: fontText})
	}
}

func (m *mydocx) addPageMeridian(jd JsonData) {
	m.addPicture(jd.PicturePath[0], properties{Breaks: true, Zoom: 1, YOffset: Inch})
	m.addTital(jd.Title, properties{FontSize: fontTitle})
	m.addParagraph(jd.Paragraphs[0], properties{Tndent: true, FontSize: fontText})
	m.addPicture(jd.PicturePath[1], properties{Breaks: false, Zoom: 3, YOffset: 3.05 * Inch})
	for i, text := range jd.Paragraphs {
		if i == 0 {
			continue
		} else if i == 1 || i == 7 {
			m.addParagraph(text, properties{Tndent: false, Bold: true, FontSize: fontText})
		} else {
			m.addParagraph(text, properties{Tndent: true, FontSize: fontText})
		}
	}
}

func RunDocx() {
	m := mydocx{document.New()}

	// 添加文件属性
	cp := m.docx.CoreProperties
	cp.SetAuthor("meridian")
	cp.SetCategory("文档")
	cp.SetContentStatus("修订版")
	cp.SetDescription("test")
	cp.SetLastModifiedBy("meridian")
	cp.SetTitle("Test docx")

	// 添加页脚
	ftr := m.docx.AddFooter()
	para := ftr.AddParagraph()
	para.Properties().SetAlignment(wml.ST_JcCenter)
	run := para.AddRun()
	run.AddText(" ")
	run.AddField(document.FieldCurrentPage)
	run.AddText(" / ")
	run.AddField(document.FieldNumberOfPages)
	m.docx.BodySection().SetFooter(ftr, wml.ST_HdrFtrDefault)

	// 添加大标题
	m.addTital("Test docx", properties{FontSize: 16})

	//table test
	m.addTable(3)

	//var datas map[string]JsonData
	//bytes, _ := ioutil.ReadFile("T:/go-io/jd.json")
	//if err := json.Unmarshal(bytes, &datas); err != nil {
	//	panic(err)
	//}
	//m.addPageJjh(datas["基金会简介"])
	//m.addPageMeridian(datas["经纶世纪简介"])

	err := m.docx.SaveToFile("T:/go-io/docx_test.docx")
	if err != nil {
		panic(err)
	}
}
