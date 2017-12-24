package utils

import (
	"github.com/fatih/color"
	"log"
	"os"
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
func LogInfo(message string, colorDisabled bool) {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	if !colorDisabled {
		color.Set(color.FgGreen)
	}
	Info.Println(message)
	if !colorDisabled {
		color.Unset()
	}
}

// LogDebug displays debug type output
func LogDebug(message string, colorDisabled bool) {
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
func LogError(message string, colorDisabled bool) {
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	if !colorDisabled {
		color.Set(color.FgRed)
	}
	Error.Println(message)
	if !colorDisabled {
		color.Unset()
	}
}

func DisplayError(message string, err error, DisableColor bool) {
	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}
	LogError(message+" "+errorMessage, DisableColor)
}

func DisplayDebug(message string, Debug bool, DisableColor bool) {
	if Debug {
		LogDebug(message, DisableColor)
	}
}

func DisplayInfo(message string, DisableColor bool) {
	LogInfo(message, DisableColor)
}
