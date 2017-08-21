package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"

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
		Tasks:  make(map[string]string),
	}
	// set a few defaults
	m.Config["del"] = " "

	// set a shell default task
	m.Tasks["shell"] = "ssh {{ .host }} {{ .args }}"
	m.Tasks["ls"] = "echo {{ .raw }}"

	// set some default inventory
	m.Inventory.Dynamic = make(map[string]string)
	m.Inventory.Dynamic["files"] = "ls -R -1"

	err := yaml.Unmarshal([]byte(config), &m)
	if err != nil {
		speak(err.Error(), false)
		speak("marvin.yml> invalid", true)
	}
	// stdin to static inventory
	if !terminal.IsTerminal(0) {
		stdin, _ := ioutil.ReadAll(os.Stdin)
		m.Inventory.StaticRaw += "\n" + string(stdin)
	}

	m.filter(query)

	if len(m.filtered) == 0 {
		speak("inventory> none matched your criteria", true)
	}
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
					i["id"] = kv[0]
				} else {
					i[kv[0]] = strings.Join(kv[1:], m.Config["del"])
					// always set an id
					if id == 0 {
						i["id"] = i[kv[0]]
					}
				}

				// regardless, lets add a raw
				i["raw"] = row
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
		if sKey == "" {
			sKey = "id"
		}
		if sValue == "" {
			sValue = "*"
		}
	} else {
		sKey, sValue = "id", q[0]
	}

	if dynamicCmd, dynamicCmdExists := m.Inventory.Dynamic[sKey]; dynamicCmdExists {
		dInventoryOutput, dInventoryError := m.exec(dynamicCmd)
		if dInventoryError == nil {
			// if no errors, add it
			m.Inventory.StaticRaw += "\n"
			for _, row := range strings.Split(dInventoryOutput, "\n") {
				// add each record, with the inventory name
				m.Inventory.StaticRaw += sKey + ":" + row + "\n"
			}
		}
	}

	if sKey == "id" || sKey == "raw" {
		// sigh, we need to generate ALL inventory because we just don't know
		for dInventoryName, dInventoryCmd := range m.Inventory.Dynamic {
			m.Inventory.StaticRaw += "\n"
			if dInventoryOutput, dInventoryError := m.exec(dInventoryCmd); dInventoryError == nil {
				// no errors, sweet
				for _, row := range strings.Split(dInventoryOutput, "\n") {
					// add each record, with the inventory name
					m.Inventory.StaticRaw += dInventoryName + ":" + row + "\n"
				}
			}
		}
	}

	// inventory should now be set ...
	m.rawToInventory(m.Inventory.StaticRaw)

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
	rawTask, exists := m.Tasks[taskName]
	var wg sync.WaitGroup
	if !exists {
		rawTask = "sh -c '" + taskName + " " + args + "'"
	}

	// first lets do an inventory check
	rawCommands := m.inventoryCheck(rawTask, m.filtered, args)

	// sweet, moving along
	// TODO: we are just going to spin them all up(we will figure out batching later)
	for cmdID, command := range rawCommands {
		wg.Add(1)
		go func(cmdID, command string) {

			output, execError := m.exec(command)

			// hopefully only one line :fingers_crossed
			if len(strings.Split(output, "\n")) > 1 {
				output = "\n" + output
			}

			if execError != nil {
				color.Red(cmdID + "> " + output)
			} else {
				color.Green(cmdID + "> " + output)
			}
			wg.Done()
		}(cmdID, command)
	}
	wg.Wait()
}

func (m *marvin) inventoryCheck(task string, filtered []map[string]string, args string) map[string]string {
	task = strings.Replace(task, "&lt;", "<", -1)
	rawCommands := make(map[string]string, 0)
	for _, inventory := range filtered {
		// add defaults to filtered inventory
		inventory["args"] = args
		inventory["time"] = time.Now().String()
		cmd, cmdParseError := m.template(task, inventory)
		if cmdParseError != nil {
			speak(inventory["id"]+"> missing/invalid arguments", true)
		}
		rawCommands[inventory["id"]] = cmd
	}
	return rawCommands
}

func (m *marvin) template(task string, inventory map[string]string) (string, error) {
	template := template.Must(template.New("translate").Parse(task))
	b := new(bytes.Buffer)
	err := template.Execute(b, inventory)
	if err == nil {
		return b.String(), nil
	}
	return "", fmt.Errorf("id: " + inventory["id"])
}

func (m *marvin) exec(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	execOutput, execError := cmd.CombinedOutput()
	output := strings.TrimSpace(string(execOutput))
	output = strings.Trim(output, "\r\n")
	return output, execError
}