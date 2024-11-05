package tools

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

func Snowflake() int64 {
	node, err := snowflake.NewNode(204)
	if err != nil {
		log.Fatalf("Snowflake Node cannot be initialized: %v", err)
	}
	uid := node.Generate().Int64()
	return uid
}
