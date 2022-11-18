package qt

import (	
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
	"github.com/go-gota/gota/series"
	"github.com/go-gota/gota/dataframe"
	"os"
)
func Test(){
	widgets.NewQApplication(len(os.Args), os.Args)

	// left sider
	splitterLeft := widgets.NewQSplitter2(core.Qt__Horizontal, nil)
	textTop := widgets.NewQTextEdit2("左部文字", splitterLeft)
	splitterLeft.AddWidget(textTop)

	// right sider
	splitterRight := widgets.NewQSplitter2(core.Qt__Vertical, splitterLeft)
	textRight := widgets.NewQTextEdit2("右部文字", splitterRight)
	textbuttom := widgets.NewQTextEdit2("下部文字", splitterLeft)
	splitterRight.AddWidget(textRight)
	splitterRight.AddWidget(textbuttom)

	splitterLeft.SetWindowTitle("splitter")
	splitterLeft.Show()

	widgets.QApplication_Exec()
}
func Test2(){
	series.New([]string{"z", "y", "d", "e"}, series.String, "col")
	dataFrame := dataframe.New(
		series.New([]string{"a", "b", "c", "d", "e"}, series.String, "alphas"),
		series.New([]int{5, 4, 2, 3, 1}, series.Int, "numbers"),
		series.New([]string{"a1", "b2", "c3", "d4", "e5"}, series.String, "alnums"),
		series.New([]bool{true, false, true, true, false}, series.Bool, "state"),
	)

	fmt.Println(dataFrame)
}
