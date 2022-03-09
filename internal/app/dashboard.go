package app

import (
	"fmt"
	"strings"
	"sync"
	"unicode/utf8"
)

const (
	widthDashboard    = 3
	numOtherRune      = 5
	numRuneForNewLine = 1
	numRuneForArrow   = 8 - numRuneForNewLine
	recordArrow       = " ─────> "
	charSelectedService   = "[*]"
	charNoSelectedService = "[ ]"
)

var (
	once                        sync.Once
	serviceNumbers              = make(map[string]int)
	selectedServices            = make(map[string]struct{})
	services                    = make([]string, 0)
	maxLenServiceName           int
	dashboardTopLine            string
	dashboardDownLine           string
	runeCountInDashboardTopLine int
)

func generateDashboard(dc DockerCompose, serviceSelected bool, serviceNumber int) string {
	var builder strings.Builder

	once.Do(func() {
		var counter int
		for service := range dc.Services {
			counter++
			serviceNumbers[service] = counter
			services = append(services, service)

			lenServiceName := utf8.RuneCountInString(service)
			if lenServiceName > maxLenServiceName {
				maxLenServiceName = lenServiceName
			}
		}

		dashboardTopLine = "┌"
		dashboardTopLine += strings.Repeat("─", maxLenServiceName+digitsCount(len(services))+numOtherRune+widthDashboard)
		dashboardTopLine += "┐"

		runeCountInDashboardTopLine = utf8.RuneCountInString(dashboardTopLine)

		dashboardDownLine = "\t└" + strings.Repeat("─", utf8.RuneCountInString(dashboardTopLine)-2) + "┘\n"
	})

	builder.WriteString("\n\t")
	builder.WriteString(dashboardTopLine + "\n")

	for _, service := range services {
		var (
			record         string
			recordHasArrow bool
		)

		if _, ok := selectedServices[service]; !ok {
			if serviceNumber > 0 {
				if serviceNumbers[service] == serviceNumber {
					if serviceSelected {
						record = fmt.Sprintf("%s│ %s %s", recordArrow, charSelectedService, service)
						selectedServices[service] = struct{}{}
						recordHasArrow = true
					} else {
						record = fmt.Sprintf("%s│ %s %s", recordArrow, charNoSelectedService, service)
						recordHasArrow = true
					}
				} else {
					record = fmt.Sprintf("\t│ %s %s", charNoSelectedService, service)
				}
			} else {
				record = fmt.Sprintf("\t│ %s %s", charNoSelectedService, service)
			}
		} else {
			if serviceNumber > 0 {
				if serviceNumbers[service] == serviceNumber {
					if serviceSelected {
						record = fmt.Sprintf("%s│ %s %s", recordArrow, charNoSelectedService, service)
						delete(selectedServices, service)
						recordHasArrow = true
					} else {
						record = fmt.Sprintf("%s│ %s %s", recordArrow, charSelectedService, service)
						recordHasArrow = true
					}
				} else {
					record = fmt.Sprintf("\t│ %s %s", charSelectedService, service)
				}
			} else {
				record = fmt.Sprintf("\t│ %s %s", charSelectedService, service)
			}
		}

		runeCountInRecord := utf8.RuneCountInString(record)
		numRepeatedEmptyChars := runeCountInDashboardTopLine - runeCountInRecord

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
