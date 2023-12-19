package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	i "github.com/Snowflake-Labs/terraform-provider-snowflake/pkg/scripts/issues"
)

func main() {
	issues := loadIssues()
	processedIssues := processIssues(issues)
	saveCsv(processedIssues)
}

func loadIssues() []i.Issue {
	bytes, err := os.ReadFile("../gh/issues.json")
	if err != nil {
		panic(err)
	}
	var issues []i.Issue
	err = json.Unmarshal(bytes, &issues)
	if err != nil {
		panic(err)
	}
	return issues
}

func processIssues(issues []i.Issue) []ProcessedIssue {
	processedIssues := make([]ProcessedIssue, 0)
	for idx, issue := range issues {
		fmt.Printf("Processing issue (%d): %d\n", idx+1, issue.Number)
		labels := make([]string, 0)
		for _, label := range issue.Labels {
			labels = append(labels, label.Name)
		}
		providerVersion := getProviderVersion()
		terraformVersion := getTerraformVersion()
		processed := ProcessedIssue{
			ID:               issue.Number,
			URL:              issue.HtmlUrl,
			Title:            issue.Title,
			ProviderVersion:  providerVersion,
			TerraformVersion: terraformVersion,
			IsBug:            slices.Contains(labels, "bug"),
			IsFeatureRequest: slices.Contains(labels, "feature-request"),
			CommentsCount:    issue.Comments,
			ReactionsCount:   issue.Reactions.TotalCount,
			Labels:           labels,
		}
		processedIssues = append(processedIssues, processed)
	}
	return processedIssues
}

func getProviderVersion() string {
	return "NONE"
}

func getTerraformVersion() string {
	return "NONE"
}

func saveCsv(issues []ProcessedIssue) {
	file, err := os.Create("issues.csv")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	w := csv.NewWriter(file)
	w.Comma = ';'

	var data [][]string
	for _, issue := range issues {
		row := []string{
			strconv.Itoa(issue.ID),
			issue.URL,
			issue.Title,
			issue.ProviderVersion,
			issue.TerraformVersion,
			strconv.FormatBool(issue.IsBug),
			strconv.FormatBool(issue.IsFeatureRequest),
			strconv.Itoa(issue.CommentsCount),
			strconv.Itoa(issue.ReactionsCount),
			strings.Join(issue.Labels, "|"),
		}
		data = append(data, row)
	}
	w.WriteAll(data)
}

type ProcessedIssue struct {
	ID               int
	URL              string
	Title            string
	ProviderVersion  string
	TerraformVersion string
	IsBug            bool
	IsFeatureRequest bool
	CommentsCount    int
	ReactionsCount   int
	Labels           []string
}
