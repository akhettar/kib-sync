/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	es "github.com/akhettar/odfe-kibana-sync/client"

	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

// creates cobra command that represents the push command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create all kiban objects (monitors, dashbaor, etc) present in the config folder",
	Long:  `create monitors, destinations, dashboards, email accoutns, group emails, search in the given kibana cluster`,
	Run:   createKibanaConfig(newCreateHandler()),
}

// CreateHandler function implementation of the Create Kibana config
type CreateHandler func(url string, body []byte) (*http.Response, error)

func newCreateHandler() CreateHandler {
	return func(path string, body []byte) (*http.Response, error) {
		// create http client
		client := es.NewClient(getValue(URL), getValue(UserName), getValue(Password))

		// invoke create
		return client.Do(path, http.MethodPost, body)
	}
}

func createKibanaConfig(handler CreateHandler) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		// crete parent config directory
		log.Println("invoking create kibana config")
		configs, err := ioutil.ReadDir(getValue(WorkDir))
		if err != nil {
			log.Fatal(err)
		}

		for _, config := range configs {
			fileInfos, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", getValue(WorkDir), config.Name()))
			if err != nil {
				ErrorLog.Printf("failed to read dir: %s", config.Name())
				continue
			}

			// Process all the kibana monitor config and push the update to Kibana
			for _, fileInfo := range fileInfos {

				log.Printf("creating kibana %s: %s to kiban cluster", config.Name(), fileInfo.Name())
				filename := fmt.Sprintf("%s/%s/%s", getValue(WorkDir), config.Name(), fileInfo.Name())
				bytes, err := ioutil.ReadFile(filename)
				if err != nil {
					ErrorLog.Printf("failed to read the content of the file %s", filename)
					continue
				}

				var m map[string]interface{}
				if err := json.NewDecoder(strings.NewReader(string(bytes))).Decode(&m); err != nil {
					ErrorLog.Printf("failed to decode the config file: %s", filename)
					continue
				}
				index := m["_index"].(string)
				documentID := m["_id"].(string)
				source := m["_source"].(map[string]interface{})

				log.Printf("pushing kibana config for document Id: %s", documentID)

				body, err := requestBody(source, config.Name())
				if err != nil {
					ErrorLog.Printf("failed to decode the request body for config: %s", filename)
					continue
				}
				// create kibana object
				res, err := handler(path(config.Name(), index, documentID), body)

				if err != nil {
					ErrorLog.Printf("Error getting response: %s", err.Error())
					continue
				}
				defer res.Body.Close()

				if res.StatusCode >= 300 {
					ErrorLog.Printf("[%d] Error indexing document", res.StatusCode)
				} else {
					// Deserialize the response into a map.
					var r map[string]interface{}
					if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
						ErrorLog.Printf("Error parsing the response body: %s", err)
					} else {
						// Print the response status and indexed document version.
						InfoLog.Printf("[%d] %s; version=%d; config=%s", res.StatusCode, r["result"], int(r["_version"].(float64)), strings.SplitAfter(fileInfo.Name(), ".")[0])
					}
				}
			}
		}
	}
}

func requestBody(content map[string]interface{}, docType string) ([]byte, error) {
	if docType == "destination" {
		content = content["destination"].(map[string]interface{})
	} else if docType == "email_group" {
		content = content["email_group"].(map[string]interface{})
	}
	body, err := json.Marshal(&content)
	if err != nil {
		ErrorLog.Println("failed to marshall the monitor doc")
		return nil, err
	}
	return body, nil
}

func path(docType string, index, documentID string) string {
	switch docType {
	case "monitor":
		return "_opendistro/_alerting/monitors"
	case "destination":
		return "_opendistro/_alerting/destinations"
	case "email_account":
		return "_opendistro/_alerting/destinations/email_accounts"
	case "email_group":
		return "_opendistro/_alerting/destinations/email_groups"
	default:
		return fmt.Sprintf("%s/_create/%s", index, documentID)
	}
}
