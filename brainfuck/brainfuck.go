package brainfuck

import (
	"errors"
	esolangs "github.com/MakeNowJust/esolangs/util"
	"io"
)

// CommandType means a kind of brainf*ck commands.
type CommandType int

// There are all brainf*ck commands.
const (
	// +
	Incr CommandType = iota
	// -
	Decr
	// >
	Next
	// <
	Prev
	// .
	Write
	// ,
	Read
	// [...]
	Loop
)

// Command means a brainf*ck command.
type Command struct {
	Type  CommandType
	Block Program
}

// Program means a brainf*ck program.
type Program []Command

// There are brainf*ck errors.
var (
	ParseError = errors.New("failed to parse string as brainf*ck program")
	ExecError  = errors.New("failed to execute brainf*ck program")
)

// There are a configuration variables.
var (
	MaxBufferSize = 30000
	EOF           = byte(0)
)

// Exec means brainf*ck execution status.
type Exec struct {
	Buffer []byte
	Index  int
}

// Parse is to parse string as brainf*ck program.
func Parse(src []byte) (Program, error) {
	if pgrm, src, err := parse([]byte(src)); err != nil {
		return nil, err
	} else {
		if len(src) == 0 {
			return pgrm, nil
		}
		return nil, ParseError
	}
}

func parse(src []byte) (Program, []byte, error) {
	pgrm := []Command{}
	for i := 0; i < len(src); i++ {
		switch src[i] {
		case '+':
			pgrm = append(pgrm, Command{Type: Incr})
		case '-':
			pgrm = append(pgrm, Command{Type: Decr})
		case '>':
			pgrm = append(pgrm, Command{Type: Next})
		case '<':
			pgrm = append(pgrm, Command{Type: Prev})
		case '.':
			pgrm = append(pgrm, Command{Type: Write})
		case ',':
			pgrm = append(pgrm, Command{Type: Read})
		case '[':
			if block, src2, err := parse(src[i+1:]); err != nil {
				return nil, nil, err
			} else {
				if len(src2) == 0 || src2[0] != ']' {
					return nil, nil, ParseError
				}
				src = src2
				i = 0
				pgrm = append(pgrm, Command{Type: Loop, Block: block})
			}
		case ']':
			return Program(pgrm), src[i:], nil
		}
	}

	return Program(pgrm), []byte{}, nil
}

func New() *Exec {
	return &Exec{
		Buffer: make([]byte, MaxBufferSize),
		Index:  0,
	}
}

func (state *Exec) value() byte {
	return state.Buffer[state.Index]
}

func (state *Exec) incr() {
	state.Buffer[state.Index] += 1
}

func (state *Exec) decr() {
	state.Buffer[state.Index] -= 1
}

func (state *Exec) next() error {
	state.Index += 1
	if state.Index >= MaxBufferSize {
		return ExecError
	}
	return nil
}

func (state *Exec) prev() error {
	state.Index -= 1
	if state.Index < 0 {
		return ExecError
	}
	return nil
}

func (state *Exec) write() error {
	return esolangs.WriteByte(state.value())
}

func (state *Exec) read() error {
	if ch, err := esolangs.ReadByte(); err != nil {
		if err == io.EOF {
			state.Buffer[state.Index] = EOF
			return nil
		}
		return err
	} else {
		state.Buffer[state.Index] = ch
		return nil
	}
}

func (state *Exec) Exec(pgrm Program) error {
	for _, cmd := range pgrm {
		switch cmd.Type {
		case Incr:
			state.incr()
		case Decr:
			state.decr()
		case Next:
			if err := state.next(); err != nil {
				return err
			}
		case Prev:
			if err := state.prev(); err != nil {
				return err
			}
		case Write:
			if err := state.write(); err != nil {
				return err
			}
		case Read:
			if err := state.read(); err != nil {
				return err
			}
		case Loop:
			for state.value() != 0 {
				if err := state.Exec(cmd.Block); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
