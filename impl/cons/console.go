package cons

import (
	"bufio"
	"io"
	"os"

	"minecraft-server/apis/base"
	"minecraft-server/apis/logs"
	"minecraft-server/impl/data/system"
)

type Console struct {
	i *io.Reader
	o *io.Writer

	logger *logs.Logging

	IChannel chan string
	OChannel chan string

	report chan system.Message
}

func NewConsole(report chan system.Message) *Console {
	file, err := os.OpenFile("latest.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		report <- system.Make(system.FAIL, err)
		return nil
	}

	i := io.MultiReader(os.Stdin)
	o := io.MultiWriter(os.Stdout, file)

	return &Console{
		i: &i,
		o: &o,

		IChannel: make(chan string),
		OChannel: make(chan string),

		report: report,
		logger: logs.NewLoggingWith("console", o, logs.EveryLevel...),
	}
}

func (c *Console) Load() {
	// handle i channel
	go func() {
		scanner := bufio.NewScanner(*c.i)

		for scanner.Scan() {
			func() {
				defer func() {
					if err := recover(); err != nil {
						c.report <- system.Make(system.FAIL, err)
					}
				}()

				c.IChannel <- scanner.Text()
			}()
		}
	}()

	// handle o channel
	go func() {
		for line := range c.OChannel {
			c.logger.Info(line)
		}
	}()
}

func (c *Console) Kill() {
	defer func() {
		_ = recover() // ignore panic with closing closed channel
	}()

	close(c.IChannel)
	close(c.OChannel)
}

func (c *Console) Name() string {
	return "ConsoleSender"
}

func (c *Console) SendMessage(message ...interface{}) {
	defer func() {
		if err := recover(); err != nil {
			c.report <- system.Make(system.FAIL, err)
		}
	}()

	c.OChannel <- base.ConvertToString(message...)
}
