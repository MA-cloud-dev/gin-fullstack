package unitrade

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"time"
)

func printJSON(out io.Writer, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, string(data))
	return err
}

func printTable(out io.Writer, headers []string, rows [][]string) {
	w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	fmt.Fprintln(w, strings.Join(headers, "\t"))
	for _, row := range rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	_ = w.Flush()
}

func formatTimePtr(value *time.Time) string {
	if value == nil || value.IsZero() {
		return "-"
	}
	return value.Local().Format(time.DateTime)
}

func formatTimeValue(value time.Time) string {
	if value.IsZero() {
		return "-"
	}
	return value.Local().Format(time.DateTime)
}

func stringValue(value *string) string {
	if value == nil || *value == "" {
		return "-"
	}
	return *value
}

func chooseOptionalString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func uintPtrValue(value *uint) string {
	if value == nil {
		return "-"
	}
	return fmt.Sprintf("%d", *value)
}

func floatPtrString(value *float64) string {
	if value == nil {
		return "-"
	}
	return fmt.Sprintf("%.2f", *value)
}
