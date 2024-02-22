/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"

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
		data, err := aws.Execute(context.Background(), owner, project, currency)
		if err != nil {
			panic(err)
		}
		for _, d := range data {
			fmt.Println(d)
		}
		s, err := store.NewBigQueryStore(context.Background(), store.BigQueryStoreConfig{})
		if err != nil {
			panic(err)
		}
		_ = s
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
}
