package app

import (
	"fmt"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/fatih/color"
)

const (
	numOtherRune         = 16
	numRuneForGreenColor = 9
	numRuneForNewLine    = 1
	numRuneForArrow      = 8 - numRuneForNewLine
	recordArrow          = " ─────> "
)

var (
	once                        sync.Once
	serviceNumbers              = make(map[string]int)
	selectedServices            = make(map[string]struct{})
	services                    = make([]string, 0)
	maxKeyLen                   int
	dashboardTopLine            string
	dashboardDownLine           string
	runeCountInDashboardTopLine int
)

func generateDashboard(dc DockerCompose, serviceSelected bool, serviceNumber int) string {
	var builder strings.Builder

	once.Do(func() {
		var counter int
		for key := range dc.Services {
			counter++
			serviceNumbers[key] = counter
			services = append(services, key)

			lenKey := utf8.RuneCountInString(key)
			if lenKey > maxKeyLen {
				maxKeyLen = lenKey
			}
		}

		dashboardTopLine = "┌"
		dashboardTopLine += strings.Repeat("─", maxKeyLen+digitsCount(len(services))+numOtherRune)
		dashboardTopLine += "┐"

		runeCountInDashboardTopLine = utf8.RuneCountInString(dashboardTopLine)

		dashboardDownLine = "\t└" + strings.Repeat("─", utf8.RuneCountInString(dashboardTopLine)-2) + "┘\n"
	})

	builder.WriteString("\n\t")
	builder.WriteString(dashboardTopLine + "\n")

	for _, service := range services {
		var (
			record             string
			recordHasArrow     bool
			recordHasCheckmark bool
		)

		if _, ok := selectedServices[service]; !ok {
			if serviceNumber > 0 {
				if serviceNumbers[service] == serviceNumber {
					if serviceSelected {
						record = fmt.Sprintf("%s│ %d. %s%s", recordArrow, serviceNumbers[service], service, color.GreenString(charSelectedService))
						selectedServices[service] = struct{}{}
						recordHasCheckmark = true
						recordHasArrow = true
					} else {
						record = fmt.Sprintf("%s│ %d. %s", recordArrow, serviceNumbers[service], service)
						recordHasArrow = true
					}
				} else {
					record = fmt.Sprintf("\t│ %d. %s", serviceNumbers[service], service)
				}
			} else {
				record = fmt.Sprintf("\t│ %d. %s", serviceNumbers[service], service)
			}
		} else {
			if serviceNumber > 0 {
				if serviceNumbers[service] == serviceNumber {
					if serviceSelected {
						record = fmt.Sprintf("%s│ %d. %s", recordArrow, serviceNumbers[service], service)
						delete(selectedServices, service)
						recordHasArrow = true
					} else {
						record = fmt.Sprintf("%s│ %d. %s%s", recordArrow, serviceNumbers[service], service, color.GreenString(charSelectedService))
						recordHasCheckmark = true
						recordHasArrow = true
					}
				} else {
					record = fmt.Sprintf("\t│ %d. %s%s", serviceNumbers[service], service, color.GreenString(charSelectedService))
					recordHasCheckmark = true
				}
			} else {
				record = fmt.Sprintf("\t│ %d. %s%s", serviceNumbers[service], service, color.GreenString(charSelectedService))
				recordHasCheckmark = true
			}
		}

		runeCountInRecord := utf8.RuneCountInString(record)
		numRepeatedEmptyChars := runeCountInDashboardTopLine - runeCountInRecord

		if recordHasCheckmark {
			numRepeatedEmptyChars += numRuneForGreenColor
		}
		if recordHasArrow {
			numRepeatedEmptyChars += numRuneForArrow
		}

		record += strings.Repeat(" ", numRepeatedEmptyChars) + "│\n"

		builder.WriteString(record)
	}

	builder.WriteString(dashboardDownLine)

	return builder.String()
}

func digitsCount(number int) int {
	var counter int
	for number != 0 {
		number /= 10
		counter++
	}

	return counter
}
