package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh/terminal"
	yaml "gopkg.in/yaml.v2"
)

type marvin struct {
	Config    map[string]string `yaml:"config"`
	Tasks     map[string]string `yaml:"tasks"`
	Inventory struct {
		static    []map[string]string
		StaticRaw string            `yaml:"static"`
		Dynamic   map[string]string `yaml:"dynamic"`
	} `yaml:"inventory"`
	filtered []map[string]string
}

func newMarvin(config, query, task, args string) marvin {
	m := marvin{
		Config: make(map[string]string),
	}
	m.Config["del"] = " "
	err := yaml.Unmarshal([]byte(config), &m)
	if err != nil {
		speak(err.Error(), false)
		speak("Unable to read in marvin.yml", true)
	}
	// stdin to static inventory
	if !terminal.IsTerminal(0) {
		stdin, _ := ioutil.ReadAll(os.Stdin)
		m.Inventory.StaticRaw += "\n" + string(stdin)
	}

	m.rawToInventory(m.Inventory.StaticRaw)
	m.filter(query)
	m.task(task, args)
	return m
}

func (m *marvin) rawToInventory(raw string) []map[string]string {
	inventory := make([]map[string]string, 0)
	// convert strings to actual inventory
	for _, iRecord := range strings.Split(strings.TrimSpace(m.Inventory.StaticRaw), "\n") {
		for _, row := range strings.Split(iRecord, "\n") {
			i := make(map[string]string)
			if row == "" {
				continue
			}
			for id, kvs := range strings.Split(row, m.Config["del"]) {
				kv := strings.Split(kvs, ":")
				if len(kv) == 1 {
					i["_id"] = kv[0]
				} else {
					i[kv[0]] = strings.Join(kv[1:], m.Config["del"])
					// always set an _id
					if id == 0 {
						i["_id"] = i[kv[0]]
					}
				}
			}

			// add the inventory
			inventory = append(inventory, i)
		}
	}
	m.Inventory.static = inventory
	return m.Inventory.static
}

func (m *marvin) filter(queryString string) []map[string]string {
	var sKey, sValue string
	// kiss for now, replace * with .* and then regex match
	filtered := make([]map[string]string, 0)
	// figure out what we are searching for.
	q := strings.Split(queryString, ":")
	if len(q) == 2 {
		sKey, sValue = q[0], q[1]
	} else {
		sKey, sValue = "_id", q[0]
	}

	// replace
	sKey = strings.Replace(sKey, "*", ".*", -1)
	sValue = strings.Replace(sValue, "*", ".*", -1)

	for _, i := range m.Inventory.static {
		for k, v := range i {
			keyMatch, _ := regexp.MatchString(sKey, k)
			valueMatch, _ := regexp.MatchString(sValue, v)
			if keyMatch && valueMatch {
				filtered = append(filtered, i)
			}
		}
	}
	m.filtered = filtered
	return filtered
}

func (m *marvin) task(taskName, args string) {
	task, exists := m.Tasks[taskName]
	var wg sync.WaitGroup
	if !exists {
		speak("Invaild task", true)
	}

	for _, inventory := range m.Inventory.static {
		wg.Add(1)
		command := task
		for k, v := range inventory {
			command = strings.Replace(command, ":"+k, v, -1)
		}
		command = strings.Replace(command, ":args", args, -1)

		go func(inventory map[string]string, command string) {
			cmd := exec.Command("sh", "-c", command)
			execOutput, execError := cmd.CombinedOutput()
			if execError != nil {
				speak(inventory["_id"]+" failed.", false)
			} else {
				speak(inventory["_id"]+" ok", false)
			}
			fmt.Println(string(execOutput))
			wg.Done()
		}(inventory, command)
	}
	wg.Wait()
}
