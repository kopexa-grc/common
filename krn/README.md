# Kopexa Resource Name (KRN)

The KRN package implements a standardized resource naming system for the Kopexa platform, following Google's resource naming design principles. It provides a consistent way to identify and reference resources across different services.

## Features

- Canonical resource naming format: `//<service-name>/<relative-resource-name>`
- Support for JSON and YAML serialization
- Database integration via `sql.Scanner` and `driver.Valuer`
- Legacy format support
- Resource ID validation
- Comprehensive test coverage

## Installation

```bash
go get github.com/kopexa-grc/common/krn
```

## Usage

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/kopexa-grc/common/krn"
)

func main() {
    // Create a new KRN
    krn, err := krn.New("//kopexa.com/frameworks/iso-27001-2022")
    if err != nil {
        panic(err)
    }

    // Get the service name
    fmt.Println(krn.ServiceName) // Output: kopexa.com

    // Get the resource path
    fmt.Println(krn.RelativeResourceName) // Output: frameworks/iso-27001-2022

    // Get the canonical string representation
    fmt.Println(krn.String()) // Output: //kopexa.com/frameworks/iso-27001-2022
}
```

### Working with Resource IDs

```go
// Get a resource ID from a collection
resourceID, err := krn.ResourceID("frameworks")
if err != nil {
    panic(err)
}
fmt.Println(resourceID) // Output: iso-27001-2022

// Create a child KRN
childKRN, err := krn.NewChildKRN(
    "//kopexa.com/frameworks/iso-27001-2022",
    "controls",
    "a.1.1",
)
if err != nil {
    panic(err)
}
fmt.Println(childKRN.String()) // Output: //kopexa.com/frameworks/iso-27001-2022/controls/a.1.1
```

### Database Integration

```go
// The KRN type implements sql.Scanner and driver.Valuer
type Resource struct {
    ID   int
    KRN  krn.KRN
    Name string
}

// Can be used directly in database operations
db.QueryRow("SELECT id, krn, name FROM resources WHERE id = ?", 1).Scan(&resource.ID, &resource.KRN, &resource.Name)
```

### JSON/YAML Support

```go
// Marshal to JSON
data, err := json.Marshal(krn)
if err != nil {
    panic(err)
}
fmt.Println(string(data)) // Output: "//kopexa.com/frameworks/iso-27001-2022"

// Unmarshal from JSON
var newKRN krn.KRN
err = json.Unmarshal(data, &newKRN)
if err != nil {
    panic(err)
}
```

## Resource ID Format

Resource IDs must follow these rules:
- 4-200 characters long
- Contains only lowercase letters, digits, dots, or hyphens
- Example: `1.1.2-tmp-configured`

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the BUSL-1.1 License - see the LICENSE file for details. 