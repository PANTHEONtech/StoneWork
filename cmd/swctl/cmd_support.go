package main

import (
	"archive/zip"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.pantheon.tech/stonework/client"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type SupportCmdOptions struct {
	OutputDirectory string
}

func NewSupportCmd(cli Cli) *cobra.Command {
	var opts SupportCmdOptions
	cmd := &cobra.Command{
		Use:                "support [flags]",
		Short:              "Export support data",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runSupportCmd(cli, opts, args)
		},
	}
	return cmd
}

func runSupportCmd(cli Cli, opts SupportCmdOptions, args []string) error {
	// create report time and dependent variables
	reportTime := time.Now()
	reportName := fmt.Sprintf("swctl-report--%s",
		strings.ReplaceAll(reportTime.UTC().Format("2006-01-02--15-04-05-.000"), ".", ""))

	// create temporal directory
	dirNamePattern := fmt.Sprintf("%v--*", reportName)
	dirName, err := os.MkdirTemp("", dirNamePattern)
	if err != nil {
		return fmt.Errorf("can't create tmp directory with name pattern %s due to %v", dirNamePattern, err)
	}
	defer os.RemoveAll(dirName)

	fullName := filepath.Join(dirName, "Interfaces.txt")
	f, err := os.OpenFile(fullName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		err = fmt.Errorf("can't open file %v due to: %v", fullName, err)
		return err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = fmt.Errorf("can't close file %v due to: %v", fullName, closeErr)
		}
	}()

	components, err := cli.Client().GetComponents()
	if err != nil {
		return err
	}

	errors := []error{
		writeReportData(cli, "Interfaces.txt", dirName, components, writeInterfaces),
		writeReportData(cli, "Status.txt", dirName, components, writeStatus),
		writeReportData(cli, "Status.json", dirName, components, writeStatusAsJson),
	}

	for _, err := range errors {
		if err != nil {
			return err
		}
	}

	for _, comp := range components {
		if comp.GetMode() != client.ComponentAuxiliary && comp.GetMode() != client.ComponentUnknown {
			info := comp.GetInfo()
			if serviceName, ok := comp.GetMetadata()["containerServiceName"]; ok {
				writeReportData(cli, "docker-logs-"+serviceName, dirName, components, writeDockerLogs, serviceName)
			}
			alias := fmt.Sprintf("%s-%s", strings.Replace(comp.GetMode().String(), " ", "-", -1), comp.GetName())
			if err = writeReportData(cli, alias+".zip", dirName, components, writeAgentCtlInfo, info.IPAddr, info.HTTPPort); err != nil {
				return err
			}
		}
	}

	if err != nil {
		err = fmt.Errorf("can't open file %v due to: %v", fullName, err)
		return err
	}

	// resolve zip file name
	simpleZipFileName := reportName + ".zip"
	zipFileName := filepath.Join(opts.OutputDirectory, simpleZipFileName)
	if opts.OutputDirectory == "" {
		zipFileName, err = filepath.Abs(simpleZipFileName)
		if err != nil {
			return fmt.Errorf("can't find out absolute path for output zip file due to: %v\n\n", err)
		}
	}

	// combine report files into one zip file
	if _, err := cli.Out().Write([]byte("Creating report zip file... ")); err != nil {
		return err
	}
	if err := createZipFile(zipFileName, dirName); err != nil {
		return fmt.Errorf("can't create zip file(%v) due to: %v", zipFileName, err)
	}
	if _, err := cli.Out().Write([]byte(fmt.Sprintf("Done.\nReport file: %v\n", zipFileName))); err != nil {
		return err
	}

	return nil
}

func writeReportData(cli Cli, fileName string, dirName string, components []client.Component, writeFunc func(Cli, io.Writer, []client.Component, ...interface{}) error,
	args ...interface{}) (err error) {
	fullName := filepath.Join(dirName, fileName)
	f, err := os.OpenFile(fullName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		err = fmt.Errorf("can't open file %v due to: %v", fullName, err)
		return
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = fmt.Errorf("can't close file %v due to: %v", fullName, closeErr)
		}
	}()

	// append some report to file
	err = writeFunc(cli, f, components, args...)
	return
}

