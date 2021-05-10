package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stackpulse/steps-sdk-go/log"
	"github.com/stackpulse/steps-sdk-go/upload"
	"sort"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/stackpulse/steps-sdk-go/step"
	events "github.com/stackpulse/steps/kubectl/base/events/get"
	"github.com/stackpulse/steps/kubectl/base/troubleshoot"
)

type Troubleshoot struct {
	troubleshoot *troubleshoot.Troubleshoot
}

type Outputs struct {
	Raw      troubleshoot.TroubleshootOutput `json:"raw"`
	Markdown string                          `json:"markdown"`
}

func (t *Troubleshoot) Init() error {
	troubleshootGet, err := troubleshoot.NewTroubleshoot(nil)
	if err != nil {
		return err
	}
	t.troubleshoot = troubleshootGet
	return nil
}

func (t *Troubleshoot) Run() (int, []byte, error) {
	troubleshootOutput, exitCode, err := t.troubleshoot.Run()
	if err != nil {
		return exitCode, troubleshootOutput, err
	}

	parsedTroubleshoot, err := t.troubleshoot.ParseObject(troubleshootOutput)
	if err != nil {
		return step.ExitCodeFailure, troubleshootOutput, fmt.Errorf("parse troubleshoot output: %w", err)
	}

	output := &Outputs{
		Raw:      parsedTroubleshoot,
		Markdown: produceMarkdown(parsedTroubleshoot),
	}

	if err := upload.RichOutput(context.Background(), strings.NewReader(output.Markdown), upload.ContentTypeMarkdown); err != nil {
		log.Logln("Error uploading rich output:", err.Error())
	}

	res, err := json.Marshal(output)
	if err != nil {
		return step.ExitCodeFailure, troubleshootOutput, fmt.Errorf("encode json output: %w", err)
	}

	return step.ExitCodeOK, res, nil
}

func sortKindEvents(events []events.Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].LastTimestamp.After(events[j].LastTimestamp)
	})
}

func renderFindingsMarkdown(hasFindingsKinds []string, troubleshootOutput troubleshoot.TroubleshootOutput) string {
	var finalMarkdown strings.Builder
	finalMarkdown.WriteString("## Findings\n---\n&nbsp;\n")

	for _, kind := range hasFindingsKinds {
		finalMarkdown.WriteString(fmt.Sprintf("### %s\n", kind))
		var curnntTable strings.Builder
		table := tablewriter.NewWriter(&curnntTable)
		// Markdown table settings
		table.SetRowLine(false)
		table.SetAutoWrapText(false)
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")

		table.SetHeader([]string{"Last_Seen", "Object_Name", "Namespace", "Reason", "Message"})
		foundEvents := troubleshootOutput[kind]
		sortKindEvents(foundEvents)

		rows := make([][]string, 0, len(foundEvents))
		for _, evt := range foundEvents {
			rows = append(rows, []string{
				time.Now().Sub(evt.LastTimestamp).Round(time.Second).String(),
				evt.ObjectName,
				evt.ObjectNamespace,
				evt.Reason,
				evt.Message,
			})
		}
		table.AppendBulk(rows)
		table.Render()
		finalMarkdown.WriteString(curnntTable.String())
		finalMarkdown.WriteString("\n&nbsp;\n")
	}

	return finalMarkdown.String()
}

func sortKinds(troubleshootOutput troubleshoot.TroubleshootOutput) (hasFindings, noFindings []string) {
	hasFindings = make([]string, 0)
	noFindings = make([]string, 0)

	for kind, foundEvents := range troubleshootOutput {
		if len(foundEvents) > 0 {
			hasFindings = append(hasFindings, kind)
		} else {
			noFindings = append(noFindings, kind)
		}
	}

	sort.Strings(hasFindings)
	sort.Strings(noFindings)
	return
}

func produceMarkdown(troubleshootOutput troubleshoot.TroubleshootOutput) string {
	var finalMarkdown strings.Builder
	finalMarkdown.WriteString("## Summary\n&nbsp;\n")

	var summaryTable strings.Builder
	table := tablewriter.NewWriter(&summaryTable)

	// Markdown table settings
	table.SetRowLine(false)
	table.SetAutoWrapText(false)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")

	table.SetHeader([]string{"", "Kind"})
	rows := make([][]string, 0, len(troubleshootOutput))

	hasFindings, noFindings := sortKinds(troubleshootOutput)

	for _, kind := range hasFindings {
		rows = append(rows, []string{":x:", kind})
	}
	for _, kind := range noFindings {
		rows = append(rows, []string{":white_check_mark:", kind})
	}

	table.AppendBulk(rows)
	table.Render()
	finalMarkdown.WriteString(summaryTable.String())

	if len(hasFindings) == 0 {
		return finalMarkdown.String()
	}

	finalMarkdown.WriteString("\n&nbsp;\n&nbsp;\n")
	finalMarkdown.WriteString(renderFindingsMarkdown(hasFindings, troubleshootOutput))
	return finalMarkdown.String()
}

func main() {
	step.Run(&Troubleshoot{})
}
