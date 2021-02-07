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
	"kib-sync/model"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetch all configured monitors from Kibana cluster",
	Long:  `The configured monitors are fetched from given kiban cluster and stored in json format in the config folder`,
	Run:   SyncMonitors(NewQueryHandler()),
}

// QueryHandler type
type QueryHandler func(object string) map[string]interface{}

// NewQueryHandler creates QueryHandler function
func NewQueryHandler() QueryHandler {
	return func(object string) map[string]interface{} {
		// create elk client
		client := newClient(getValue(UserName), getValue(Password), getValue(URL))

		// download all the monitors
		request, err := json.Marshal(model.SearchQuery{Size: 1000, Query: model.Query{model.Bool{model.Must{model.Exists{"monitor"}}}}})
		if err != nil {
			log.Fatal(err)
		}
		resp, err := client.Search(client.Search.WithBody(strings.NewReader(string(request))))
		if err != nil {
			ErrorLogger.Fatal(err)
		}

		var r map[string]interface{}

		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
			ErrorLogger.Fatalf("Error parsing the response body: %s", err)
		}

		// Print the response status, number of results, and request duration.
		log.Printf(
			"[%s] %d results; took: %dms",
			resp.Status(),
			int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
			int(r["took"].(float64)),
		)
		return r
	}
}

// SyncMonitors download all the configured monitors form Kibana cluster
func SyncMonitors(handler QueryHandler) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		configs := []string{"monitor", "dashboard", "search", "destination"}

		// crete parent config directory
		createDir("config")

		for _, config := range configs {

			// Perform the query
			r := handler(config)

			// Create folder if not exist for the configuration files for the given condfig
			createDir(fmt.Sprintf("%s/%s", "config", config))

			// Print the ID and document source for each hit.
			counter := 0
			var ids []string
			for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
				id := hit.(map[string]interface{})["_id"]
				ids = append(ids, id.(string))
				InfoLogger.Printf("successfully fetched %s id= %s", config, id)
				file, err := os.Create(fmt.Sprintf("%s/%s/%s.json", "config", config, id))
				if err != nil {
					ErrorLogger.Println(err)
					continue
				}
				b, _ := json.MarshalIndent(hit, "", "\t")
				file.WriteString(string(b))
				file.Close()
				counter++
			}

			InfoLogger.Printf("all of the %d kiban %s configs successfully synched", counter, config)

			// remove the redundant configs
			fileInfos, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", "config", config))
			if err != nil {
				log.Fatal(err)
			}

			// Remove redundant files
			for _, fileInfo := range fileInfos {
				id := strings.Split(fileInfo.Name(), ".")[0]
				if !find(ids, id) {
					WarningLogger.Printf("removing kiban config with id: %s", fileInfo.Name())
					os.Remove(fmt.Sprintf("%s/%s", "config", fileInfo.Name()))
				}
			}
		}

	}
}

func createDir(name string) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(name, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}
}

func find(ids []string, id string) bool {
	for _, val := range ids {
		if id == val {
			return true
		}
	}
	return false
}
