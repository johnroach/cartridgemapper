package utils

import (
	"log"
	"os"

	"github.com/fatih/color"
)

var (
	// Info logger
	Info *log.Logger
	// Warning logger
	Warning *log.Logger
	// Error logger
	Error *log.Logger
	// Debug logger
	Debug *log.Logger
)

// LogInfo displays info type output
func logInfo(message string, colorDisabled bool) {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	if !colorDisabled {
		color.Set(color.FgGreen)
	}
	Info.Println(message)
	if !colorDisabled {
		color.Unset()
	}
}

// LogWarning displays warning type output
func logWarning(message string, colorDisabled bool) {
	Info = log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime)
	if !colorDisabled {
		color.Set(color.FgYellow)
	}
	Info.Println(message)
	if !colorDisabled {
		color.Unset()
	}
}

// LogDebug displays debug type output
func logDebug(message string, colorDisabled bool) {
	Debug = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime)
	if !colorDisabled {
		color.Set(color.FgHiBlue)
	}
	Debug.Println(message)
	if !colorDisabled {
		color.Unset()
	}
}

// LogError displays error type output
func logError(message string, colorDisabled bool) {
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	if !colorDisabled {
		color.Set(color.FgRed)
	}
	Error.Println(message)
	if !colorDisabled {
		color.Unset()
	}
}

//DisplayError gets the error and the message to show to the user and displays it
func DisplayError(message string, err error, DisableColor bool) {
	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}
	logError(message+" "+errorMessage, DisableColor)
}

//DisplayDebug displays debug messages
func DisplayDebug(message string, Debug bool, DisableColor bool) {
	if Debug {
		logDebug(message, DisableColor)
	}
}

//DisplayWarning displays warning messages
func DisplayWarning(message string, DisableColor bool) {
	logWarning(message, DisableColor)
}

//DisplayInfo displays info messages
func DisplayInfo(message string, DisableColor bool) {
	logInfo(message, DisableColor)
}
