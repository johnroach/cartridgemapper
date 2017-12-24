package utils

import (
	"github.com/fatih/color"
	//"io/ioutil"
	"log"
	"os"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Debug   *log.Logger
)

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
