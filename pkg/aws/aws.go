package aws

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/sirupsen/logrus"
	"github.com/takutakahashi/billcap-schema/pkg/schema"
)

func Execute(ctx context.Context, baseCurrency string) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	ceClient := costexplorer.NewFromConfig(cfg)
	targetGranularity := types.GranularityDaily
	format := map[types.Granularity]string{
		types.GranularityHourly: "2006-01-02T15:04:05Z",
		types.GranularityDaily:  "2006-01-02",
	}
	now := time.Now()
	today := now.Format(format[targetGranularity])
	yesterday := now.AddDate(0, 0, -1).Format(format[targetGranularity])
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
		Granularity: targetGranularity,
		Metrics:     []string{"UnblendedCost", "UsageQuantity"},
	}

	output, err := ceClient.GetCostAndUsage(ctx, input)
	if err != nil {
		log.Fatalf("failed to get cost and usage, %v", err)
	}

	for _, result := range output.ResultsByTime {
		for _, group := range result.Groups {
			ret := schema.TransformedData{
				Time:              now,
				SchemaVersion:     schema.SchemaVersionTransformedData,
				Owner:             "takutakahashi",
				Project:           "lab",
				Provider:          "AWS",
				Service:           group.Keys[0],
				SKU:               group.Keys[1],
				CostAmount:        parseSize(*group.Metrics["UnblendedCost"].Amount),
				CostAmountUnit:    *group.Metrics["UnblendedCost"].Unit,
				UsageQuantity:     parseSize(*group.Metrics["UsageQuantity"].Amount),
				UsageQuantityUnit: *group.Metrics["UsageQuantity"].Unit,
				ExchangeRate:      150,
				TotalCost:         parseSize(*group.Metrics["UnblendedCost"].Amount) * 150,
				TotalUnit:         baseCurrency,
			}
			logrus.Info(ret)
		}
	}
}
func parseSize(str string) float64 {
	// 数値以外の文字を取り除く
	trimmed := strings.TrimFunc(str, func(r rune) bool {
		return !unicode.IsNumber(r) && r != '.'
	})

	if trimmed == "" {
		return -1 // 数値が見つからない場合
	}

	// float64 にパースする
	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return -1 // パースに失敗した場合
	}
	return value
}
