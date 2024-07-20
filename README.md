# Commentary

Commentary is a highly opinionated CLI tool for processing and formatting comments in Go source files. It enforces a consistent style for comments based on the export status of the entities they document.

## Features

- **Ignore comment blocks**: Ignores comments that begin with `/*` and end with `*/`.
- **Capitalization enforcement**:
  - Capitalizes comments for exported functions, types, structs, and interfaces.
  - Lowercases comments for non-exported functions, types, structs, and interfaces.
- **Internal comment styling**:
  - Lowercases comments within the body of functions, types, structs, and interfaces, whether they are exported or internal.

## Installation

You can install the `commentary` tool using `go install`:

```sh
go install github.com/johnmikee/commentary/cmd/commentary@latest
```

## Usage

Navigate to the directory containing the Go files you want to process and run:

```sh
commentary -dir path/to/your/go/files [-write]
```

- `-dir`: Directory to scan for `.go` files (default is current directory).
- `-write`: Write changes to files if set to `true` (default is `false`).

## Example

```sh
# Scan the current directory and print changes without writing to files
commentary

# Scan the specified directory and write changes to files
commentary -dir /path/to/go/files -write
```

## Opinionated Rules

Commentary enforces the following rules:

1. **Comment Blocks**: Comments beginning with `/*` and ending with `*/` are ignored.
2. **Exported Entities**: Comments for exported functions, types, structs, and interfaces must start with a capital letter.
3. **Non-Exported Entities**: Comments for non-exported functions, types, structs, and interfaces must start with a lowercase letter.
4. **Internal Comments**: Comments within the body of functions, types, structs, and interfaces must start with a lowercase letter, regardless of export status.

## Contribution

Feel free to open issues or submit pull requests if you have suggestions or improvements. This tool is designed to enforce a specific style, and contributions that enhance its ability to do so are welcome.
