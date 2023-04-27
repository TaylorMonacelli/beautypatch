package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	filePaths := []string{
		"C:\\Windows\\System32\\control.exe",
		"C:\\Windows\\System32\\windowspowershell\\v1.0\\powershell.exe",
		"C:\\Windows\\explorer.exe",
	}

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
      </taskbar:TaskbarPinList>
    </defaultlayout:TaskbarLayout>
  </CustomTaskbarLayoutCollection>
</LayoutModificationTemplate>`

	for _, path := range filePaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.WithFields(log.Fields{
				"path": path,
			}).Info("File does not exist")
		} else {
			log.WithFields(log.Fields{
				"path": path,
			}).Info("File exists")

			xmlTemplate = fmt.Sprintf("%s\n<taskbar:DesktopApp DesktopApplicationLinkPath=\"%s\" />", xmlTemplate, path)

			docStartIndex := len(xmlTemplate) - len("</taskbar:TaskbarPinList>\n      </defaultlayout:TaskbarLayout>\n    </CustomTaskbarLayoutCollection>\n  </LayoutModificationTemplate>")
			xmlTemplate = xmlTemplate[:docStartIndex] + fmt.Sprintf("\n<taskbar:DesktopApp DesktopApplicationLinkPath=\"%s\" />", path) + xmlTemplate[docStartIndex:]

			log.WithFields(log.Fields{
				"path": path,
			}).Info("Path added to XML template")
		}
	}

	fmt.Println(xmlTemplate)
}
