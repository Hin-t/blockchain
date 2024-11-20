package BLC

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

//设置端口号（环境变量）
func (cli *CLI) SetNodeID(nodeID string) {
	if nodeID == "" {
		fmt.Println("nodeID is empty, please set the node id...")
		os.Exit(1)
	}
	switch runtime.GOOS {
	case "darwin":
		fmt.Println("OS X.")
		err := os.Setenv("Node_ID", nodeID)
		if err != nil {
			log.Fatalf("set env failed! %v\n", err)
		}
	case "linux":
	case "windows":
		fmt.Printf("Windows. node id is %s\n", nodeID)
		err := os.Setenv("Node_ID", nodeID)
		if err != nil {
			log.Fatalf("set env failed! %v\n", err)
		}
	default:
		return
	}
}
