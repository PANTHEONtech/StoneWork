package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"go.pantheon.tech/stonework/client"
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

	defer func(path string) {
		err = os.RemoveAll(path)
		if err != nil {
			err = fmt.Errorf("can't remove all files in temporary directory %s due to %v", path, err)
		}
	}(dirName)

	fullName := filepath.Join(dirName, "interfaces.txt")
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
		writeReportData(cli, "interfaces.txt", dirName, components, writeInterfaces),
		writeReportData(cli, "status.txt", dirName, components, writeStatus),
		writeReportData(cli, "status.json", dirName, components, writeStatusAsJson),
		writeReportData(cli, "docker-compose.yaml", dirName, components, writeDockerComposeConfig),
		writeReportData(cli, "docker-ps.txt", dirName, components, writeDockerContainers),
	}

	for _, comp := range components {
		if comp.GetMode() != client.ComponentAuxiliary && comp.GetMode() != client.ComponentUnknown {
			info := comp.GetInfo()
			alias := fmt.Sprintf("%s-", comp.GetName())

			if serviceName, ok := comp.GetMetadata()["containerServiceName"]; ok {
				err = writeReportData(cli, strings.ToLower(alias)+"docker-logs"+".log", dirName, components, writeDockerLogs, serviceName)
				if err != nil {
					fmt.Println(alias, err)
				}
			}
			buffer := strings.ToLower(alias) + "vppagent-report"
			err = writeReportData(cli, buffer+".zip", dirName, components, writeAgentCtlInfo, info.IPAddr, info.HTTPPort)
			errors = append(errors, err)
			err := os.Mkdir(path.Join(dirName, buffer), 0777)
			_ = err

			err = extractZip(dirName+"/"+buffer+".zip", path.Join(dirName, buffer))
			if err != nil {
				return err
			}

			err = os.Remove(dirName + "/" + buffer + ".zip")
			if err != nil {
				return err
			}

			for _, comp := range components {
				if sn, ok := comp.GetMetadata()["containerID"]; ok {
					err = writeReportData(cli, alias+"docker-inspect.txt", dirName, nil, writeDockerInspect, sn)
					if err != nil {
						errors = append(errors, err)
					}
				}

			}

		}
	}

	for _, err2 := range errors {
		if err2 != nil {
			logrus.Warnln(err2)
		}
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
	file, err := os.OpenFile(fullName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		err = fmt.Errorf("can't open file %v due to: %v", fullName, err)
		return
	}
	defer file.Close()

	err = writeFunc(cli, file, components, args...)
	return err
}

func extractZip(sourceZip string, destinationFolder string) error {
	zipReader, err := zip.OpenReader(sourceZip)
	if err != nil {
		return err
	}

	for _, file := range zipReader.File {
		zippedFile, _ := file.Open()

		destinationFile, err := os.OpenFile(destinationFolder+"/"+file.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			return err
		}

		_, err = io.Copy(destinationFile, zippedFile)
		if err != nil {
			return err
		}
		destinationFile.Close()
		zippedFile.Close()
	}
	return nil
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

func writeDockerComposeConfig(cli Cli, w io.Writer, components []client.Component, otherArgs ...interface{}) error {
	stdout, _, err := cli.Exec("docker compose config", []string{})
	if err != nil {
		return err
	}
	fmt.Fprintln(w, stdout)
	return nil
}

func writeDockerContainers(cli Cli, w io.Writer, components []client.Component, otherArgs ...interface{}) error {
	stdout, _, err := cli.Exec("docker compose ps --all", []string{})
	if err != nil {
		return err
	}
	fmt.Fprintln(w, stdout)
	return nil
}

func writeDockerInspect(cli Cli, w io.Writer, components []client.Component, otherArgs ...interface{}) error {
	cmd := fmt.Sprintf("docker inspect %s", fmt.Sprintf("%s", otherArgs[0]))
	stdout, _, err := cli.Exec(cmd, []string{})
	if err != nil {
		return err
	}
	fmt.Fprintln(w, stdout)
	return nil
}

func writeAgentCtlInfo(cli Cli, w io.Writer, components []client.Component, args ...interface{}) error {
	tempDirName, err := os.MkdirTemp("", "agentctl-reports-*")
	if err != nil {
		return fmt.Errorf("can't create tmp directory with namedue to %v", err)
	}
	defer os.RemoveAll(tempDirName)

	host := args[0]
	port := args[1]

	cmd := fmt.Sprintf("agentctl report --host %s --http-port %d -o %s -i", host, port, tempDirName)
	_, _, err = cli.Exec(cmd, []string{})
	if err != nil {
		return err
	}
	// extract ,delete zip, read files
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

	_, err = w.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func writeDockerLogs(cli Cli, w io.Writer, components []client.Component, args ...interface{}) error {
	serviceName := args[0]
	cmd := fmt.Sprintf("docker compose logs --no-color -n 10000 %s", serviceName)
	stdout, _, err := cli.Exec(cmd, []string{})
	if err != nil {
		return err
	}
	fmt.Fprintln(w, stdout)
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

	err = addWholeFolderToZip(zipWriter, dirName)
	if err != nil {
		return err
	}

	return nil
}

func addWholeFolderToZip(zipWriter *zip.Writer, dirName string) error {

	err := filepath.Walk(dirName, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		relPath := strings.TrimPrefix(filePath, dirName)
		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
