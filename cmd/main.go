package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/deqdev/nodeparse/pkg/manager"
)

func main() {
	var (
		url      string
		filename string
	)

	flag.StringVar(&url, "url", "", "Node subscription URL")
	flag.StringVar(&filename, "file", "", "Node file path")
	flag.Parse()

	nm := manager.NewNodeManager()

	if url != "" {
		if err := nm.LoadFromURL(url); err != nil {
			log.Fatal(err)
		}
	}

	if filename != "" {
		if err := nm.LoadFromFile(filename); err != nil {
			log.Fatal(err)
		}
	}

	err := nm.LoadFromFile("./conf/nodes")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("node: %v\n", nm)

	configs := nm.ExportToClash()
	fmt.Println(configs)
	// 处理配置...
}
