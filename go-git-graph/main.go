package main

import "flag"

func main() {
	var folder string
	var email string
	var clear bool

	flag.StringVar(&folder, "add", "", "Add a folder to scan for Git repositories")
	flag.StringVar(&email, "email", "hthant00chk@gmail.com", "Your email to scan")
	flag.BoolVar(&clear, "clear", false, "Clear repo histories")
	flag.Parse()

	if clear {
		clearGraph()
	}

	if folder != "" {
		scan(folder)
		return
	}

	graph(email)
}
