package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type colours struct {
	bg tcell.Color
	fg tcell.Color
}

func Start() {
	defaultColours := colours{
		bg: tcell.ColorNone,
		fg: tcell.ColorWhite,
	}

	app := tview.NewApplication()

	// Create components
	msgView := NewMessageView(defaultColours)

	msgInput := NewMessageInput(defaultColours, msgView)

	usrList := NewUserList(defaultColours)

	messageFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(msgView, 0, 1, false).
		AddItem(msgInput, 3, 1, true)
	mainFlex := tview.NewFlex().
		AddItem(messageFlex, 0, 1, true).
		AddItem(usrList, 40, 1, false)

	mainFlex.SetBackgroundColor(defaultColours.bg)
	messageFlex.SetBackgroundColor(defaultColours.bg)
	app.SetRoot(mainFlex, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func NewMessageInput(col colours, msgView *tview.TextView) *tview.TextArea {
	input := tview.NewTextArea()
	input.SetPlaceholder("Send a message...")
	input.SetPlaceholderStyle(tcell.StyleDefault.Foreground(col.bg))
	input.SetTextStyle(tcell.StyleDefault.Background(col.bg))
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			msgView.SetText(input.GetText())
			input.SetText("", true)
			// Clear event to not send newline
			event = tcell.NewEventKey(tcell.KeyBackspace, 'a', tcell.ModNone)
		}
		return event
	})
	input.SetBorder(true)
	input.SetBackgroundColor(col.bg)
	return input
}

func NewMessageView(col colours) *tview.TextView {
	msgView := tview.NewTextView()
	msgView.SetBackgroundColor(col.bg)
	msgView.SetTitle("Messages").SetTitleAlign(tview.AlignLeft)
	msgView.SetBorder(true)
	return msgView
}

func NewUserList(col colours) *tview.TextView {
	usrList := tview.NewTextView()
	usrList.SetBackgroundColor(col.bg)
	usrList.SetTitle("Connected Users").SetTitleAlign(tview.AlignLeft)
	usrList.SetBorder(true)
	return usrList
}
