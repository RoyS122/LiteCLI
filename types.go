package litecli

type Command struct {
	Use   string
	Short string
	Long  string
	Run   func(cmd *Command, args []string)
	Flags []Flag
}

type Flag struct {
	Name         string
	Short        rune
	Value        string
	DefaultValue string
	Target       any
	Type         uint8
}

type App struct {
	Name        string
	Description string
	Version     string
	RootCmd     *Command
}

type FlagList []Flag

const (
	STRING = iota
	INT
	UINT8
)
