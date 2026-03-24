package main

import "fmt"

// ================================
// INTERFACES - The power of Go
// ================================

// Interface: any struct that has a "GetInfo" and "CheckStatus" method is a Service
type Service interface {
	GetInfo() string
	CheckStatus() string
}

// WebServer struct
type WebServer struct {
	Name string
	Port int
	Up   bool
}

func (w WebServer) GetInfo() string {
	return fmt.Sprintf("Web: %s (port %d)", w.Name, w.Port)
}

func (w WebServer) CheckStatus() string {
	if w.Up {
		return "✅ running"
	}
	return "❌ down"
}

// Database struct
type Database struct {
	Name    string
	Port    int
	Engine  string
	Running bool
}

func (d Database) GetInfo() string {
	return fmt.Sprintf("DB: %s [%s] (port %d)", d.Name, d.Engine, d.Port)
}

func (d Database) CheckStatus() string {
	if d.Running {
		return "✅ running"
	}
	return "❌ down"
}

// Cache struct
type Cache struct {
	Name string
	Port int
	Up   bool
}

func (c Cache) GetInfo() string {
	return fmt.Sprintf("Cache: %s (port %d)", c.Name, c.Port)
}

func (c Cache) CheckStatus() string {
	if c.Up {
		return "✅ running"
	}
	return "❌ down"
}

// This function works with ANY Service - that's the power of interfaces!
// It doesn't care if it's a WebServer, Database, or Cache
func printServiceStatus(s Service) {
	fmt.Printf("  %-35s %s\n", s.GetInfo(), s.CheckStatus())
}

func main() {
	// ================================
	// 1. BASIC INTERFACE USAGE
	// ================================

	fmt.Println("=== Interface - One Function for All Types ===")

	web := WebServer{Name: "nginx", Port: 80, Up: true}
	db := Database{Name: "postgres", Port: 5432, Engine: "PostgreSQL", Running: true}
	cache := Cache{Name: "redis", Port: 6379, Up: false}

	// All three different structs, but one function handles them all!
	printServiceStatus(web)
	printServiceStatus(db)
	printServiceStatus(cache)

	// ================================
	// 2. SLICE OF INTERFACES
	// ================================

	fmt.Println("\n=== Slice of Services (mixed types) ===")

	// We can put different types in the same slice because they all satisfy Service
	services := []Service{
		WebServer{Name: "nginx", Port: 80, Up: true},
		WebServer{Name: "apache", Port: 8080, Up: true},
		Database{Name: "postgres", Port: 5432, Engine: "PostgreSQL", Running: true},
		Database{Name: "mysql", Port: 3306, Engine: "MySQL", Running: false},
		Cache{Name: "redis", Port: 6379, Up: true},
		Cache{Name: "memcached", Port: 11211, Up: false},
	}

	upCount := 0
	downCount := 0

	for _, svc := range services {
		printServiceStatus(svc)
		if svc.CheckStatus() == "✅ running" {
			upCount++
		} else {
			downCount++
		}
	}

	fmt.Printf("\nTotal: %d services | Up: %d | Down: %d\n", len(services), upCount, downCount)

	// ================================
	// 3. EMPTY INTERFACE (any)
	// ================================

	fmt.Println("\n=== Empty Interface (any) ===")

	// 'any' is the same as 'interface{}' - it can hold any type
	// Useful but use carefully - you lose type safety
	var data any

	data = "hello"
	fmt.Println("String:", data)

	data = 42
	fmt.Println("Int:", data)

	data = WebServer{Name: "test", Port: 80, Up: true}
	fmt.Println("Struct:", data)

	// Type assertion - getting the actual value back
	if ws, ok := data.(WebServer); ok {
		fmt.Println("It's a WebServer:", ws.Name)
	}

	// ================================
	// 4. TYPE SWITCH
	// ================================

	fmt.Println("\n=== Type Switch ===")

	things := []any{"hello", 42, true, WebServer{Name: "test", Port: 80, Up: true}}

	for _, thing := range things {
		switch v := thing.(type) {
		case string:
			fmt.Printf("  String: %q\n", v)
		case int:
			fmt.Printf("  Int: %d\n", v)
		case bool:
			fmt.Printf("  Bool: %t\n", v)
		case WebServer:
			fmt.Printf("  WebServer: %s on port %d\n", v.Name, v.Port)
		default:
			fmt.Printf("  Unknown type: %T\n", v)
		}
	}
}