func writeInterfaces(cli Cli, w io.Writer, components []client.Component, otherArgs ...interface{}) error {
	for _, compo := range components {
		if sn, ok := compo.GetMetadata()["containerServiceName"]; ok {
			cmd := fmt.Sprintf("vpp-probe --color never --env=%s --query label=%s=%s discover", defaultVppProbeEnv, client.DockerComposeServiceLabel, sn)
			stdout, _, err := cli.Exec(cmd, []string{})
			if err != nil {
				if ee, ok := err.(*exec.ExitError); ok {
					logrus.Tracef("vpp-probe discover failed for service %s with error: %v: %s", sn, ee.String(), ee.Stderr)
					continue
				}
			}
			fmt.Fprintln(w, stdout)
		}
	}
	return nil
}

func writeStatus(cli Cli, w io.Writer, components []client.Component, otherArgs ...interface{}) error {
	infos, err := getStatusInfo(components)
	if err != nil {
		return err
	}
	printStatusTable(w, infos, false)

	return nil
}

func writeStatusAsJson(cli Cli, w io.Writer, components []client.Component, otherArgs ...interface{}) error {
	infos, err := getStatusInfo(components)
	if err != nil {
		return err
	}

	if err := formatAsTemplate(w, "json", infos); err != nil {
		return err
	}

	return nil
}

func writeAgentCtlInfo(cli Cli, w io.Writer, components []client.Component, args ...interface{}) error {
	tempDirName, err := os.MkdirTemp("", "agentctl-reports-*")
	defer os.RemoveAll(tempDirName)

	host := args[0]
	port := args[1]

	cmd := fmt.Sprintf("agentctl report --host %s --http-port %d -o %s -i", host, port, tempDirName)
	_, _, err = cli.Exec(cmd, []string{})
	if err != nil {
		return err
	}
	// fmt.Println(stdout)

	files, err := os.ReadDir(tempDirName)
	if err != nil {
		return err
	}
	zipFilename := filepath.Join(tempDirName, files[0].Name())
	file, err := os.Open(zipFilename)
	if err != nil {
		return err
	}
	defer file.Close()

	stats, err := file.Stat()
	data := make([]byte, stats.Size())
	if _, err = file.Read(data); err != nil {
		return err
	}

	w.Write(data)
	return nil
}

func writeDockerLogs(cli Cli, w io.Writer, components []client.Component, args ...interface{}) error {
	// container := fmt.Sprintf("%s", args[0])
	// logs, err := cli.Client().GetLogs(container)
	// if err != nil {
	//	return err
	// }
	// fmt.Println(logs)
	return nil
}

func createZipFile(zipFileName string, dirName string) (err error) {
	// create zip writer
	zipFile, err := os.Create(zipFileName)
	if err != nil {
		return fmt.Errorf("can't create empty zip file(%v) due to: %v", zipFileName, err)
	}
	defer func() {
		if closeErr := zipFile.Close(); closeErr != nil {
			err = fmt.Errorf("can't close zip file %v due to: %v", zipFileName, closeErr)
		}
	}()
	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if closeErr := zipWriter.Close(); closeErr != nil {
			err = fmt.Errorf("can't close zip file writer for zip file %v due to: %v", zipFileName, closeErr)
		}
	}()

	// Add files to zip
	dirItems, err := os.ReadDir(dirName)
	if err != nil {
		return fmt.Errorf("can't read report directory(%v) due to: %v", dirName, err)
	}
	for _, dirItem := range dirItems {
		if !dirItem.IsDir() {
			if err = addFileToZip(zipWriter, filepath.Join(dirName, dirItem.Name())); err != nil {
				return fmt.Errorf("can't add file dirItem.Name() to report zip file due to: %v", err)
			}
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	// open file for addition
	fileToZip, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("can't open file %v due to: %v", filename, err)
	}
	defer func() {
		if closeErr := fileToZip.Close(); closeErr != nil {
			err = fmt.Errorf("can't close zip file %v opened "+
				"for file appending due to: %v", filename, closeErr)
		}
	}()

	// get information from file for addition
	info, err := fileToZip.Stat()
	if err != nil {
		return fmt.Errorf("can't get information about file (%v) "+
			"that should be added to zip file due to: %v", filename, err)
	}

	// add file to zip file
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("can't create zip file info header for file %v due to: %v", filename, err)
	}
	header.Method = zip.Deflate // enables compression
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("can't create zip header for file %v due to: %v", filename, err)
	}
	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		return fmt.Errorf("can't copy content of file %v to zip file due to: %v", filename, err)
	}
	return nil
}
