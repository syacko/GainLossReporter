package main

import (
	"GainLossReporter/packages"
	"github.com/integrii/flaggy"
)

const (
	LOGMESSAGEPREFIX = "GainLossReporter"
	VERSION          = "1.0"
)

var (
	fileName = ""
	action   = ""
)

func init() {
	// Set your program's name and description.  These appear in help output.
	flaggy.SetName(LOGMESSAGEPREFIX)
	flaggy.SetVersion(VERSION)
	flaggy.SetDescription(`
	`)

	// You can disable various things by changing bool on the default parser
	// (or your own parser if you have created one).
	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	// You can set a help prepend or append on the default parser.
	//flaggy.DefaultParser.AdditionalHelpPrepend = "https://gitlab.com/getsote/utilities"

	flaggy.AddPositionalValue(&fileName, "fileName", 1, true, "The full qualified file name for the input csv file")
	flaggy.AddPositionalValue(&action, "action", 2, true, "What you want Gain/Loss Reporter to output")
	flaggy.Parse()
}

func main() {

	if err := packages.Run(fileName); err != nil {
		panic("Application Failed: " + err.Error())
	}

}
