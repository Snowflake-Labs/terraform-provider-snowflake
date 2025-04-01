package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

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
	versionsBuffer := GetAndReturnBody("https://registry.terraform.io/v1/providers/Snowflake-Labs/snowflake/versions")
	var versionsResponse VersionsResponse
	if err := json.Unmarshal(versionsBuffer.Bytes(), &versionsResponse); err != nil {
		panic(err)
	}

	versionsMap := make(map[string]int)

	for _, version := range versionsResponse.Versions {
		versionUrl := fmt.Sprintf("https://registry.terraform.io/v1/providers/Snowflake-Labs/snowflake/%s/download/darwin/amd64", version.Version)
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

	resp, err := http.Get(url)
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
