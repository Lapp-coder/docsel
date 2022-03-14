package app

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/eiannone/keyboard"

	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const helpMessage = `Select the desired service by pressing the Space bar when the arrow on the left points to it. 
You can move the arrow to the keys j (down), k (up) or the arrows ↓ ↑. 
Press the Enter key to start your selected services.
`

const (
	NameRootCmd      = "docsel"
	ShortDescRootCmd = "Docsel is a utility that allows you to run your selected services in Docker based on the docker-compose file"
)

const (
	NameRunCmd      = "run"
	ShortDescRunCmd = "Launches the service selection panel"

	FlagPath             = "path"
	ShorthandFlagPath    = "p"
	DescFlagPath         = "Path to docker-compose file"
	DefaultFlagPathValue = "docker-compose.yaml"

	FlagDetach             = "detach"
	ShorthandFlagDetach    = "d"
	DescFlagDetach         = "Running selected services in the background"
	DefaultFlagDetachValue = false

	FlagRemove             = "rm"
	DescFlagRemove         = "Removes stopped service containers. Ignored in detached mode"
	DefaultFlagRemoveValue = false

	FlagSave             = "save"
	ShorthandFlagSave    = "s"
	DescFlagSave         = "Saves the docker-compose file generated with the selected services"
	DefaultFlagSaveValue = false

	FlagBuild             = "build"
	DescFlagBuild         = "Build images before starting containers"
	DefaultFlagBuildValue = false
)

var RootCmd = &cobra.Command{
	Use:   NameRootCmd,
	Short: ShortDescRootCmd,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := cmd.Help(); err != nil {
			return err
		}

		return nil
	},
}

const (
	firstNameGeneratedFile = "docker-compose.yaml"
	lastNameGeneratedFile  = "docsel.yaml"
)

var errEmptyServices = errors.New("empty services in docker-compose file")

var RunCmd = &cobra.Command{
	Use:   NameRunCmd,
	Short: ShortDescRunCmd,
	RunE: func(cmd *cobra.Command, args []string) error {
		cleanConsole()

		filepath := cmd.Flag(FlagPath).Value.String()
		fileContent, err := os.ReadFile(filepath)
		if err != nil {
			return err
		}

		var dc DockerCompose
		if err := yaml.Unmarshal(fileContent, &dc); err != nil {
			return err
		}

		if len(dc.Services) == 0 {
			return errEmptyServices
		}

		recordNumber := 1
		fmt.Println(helpMessage, generateDashboard(dc, false, recordNumber))

		var filename string
		switch filepath {
		case "docker-compose.yaml", "docker-compose.yml", "compose.yaml", "compose.yml":
			filename = lastNameGeneratedFile
		default:
			filename = firstNameGeneratedFile
		}

	Label:
		file, err := os.Create(filename)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				if err = os.Remove(filename); err != nil {
					return err
				}
				goto Label
			}

			return err
		}
		defer func() {
			if cmd.Flag(FlagSave).Value.String() == "false" {
				if err = file.Close(); err != nil {
					log.Fatalln(err)
				}
				if err = os.Remove(filename); err != nil {
					log.Fatalln(err)
				}
			} else {
				fileInfo, err := file.Stat()
				if err != nil {
					log.Fatalln(err)
				}

				if fileInfo.Size() > 0 {
					if err = file.Close(); err != nil {
						log.Fatalln(err)
					}

					if filename == lastNameGeneratedFile {
						return
					} else if filename == firstNameGeneratedFile {
						if err = os.Rename(filename, lastNameGeneratedFile); err != nil {
							log.Fatalln(err)
						}
						return
					}
				}

				if err = os.Remove(filename); err != nil {
					log.Fatalln(err)
				}
			}
		}()

		for {
			char, key, err := keyboard.GetSingleKey()
			if err != nil {
				return err
			}

			switch key {
			case keyboard.KeyEnter:
				if len(selectedServices) > 0 {
					for service := range dc.Services {
						if _, ok := selectedServices[service]; !ok {
							delete(dc.Services, service)
						}
					}

					b, err := yaml.Marshal(dc)
					if err != nil {
						return err
					}

					if _, err = file.WriteString(string(b)); err != nil {
						return err
					}

					if cmd.Flag(FlagDetach).Value.String() == "false" {
						var dockerCmd *exec.Cmd
						if cmd.Flag(FlagBuild).Value.String() == "true" {
							dockerCmd = exec.Command("docker-compose", "-f", filename, "up", "--build")
						} else {
							dockerCmd = exec.Command("docker-compose", "-f", filename, "up")
						}

						dockerCmd.Stderr = os.Stderr
						dockerCmd.Stdout = os.Stdout

						if err := dockerCmd.Start(); err != nil {
							return err
						}

						quit := make(chan os.Signal, 1)
						signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
						<-quit

						if cmd.Flag(FlagRemove).Value.String() == "false" {
							dockerCmd = exec.Command("docker-compose", "stop")
						} else {
							dockerCmd = exec.Command("docker-compose", "rm", "--stop", "-f")
						}

						dockerCmd.Stderr = os.Stderr
						dockerCmd.Stdout = os.Stdout

						return dockerCmd.Run()
					}

					var dockerCmd *exec.Cmd
					if cmd.Flag(FlagBuild).Value.String() == "true" {
						dockerCmd = exec.Command("docker-compose", "-f", filename, "up", "-d", "--build")
					} else {
						dockerCmd = exec.Command("docker-compose", "-f", filename, "up", "-d")
					}

					dockerCmd.Stderr = os.Stderr
					dockerCmd.Stdout = os.Stdout

					return dockerCmd.Run()
				}

				fmt.Println("none of the services has been selected")

				return nil
			case keyboard.KeyCtrlC:
				return nil
			case keyboard.KeySpace:
				cleanConsole()
				fmt.Println(helpMessage, generateDashboard(dc, true, recordNumber))
			default:
				if key == keyboard.KeyArrowUp || char == 'k' {
					if recordNumber > 1 {
						recordNumber--

						cleanConsole()
						fmt.Println(helpMessage, generateDashboard(dc, false, recordNumber))
					}
				} else if key == keyboard.KeyArrowDown || char == 'j' {
					if recordNumber < len(services) {
						recordNumber++

						cleanConsole()
						fmt.Println(helpMessage, generateDashboard(dc, false, recordNumber))
					}
				}
			}
		}
	},
}
