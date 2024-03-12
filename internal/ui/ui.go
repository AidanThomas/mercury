package ui

import (
	"fmt"
	"strconv"
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
)

func Start() error {
	app := tview.NewApplication()

	out := make(chan string)

	defaultTheme := colours{
		bg: tcell.ColorNone,
		fg: tcell.ColorNone,
	}

	msgInput := newMessageInput(defaultTheme, out)
	msgView := newMessageView(defaultTheme)
	usrList := newUserList(defaultTheme)

	messageFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(msgView, 0, 1, false).
		AddItem(msgInput, 3, 1, true)
	mainFlex := tview.NewFlex().
		AddItem(messageFlex, 0, 1, true).
		AddItem(usrList, 40, 1, false)

	mainFlex.SetBackgroundColor(defaultTheme.bg)
	messageFlex.SetBackgroundColor(defaultTheme.bg)

	go waitForOutgoing(out, msgView)

	app.SetRoot(mainFlex, true)
	if err := app.Run(); err != nil {
		return err
	}
	return nil
}

func newMessageInput(theme colours, out chan string) *tview.TextArea {
	input := tview.NewTextArea()
	input.SetPlaceholder("Send a message...")
	input.SetPlaceholderStyle(tcell.StyleDefault.Foreground(theme.bg))
	input.SetTextStyle(tcell.StyleDefault.Background(theme.bg))
	input.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			out <- fmt.Sprintf(" > %s", input.GetText())
			input.SetText("", false)
			return nil
		case tcell.KeyCtrlC:
			out <- "Control C"
			input.SetText("", false)
			return nil
		}
		return event
	})
	input.SetBorder(true)
	input.SetBackgroundColor(theme.bg)
	return input
}

func newMessageView(theme colours) *tview.TextView {
	msgView := tview.NewTextView()
	msgView.SetScrollable(false)
	msgView.SetTextAlign(tview.AlignLeft)
	msgView.SetBackgroundColor(theme.bg)
	msgView.SetTitle("Messages").SetTitleAlign(tview.AlignLeft)
	msgView.SetBorder(true)
	return msgView
}

func newUserList(theme colours) *tview.TextView {
	usrList := tview.NewTextView()
	usrList.SetBackgroundColor(theme.bg)
	usrList.SetTitle("Connected Users").SetTitleAlign(tview.AlignLeft)
	usrList.SetBorder(true)
	return usrList
}

func waitForOutgoing(out chan string, msgView *tview.TextView) {
	for {
		msg := <-out

		messages = append(messages, message.Message{Body: msg})
		_, _, _, height := msgView.GetInnerRect()
		diff := height - len(messages)
		var msgs []string
		if diff > 0 {
			for i := 0; i < diff; i++ {
				msgs = append(msgs, "")
			}
		}
		for _, msg := range messages {
			msgs = append(msgs, msg.Body)
		}
		renderMessages(msgs, msgView)
	}
}

func renderMessages(msgs []string, msgView *tview.TextView) {
	msgView.SetText(strings.Join(msgs, "\n"))
}
