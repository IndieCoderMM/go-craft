# Git Contribution Grapher

This program is to visualize commits of local Git repositories. 

It scans a specified folder for local Git repositories and graphs contributions for a given email.

## Features

- **Scan Git Repos**: Scans the specified folder for local Git repositories and graphs contributions.
- **Graph Contributions**: Visualizes contributions for a given email.
- **List Repositories**: Displays all detected repositories.
- **Clear Histories**: Option to clear contribution histories.
- **View/Save Config**: Manage user email configuration for contribution tracking.

## Usage

1. **Build the Program**: 
```go
go build -o ggg.exe
```
   
2. **Scan Repositories**: 
```go
./ggg.exe -scan <projects> -email <email>
```
3. **List Repositories**: 
```go
./ggg.exe -list
```
   
3. **Graph Contributions**: 
```go
./ggg.exe
```

### Options:

- `-scan <folder>`: Specify a folder to scan for Git repositories.
- `-email <email>`: Set your email for contribution tracking.
- `-clean`: Clear repository contribution histories.
- `-list`: List all registered repositories.
- `-config`: Show current configuration.

### Resources

- [Tutorial](https://flaviocopes.com/go-git-contributions/)
- [Go-Git](https://github.com/go-git/go-git)
