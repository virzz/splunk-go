package main

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func generateSPL(cmd *cobra.Command, args []string) string {
	earliest, _ := cmd.Flags().GetString("earliest")
	latest, _ := cmd.Flags().GetString("latest")
	dedups, _ := cmd.Flags().GetStringArray("dedup")
	rename, _ := cmd.Flags().GetStringToString("rename")
	table, _ := cmd.Flags().GetStringArray("table")
	spl, _ := cmd.Flags().GetString("spl")

	buf := strings.Builder{}
	buf.WriteString(strings.Join(args, " "))
	if earliest != "" {
		buf.WriteString(" earliest=" + earliest)
	}
	if latest != "" {
		buf.WriteString(" latest=" + latest)
	}
	if spl != "" {
		buf.WriteString(spl)
	}
	if len(dedups) > 0 {
		buf.WriteString(" | dedup " + strings.Join(dedups, ","))
	}
	for k, v := range rename {
		buf.WriteString(" | rename " + k + " TO " + v)
	}
	buf.WriteString(" | table " + strings.Join(table, " "))
	return buf.String()
}

var searchPFlag = pflag.NewFlagSet("search", pflag.ContinueOnError)

func init() {
	searchPFlag.StringP("spl", "a", "", "other spl")
	searchPFlag.StringP("earliest", "e", "-1h", "earliest")
	searchPFlag.StringP("latest", "l", "now", "latest")
	searchPFlag.StringToStringP("rename", "r", map[string]string{}, "rename")
	searchPFlag.StringArrayP("dedup", "d", []string{}, "dedup fields")
	searchPFlag.StringArrayP("table", "t", []string{}, "table fields")
}
