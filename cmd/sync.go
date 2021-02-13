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
	es "kib-sync/client"
	"log"
	"net/http"
	"os"
	"strings"

	"kib-sync/model"

	"github.com/spf13/cobra"
)

const (
	monitorQueryPath string = "_opendistro/_alerting/monitors/_search"
	searchQueryPath  string = "_search"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Fetch all configured monitors from Kibana cluster",
	Long:  `The configured monitors are fetched from given kiban cluster and stored in json format in the config folder`,
	Run:   SyncMonitors(NewQueryHandler()),
}

type QueryHandler func(path string, body []byte) (map[string]interface{}, error)

func NewQueryHandler() QueryHandler {
	return func(path string, body []byte) (map[string]interface{}, error) {
		c := es.NewClient(getValue(URL), getValue(UserName), getValue(Password))

		// query all monitors
		res, err := c.Do(path, http.MethodGet, body)

		if err != nil {
			ErrorLog.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode >= 300 {
			ErrorLog.Fatalf("Got response with status code: %d", res.StatusCode)
		}
		var r map[string]interface{}

		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			ErrorLog.Fatalf("Error parsing the response body: %s", err)
		}

		// Print the response status, number of results, and request duration.
		log.Printf(
			"[%d] %d results; took: %dms",
			res.StatusCode,
			int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
			int(r["took"].(float64)),
		)
		return r, nil
	}
}

// SyncMonitors download all the configured monitors form Kibana cluster
func SyncMonitors(handler QueryHandler) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		// crete parent config directory
		createDir("config")

		// Query all the config of the following kiban objects
		objects := []string{"monitor", "dashboard", "search", "destination"}
		for _, obj := range objects {
			query, err := json.Marshal(model.QueryRequest{Size: 1000, Query: model.Query{model.Bool{model.Must{model.Exists{obj}}}}})

			if err != nil {
				ErrorLog.Fatal("failed to marshal the query request")
			}
			var path string

			if obj == "monitor" {
				path = monitorQueryPath
			} else {
				path = searchQueryPath
			}

			// Run the query
			r, err := handler(path, query)

			if err != nil {
				ErrorLog.Fatal(err)
			}

			// Create folder if not exist for the configuration files for the given condfig
			createDir(fmt.Sprintf("%s/%s", "config", obj))

			// Print the ID and document source for each hit.
			counter := 0
			var ids []string
			for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
				id := hit.(map[string]interface{})["_id"]
				ids = append(ids, id.(string))
				InfoLog.Printf("successfully fetched %s id= %s", obj, id)
				file, err := os.Create(fmt.Sprintf("%s/%s/%s.json", "config", obj, id))
				if err != nil {
					ErrorLog.Println(err)
					continue
				}

				b, err := json.MarshalIndent(hit, "", "\t")
				if err != nil {
					ErrorLog.Printf("failed to marshall the resposne: %v", err.Error())
				}
				file.WriteString(string(b))
				file.Close()
				counter++
			}

			InfoLog.Printf("all of the %d kiban monitor configs successfully synched", counter)
			// remove the redundant configs
			fileInfos, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", "config", obj))
			if err != nil {
				log.Fatal(err)
			}

			// Remove redundant files
			for _, fileInfo := range fileInfos {
				id := strings.Split(fileInfo.Name(), ".")[0]
				if !find(ids, id) {
					WarnLog.Printf("removing kiban config with id: %s", fileInfo.Name())
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
