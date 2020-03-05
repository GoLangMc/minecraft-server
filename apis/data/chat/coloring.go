package chat

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type ChatColor int

type ColorCode struct {
	Chat string
	Motd string
	Json string

	Dec string
	Hex string
}

const (
	DarkRed ChatColor = iota
	Red

	Gold
	Yellow

	DarkGreen
	Green

	DarkAqua
	Aqua

	DarkBlue
	Blue

	DarkPurple
	Purple

	White
	Black

	DarkGray
	Gray

	Obfuscated
	Bold
	Strikethrough
	Underline
	Italic
	Reset
)

const ColorCChar = '§'
const ColorAChar = '&'

var codeToCode = map[ChatColor]*ColorCode{
	DarkRed: {
		Chat: `§4`,
		Motd: `\u00A74`,
		Json: `dark_red`,
		Dec:  `11141120`,
		Hex:  `AA0000`,
	},
	Red: {
		Chat: `§c`,
		Motd: `\u00A7c`,
		Json: `red`,
		Dec:  `16733525`,
		Hex:  `FF5555`,
	},
	Gold: {
		Chat: `§6`,
		Motd: `\u00A76`,
		Json: `gold`,
		Dec:  `16755200`,
		Hex:  `FFAA00`,
	},
	Yellow: {
		Chat: `§e`,
		Motd: `\u00A7e`,
		Json: `yellow`,
		Dec:  `16777045`,
		Hex:  `FFFF55`,
	},
	DarkGreen: {
		Chat: `§2`,
		Motd: `\u00A72`,
		Json: `dark_green`,
		Dec:  `43520`,
		Hex:  `00AA00`,
	},
	Green: {
		Chat: `§a`,
		Motd: `\u00A7a`,
		Dec:  `5635925`,
		Hex:  `55FF55`,
	},
	DarkAqua: {
		Chat: `§3`,
		Motd: `\u00A73`,
		Json: `dark_aqua`,
		Dec:  `43690`,
		Hex:  `00AAAA`,
	},
	Aqua: {
		Chat: `§b`,
		Motd: `\u00A7b`,
		Json: `aqua`,
		Dec:  `5636095`,
		Hex:  `55FFFF`,
	},
	DarkBlue: {
		Chat: `§1`,
		Motd: `\u00A71`,
		Json: `dark_blue`,
		Dec:  `170`,
		Hex:  `0000AA`,
	},
	Blue: {
		Chat: `§9`,
		Motd: `\u00A79`,
		Json: `blue`,
		Dec:  `5592575`,
		Hex:  `5555FF`,
	},
	DarkPurple: {
		Chat: `§5`,
		Motd: `\u00A75`,
		Json: `dark_purple`,
		Dec:  `11141290`,
		Hex:  `AA00AA`,
	},
	Purple: {
		Chat: `§d`,
		Motd: `\u00A7d`,
		Json: `light_purple`,
		Dec:  `16733695`,
		Hex:  `FF55FF`,
	},
	White: {
		Chat: `§f`,
		Motd: `\u00A7f`,
		Json: `white`,
		Dec:  `16777215`,
		Hex:  `FFFFFF`,
	},
	Black: {
		Chat: `§0`,
		Motd: `\u00A70`,
		Json: `black`,
		Dec:  `0`,
		Hex:  `000000`,
	},
	DarkGray: {
		Chat: `§8`,
		Motd: `\u00A78`,
		Json: `dark_gray`,
		Dec:  `5592405`,
		Hex:  `555555`,
	},
	Gray: {
		Chat: `§7`,
		Motd: `\u00A77`,
		Json: `gray`,
		Dec:  `11184810`,
		Hex:  `AAAAAA`,
	},

	Obfuscated: {
		Chat: `§k`,
		Motd: `\u00A7k`,
		Json: `obfuscated`,
	},
	Bold: {
		Chat: `§l`,
		Motd: `\u00A7l`,
		Json: `bold`,
	},
	Strikethrough: {
		Chat: `§m`,
		Motd: `\u00A7m`,
		Json: `strikethrough`,
	},
	Underline: {
		Chat: `§n`,
		Motd: `\u00A7n`,
		Json: `underline`,
	},
	Italic: {
		Chat: `§o`,
		Motd: `\u00A7o`,
		Json: `italic`,
	},
	Reset: {
		Chat: `§r`,
		Motd: `\u00A7r`,
		Json: `reset`,
	},
}
var codeToForm = map[ChatColor]color.Attribute{
	DarkRed:    color.FgHiRed,
	Red:        color.FgRed,
	Gold:       color.FgYellow,
	Yellow:     color.FgHiYellow,
	DarkGreen:  color.FgGreen,
	Green:      color.FgHiGreen,
	DarkAqua:   color.FgCyan,
	Aqua:       color.FgHiCyan,
	DarkBlue:   color.FgBlue,
	Blue:       color.FgHiBlue,
	DarkPurple: color.FgMagenta,
	Purple:     color.FgHiMagenta,
	White:      color.FgHiWhite,
	Black:      color.FgBlack,
	DarkGray:   color.FgHiBlack,
	Gray:       color.FgWhite,

	Reset:         color.Reset,
	Obfuscated:    color.BlinkRapid,
	Bold:          color.Bold,
	Strikethrough: color.CrossedOut,
	Underline:     color.Underline,
	Italic:        color.Italic,
}

