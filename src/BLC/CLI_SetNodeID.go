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
		fmt.Println("OS X detected.")
	case "linux":
		fmt.Println("Linux detected.")
	case "windows":
		fmt.Println("Windows detected.")
	default:
		fmt.Printf("Unsupported OS: %s\n", runtime.GOOS)
		os.Exit(1)
	}

	fmt.Printf("node id is %s\n", nodeID)
	err := os.Setenv("NODE_ID", nodeID)
	if err != nil {
		log.Fatalf("set env failed! %v\n", err)
	}
}
