package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

/*
This script is used to download all the public GPG keys from the Terraform registry for the Snowflake provider.
It uses Terraform Registry API to get the list of all Snowflake provider versions.
Then, it iterates over them to find all GPG keys and presents them in the terminal output at the end.

Note: Before running the script, make sure the SnowflakeRegistryUrl is set to the correct value.
*/

const SnowflakeRegistryUrl = "https://registry.terraform.io/v1/providers/Snowflake-Labs/snowflake"

type VersionsResponse struct {
	Versions []Version `json:"versions"`
}

type Version struct {
	Version string `json:"version"`
}

type DownloadResponse struct {
	SigningKeys SigningKeys `json:"signing_keys"`
}

type SigningKeys struct {
	GpgPublicKeys []GpgPublicKeys `json:"gpg_public_keys"`
}

type GpgPublicKeys struct {
	AsciiArmor string `json:"ascii_armor"`
}

func main() {
	versionsBuffer := GetAndReturnBody(fmt.Sprintf("%s/versions", SnowflakeRegistryUrl))
	var versionsResponse VersionsResponse
	if err := json.Unmarshal(versionsBuffer.Bytes(), &versionsResponse); err != nil {
		panic(err)
	}

	versionsMap := make(map[string]int)

	for _, version := range versionsResponse.Versions {
		versionUrl := fmt.Sprintf("%s/%s/download/darwin/amd64", SnowflakeRegistryUrl, version.Version)
		versionBuffer := GetAndReturnBody(versionUrl)

		var downloadResponse DownloadResponse
		if err := json.Unmarshal(versionBuffer.Bytes(), &downloadResponse); err != nil {
			panic(err)
		}

		for _, gpgKey := range downloadResponse.SigningKeys.GpgPublicKeys {
			if count, ok := versionsMap[gpgKey.AsciiArmor]; ok {
				versionsMap[gpgKey.AsciiArmor] = count + 1
			} else {
				versionsMap[gpgKey.AsciiArmor] = 1
			}
		}
	}

	for key, count := range versionsMap {
		log.Printf("\nKey count: %d\n%s", count, key)
	}
}

func GetAndReturnBody(url string) *bytes.Buffer {
	log.Printf("Calling %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		panic(err)
	}

	return buf
}
