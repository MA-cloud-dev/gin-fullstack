package unitrade

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type actionHelpSpec struct {
	Name        string
	Description string
}

type commandHelpSpec struct {
	Name        string
	Description string
	Usage       string
	Actions     []actionHelpSpec
}

var commandHelpSpecs = []commandHelpSpec{
	{
		Name:        "login",
		Description: "Save a PAT-backed session",
		Usage:       "unitrade login [flags]",
	},
	{
		Name:        "logout",
		Description: "Remove a saved profile",
		Usage:       "unitrade logout [flags]",
	},
	{
		Name:        "whoami",
		Description: "Show the active session",
		Usage:       "unitrade whoami [flags]",
	},
	{
		Name:        "auth",
		Description: "Campus authentication workflows",
		Usage:       "unitrade auth <action> [flags]",
		Actions: []actionHelpSpec{
			{Name: "list", Description: "List campus authentication requests"},
			{Name: "show", Description: "Show campus authentication details"},
			{Name: "approve", Description: "Approve a campus authentication request"},
			{Name: "reject", Description: "Reject a campus authentication request"},
			{Name: "revoke", Description: "Revoke a campus authentication review"},
		},
	},
	{
		Name:        "announcement",
		Description: "Campus announcement workflows",
		Usage:       "unitrade announcement <action> [flags]",
		Actions: []actionHelpSpec{
			{Name: "list", Description: "List campus announcements"},
			{Name: "show", Description: "Show campus announcement details"},
			{Name: "create", Description: "Create a campus announcement"},
			{Name: "update", Description: "Update a campus announcement"},
			{Name: "publish", Description: "Publish a campus announcement"},
			{Name: "unpublish", Description: "Unpublish a campus announcement"},
			{Name: "delete", Description: "Delete a campus announcement"},
		},
	},
	{
		Name:        "product",
		Description: "Campus product workflows",
		Usage:       "unitrade product <action> [flags]",
		Actions: []actionHelpSpec{
			{Name: "list", Description: "List campus products"},
			{Name: "show", Description: "Show campus product details"},
			{Name: "set-status", Description: "Update a campus product status"},
		},
	},
	{
		Name:        "report",
		Description: "Campus report workflows",
		Usage:       "unitrade report <action> [flags]",
		Actions: []actionHelpSpec{
			{Name: "list", Description: "List campus reports"},
			{Name: "show", Description: "Show campus report details"},
			{Name: "handle", Description: "Handle a campus report"},
		},
	},
	{
		Name:        "category",
		Description: "Campus category workflows",
		Usage:       "unitrade category <action> [flags]",
		Actions: []actionHelpSpec{
			{Name: "list", Description: "List campus categories"},
			{Name: "show", Description: "Show campus category details"},
			{Name: "create", Description: "Create a campus category"},
			{Name: "update", Description: "Update a campus category"},
			{Name: "enable", Description: "Enable a campus category"},
			{Name: "disable", Description: "Disable a campus category"},
		},
	},
	{
		Name:        "user",
		Description: "Campus user workflows",
		Usage:       "unitrade user <action> [flags]",
		Actions: []actionHelpSpec{
			{Name: "list", Description: "List campus users"},
			{Name: "show", Description: "Show campus user details"},
			{Name: "enable", Description: "Enable a campus user"},
			{Name: "disable", Description: "Disable a campus user"},
		},
	},
	{
		Name:        "staff",
		Description: "Campus staff workflows",
		Usage:       "unitrade staff <action> [flags]",
		Actions: []actionHelpSpec{
			{Name: "list", Description: "List campus staff"},
			{Name: "show", Description: "Show campus staff details"},
			{Name: "enable", Description: "Enable a campus staff account"},
			{Name: "disable", Description: "Disable a campus staff account"},
		},
	},
	{
		Name:        "order",
		Description: "Campus order queries",
		Usage:       "unitrade order <action> [flags]",
		Actions: []actionHelpSpec{
			{Name: "list", Description: "List campus orders"},
			{Name: "show", Description: "Show campus order details"},
		},
	},
}

func isHelpToken(value string) bool {
	switch strings.TrimSpace(value) {
	case "help", "-h", "--help":
		return true
	default:
		return false
	}
}

func commandHelpSpecByName(name string) (commandHelpSpec, bool) {
	for _, spec := range commandHelpSpecs {
		if spec.Name == name {
			return spec, true
		}
	}
	return commandHelpSpec{}, false
}

func (a *App) printRootUsage() {
	fmt.Fprintln(a.stdout, "Unitrade campus CLI")
	fmt.Fprintln(a.stdout, "")
	fmt.Fprintln(a.stdout, "Usage:")
	fmt.Fprintln(a.stdout, "  unitrade <command> <action> [flags]")
	fmt.Fprintln(a.stdout, "")
	fmt.Fprintln(a.stdout, "Commands:")

	w := tabwriter.NewWriter(a.stdout, 0, 4, 2, ' ', 0)
	for _, spec := range commandHelpSpecs {
		fmt.Fprintf(w, "  %s\t%s\n", spec.Name, spec.Description)
		if len(spec.Actions) == 0 {
			continue
		}
		for _, action := range spec.Actions {
			fmt.Fprintf(w, "    %s\t%s\n", action.Name, action.Description)
		}
	}
	_ = w.Flush()

	fmt.Fprintln(a.stdout, "")
	fmt.Fprintln(a.stdout, "Help:")
	fmt.Fprintln(a.stdout, "  unitrade <command> -h           Show command help")
	fmt.Fprintln(a.stdout, "  unitrade <command> <action> -h  Show action flags")
}

func (a *App) printCommandUsage(commandName string) bool {
	spec, ok := commandHelpSpecByName(commandName)
	if !ok {
		return false
	}

	fmt.Fprintf(a.stdout, "Unitrade %s command\n\n", spec.Name)
	fmt.Fprintln(a.stdout, "Usage:")
	fmt.Fprintf(a.stdout, "  %s\n", spec.Usage)
	fmt.Fprintln(a.stdout, "")
	fmt.Fprintln(a.stdout, "Description:")
	fmt.Fprintf(a.stdout, "  %s\n", spec.Description)

	if len(spec.Actions) > 0 {
		fmt.Fprintln(a.stdout, "")
		fmt.Fprintln(a.stdout, "Actions:")
		printIndentedActionTable(a.stdout, spec.Actions)
		fmt.Fprintln(a.stdout, "")
		fmt.Fprintln(a.stdout, "Help:")
		fmt.Fprintf(a.stdout, "  unitrade %s <action> -h  Show action flags\n", spec.Name)
	}
	return true
}

func printIndentedActionTable(out io.Writer, actions []actionHelpSpec) {
	w := tabwriter.NewWriter(out, 0, 4, 2, ' ', 0)
	for _, action := range actions {
		fmt.Fprintf(w, "  %s\t%s\n", action.Name, action.Description)
	}
	_ = w.Flush()
}
