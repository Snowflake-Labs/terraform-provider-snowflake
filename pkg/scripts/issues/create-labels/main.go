package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

var labels = []string{
	"category:resource",
	"category:data_source",
	"category:import",
	"category:sdk",
	"category:identifiers",
	"category:provider_config",
	"category:grants",
	"category:other",
	"resource:account",
	"resource:account_parameter",
	"resource:account_password_policy",
	"resource:alert",
	"resource:api_integration",
	"resource:database",
	"resource:database_role",
	"resource:dynamic_table",
	"resource:email_notification_integration",
	"resource:external_function",
	"resource:external_oauth_integration",
	"resource:external_table",
	"resource:failover_group",
	"resource:file_format",
	"resource:function",
	"resource:grant_account_role",
	"resource:grant_database_role",
	"resource:grant_ownership",
	"resource:grant_privileges_to_account_role",
	"resource:grant_privileges_to_database_role",
	"resource:grant_privileges_to_share",
	"resource:managed_account",
	"resource:masking_policy",
	"resource:materialized_view",
	"resource:network_policy",
	"resource:network_policy_attachment",
	"resource:notification_integration",
	"resource:oauth_integration",
	"resource:object_parameter",
	"resource:password_policy",
	"resource:pipe",
	"resource:procedure",
	"resource:resource_monitor",
	"resource:role",
	"resource:row_access_policy",
	"resource:saml_integration",
	"resource:schema",
	"resource:scim_integration",
	"resource:sequence",
	"resource:session_parameter",
	"resource:share",
	"resource:stage",
	"resource:storage_integration",
	"resource:stream",
	"resource:table",
	"resource:table_column_masking_policy_application",
	"resource:table_constraint",
	"resource:tag",
	"resource:tag_association",
	"resource:tag_masking_policy_association",
	"resource:task",
	"resource:unsafe_execute",
	"resource:user",
	"resource:user_password_policy_attachment",
	"resource:user_public_keys",
	"resource:view",
	"resource:warehouse",
	"data_source:accounts",
	"data_source:alerts",
	"data_source:current_account",
	"data_source:current_role",
	"data_source:database",
	"data_source:database_roles",
	"data_source:databases",
	"data_source:dynamic_tables",
	"data_source:external_functions",
	"data_source:external_tables",
	"data_source:failover_groups",
	"data_source:file_formats",
	"data_source:functions",
	"data_source:grants",
	"data_source:masking_policies",
	"data_source:materialized_views",
	"data_source:parameters",
	"data_source:pipes",
	"data_source:procedures",
	"data_source:resource_monitors",
	"data_source:roles",
	"data_source:row_access_policies",
	"data_source:schemas",
	"data_source:sequences",
	"data_source:shares",
	"data_source:stages",
	"data_source:storage_integrations",
	"data_source:streams",
	"data_source:system_generate_scim_access_token",
	"data_source:system_get_aws_sns_iam_policy",
	"data_source:system_get_privatelink_config",
	"data_source:system_get_snowflake_platform_info",
	"data_source:tables",
	"data_source:tasks",
	"data_source:users",
	"data_source:views",
	"data_source:warehouses",
}

func main() {
	accessToken := getAccessToken()
	repoLabels := loadRepoLabels(accessToken)
	jsonRepoLabels, _ := json.MarshalIndent(repoLabels, "", "\t")
	log.Println(string(jsonRepoLabels))
	successful, failed := createLabelsIfNotPresent(accessToken, repoLabels, labels)
	fmt.Printf("\nSuccessfully created labels:\n")
	for _, label := range successful {
		fmt.Println(label)
	}
	fmt.Printf("\nUnsuccessful label creation:\n")
	for _, label := range failed {
		fmt.Println(label)
	}
}

type ReadLabel struct {
	ID          int    `json:"id"`
	NodeId      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

func loadRepoLabels(accessToken string) []ReadLabel {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/Snowflake-Labs/terraform-provider-snowflake/labels", nil)
	if err != nil {
		panic("failed to create list labels request: " + err.Error())
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic("failed to retrieve repository labels: " + err.Error())
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("failed to read list labels response body: " + err.Error())
	}

	var readLabels []ReadLabel
	err = json.Unmarshal(bodyBytes, &readLabels)
	if err != nil {
		panic("failed to unmarshal read labels: " + err.Error())
	}

	return readLabels
}

type CreateLabelRequestBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func createLabelsIfNotPresent(accessToken string, repoLabels []ReadLabel, labels []string) (successful []string, failed []string) {
	repoLabelNames := make([]string, len(repoLabels))
	for i, label := range repoLabels {
		repoLabelNames[i] = label.Name
	}

	for _, label := range labels {
		if slices.Contains(repoLabelNames, label) {
			continue
		}

		time.Sleep(3 * time.Second)
		log.Println("Processing:", label)

		var requestBody []byte
		var err error
		parts := strings.Split(label, ":")
		labelType := parts[0]
		labelValue := parts[1]

		switch labelType {
		// Categories will be created by hand
		case "resource", "data_source":
			requestBody, err = json.Marshal(&CreateLabelRequestBody{
				Name:        label,
				Description: fmt.Sprintf("Issue connected to the snowflake_%s resource", labelValue),
			})
		default:
			log.Println("Unknown label type:", labelType)
			continue
		}

		if err != nil {
			log.Println("Failed to marshal create label request body:", err)
			failed = append(failed, label)
			continue
		}

		req, err := http.NewRequest(http.MethodPost, "https://api.github.com/repos/Snowflake-Labs/terraform-provider-snowflake/labels", bytes.NewReader(requestBody))
		if err != nil {
			log.Println("failed to create label request:", err)
			failed = append(failed, label)
			continue
		}
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println("failed to create a new label: ", label, err)
			failed = append(failed, label)
			continue
		}

		if resp.StatusCode != http.StatusCreated {
			log.Println("incorrect status code, expected 201, and got:", resp.StatusCode)
			failed = append(failed, label)
			continue
		}

		successful = append(successful, label)
	}

	return
}

func getAccessToken() string {
	token := os.Getenv("SF_TF_SCRIPT_GH_ACCESS_TOKEN")
	if token == "" {
		panic(errors.New("GitHub access token missing"))
	}
	return token
}
