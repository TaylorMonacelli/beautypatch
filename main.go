package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

func main() {
	paths := []string{
		"C:\\Windows\\System32\\control.exe",
		"C:\\Windows\\System32\\windowspowershell\\v1.0\\powershell.exe",
		"C:\\Windows\\explorer.exe",
	}

	// create a new logger instance
	logger := logrus.New()

	// set logger output to stdout
	logger.SetOutput(os.Stdout)

	// create a new taskbar pin list
	taskbarPinList := ""

	// iterate over paths
	for _, path := range paths {
		// check if the file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			logger.Warnf("%s does not exist", path)
			continue
		}

		// add the path to the taskbar pin list
		taskbarPinList += fmt.Sprintf(`<taskbar:DesktopApp DesktopApplicationLinkPath="%s" />\n`, filepath.ToSlash(path))
		logger.Infof("%s added to the taskbar pin list", path)
	}

	// create the layout modification template
	layoutModificationTemplate := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
	<LayoutModificationTemplate xmlns="http://schemas.microsoft.com/Start/2014/LayoutModification"
	    xmlns:defaultlayout="http://schemas.microsoft.com/Start/2014/FullDefaultLayout"
	    xmlns:start="http://schemas.microsoft.com/Start/2014/StartLayout"
	    xmlns:taskbar="http://schemas.microsoft.com/Start/2014/TaskbarLayout"
	    Version="1">
	  <CustomTaskbarLayoutCollection PinListPlacement="Replace">
	    <defaultlayout:TaskbarLayout>
	      <taskbar:TaskbarPinList>
	        %s
	      </taskbar:TaskbarPinList>
	    </defaultlayout:TaskbarLayout>
	  </CustomTaskbarLayoutCollection>
	</LayoutModificationTemplate>`, taskbarPinList)

	// log the layout modification template
	logger.Infof("Layout modification template:\n%s", layoutModificationTemplate)
}
