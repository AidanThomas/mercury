package ui

import (
	"fmt"
	"strings"

	"github.com/AidanThomas/mercury/internal/message"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type colours struct {
	bg tcell.Color
	fg tcell.Color
}

var (
	messages []message.Message
	users    []string
)

func Start() {
	defaultColours := colours{
		bg: tcell.ColorNone,
		fg: tcell.ColorWhite,
	}

	out := make(chan string)

	app := tview.NewApplication()

	// Create components
	msgView := NewMessageView(defaultColours)
	msgInput := NewMessageInput(defaultColours, out)
	usrList := NewUserList(defaultColours)

	messageFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(msgView, 0, 1, false).
		AddItem(msgInput, 3, 1, true)
	mainFlex := tview.NewFlex().
		AddItem(messageFlex, 0, 1, true).
		AddItem(usrList, 40, 1, false)

	mainFlex.SetBackgroundColor(defaultColours.bg)
	messageFlex.SetBackgroundColor(defaultColours.bg)

	go waitForOutgoing(out, msgView)

	app.SetRoot(mainFlex, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func NewMessageInput(col colours, out chan string) *tview.TextArea {
	input := tview.NewTextArea()
	input.SetPlaceholder("Send a message...")
	input.SetPlaceholderStyle(tcell.StyleDefault.Foreground(col.bg))
	input.SetTextStyle(tcell.StyleDefault.Background(col.bg))
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			out <- fmt.Sprintf(" > %s", input.GetText())
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
	msgView.SetScrollable(true)
	msgView.SetTextAlign(tview.AlignLeft)
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

func waitForOutgoing(out chan string, msgView *tview.TextView) {
	for {
		msg := <-out
		messages = append(messages, message.Message{Body: msg})
		var msgs []string
		for i := len(messages) - 1; i >= 0; i-- {
			msgs = append(msgs, messages[i].Body)
		}
		msgView.SetText(strings.Join(msgs, "\n"))
	}
}
