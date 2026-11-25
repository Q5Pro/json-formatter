// JSON Formatter - Pretty-print, minify, validate and query JSON
// Made in discord.gg/codeshare · astra-dev.com.tr
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readInput(file string) ([]byte, error) {
	if file != "" {
		return os.ReadFile(file)
	}
	var sb strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		sb.WriteString(scanner.Text())
		sb.WriteString("\n")
	}
	return []byte(sb.String()), scanner.Err()
}

func prettyPrint(data []byte, indent string) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}
	out, err := json.MarshalIndent(v, "", indent)
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func minify(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}
	out, err := json.Marshal(v)
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func validate(data []byte) {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Valid JSON ✓")
	printStats(v)
}

func printStats(v interface{}) {
	switch val := v.(type) {
	case map[string]interface{}:
		fmt.Printf("Type: object  Keys: %d\n", len(val))
	case []interface{}:
		fmt.Printf("Type: array  Length: %d\n", len(val))
	case string:
		fmt.Printf("Type: string  Length: %d\n", len(val))
	case float64:
		fmt.Printf("Type: number  Value: %v\n", val)
	case bool:
		fmt.Printf("Type: boolean  Value: %v\n", val)
	case nil:
		fmt.Println("Type: null")
	}
}

func queryPath(data []byte, path string) error {
	var root interface{}
	if err := json.Unmarshal(data, &root); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}

	parts := strings.Split(strings.TrimPrefix(path, "."), ".")
	current := root

	for _, part := range parts {
		if part == "" {
			continue
		}
		switch c := current.(type) {
		case map[string]interface{}:
			val, ok := c[part]
			if !ok {
				return fmt.Errorf("key '%s' not found", part)
			}
			current = val
		case []interface{}:
			idx, err := strconv.Atoi(part)
			if err != nil {
				return fmt.Errorf("array index must be integer, got '%s'", part)
			}
			if idx < 0 || idx >= len(c) {
				return fmt.Errorf("index %d out of range (len=%d)", idx, len(c))
			}
			current = c[idx]
		default:
			return fmt.Errorf("cannot navigate into %T at '%s'", current, part)
		}
	}

	out, err := json.MarshalIndent(current, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func keys(data []byte) error {
	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}
	switch obj := v.(type) {
	case map[string]interface{}:
		for k := range obj {
			fmt.Println(k)
		}
	case []interface{}:
		fmt.Printf("Array with %d elements (indices 0-%d)\n", len(obj), len(obj)-1)
	default:
		fmt.Printf("Not an object or array: %T\n", v)
	}
	return nil
}

func usage() {
	fmt.Println(`JSON Formatter & Query Tool
Usage:
  json-formatter <command> [file] [options]
  cat file.json | json-formatter <command>

Commands:
  fmt [--indent N]   Pretty-print JSON (default indent: 2)
  min                Minify JSON
  validate           Validate JSON and show stats
  keys               List top-level keys
  get <path>         Extract value at path (e.g. .user.name or .items.0)

Examples:
  json-formatter fmt data.json
  json-formatter fmt data.json --indent 4
  json-formatter min data.json
  json-formatter validate data.json
  json-formatter keys data.json
  json-formatter get data.json .user.email
  cat data.json | json-formatter get .items.0`)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
		return
	}

	cmd := args[0]
	rest := args[1:]

	// Parse file and flags
	file := ""
	extra := []string{}
	for _, a := range rest {
		if !strings.HasPrefix(a, "--") && !strings.HasPrefix(a, ".") && file == "" {
			// check if it looks like a file (not a path query)
			if _, err := os.Stat(a); err == nil {
				file = a
				continue
			}
		}
		extra = append(extra, a)
	}

	data, err := readInput(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	data = []byte(strings.TrimSpace(string(data)))
	if len(data) == 0 {
		fmt.Fprintln(os.Stderr, "Error: empty input")
		os.Exit(1)
	}

	switch cmd {
	case "fmt", "format":
		indent := "  "
		for i, a := range extra {
			if a == "--indent" && i+1 < len(extra) {
				n, err := strconv.Atoi(extra[i+1])
				if err == nil && n > 0 {
					indent = strings.Repeat(" ", n)
				}
			}
		}
		if err := prettyPrint(data, indent); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case "min", "minify":
		if err := minify(data); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case "validate", "valid":
		validate(data)

	case "keys":
		if err := keys(data); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case "get":
		path := ""
		for _, a := range extra {
			if strings.HasPrefix(a, ".") || (!strings.HasPrefix(a, "--") && a != file) {
				path = a
				break
			}
		}
		if path == "" && len(rest) > 0 {
			path = rest[len(rest)-1]
		}
		if path == "" {
			fmt.Fprintln(os.Stderr, "Error: path required (e.g. .user.name)")
			os.Exit(1)
		}
		if err := queryPath(data, path); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	case "help", "--help", "-h":
		usage()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		usage()
		os.Exit(1)
	}
}
