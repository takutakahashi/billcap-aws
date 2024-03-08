package aws

import (
	"context"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"github.com/takutakahashi/billcap-schema/pkg/schema"
)

func Execute(ctx context.Context, owner, project, baseCurrency string, from, to string) ([]schema.TransformedData, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	ceClient := costexplorer.NewFromConfig(cfg)
	input := &costexplorer.GetCostAndUsageInput{
		TimePeriod: &types.DateInterval{
			Start: aws.String(from),
			End:   aws.String(to),
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

	output, err := ceClient.GetCostAndUsage(ctx, input)
	if err != nil {
		return nil, err
	}
	ret := []schema.TransformedData{}
	for _, result := range output.ResultsByTime {
		for _, group := range result.Groups {
			t, err := time.Parse("2006-01-02", *result.TimePeriod.Start)
			if err != nil {
				return nil, err
			}
			ret = append(ret, schema.TransformedData{
				Time:              t,
				SchemaVersion:     schema.SchemaVersionTransformedData,
				Owner:             owner,
				Project:           project,
				Provider:          "AWS",
				Service:           group.Keys[0],
				SKU:               group.Keys[1],
				CostAmount:        parseSize(*group.Metrics["UnblendedCost"].Amount),
				CostAmountUnit:    *group.Metrics["UnblendedCost"].Unit,
				UsageQuantity:     parseSize(*group.Metrics["UsageQuantity"].Amount),
				UsageQuantityUnit: *group.Metrics["UsageQuantity"].Unit,
				ExchangeRate:      150,
				TotalCost:         parseSize(*group.Metrics["UnblendedCost"].Amount) * 150,
				TotalCostUnit:     baseCurrency,
			})
		}
	}
	return ret, nil
}
func parseSize(str string) float64 {
	trimmed := strings.TrimFunc(str, func(r rune) bool {
		return !unicode.IsNumber(r) && r != '.'
	})

	if trimmed == "" {
		return -1
	}

	value, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return -1
	}
	return value
}
