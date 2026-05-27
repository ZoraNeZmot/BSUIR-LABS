// package main

// import (
// 	"chat/iup"
// 	"chat/node"
// 	"chat/startup"
// 	"log"
// 	"time"
// )

// func main() {
// 	localNode := iup.NewIUP(startup.Name)
// 	if localNode == nil {
// 		log.Fatal("Failed to create IUP")
// 	}

// 	localListener := node.NewListenerTCP(localNode.GetNode())
// 	if localListener == nil {
// 		log.Fatal("Failed to create listener")
// 	}

// 	localNode.ListenIncoming()
// 	localNode.SendGreetingMessage()

// 	time.Sleep(10 * time.Second)
// }

package main

import (
	"chat/iup"
	"chat/node"
	"chat/startup"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ChatUI struct {
	window      fyne.Window
	messages    *widget.List
	messageData []Message
	input       *widget.Entry
	sendBtn     *widget.Button
}

type Message struct {
	sender  string
	content string
	isOwn   bool
}

// NewChatUI creates and returns a new chat UI instance
func NewChatUI() *ChatUI {
	myApp := app.New()
	window := myApp.NewWindow("Chat Application")

	chat := &ChatUI{
		window:      window,
		messageData: []Message{},
	}

	// Create messages list
	chat.messages = widget.NewList(
		func() int {
			return len(chat.messageData)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			msg := chat.messageData[i]
			if msg.isOwn {
				label.SetText("Me: " + msg.content)
			} else {
				label.SetText(msg.sender + ": " + msg.content)
			}
		},
	)

	// Create input field
	chat.input = widget.NewEntry()
	chat.input.SetPlaceHolder("Type your message...")

	// Create send button
	chat.sendBtn = widget.NewButton("Send", func() {
		// This will be handled by the processing function
	})

	// Layout
	content := container.NewBorder(
		nil,
		container.NewBorder(nil, nil, nil, chat.sendBtn, chat.input),
		nil,
		nil,
		chat.messages,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(500, 600))

	return chat
}

// ProcessMyMessage handles sending messages from the current user
func (c *ChatUI) ProcessMyMessage(message string) {
	if message == "" {
		return
	}

	// Add own message to the list
	c.messageData = append(c.messageData, Message{
		sender:  "Me",
		content: message,
		isOwn:   true,
	})

	// Refresh the messages list
	c.messages.Refresh()

	// Clear the input field
	c.input.SetText("")
}

// ProcessOtherNodeMessage handles receiving messages from other nodes
func (c *ChatUI) ProcessOtherNodeMessage(senderName, message string) {
	if message == "" {
		return
	}

	// Add other node's message to the list
	c.messageData = append(c.messageData, Message{
		sender:  senderName,
		content: message,
		isOwn:   false,
	})

	// Refresh the messages list
	c.messages.Refresh()
}

// Run starts the chat UI application
func (c *ChatUI) Run() {
	c.window.ShowAndRun()
}

// SetSendHandler allows setting custom send button behavior
func (c *ChatUI) SetSendHandler(handler func(string)) {
	c.sendBtn.OnTapped = func() {
		message := c.input.Text
		if message != "" {
			handler(message)
		}
	}
}

func main() {
	chat := NewChatUI()

	localPeer := iup.NewIUP(startup.Name)
	if localPeer == nil {
		log.Fatal("Failed to create local peer")
	}

	localListener := node.NewListenerTCP(localPeer.GetNode())
	if localListener == nil {
		log.Fatal("Failed to create listener")
	}
	if err := localListener.Start(); err != nil {
		log.Fatal("Failed to start listener:", err)
	}
	localPeer.SetChanel(localListener.NodeCh)
	localPeer.ListenIncoming()
	localPeer.SendGreetingMessage()

	localListener.SetChatHandler(func(sender string, text string) {
		fyne.Do(func() {
			chat.ProcessOtherNodeMessage(sender, text)
		})
	})

	chat.SetSendHandler(func(msg string) {
		chat.ProcessMyMessage(msg)
		localListener.SendChat(startup.Name, msg)
	})

	chat.Run()

}
