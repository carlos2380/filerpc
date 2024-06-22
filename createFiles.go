package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	basePath := "core"
	content := `{
    "reward": "100000"
}`

	if err := os.MkdirAll(basePath, os.ModePerm); err != nil {
		fmt.Printf("Error creating directory %s: %v\n", basePath, err)
		return
	}

	//Number of threads
	for i := 0; i <= 8; i++ {
		//Number of files per thread
		for j := 0; j <= 1000; j++ {
			fileName := fmt.Sprintf("1.%d.%d.json", i, j)
			filePath := filepath.Join(basePath, fileName)
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				fmt.Printf("Error writing file %s: %v\n", filePath, err)
			}
		}
	}
	fmt.Println("All files created successfully.")
}
