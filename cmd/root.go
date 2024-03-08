/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/takutakahashi/billcap-aws/pkg/aws"
	"github.com/takutakahashi/billcap-schema/pkg/store"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "billcap-aws",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		owner := cmd.Flag("owner").Value.String()
		project := cmd.Flag("project").Value.String()
		currency := cmd.Flag("currency").Value.String()
		ctx := context.Background()
		from := cmd.Flag("date-from").Value.String()
		if from == "" {
			from = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		}
		to := cmd.Flag("date-to").Value.String()
		if to == "" {
			to = time.Now().Format("2006-01-02")
		}

		data, err := aws.Execute(ctx, owner, project, currency, from, to)
		if err != nil {
			logrus.Fatal(err)
		}
		for _, d := range data {
			fmt.Println(d)
		}

		googleServiceAccountJSONPath := cmd.Flag("google-serviceaccount-json-path").Value.String()
		googleProjectID := cmd.Flag("google-project-id").Value.String()
		if googleProjectID == "" {
			logrus.Fatal("google-project-id is required")
		}
		googleDatasetID := cmd.Flag("google-dataset-id").Value.String()
		googleTableID := cmd.Flag("google-table-id").Value.String()
		s, err := store.NewBigQueryStore(context.Background(), store.BigQueryStoreConfig{
			CredentialsPath: googleServiceAccountJSONPath,
			ProjectID:       googleProjectID,
			Transformed: store.BigQueryDatabase{
				DatasetID: googleDatasetID,
				TableID:   googleTableID,
			},
		})
		if err != nil {
			logrus.Fatal(err)
		}
		if err := s.Setup(ctx); err != nil {
			logrus.Fatal(err)
		}
		if err := s.LoadTransformed(ctx, data); err != nil {
			logrus.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.billcap-aws.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("owner", "o", "", "Owner of data")
	rootCmd.Flags().StringP("project", "p", "", "Project name")
	rootCmd.Flags().StringP("currency", "c", "USD", "Base currency unit")
	rootCmd.Flags().StringP("store", "s", "bigquery", "Store driver")
	rootCmd.Flags().StringP("google-serviceaccount-json-path", "j", "/secret/serviceaccount.json", "Google Cloud Service Account JSON Path")
	rootCmd.Flags().StringP("google-project-id", "i", "", "Google Cloud Project ID")
	rootCmd.Flags().StringP("google-dataset-id", "d", "billcap", "Google Cloud Dataset ID")
	rootCmd.Flags().StringP("google-table-id", "t", "billcap", "Google Cloud Table ID")
	rootCmd.Flags().String("date-from", "", "Date from. Format: 2006-01-02")
	rootCmd.Flags().String("date-to", "", "Date to Format: 2006-01-02")

}
