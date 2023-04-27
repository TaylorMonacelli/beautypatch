package main

import (
	"fmt"
	"os"
	"github.com/sirupsen/logrus"
)

func main() {
	paths := []string{
		"C:\\Windows\\System32\\control.exe",
		"C:\\Windows\\System32\\windowspowershell\\v1.0\\powershell.exe",
		"C:\\Windows\\explorer.exe",
	}

	// Initialize logrus
	logger := logrus.New()
	logger.SetOutput(os.Stdout)

	xmlTemplate := `<?xml version="1.0" encoding="utf-8"?>
<LayoutModificationTemplate
    xmlns="http://schemas.microsoft.com/Start/2014/LayoutModification"
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
</LayoutModificationTemplate>`

	var taskbarApps string
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			logger.Infof("File exists: %s", path)
			taskbarApps += "        <taskbar:DesktopApp DesktopApplicationLinkPath=\"" + path + "\" />\n"
		} else {
			logger.Warnf("File does not exist: %s", path)
		}
	}

	xmlContent := []byte(fmt.Sprintf(xmlTemplate, taskbarApps))

	// Write to file
	file, err := os.Create("Taskbar.xml")
	if err != nil {
		logger.Errorf("Failed to create file: %s", err.Error())
		return
	}
	defer file.Close()

	if _, err := file.Write(xmlContent); err != nil {
		logger.Errorf("Failed to write to file: %s", err.Error())
		return
	}

	logger.Infof("LayoutModificationTemplate written to Taskbar.xml file")
}
