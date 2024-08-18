package main

import "flag"

func main() {
	var folder string
	var email string

	flag.StringVar(&folder, "add", "", "Add a folder to scan for Git repositories")
	flag.StringVar(&email, "email", "your@email.com", "Your email to scan")
	flag.Parse()

	if folder != "" {
		scan(folder)
		return
	}

	graph(email)
}
