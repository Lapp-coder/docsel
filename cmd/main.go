package main

import (
	"log"

	"github.com/Lapp-coder/docsel/internal/app"
)

func main() {
	var filepath string
	app.RunCmd.
		PersistentFlags().
		StringVarP(&filepath, app.FlagPath, app.ShorthandFlagPath, app.DefaultFlagPathValue, app.DescFlagPath)
	app.RunCmd.
		PersistentFlags().
		BoolP(app.FlagDetach, app.ShorthandFlagDetach, app.DefaultFlagDetachValue, app.DescFlagDetach)
	app.RunCmd.
		PersistentFlags().
		Bool(app.FlagRemove, app.DefaultFlagRemoveValue, app.DescFlagRemove)
	app.RunCmd.
		PersistentFlags().
		BoolP(app.FlagSave, app.ShorthandFlagSave, app.DefaultFlagSaveValue, app.DescFlagSave)

	app.RootCmd.AddCommand(app.RunCmd)
	if err := app.RootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
