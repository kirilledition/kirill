package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	GitCommit string = "0000000"
	GitTag    string = "nonono"
)

func LongDescription() string {
	return fmt.Sprintf(`
┌──────────────────────────────────────────┐
│                             .---..---.   │
│       .     .--.        .--.|   ||   |   │
│     .'|     |__|        |__||   ||   |   │
│   .'  |     .--..-,.--. .--.|   ||   |   │
│  <    |     |  ||  .-. ||  ||   ||   |   │     Q
│   |   | ____|  || |  | ||  ||   ||   |   │  ___|\_.-,
│   |   | \ .'|  || |  | ||  ||   ||   |   S\ Q~\___ \|
│   |   |/  . |  || |  '- |  ||   ||   |   │(   )o 5) Q
│   |    /\  \|__|| |     |__||   ||   |   │\\  \_ ()
│   |   |  \  \   | |         '---''---'   │ \'. _'/'.
│   '    \  \  \  |_|                     .-. '-(  x< \
│  '------'  '---'            ,o         /\, '.  )  /'\\
└─────────────────────────────\'.__.----/ .'\  '.-'/   \\
 by kirilledition              '---'q__/.'__ ;    /     \\_
                                    '---'     '--'       '"'
Version: %s
Git Commit: %s

kirill is a toolbox for bioinformatics`,
		GitTag, GitCommit[:7])
}

var rootCmd = &cobra.Command{
	Use:   "kirill",
	Short: "Yet another bioinformatics toolbox",
	Long:  LongDescription(),
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
