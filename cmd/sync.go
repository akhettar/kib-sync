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
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetch all configured monitors from Kibana cluster",
	Long:  `The configured monitors are fetched from given kiban cluster and stored in json format in the config folder`,
	Run:   SyncMonitors(),
}

// SyncMonitors download all the configured monitors form Kibana cluster
func SyncMonitors() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		// create elk client
		client := newClient(getValue(UserName), getValue(Password), getValue(URL))

		// download all the monitors
		resp, err := client.Search(client.Search.WithBody(strings.NewReader("{\"size\": 10000, \"query\":{ \"bool\": {\"must\": { \"exists\": { \"field\" : \"monitor\" }}}}}")))
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

		// Print the ID and document source for each hit.
		counter := 0
		var ids []string
		for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
			source := hit.(map[string]interface{})["_source"]
			id := hit.(map[string]interface{})["_id"]
			ids = append(ids, id.(string))
			name := source.(map[string]interface{})["monitor"].(map[string]interface{})["name"].(string)
			InfoLogger.Printf("successfully fetched Monitor name= %s", name)
			file, err := os.Create(fmt.Sprintf("%s/%s.json", "config", id))
			if err != nil {
				ErrorLogger.Println(err)
				continue
			}
			b, _ := json.MarshalIndent(hit, "", "\t")
			file.WriteString(string(b))
			file.Close()
			counter++
		}

		InfoLogger.Printf("all of the %d kiban monitor configs successfully synched", counter)

		// remove the redundant configs
		fileInfos, err := ioutil.ReadDir("config")
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

func find(ids []string, id string) bool {
	for _, val := range ids {
		if id == val {
			return true
		}
	}
	return false
}
