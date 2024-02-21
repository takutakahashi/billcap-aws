package aws

import (
	"context"
	"encoding/csv"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
)

func Execute(ctx context.Context) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	ceClient := costexplorer.NewFromConfig(cfg)

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// Cost Explorer APIリクエスト用の入力を作成
	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(yesterday),
			End:   aws.String(today),
		},
		GroupBy: []types.GroupDefinition{
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String("SERVICE"),
			},
			{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String("USAGE_TYPE"),
			},
		},
		Granularity: types.GranularityDaily,
		Metrics:     []string{"UnblendedCost", "UsageQuantity"},
	}

	// Cost Explorer APIを実行
	output, err := ceClient.GetCostAndUsage(ctx, input)
	if err != nil {
		log.Fatalf("failed to get cost and usage, %v", err)
	}

	// CSVライターを準備
	csvWriter := csv.NewWriter(os.Stdout)
	defer csvWriter.Flush()

	// CSVヘッダを書き出し
	headers := []string{"Start Date", "End Date", "Usage Type", "Unblended Cost", "Usage Quantity"}
	if err := csvWriter.Write(headers); err != nil {
		log.Fatalf("failed to write headers to CSV, %v", err)
	}

	// 結果をCSVに書き出し
	for _, result := range output.ResultsByTime {
		for _, group := range result.Groups {

			record := []string{
				*result.TimePeriod.Start,
				*result.TimePeriod.End,
				group.Keys[0],
				group.Keys[1],
				*group.Metrics["UsageQuantity"].Amount,
				*group.Metrics["UsageQuantity"].Unit,
				*group.Metrics["UnblendedCost"].Amount,
				*group.Metrics["UnblendedCost"].Unit,
			}
			if err := csvWriter.Write(record); err != nil {
				log.Fatalf("failed to write record to CSV, %v", err)
			}
		}
	}
}
