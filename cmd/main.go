package main

import (
	"fmt"
	"log"

	"github.com/samarth8765/snowflake-id/snowflakeId"
)

func main() {
	node, err := snowflakeId.NewNode(1)
	if err != nil {
		log.Fatalf("Fail to create node: %v", err)
	}

	id := node.GenerateID()
	snowflakeID := snowflakeId.ID(id)

	s, err := snowflakeID.String(16)
	if err != nil {
		log.Fatal("Fail to represent string format", err)
	}

	fmt.Printf("%s", s)
}
