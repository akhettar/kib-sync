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
	"testing"

	"github.com/spf13/cobra"
)

func TestSyncConfig(t *testing.T) {

	// 1. Create the sync command
	cmd := &cobra.Command{
		Use:   "sync",
		Short: "Fetch all configured monitors from Kibana cluster",
		Long:  `The configured monitors are fetched from given kiban cluster and stored in json format in the config folder`,
		Run:   SyncConfig(queryHandlerSuccess(), []string{"monitor"}),
	}

	// 2. Remove the config file
	filename := "13otoHcBbX-aeATowSlk.json"
	if err := os.Remove(fmt.Sprintf("config/monitor/%s", filename)); err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	// 3. Execute the command
	cmd.Execute()

	// 4. Assert monitor file created in the config folder
	if _, err := os.Stat(fmt.Sprintf("config/monitor/%s", filename)); os.IsNotExist(err) {
		t.Errorf("failed to create file successfully in the config/monitor folder")
	}

	// 5. clean up
	os.RemoveAll(fmt.Sprintf("config"))
}

func queryHandlerSuccess() QueryHandler {
	return func(path string, body []byte) (map[string]interface{}, error) {

		file, err := ioutil.ReadFile("../data/monitors.json")
		if err != nil {
			log.Fatal(err)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(file, &result); err != nil {
			ErrorLog.Printf("failed to unmarshal the result")
		}
		return result, nil
	}
}
