package msgs

import (
	"encoding/json"
	"strings"

	"github.com/golangmc/minecraft-server/apis/data/chat"
)

type MessagePosition byte

const (
	NormalChat MessagePosition = iota
	SystemChat
	HotBarText
)

type Message struct {
	Text  string          `json:"text"`
	Color *chat.ChatColor `json:"color,string,omitempty"`

	Bold          *bool `json:"bold,boolean,omitempty"`
	Italic        *bool `json:"italic,boolean,omitempty"`
	Underlined    *bool `json:"underlined,boolean,omitempty"`
	Strikethrough *bool `json:"strikethrough,boolean,omitempty"`
	Obfuscated    *bool `json:"obfuscated,boolean,omitempty"`

	Extra []*Message `json:"extra,omitempty"`

	head *Message
}

func New(text string) *Message {
	return &Message{
		Text: text,
	}
}

func (c *Message) SetColor(code chat.ChatColor) *Message {
	c.Color = &code
	return c
}

func (c *Message) SetBold(value bool) *Message {
	c.Bold = &value
	return c
}

func (c *Message) SetItalic(value bool) *Message {
	c.Italic = &value
	return c
}

func (c *Message) SetUnderlined(value bool) *Message {
	c.Underlined = &value
	return c
}

func (c *Message) SetStrikethrough(value bool) *Message {
	c.Strikethrough = &value
	return c
}

func (c *Message) SetObfuscated(value bool) *Message {
	c.Obfuscated = &value
	return c
}

// creates and returns a new Chat object, and adds it to the caller's extra slice
func (c *Message) Add(text string) *Message {
	chat := New(text)
	chat.head = c

	c.Extra = append(c.Extra, chat)

	return chat
}

func (c *Message) Reset() *Message {

	next := c.Add("").SetColor(chat.Reset)

	if c.Bold != nil && *c.Bold == true {
		next.SetBold(false)
	}

	if c.Italic != nil && *c.Italic == true {
		next.SetItalic(false)
	}

	if c.Underlined != nil && *c.Underlined == true {
		next.SetUnderlined(false)
	}

	if c.Strikethrough != nil && *c.Strikethrough == true {
		next.SetStrikethrough(false)
	}

	if c.Obfuscated != nil && *c.Obfuscated == true {
		next.SetObfuscated(false)
	}

	return next
}

func (c *Message) AsJson() string {

	chat := c

	for chat.head != nil {
		chat = chat.head
	}

	if text, err := json.Marshal(chat); err != nil {
		panic(err)
	} else {
		return string(text)
	}
}

func (c *Message) AsText() string {
	builder := strings.Builder{}

	curr := c

	for curr.head != nil {
		curr = curr.head
	}

	builder.WriteString(curr.asText())

	return builder.String()
}

func (c *Message) asText() string {
	builder := strings.Builder{}

	if c.Color != nil {
		builder.WriteString(c.Color.String())
	}

	if c.Bold != nil && *c.Bold == true {
		builder.WriteString(chat.Bold.String())
	}

	if c.Italic != nil && *c.Italic == true {
		builder.WriteString(chat.Italic.String())
	}

	if c.Underlined != nil && *c.Underlined == true {
		builder.WriteString(chat.Underline.String())
	}

	if c.Strikethrough != nil && *c.Strikethrough == true {
		builder.WriteString(chat.Strikethrough.String())
	}

	if c.Obfuscated != nil && *c.Obfuscated == true {
		builder.WriteString(chat.Obfuscated.String())
	}

	builder.WriteString(c.Text)

	for _, extra := range c.Extra {
		builder.WriteString(extra.asText())
	}

	return builder.String()
}

func (c *Message) String() string {
	return c.AsJson()
}
