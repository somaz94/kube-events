package report

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/somaz94/kube-events/internal/event"
)

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
	colorReset  = "\033[0m"
)

// Summary holds grouped events and statistics.
type Summary struct {
	Groups       []event.ResourceGroup
	Events       []event.Event
	TotalEvents  int
	WarningCount int
	NormalCount  int
	Resources    int
}

// NewSummary creates a Summary from grouped and filtered events.
func NewSummary(groups []event.ResourceGroup, events []event.Event) *Summary {
	s := &Summary{
		Groups:      groups,
		Events:      events,
		TotalEvents: len(events),
		Resources:   len(groups),
	}
	for _, e := range events {
		switch e.Type {
		case "Warning":
			s.WarningCount++
		default:
			s.NormalCount++
		}
	}
	return s
}

// PrintColor outputs events grouped by resource with ANSI colors.
func (s *Summary) PrintColor(w io.Writer, summaryOnly bool) error {
	if summaryOnly {
		return s.printSummaryLine(w, true)
	}

	if s.TotalEvents == 0 {
		fmt.Fprintf(w, "%sNo events found.%s\n", colorGray, colorReset)
		return nil
	}

	for i, g := range s.Groups {
		if i > 0 {
			fmt.Fprintln(w)
		}

		nsLabel := ""
		if g.Key.Namespace != "" {
			nsLabel = fmt.Sprintf(" %s[%s]%s", colorCyan, g.Key.Namespace, colorReset)
		}

		fmt.Fprintf(w, "%s%s/%s%s%s (%d events)\n",
			colorBold, g.Key.Kind, g.Key.Name, colorReset, nsLabel, len(g.Events))

		for _, e := range g.Events {
			typeColor := colorGreen
			typeIcon := "  "
			if e.Type == "Warning" {
				typeColor = colorYellow
				typeIcon = "! "
			}

			age := formatAge(e.Age)
			fmt.Fprintf(w, "  %s%s%-18s%s %s%-8s%s %s\n",
				typeColor, typeIcon, e.Reason, colorReset,
				colorGray, age, colorReset,
				e.Message)
		}
	}

	fmt.Fprintln(w)
	return s.printSummaryLine(w, true)
}

// PrintPlain outputs events without ANSI colors.
func (s *Summary) PrintPlain(w io.Writer, summaryOnly bool) error {
	if summaryOnly {
		return s.printSummaryLine(w, false)
	}

	if s.TotalEvents == 0 {
		fmt.Fprintln(w, "No events found.")
		return nil
	}

	for i, g := range s.Groups {
		if i > 0 {
			fmt.Fprintln(w)
		}

		nsLabel := ""
		if g.Key.Namespace != "" {
			nsLabel = fmt.Sprintf(" [%s]", g.Key.Namespace)
		}

		fmt.Fprintf(w, "%s/%s%s (%d events)\n", g.Key.Kind, g.Key.Name, nsLabel, len(g.Events))

		for _, e := range g.Events {
			icon := "  "
			if e.Type == "Warning" {
				icon = "! "
			}
			age := formatAge(e.Age)
			fmt.Fprintf(w, "  %s%-18s %-8s %s\n", icon, e.Reason, age, e.Message)
		}
	}

	fmt.Fprintln(w)
	return s.printSummaryLine(w, false)
}

