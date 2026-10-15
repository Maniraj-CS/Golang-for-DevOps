 <div align="center">

<img src="https://readme-typing-svg.demolab.com?font=JetBrains+Mono&weight=600&size=28&duration=3000&pause=1000&color=00ADD8&center=true&vCenter=true&width=700&lines=Golang+for+DevOps;Build+Tools+That+Solve+Real+Problems;Learn+Go+Through+Practical+Projects" alt="Go for DevOps animation" />

<p>
  Practical Go notes for Linux automation, cloud tooling, APIs, and infrastructure operations.
</p>

<a href="https://go.dev/">
  <img src="https://img.shields.io/badge/Language-Go-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go" />
</a>
<a href="https://pkg.go.dev/">
  <img src="https://img.shields.io/badge/Documentation-Go_Packages-64748B?style=flat-square&logo=go&logoColor=white" alt="Go Packages" />
</a>
<a href="https://github.com/spf13/cobra">
  <img src="https://img.shields.io/badge/CLI-Cobra-475569?style=flat-square" alt="Cobra CLI" />
</a>

</div>

---

## About This Repository

This repository focuses on learning Go by building practical tools for DevOps work.

Instead of studying every Go concept, the focus is on writing command-line tools, automating system tasks, working with APIs, checking service health, and interacting with cloud-native infrastructure.

## Who Is This For?

- DevOps learners who want to move beyond Bash and Python scripts.
- SRE learners who want to build reliable automation tools.
- Developers interested in infrastructure and cloud automation.

## Learning Approach

- Learn concepts through practical examples.
- Build a working tool for each major topic.
- Practice error handling and maintainable code.
- Combine the concepts into a final DevOps project.

## Table of Contents

