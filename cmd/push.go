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
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/spf13/cobra"
)

// creates cobra command that represents the push command

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "pushes the monitor's config to Kibana cluster",
	Long:  `Read monitor configuration and pushes them kibana cluster`,
	Run:   PushMonitorsConfig(),
}

// PushMonitorsConfig pushes the monitor config to kibana cluster
func PushMonitorsConfig() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		// Create an instance of the ELk client
		client := newClient(getValue(UserName), getValue(Password), getValue(URL))

		fileInfos, err := ioutil.ReadDir("config")
		if err != nil {
			log.Fatal(err)
		}

		// Process all the kibana monitor config and push the update to Kibana
		for _, fileInfo := range fileInfos {
			log.Printf("pushing kibana config: %s to kiban cluster", strings.SplitAfter(fileInfo.Name(), ".")[0])
			bytes, _ := ioutil.ReadFile(fmt.Sprintf("config/%s", fileInfo.Name()))

			var m map[string]interface{}
			if err := json.NewDecoder(strings.NewReader(string(bytes))).Decode(&m); err != nil {
				ErrorLogger.Printf("failed to decode the monitor config: %s", err.Error())
				continue
			}
			index := m["_index"].(string)
			documentID := m["_id"].(string)
			monitorDoc := m["_source"].(map[string]interface{})
			docBytes, err := json.Marshal(&monitorDoc)
			if err != nil {
				log.Fatal("failed to marshall the monitor doc")
			}
			log.Printf("pushing monitor config for document id: %s", documentID)
			req := esapi.IndexRequest{
				Index:        index,
				DocumentID:   documentID,
				DocumentType: "_doc",
				Body:         strings.NewReader(string(docBytes)),
				Refresh:      "true",
			}
			if err != nil {
				ErrorLogger.Printf("Error getting response: %s", err.Error())
				continue
			}
			res, err := req.Do(context.Background(), client)

			if err != nil {
				ErrorLogger.Printf("Error getting response: %s", err.Error())
				continue
			}

			defer res.Body.Close()

			if res.IsError() {
				log.Printf("[%s] Error indexing document", res.Status())
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					ErrorLogger.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					InfoLogger.Printf("[%s] %s; version=%d; monitor=%s", res.Status(), r["result"], int(r["_version"].(float64)), strings.SplitAfter(fileInfo.Name(), ".")[0])
				}
			}
		}

	}
}
