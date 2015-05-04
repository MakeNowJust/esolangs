package main

import (
	"github.com/MakeNowJust/esolangs/brainfuck"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "brainfuck"
	app.Version = "0.1.0"
	app.Usage = "brainf*ck interpreter written in Go"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "buffer, b",
			Value: 30000,
			Usage: "buffer size to execute brainf*ck program",
		},
		cli.IntFlag{
			Name:  "eof, e",
			Value: 0,
			Usage: "eof value to return in breinf*ck execution",
		},
	}
	app.Action = func(c *cli.Context) {
		brainfuck.MaxBufferSize = c.Int("buffer")
		brainfuck.EOF = byte(c.Int("eof"))

		exec := brainfuck.New()

		if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
			os.Exit(1)
		} else {

			for _, file := range c.Args() {
				if src, err := ioutil.ReadFile(file); err != nil {
					log.Fatal(err)
				} else if pgrm, err := brainfuck.Parse(src); err != nil {
					log.Fatal(err)
				} else if err := exec.Exec(pgrm); err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	app.Run(os.Args)
}