| Section | Topics |
|---|---|
| [Go Fundamentals](#go-fundamentals) | Setup, variables, data types, control flow, functions |
| [Core Go Concepts](#core-go-concepts) | Slices, maps, pointers, structs, interfaces, JSON |
| [System Automation](#system-automation) | Files, errors, packages, system commands |
| [Networking and APIs](#networking-and-apis) | HTTP, health checks, concurrency, REST APIs |
| [CLI Development](#cli-development) | Cobra, Viper, flags, Go binaries |
| [Cloud-Native Tools](#cloud-native-tools) | Docker SDK, Kubernetes client-go |
| [Capstone Project](#capstone-project) | Build a Kubernetes status tool |

---

## Learning Projects

<div align="center">
<table>
<tr>
<td width="50%" valign="top">

### 🟦 Go Fundamentals

Learn the basic building blocks of Go and create your first command-line tool.

**Topics covered:**

- Go installation and `go env`
- Variables, constants, and data types
- Go modules and project structure
- User input and command-line arguments
- `fmt` and `os` packages

**Project:** Build a localhost port checker.

<a href="./Day-01/">
<img src="https://img.shields.io/badge/Explore_Fundamentals-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Explore Go Fundamentals" />
</a>

</td>
<td width="50%" valign="top">

### 🔁 Control Flow and Functions

Use conditions, loops, collections, and functions to process data.

**Topics covered:**

- `if`, `else`, and `switch`
- `for` and `range`
- Arrays, slices, and maps
- Multiple return values
- Pointer basics

**Project:** Build a server log analyzer.

<a href="./Day-02/">
<img src="https://img.shields.io/badge/Explore_Control_Flow-64748B?style=for-the-badge&logo=go&logoColor=white" alt="Explore Control Flow" />
</a>

</td>
</tr>
<tr>
<td width="50%" valign="top">

### 🧩 Structs, Interfaces and JSON

Organize data and work with configuration files.

**Topics covered:**

- Structs and methods
- Interfaces
- JSON encoding and decoding
- Struct tags
- Reading configuration files

**Project:** Build a server configuration viewer supporting JSON and YAML.

<a href="./Day-03/">
<img src="https://img.shields.io/badge/Explore_Data_Handling-0D9488?style=for-the-badge&logo=go&logoColor=white" alt="Explore Data Handling" />
</a>

</td>
<td width="50%" valign="top">

### 📁 File Handling and Errors

Build tools that read, process, and write files reliably.

**Topics covered:**

- File operations with `os`
- Line-by-line reading with `bufio`
- Error wrapping and custom errors
- `errors.Is` and `errors.As`
- Packages and `init()`

**Project:** Build a log parser with severity filtering and error handling.

<a href="./Day-04/">
<img src="https://img.shields.io/badge/Explore_File_Handling-9333EA?style=for-the-badge&logo=go&logoColor=white" alt="Explore File Handling" />
</a>

</td>
</tr>
<tr>
<td width="50%" valign="top">

### 🌐 Networking and Concurrency

Check services and perform network requests efficiently.

**Topics covered:**

- HTTP requests and response codes
- Timeouts and headers
- Goroutines and channels
- `sync.WaitGroup`
- Concurrent task execution

**Project:** Build a concurrent website health checker.

<a href="./Day-05/">
<img src="https://img.shields.io/badge/Explore_Health_Checks-2563EB?style=for-the-badge&logo=go&logoColor=white" alt="Explore Health Checks" />
</a>

</td>
<td width="50%" valign="top">

### 🖥️ System Automation

Use Go to run system commands and collect machine information.

**Topics covered:**

- Running commands with `os/exec`
- Capturing standard output and errors
- Environment variables
- File paths and string processing
- `defer`, `panic`, and `recover`

**Project:** Build a system information collector that produces a JSON or table report.

<a href="./Day-06/">
<img src="https://img.shields.io/badge/Explore_System_Tools-475569?style=for-the-badge&logo=linux&logoColor=white" alt="Explore System Tools" />
</a>

</td>
</tr>
<tr>
<td width="50%" valign="top">

### 🔌 REST API Development

Connect Go tools to external services and expose your own endpoints.

**Topics covered:**

- GET, POST, PUT, and DELETE
- Authentication headers and tokens
- Nested JSON responses
- HTTP handlers and `http.ServeMux`
- Environment-based configuration

**Project:** Build a GitHub API CLI and a local `/repos` endpoint.

<a href="./Day-07/">
<img src="https://img.shields.io/badge/Explore_REST_APIs-16A34A?style=for-the-badge&logo=go&logoColor=white" alt="Explore REST APIs" />
</a>

</td>
<td width="50%" valign="top">

### ⚙️ CLI Application Development

Combine your tools into a reusable command-line application.

**Topics covered:**

- Cobra commands and subcommands
- Flags and configuration with Viper
- Persistent and local flags
- Go tool dependencies
- Building and cross-compiling binaries

**Project:** Build `devctl` with health, sysinfo, and GitHub commands.

<a href="./Day-08/">
<img src="https://img.shields.io/badge/Explore_CLI_Development-F59E0B?style=for-the-badge&logo=go&logoColor=white" alt="Explore CLI Development" />
</a>

</td>
</tr>
<tr>
<td width="50%" valign="top">

### 🐳 Docker and Kubernetes

Interact with containers and Kubernetes resources from Go.

**Topics covered:**

- Docker SDK for Go
- Containers and images
- Kubernetes `client-go`
- Pods, deployments, and services
- Contexts and timeouts

**Project:** Build a tool to inspect Docker containers and Kubernetes resources.

<a href="./Day-09/">
<img src="https://img.shields.io/badge/Explore_Cloud_Native-0EA5E9?style=for-the-badge&logo=docker&logoColor=white" alt="Explore Cloud Native" />
</a>

</td>
<td width="50%" valign="top">

### 🚀 Capstone: Kubernetes Status Tool

Combine the course concepts into a practical Kubernetes CLI tool.

**Topics covered:**

- Cobra CLI and command flags
- Kubernetes API integration
- Concurrent resource fetching
- Error handling and JSON output
- Terminal status formatting

**Project:** Build `kubectl-status` to display resources and highlight unhealthy workloads.

<a href="./Day-10/">
<img src="https://img.shields.io/badge/Explore_Capstone-DC2626?style=for-the-badge&logo=kubernetes&logoColor=white" alt="Explore Capstone Project" />
</a>

</td>
</tr>
</table>
</div>

---

## What You Will Build

By working through these topics, you will practice building:

- A port checker for localhost.
- A server log analyzer and parser.
- A configuration file viewer.
- A concurrent website health checker.
- A system information collector.
- A GitHub API command-line tool.
- A reusable `devctl` CLI.
- A Docker and Kubernetes resource inspector.
- A Kubernetes status plugin.

## Capstone Project

The final project is `kubectl-status`, a tool designed to show Kubernetes resources in a clear terminal view.

Planned features:

- Display pods, deployments, and services.
- Highlight unhealthy pods and restart issues.
- Fetch resources concurrently.
- Filter resources by namespace.
- Support table and JSON output.
- Build a standalone binary that can be installed in `PATH` and used as a kubectl plugin.

---

## Scope of This Repository

**Included:** Go fundamentals, system automation, CLI tools, HTTP APIs, concurrency, Docker, Kubernetes, and practical DevOps projects.

**Not included:** Frontend development, web frameworks as a primary focus, exhaustive coverage of every Go feature, or competitive programming.

---

<div align="center">

**Learn Go. Build Tools. Automate Infrastructure.**

<img src="https://capsule-render.vercel.app/api?type=rect&color=0:0f172a,100:164e63&height=3&section=footer" width="100%" alt="" />

</div>