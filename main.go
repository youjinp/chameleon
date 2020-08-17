package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/spf13/cobra"
	"github.com/youjinp/chameleon/pkg/copy"
	"github.com/youjinp/chameleon/pkg/paste"
	"github.com/youjinp/chameleon/pkg/utils"
)

// TODO:
// - build lambda + sfn pipeline to export data from dynamodb to s3 (csv)
// - use previously built lambda + sfn pipeline to import data from s3 to dynamodb

// OR figure out how to use datapipeline? May be expensive, have to spin up and destroy clusters?

// This program downloads all data to local storage
// => not suitable for large dataset

var accessIDCopy string
var accessKeyCopy string
var sessionCopy string
var tableCopy string
var fileCopy string

var accessIDPaste string
var accessKeyPaste string
var sessionPaste string
var tablePaste string
var filePaste string

var rootCmd = &cobra.Command{
	Use:   "chameleon",
	Short: "A tool for copying DynamoDB data",
	Long:  `Chameleon is a CLI tool that helps with copying DynamoDB data.`,
}

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy dynamodb's data to a location",
	Run: func(cmd *cobra.Command, args []string) {
		// configure aws
		cfg, err := external.LoadDefaultAWSConfig(
			external.WithCredentialsProvider{
				CredentialsProvider: aws.NewStaticCredentialsProvider(accessIDCopy, accessKeyCopy, sessionCopy),
			},
		)
		utils.CheckError("Failed to load aws config", err)

		// download data
		dynamodbClient := dynamodb.New(cfg)
		b := copy.Copy{dynamodbClient, copy.Options{tableCopy, 1000, fileCopy}}
		b.Start()
	},
}

var pasteCmd = &cobra.Command{
	Use:   "paste",
	Short: "Paste dynamodb's data from a location",
	Run: func(cmd *cobra.Command, args []string) {
		// configure aws
		cfg, err := external.LoadDefaultAWSConfig(
			external.WithCredentialsProvider{
				CredentialsProvider: aws.NewStaticCredentialsProvider(accessIDPaste, accessKeyPaste, sessionPaste),
			},
		)
		utils.CheckError("Failed to load aws config", err)

		// download data
		dynamodbClient := dynamodb.New(cfg)
		b := paste.Paste{dynamodbClient, paste.Options{tablePaste, filePaste}}
		b.Start()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Chameleon",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Chameleon v0.0.1")
	},
}

func main() {

	rootCmd.AddCommand(copyCmd)
	rootCmd.AddCommand(pasteCmd)
	rootCmd.AddCommand(versionCmd)

	// env vars
	id := os.Getenv("AWS_ACCESS_KEY_ID")
	key := os.Getenv("AWS_SECRET_ACCESS_KEY")
	token := os.Getenv("AWS_SESSION_TOKEN")

	// copy flags
	copyCmd.Flags().StringVarP(&accessIDCopy, "id", "a", id, "aws access ID or set AWS_ACCESS_KEY_ID")
	copyCmd.Flags().StringVarP(&accessKeyCopy, "key", "k", key, "aws secret access ID or set AWS_SECRET_ACCESS_KEY")
	copyCmd.Flags().StringVarP(&sessionCopy, "token", "s", token, "aws session token or set AWS_SESSION_TOKEN")
	copyCmd.Flags().StringVarP(&tableCopy, "table", "t", "", "table name")
	copyCmd.Flags().StringVarP(&fileCopy, "output", "o", "data", "output path")
	copyCmd.MarkFlagRequired("table")

	// paste flags
	pasteCmd.Flags().StringVarP(&accessIDPaste, "id", "a", id, "aws access ID or set AWS_ACCESS_KEY_ID")
	pasteCmd.Flags().StringVarP(&accessKeyPaste, "key", "k", key, "aws secret access ID or set AWS_SECRET_ACCESS_KEY")
	pasteCmd.Flags().StringVarP(&sessionPaste, "token", "s", token, "aws session token or set AWS_SESSION_TOKEN")
	pasteCmd.Flags().StringVarP(&tablePaste, "table", "t", "", "table name")
	pasteCmd.Flags().StringVarP(&filePaste, "input", "i", "data", "input path")
	pasteCmd.MarkFlagRequired("table")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
