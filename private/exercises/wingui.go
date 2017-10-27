// @Author  : Reticence (liuyang_blue@qq.com)
// @Homepage: https://github.com/Reticence
// @Date    : 2017-10-25 17:07
// @Version : 1.0
// @Software: Gogland

package exercises

import (
	"github.com/andlabs/ui"
)

func uiStructure() {
	name := ui.NewEntry()
	button := ui.NewButton("Greet")
	greeting := ui.NewLabel("")
	box := ui.NewVerticalBox()
	box.Append(ui.NewLabel("Enter your name:"), false)
	box.Append(name, false)
	box.Append(button, false)
	box.Append(greeting, false)
	window := ui.NewWindow("Hello", 300, 100, false)
	window.SetChild(box)
	button.OnClicked(func(*ui.Button) {
		greeting.SetText("Hello, " + name.Text() + "!")
	})
	window.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	window.Show()
}

func RunWinGui()  {
	err := ui.Main(uiStructure)
	if err != nil {
		panic(err)
	}
}