# Day 01 - Go Setup, Data Types and Your First CLI Tool

## What are we doing today?

Today we will:
1. Install Go on your machine
2. Understand how a Go project is structured
3. Learn about data types and variables
4. Build a port checker CLI tool

By the end of this episode, you will have a working CLI tool that checks if a port is open or closed on your machine. Pretty useful for DevOps, right?

---

## 1. Installing Go

Go to [https://go.dev/dl/](https://go.dev/dl/) and download the latest version for your OS.

After installing, open your terminal and run:

```bash
go version
```

You should see something like:

```
go version go1.26 darwin/arm64
```

Now run this to see all the Go environment variables:

```bash
go env
```

The important ones to know:
- `GOPATH` - where Go stores downloaded packages (usually `~/go`)
- `GOROOT` - where Go itself is installed
- `GOBIN` - where Go puts compiled binaries

But don't worry too much about these. Go modules handle most of the stuff for us now.

---

## 2. Creating Your First Go Project

```bash
mkdir port-checker
cd port-checker
go mod init port-checker
```

`go mod init` creates a `go.mod` file. This file tracks your project name and all the dependencies you use. Think of it like `package.json` in Node or `requirements.txt` in Python.

Now create a file called `main.go`:

```bash
touch main.go
```

Every Go program starts from the `main` package and the `main` function. This is the entry point - just like how every movie has a first scene.

---

## 3. Hello World - The Tradition

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello DevOps Engineers!")
}
```

Run it:

```bash
go run main.go
```

What is happening here:
- `package main` - tells Go this is the main package (entry point)
- `import "fmt"` - we are importing the `fmt` package for printing stuff
- `func main()` - this is the main function, Go starts running from here
- `fmt.Println` - prints a line to the terminal

---

## 4. Data Types

Go is a statically typed language. Matlab hai bhai, you have to tell Go what type of data you are storing.

### Basic types:

```go
var name string = "Abhishek"    // text
var age int = 25                // whole number
var salary float64 = 50000.50   // decimal number
var isDevOps bool = true        // true or false
```

### Short declaration (the Go way):

Most of the time, we use `:=` instead of `var`. Go is smart enough to figure out the type.

```go
name := "Abhishek"
age := 25
salary := 50000.50
isDevOps := true
```

`:=` only works inside functions. For package level variables, you still need `var`.

### Constants:

```go
const maxRetries = 3
const appName = "port-checker"
```

Constants cannot be changed once set. Use them for values that should never change.

---

## 5. The fmt Package - Your Best Friend

```go
name := "Abhishek"
age := 25

// Simple print
fmt.Println("Hello", name)

// Formatted print (like printf in C)
fmt.Printf("Name: %s, Age: %d\n", name, age)

// Store formatted string in a variable
msg := fmt.Sprintf("Server %s is running on port %d", name, 8080)
fmt.Println(msg)
```

Common format verbs:
- `%s` - string
- `%d` - integer
- `%f` - float
- `%v` - any value (Go figures it out)
- `%t` - boolean
- `\n` - new line

---

## 6. Taking User Input

### Method 1: Using os.Args (command line arguments)

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Please provide your name")
        return
    }
    name := os.Args[1]
    fmt.Printf("Hello %s!\n", name)
}
```

Run it like:
```bash
go run main.go Abhishek
```

`os.Args` is a slice (like an array) that holds all the command line arguments. `os.Args[0]` is always the program name itself, so the actual arguments start from `os.Args[1]`.

### Method 2: Using fmt.Scan

```go
var name string
fmt.Print("Enter your name: ")
fmt.Scan(&name)
fmt.Printf("Hello %s!\n", name)
```

The `&` before `name` means we are passing the memory address of the variable. Don't worry too much about this now - we will cover pointers later.

---

## 7. Project: Port Checker CLI Tool

Now let's build something real! We will create a tool that checks if a given port is open or closed on localhost.

This is useful when you want to quickly check:
- Is my Nginx running on port 80?
- Is the database running on port 5432?
- Is my app listening on port 3000?

### How does it work?

We try to make a TCP connection to `localhost:<port>`. If the connection is successful, the port is open. If it fails, the port is closed. Simple!

### The Code

Check the `main.go` file in this folder for the complete code.

### How to run

```bash
# Check a single port
go run main.go 8080

# Build and run
go build -o port-checker
./port-checker 8080
./port-checker 3000
./port-checker 22
```

### Sample output

```
Checking port 8080 on localhost...
Port 8080 is CLOSED
```

```
Checking port 22 on localhost...
Port 22 is OPEN
```

---

## 8. Build Your Binary

One of the best things about Go - you can compile your code into a single binary. No runtime needed, no dependencies. Just one file.

```bash
go build -o port-checker
```

This creates a binary called `port-checker`. You can copy this to any machine with the same OS and it will just work.

Want to build for Linux from your Mac?

```bash
GOOS=linux GOARCH=amd64 go build -o port-checker-linux
```

This is called cross-compilation. Very handy for DevOps - build on your laptop, deploy to your Linux servers.

---

## Mini Challenges

Try these on your own:

1. **Easy**: Modify the tool to accept a hostname along with port (e.g., `./port-checker google.com 443`)
2. **Medium**: Check multiple ports at once (e.g., `./port-checker 80 443 8080 3000`)
3. **Hard**: Add a timeout flag so user can control how long to wait before saying port is closed

We will cover some of these patterns in the upcoming episodes!

---

## What We Learned Today

- How to install Go and create a new project
- Data types: string, int, float64, bool
- Variables: `var` vs `:=` vs `const`
- fmt package: Println, Printf, Sprintf
- Taking input from command line using `os.Args`
- Building a real CLI tool that checks ports
- Cross-compiling Go binaries

---

Next up: **Day 02 - Conditionals, Loops and Functions** where we will build a log file analyzer!
