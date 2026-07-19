package litecli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func (c *Command) Execute() {
	var startedFlag int = -1
	var positionArg []string
	for _, arg := range os.Args[1:] {
		var flagIndex int = -1
		var pre uint8
		if strings.HasPrefix(arg, "-") {
			pre += 1
			if strings.HasPrefix(arg, "--") {
				pre += 1
			}

			var argSplited []string = strings.Split(arg[pre:], "=")

			if pre == 1 {
				flagIndex = c.checkFlagShort(rune(argSplited[0][0]))
			} else {
				flagIndex = c.checkFlagName(argSplited[0])
			}

			if flagIndex == -1 {
				continue
			}

			if len(argSplited) > 1 {
				put(&c.Flags[flagIndex], argSplited[1])
				continue
			}
			if c.Flags[flagIndex].Type == BOOL {
				put(&c.Flags[flagIndex], "1")
				continue
			}
			startedFlag = flagIndex
			continue
		}
		if startedFlag > -1 {
			put(&c.Flags[startedFlag], arg)
		} else {
			positionArg = append(positionArg, arg)
		}

		startedFlag = -1

	}

	c.Run(c, positionArg)
}

func put(f *Flag, value string) {
	switch target := f.Target.(type) {

	case *string:
		*target = value

	case *int:
		v, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Put case INT:", err)
			return
		}
		*target = v

	case *uint8:
		v, err := strconv.ParseUint(value, 10, 8)
		if err != nil {
			fmt.Println("Put case UINT8:", err)
			return
		}
		*target = uint8(v)

	case *uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			fmt.Println("Put case UINT64:", err)
			return
		}
		*target = v

	case *float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			fmt.Println("Put case FLOAT64:", err)
			return
		}
		*target = v

	case *bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			fmt.Println("Put case BOOL:", err)
			return
		}
		*target = v

	default:
		fmt.Printf("Put: type de target non supporté (%T)\n", f.Target)
	}
}

func (c *Command) checkFlagName(flagFull string) int {
	for index, f := range c.Flags {
		if f.Name == flagFull {
			return index
		}
	}
	return -1
}

func (c *Command) checkFlagShort(flagRune rune) int {
	for index, f := range c.Flags {
		if f.Short == flagRune {
			return index
		}
	}
	return -1
}

func (c *Command) StringVarP(s *string, flagFull string, flagRune rune, defaultValue string) {
	*s = defaultValue
	c.Flags = append(c.Flags, Flag{Target: s, Name: flagFull, Short: flagRune, DefaultValue: defaultValue, Type: STRING})
}

func (c *Command) IntVarP(i *int, flagFull string, flagRune rune, defaultValue int) {
	*i = defaultValue
	c.Flags = append(c.Flags, Flag{Target: i, Name: flagFull, Short: flagRune, DefaultValue: strconv.Itoa(defaultValue), Type: INT})
}

func (c *Command) Uint8VarP(ui *uint8, flagFull string, flagRune rune, defaultValue uint8) {
	*ui = defaultValue
	c.Flags = append(c.Flags, Flag{Target: ui, Name: flagFull, Short: flagRune, DefaultValue: strconv.FormatUint(uint64(defaultValue), 10), Type: UINT8})
}

func (c *Command) Uint64VarP(ui *uint64, flagFull string, flagRune rune, defaultValue uint64) {
	*ui = defaultValue
	c.Flags = append(c.Flags, Flag{Target: ui, Name: flagFull, Short: flagRune, DefaultValue: strconv.FormatUint(uint64(defaultValue), 10), Type: UINT64})
}

func (c *Command) Float64VarP(i *float64, flagFull string, flagRune rune, defaultValue float64) {
	*i = defaultValue
	strDefault := strconv.FormatFloat(defaultValue, 'f', -1, 64)

	c.Flags = append(c.Flags, Flag{
		Target:       i,
		Name:         flagFull,
		Short:        flagRune,
		DefaultValue: strDefault,
		Type:         FLOAT,
	})
}

func (c *Command) BoolVarP(i *bool, flagFull string, flagRune rune, defaultValue bool) {
	*i = defaultValue

	c.Flags = append(c.Flags, Flag{
		Target:       i,
		Name:         flagFull,
		Short:        flagRune,
		DefaultValue: strconv.FormatBool(defaultValue),
		Type:         BOOL,
	})
}
