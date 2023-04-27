package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

func main() {
	paths := []string{
		"%ALLUSERSPROFILE%\\Microsoft\\Windows\\Start Menu\\Programs\\Google Chrome.lnk",
		"%appdata%\\Microsoft\\Windows\\Start Menu\\Programs\\System Tools\\Command Prompt.lnk",
		"C:\\Program Files\\Microsoft VS Code\\code.exe",
		"C:\\Program Files\\PowerShell\\7\\pwsh.exe",
		"C:\\Windows\\System32\\control.exe",
		"C:\\Windows\\System32\\windowspowershell\\v1.0\\powershell.exe",
		"C:\\Windows\\System32\\windowspowershell\\v1.0\\powershell.lnk",
		"C:\\Windows\\explorer.exe",
	}

	// Create the log file directory if it doesn't exist
	logFilePath := "C:\\ProgramData\\Streambox\\Beautypatch\\log\\run-log.txt"
	logDirPath := filepath.Dir(logFilePath)
	if _, err := os.Stat(logDirPath); os.IsNotExist(err) {
		err := os.MkdirAll(logDirPath, os.ModePerm)
		if err != nil {
			logrus.Fatalf("Failed to create log directory: %v", err)
		}
	}

	// Create a RotateLogs object
	rotateLogs, err := rotatelogs.New(
		logFilePath+".%Y%m%d%H%M%S",
		rotatelogs.WithLinkName(logFilePath),
		rotatelogs.WithRotationTime(time.Hour),
		rotatelogs.WithMaxAge(time.Hour*24*30*6),
		rotatelogs.WithRotationSize(1024*1024),
	)
	if err != nil {
		logrus.Fatalf("Failed to create RotateLogs object: %v", err)
	}

	// Initialize logrus
	// Create the logger and set the output to the log file
	logger := logrus.New()
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}
	logger.SetOutput(rotateLogs)

	// Log a message
	logger.Info("Logging to file and stdout")

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
			logger.Infof("Path exists: %s", path)
		} else {
			logger.Warnf("Path does not exist, but adding it to : %s", path)
		}
		taskbarApps += "        <taskbar:DesktopApp DesktopApplicationLinkPath=\"" + path + "\" />\n"
	}

	xmlContent := []byte(fmt.Sprintf(xmlTemplate, taskbarApps))

	configDir, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	filename := "Taskbar.xml"
	layoutPath := filepath.Join(configDir, "BeautyPatch", filename)

	xmlDirPath := filepath.Dir(layoutPath)
	if _, err := os.Stat(xmlDirPath); os.IsNotExist(err) {
		err := os.MkdirAll(xmlDirPath, os.ModePerm)
		if err != nil {
			logrus.Fatalf("Failed to create Taskbar.xml directory: %v", err)
		}
	}

	file, err := os.Create(layoutPath)
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

	mountPath := "c:\\"
	cmd := exec.Command("powershell", "-Command", fmt.Sprintf("Import-StartLayout -LayoutPath %s -MountPath %s", layoutPath, mountPath))
	logger.Infof("LayoutModificationTemplate written to Taskbar.xml file")

	// Log the command as a string
	logrus.Infof("Running command: %s", strings.Join(cmd.Args, " "))

	err = cmd.Run()
	if err != nil {
		logrus.WithError(err).Error("Failed to import StartLayout")
	}
}
