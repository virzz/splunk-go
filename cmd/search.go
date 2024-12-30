package main

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
	"github.com/virzz/splunk-go"
)

var (
	searchCmd = &cobra.Command{Use: "search", Aliases: []string{"s"}}
	exportCmd = &cobra.Command{
		Use: "export", Aliases: []string{"e"},
		Short: "export splunk search to output",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")
			mode, _ := cmd.Flags().GetString("mode")
			wait, _ := cmd.Flags().GetInt("wait")
			var om splunk.OutputMode
			switch mode {
			case "json":
				om = splunk.OutputModeJSON
			default:
				om = splunk.OutputModeCSV
			}
			if filepath.Ext(output) == "" {
				output = output + "." + mode
			}
			if isJob, _ := cmd.Flags().GetBool("is-job"); !isJob {
				spl := generateSPL(cmd, args)
				cmd.Println("SPL: " + spl)
				dryRun, _ := cmd.Flags().GetBool("dry-run")
				if dryRun {
					return nil
				}
				return splunk.Search.QueryAndDownload(spl, output, wait, om)
			}
			r, err := splunk.Search.Status(args[0], wait)
			if err != nil {
				return err
			}
			if !r.IsDone {
				return errors.New("job is not done")
			}
			if r.ResultCount == 0 {
				return errors.New("job result count is 0")
			}
			return splunk.Search.Download(args[0], output, om)
		},
	}
)

func init() {
	exportCmd.Flags().AddFlagSet(searchPFlag)
	exportCmd.Flags().StringP("output", "o", time.Now().Format("20060102150405"), "output file")
	exportCmd.Flags().StringP("mode", "m", "csv", "output mode: csv,json")
	exportCmd.Flags().IntP("wait", "w", 3, "times to wait for search job to finish,-1 is forever")
	exportCmd.Flags().BoolP("is-job", "j", false, "use job id to export")
	exportCmd.Flags().BoolP("dry-run", "x", false, "print splunk search query")

	searchCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(searchCmd)
}
