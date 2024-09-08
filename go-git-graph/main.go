package main

import "flag"

func main() {
	var folder string
	var email string
	var config bool
	var clean bool
	var list bool

	flag.StringVar(&folder, "scan", "", "Add a folder to scan for Git repositories")
	flag.StringVar(&email, "email", "", "Your email to scan")
	flag.BoolVar(&clean, "clean", false, "Clear repo histories")
	flag.BoolVar(&list, "list", false, "List all repos")
	flag.BoolVar(&config, "config", false, "Show current configs")

	flag.Parse()

	if email != "" {
		saveConfigs(email)
	}

	switch {
	case config:
		getConfigs(true)
	case clean:
		clearGraph()
	case list:
		listRepos()
	case folder != "":
		configs, err := getConfigs(false)
		if err != nil {
			panic(err)
		}
		ignore := configs["ignore"]
		scan(folder, ignore)
	default:
		configs, err := getConfigs(false)
		if err != nil {
			panic(err)
		}
		graph(configs["email"])
	}
}
