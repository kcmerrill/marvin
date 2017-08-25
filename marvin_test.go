package main

import (
	"fmt"
	"log"
	"testing"
)

func getSampleConfig() string {
	return `
tasks:
  mysql: mysql_command	
  template: :db :port :host
  task: ignore
inventory:
  dynamic: 
    aws: aws_command
  static: |
    db:db_name host:host_name port:port_number
    db:db_name2 host:host_name2 port:port_number2
    nothingheretosee
`
}
func TestNewMarvinStruct(t *testing.T) {
	m := newMarvin(getSampleConfig(), ".", "db:", "task", "")
	if taskCommand, taskExists := m.Tasks["mysql"]; !taskExists && taskCommand != "mysql_command" {
		log.Fatalf("Expected mysql_command, Actual: " + taskCommand)
	}

	if dInventoryCommand, dInventoryExists := m.Inventory.Dynamic["aws"]; !dInventoryExists && dInventoryCommand != "aws_command" {
		log.Fatalf("Expected aws_command, Actual: " + dInventoryCommand)
	}
}

func TestRawToInventory(t *testing.T) {
	m := newMarvin(getSampleConfig(), ".", "db:", "task", "")
	if len(m.Inventory.static) != 3 {
		log.Fatalf("Expected 3 static inventories")
	}

	if m.Inventory.static[0]["db"] != "db_name" {
		log.Fatalf("1st inventory should be db:db_name")
	}

	if m.Inventory.static[1]["db"] != "db_name2" {
		log.Fatalf("2nd inventory should be db:db_name2")
	}

	fmt.Println(m.Inventory.static[2])
	if m.Inventory.static[2]["id"] != "nothingheretosee" {
		log.Fatalf("Without a ':', id should be the key")
	}
}

func TestMarvinFilter(t *testing.T) {
	m := newMarvin(getSampleConfig(), ".", "db:*", "task", "")
	if len(m.filtered) != 2 {
		log.Fatalf("Should be 2 db matches")
	}

	m.filter("db:db_name2")
	if len(m.filtered) != 1 {
		log.Fatalf("Should only be 1 db_name2 matches")
	}

	m.filter("not*")
	if len(m.filtered) != 1 {
		log.Fatalf("Should only be 1 not*")
	}

	m.filter("*here")
	if len(m.filtered) != 1 {
		log.Fatalf("Should only be 1 *here")
	}

	m.filter("bingowashisnameo")
	if len(m.filtered) != 0 {
		log.Fatalf("bingowashisnameo should have have matched anything")
	}

	m.filter("db:bingowashisnameo")
	if len(m.filtered) != 0 {
		log.Fatalf("bingowashisnameo should have have matched anything")
	}

	m.filter("db:*name")
	if len(m.filtered) != 2 {
		log.Fatalf("Should have been 2 matches")
	}
}
