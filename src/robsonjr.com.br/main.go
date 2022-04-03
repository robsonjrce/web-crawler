package main

import (
	"robsonjr.com.br/cmd"
	"robsonjr.com.br/utils/signals"
)

func main() {
	signals.InstallInterruptSignal()

	cmd.Execute()
}