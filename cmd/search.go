package main

import (
	"errors"
	"time"

	"github.com/spf13/cobra"
	"github.com/virzz/splunk-go"
)

var (
	searchCmd = &cobra.Command{Use: "search"}
	exportCmd = &cobra.Command{
		Use:   "export",
		Short: "export splunk search to output",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			output, _ := cmd.Flags().GetString("output")
			if output == "" {
				output = time.Now().Format("20060102150405.tmp")
			}
			mode, _ := cmd.Flags().GetString("mode")
			wait, _ := cmd.Flags().GetInt("wait")
			var om splunk.OutputMode
			switch mode {
			case "json":
				om = splunk.OutputModeJSON
			default:
				om = splunk.OutputModeCSV
			}
			if isJob, _ := cmd.Flags().GetBool("is-job"); !isJob {
				spl := generateSPL(cmd, args)
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
	exportCmd.Flags().StringP("output", "o", "", "output file")
	exportCmd.Flags().StringP("mode", "m", "csv", "output mode: csv,json")
	exportCmd.Flags().IntP("wait", "w", 3, "times to wait for search job to finish,-1 is forever")
	exportCmd.Flags().BoolP("is-job", "j", false, "use job id to export")

	searchCmd.AddCommand(exportCmd)
	rootCmd.AddCommand(searchCmd)
}
