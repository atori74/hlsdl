package main

import (
	"fmt"
)

func HandleCommandLine(d *Downloader) {
	defer fmt.Println("Finish.")
	for {
		cmd := prompt("command")

		switch cmd {
		case "dl":
			url := prompt("m3u8 URL?")
			filename := prompt("filename?")
			d.Enqueue(HLSInfo{URL: url, Title: filename})
		case "info":
			i := d.Info()
			fmt.Printf("Download Path: %s\n", i.DownloadPath)
			fmt.Printf("Concurrency: %d\n", i.Concurrency)
			fmt.Printf("Waiting Tasks: %d\n", i.Waiting)
		case "help":
			fmt.Println("commands: dl, info, help, quit")
		case "quit":
			fmt.Println("Waiting for all tasks done...")
			d.Quit()
			return
		}
	}
}

func prompt(title string) string {
	var input string
	fmt.Printf("%s> ", title)
	fmt.Scanln(&input)

	return input
}