// PrintJSON outputs events as JSON.
func (s *Summary) PrintJSON(w io.Writer) error {
	output := struct {
		Summary struct {
			TotalEvents  int `json:"totalEvents"`
			WarningCount int `json:"warningCount"`
			NormalCount  int `json:"normalCount"`
			Resources    int `json:"resources"`
		} `json:"summary"`
		Groups []jsonGroup `json:"groups"`
	}{}

	output.Summary.TotalEvents = s.TotalEvents
	output.Summary.WarningCount = s.WarningCount
	output.Summary.NormalCount = s.NormalCount
	output.Summary.Resources = s.Resources

	for _, g := range s.Groups {
		jg := jsonGroup{
			Kind:      g.Key.Kind,
			Name:      g.Key.Name,
			Namespace: g.Key.Namespace,
		}
		for _, e := range g.Events {
			jg.Events = append(jg.Events, jsonEvent{
				Type:    e.Type,
				Reason:  e.Reason,
				Message: e.Message,
				Age:     formatAge(e.Age),
				Count:   e.Count,
			})
		}
		output.Groups = append(output.Groups, jg)
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(output)
}

type jsonGroup struct {
	Kind      string      `json:"kind"`
	Name      string      `json:"name"`
	Namespace string      `json:"namespace,omitempty"`
	Events    []jsonEvent `json:"events"`
}

type jsonEvent struct {
	Type    string `json:"type"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Age     string `json:"age"`
	Count   int32  `json:"count"`
}

// PrintMarkdown outputs events as a markdown table.
func (s *Summary) PrintMarkdown(w io.Writer) error {
	fmt.Fprintln(w, "## Kubernetes Events Summary")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "**%d** events across **%d** resources | ", s.TotalEvents, s.Resources)
	fmt.Fprintf(w, "Warning: **%d** | Normal: **%d**\n\n", s.WarningCount, s.NormalCount)

	if s.TotalEvents == 0 {
		fmt.Fprintln(w, "No events found.")
		return nil
	}

	fmt.Fprintln(w, "| Type | Resource | Reason | Age | Message |")
	fmt.Fprintln(w, "|------|----------|--------|-----|---------|")

	for _, g := range s.Groups {
		for _, e := range g.Events {
			resource := fmt.Sprintf("%s/%s", g.Key.Kind, g.Key.Name)
			if g.Key.Namespace != "" {
				resource = fmt.Sprintf("%s/%s [%s]", g.Key.Kind, g.Key.Name, g.Key.Namespace)
			}
			fmt.Fprintf(w, "| %s | %s | %s | %s | %s |\n",
				e.Type, resource, e.Reason, formatAge(e.Age), truncate(e.Message, 80))
		}
	}
	return nil
}

// PrintTable outputs events as an ASCII table.
func (s *Summary) PrintTable(w io.Writer) error {
	fmt.Fprintf(w, "%-9s %-40s %-20s %-8s %s\n", "TYPE", "RESOURCE", "REASON", "AGE", "MESSAGE")
	fmt.Fprintln(w, strings.Repeat("-", 120))

	for _, g := range s.Groups {
		for _, e := range g.Events {
			resource := fmt.Sprintf("%s/%s", g.Key.Kind, g.Key.Name)
			if g.Key.Namespace != "" {
				resource = fmt.Sprintf("[%s] %s/%s", g.Key.Namespace, g.Key.Kind, g.Key.Name)
			}
			fmt.Fprintf(w, "%-9s %-40s %-20s %-8s %s\n",
				e.Type, truncate(resource, 38), truncate(e.Reason, 18), formatAge(e.Age), truncate(e.Message, 50))
		}
	}

	fmt.Fprintln(w, strings.Repeat("-", 120))
	fmt.Fprintf(w, "Total: %d events, %d resources (Warning: %d, Normal: %d)\n",
		s.TotalEvents, s.Resources, s.WarningCount, s.NormalCount)
	return nil
}

func (s *Summary) printSummaryLine(w io.Writer, colorize bool) error {
	if colorize {
		fmt.Fprintf(w, "%sSummary:%s %d events, %d resources",
			colorBold, colorReset, s.TotalEvents, s.Resources)
		if s.WarningCount > 0 {
			fmt.Fprintf(w, " | %sWarning: %d%s", colorYellow, s.WarningCount, colorReset)
		}
		if s.NormalCount > 0 {
			fmt.Fprintf(w, " | %sNormal: %d%s", colorGreen, s.NormalCount, colorReset)
		}
		fmt.Fprintln(w)
	} else {
		fmt.Fprintf(w, "Summary: %d events, %d resources (Warning: %d, Normal: %d)\n",
			s.TotalEvents, s.Resources, s.WarningCount, s.NormalCount)
	}
	return nil
}

func formatAge(d interface{ Seconds() float64 }) string {
	sec := d.Seconds()
	switch {
	case sec < 60:
		return fmt.Sprintf("%ds", int(sec))
	case sec < 3600:
		return fmt.Sprintf("%dm", int(sec/60))
	case sec < 86400:
		return fmt.Sprintf("%dh", int(sec/3600))
	default:
		return fmt.Sprintf("%dd", int(sec/86400))
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