var charToCode = map[rune]ChatColor{
	'4': DarkRed,
	'c': Red,
	'6': Gold,
	'e': Yellow,

	'2': DarkGreen,
	'a': Green,

	'3': DarkAqua,
	'b': Aqua,

	'1': DarkBlue,
	'9': Blue,

	'5': DarkPurple,
	'd': Purple,

	'f': White,
	'0': Black,

	'8': DarkGray,
	'7': Gray,

	'k': Obfuscated,
	'l': Bold,
	'm': Strikethrough,
	'n': Underline,
	'o': Italic,
	'r': Reset,
}

var jsonToCode = map[string]ChatColor{
	`dark_red`: DarkRed,
	`red`:      Red,
	`gold`:     Gold,
	`yellow`:   Yellow,

	`dark_green`: DarkGreen,
	`green`:      Green,

	`dark_aqua`: DarkAqua,
	`aqua`:      Aqua,

	`dark_blue`: DarkBlue,
	`blue`:      Blue,

	`dark_purple`:  DarkPurple,
	`light_purple`: Purple,

	`white`: White,
	`black`: Black,

	`dark_gray`: DarkGray,
	`gray`:      Gray,

	`obfuscated`:    Obfuscated,
	`bold`:          Bold,
	`strikethrough`: Strikethrough,
	`underline`:     Underline,
	`italic`:        Italic,
	`reset`:         Reset,
}

func (code ChatColor) String() string {
	return codeToCode[code].Chat
}

func (code *ChatColor) MarshalJSON() ([]byte, error) {
	return []byte(`"` + codeToCode[*code].Json + `"`), nil
}

func (code *ChatColor) UnmarshalJSON(bytes []byte) error {
	*code = jsonToCode[string(bytes)]
	return nil
}

func (code *ChatColor) On(text string) string {
	if len(text) == 0 {
		return ""
	}

	return fmt.Sprintf("%p%s%v", code, text, Reset)
}

func Translate(text string) string {

	build := strings.Builder{}
	chars := []rune(text)

	for i := 0; i < len(chars); i++ {

		r := chars[i]

		if r != ColorAChar || i+1 >= len(chars) {
			build.WriteRune(r)
		} else {
			c := codeToCode[charToCode[chars[i+1]]]

			if c == nil {
				build.WriteRune(r)
			} else {
				build.WriteString(c.Chat)
				i++
			}
		}
	}

	return build.String()
}

func TranslateConsole(text string) string {
	text = Translate(text)

	build := strings.Builder{}
	temps := strings.Builder{}

	chars := []rune(text)
	forms := make([]color.Attribute, 0)

	for i := 0; i < len(chars); i++ {
		r := chars[i]
		if r != ColorCChar || i+1 >= len(chars) {
			temps.WriteRune(r)
			continue
		}

		f, con := codeToForm[charToCode[chars[i+1]]]
		if !con {
			temps.WriteRune(r)
			continue
		}

		if temps.Len() > 0 {
			build.WriteString(color.New(forms...).Sprint(temps.String()))
			temps.Reset()
		}

		i++
		if f <= color.CrossedOut {
			forms = append(forms, f)
		} else {
			forms = make([]color.Attribute, 0)
			forms = append(forms, f)
		}
	}

	if temps.Len() > 0 {
		build.WriteString(color.New(forms...).Sprint(temps.String()))
	}

	return build.String()
}
