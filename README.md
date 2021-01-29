# Gonfig

[![Go Reference](https://pkg.go.dev/badge/github.com/milad-abbasi/gonfig.svg)](https://pkg.go.dev/github.com/milad-abbasi/gonfig)
![Build Status](https://github.com/milad-abbasi/gonfig/workflows/Build/badge.svg)
[![codecov](https://codecov.io/gh/milad-abbasi/gonfig/branch/master/graph/badge.svg?token=jv13V1BgIP)](https://codecov.io/gh/milad-abbasi/gonfig)
[![Go Report Card](https://goreportcard.com/badge/github.com/milad-abbasi/gonfig)](https://goreportcard.com/report/github.com/milad-abbasi/gonfig)
[![GitHub release](https://img.shields.io/github/release/milad-abbasi/gonfig.svg)](https://gitHub.com/milad-abbasi/gonfig/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Tag-based configuration parser which loads values from different providers into typesafe struct.

## Installation
This package needs go version 1.15+
```bash
go get -u github.com/milad-abbasi/gonfig
```

## Usage
```go
package main

import (
    "fmt"
    "net/url"
    "time"
    
    "github.com/milad-abbasi/gonfig"
)

// Unexported struct fields are ignored
type Config struct {
	Host     string
	Log      bool
	Expiry   int
	Port     uint
	Pi       float64
	Com      complex128
	Byte     byte
	Rune     rune
	Duration time.Duration
	Time     time.Time
	url      url.URL
	Redis    struct {
		Hosts []string
		Ports []int
	}
}

func main() {
	var c Config

	// Input argument must be a pointer to struct
	err := gonfig.Load().FromEnv().Into(&c)
	if err != nil {
		fmt.Println(err)
	}
}
```

## Tags
All the tags are optional.

### Config
`config` tag is used to change default key for fetching the value of field.
```go
type Config struct {
	HostName string `config:"HOST"`
}

func main() {
	var c Config

	os.Setenv("HOST", "golang.org")
	gonfig.Load().FromEnv().Into(&c)
}
```

### File related tags
`json`, `yaml` and `toml` tags are used to change default key for fetching the value from file.
```go
type Config struct {
	HostName string `json:"host" yaml:"host" toml:"host"`
}

func main() {
	var c Config

	gonfig.Load().FromFile("config.json").Into(&c)
}
```

### Default
`default` tag is used to declare a default value in case of missing value.
```go
type Config struct {
	Host url.URL `default:"golang.org"`
}
```

### Required
`required` tag is used to make sure a value is present for corresponding field.  
Fields are optional by default.
```go
type Config struct {
	Host url.URL `required:"true"`
}
```

### Ignore
`ignore` tag is used to skip populating a field.
Ignore is `false` by default.
```go
type Config struct {
	Ignored     int `ignore:"true"`
	AlsoIgnored int `config:"-"`
}
```

### Expand
`expand` tag is used to expand value from OS environment variables.  
Expand is `false` by default.
```go
type Config struct {
	Expanded int `expand:"true" default:"${ENV_VALUE}"`
}

func main() {
	var c Config

	os.Setenv("ENV_VALUE", "123")
	gonfig.Load().FromEnv().Into(&c)
	fmt.Println(c.Expanded) // 123
}
```

### Separator
`separator` tag is used to separate slice/array items.  
Default separator is a single space.
```go
type Config struct {
	List []int `separator:"," default:"1, 2, 3"`
}
```

### Format
`format` tag is used for parsing time strings.    
Default format is `time.RFC3339`.
```go
type Config struct {
	Time time.Time `format:"2006-01-02T15:04:05.999999999Z07:00"`
}
```

## Providers
Providers can be chained together and they are applied in the specified order.  
If multiple values are provided for a field, last one will get applied.

### Supported providers
- Environment variables
- files
    - .json
    - .yaml (.yml)
    - .toml
    - .env

```go
func main() {
	var c Config

	gonfig.
		Load().
		FromEnv().
		FromFile("config.json").
		FromFile("config.yaml").
		FromFile("config.toml").
		FromFile(".env").
		AddProvider(CustomProvider).
		Into(&c)
}
```

### Env Provider
Env provider will populate struct fields based on the hierarchy of struct.

```go
type Config struct {
	PrettyLog bool
	Redis     struct {
		Host string
		Port int
	}
}

func main() {
	var c Config

	gonfig.
		Load().
		FromEnv().
		Into(&c)
}
```
It will check for following keys:
- `PRETTY_LOG`
- `REDIS_HOST`
- `REDIS_PORT`

To change default settings, make an `EnvProvider` and add it to the providers list manually:
```go
type Config struct {
	PrettyLog bool
	Redis     struct {
		Host string
		Port int
	}
}

func main() {
	var c Config

	ep := gonfig.EnvProvider{
		Prefix:         "APP_", // Default to ""
		SnakeCase:      false,  // Defaults to true
		UpperCase:      false,  // Defaults to true
		FieldSeparator: "__",   // Defaults to "_"
		Source:         ".env", // Defaults to OS env vars
		Required:       true,   // Defaults to false
	}

	gonfig.
		Load().
		AddProvider(&ep).
		Into(&c)
}
```
It will check for following keys in `.env` file:
- `APP_PrettyLog`
- `APP_Redis__Host`
- `APP_Redis__Port`

### File Provider
File provider uses third party parsers for parsing files, read their documentation for more info.
- [json](https://golang.org/pkg/encoding/json)
- [yaml](https://github.com/go-yaml/yaml/tree/v3)
- [toml](https://github.com/BurntSushi/toml)
- [env](https://github.com/joho/godotenv)

### Custom Provider
You can use your own provider by implementing `Provider` interface and one or both `Unmarshaler` and `Filler` interfaces.

```go
type CustomProvider struct{}

func (cp *CustomProvider) Name() string {
	return "custom provider"
}

func (cp *CustomProvider) UnmarshalStruct(i interface{}) error {
	// UnmarshalStruct receives a struct pointer and unmarshalls values into it.
	return nil
}

func (cp *CustomProvider) Fill(in *gonfig.Input) error {
	// Fill receives struct fields and set their values.
	return nil
}

func main() {
	var c Config

	gonfig.
		Load().
		AddProvider(new(CustomProvider)).
		Into(&c)
}
```

## Supported types
Any other type except the followings, results an error
- `string`
- `bool`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `complex64`, `complex128`
- `byte`
- `rune`
- [time.Duration](https://golang.org/pkg/time/#Duration)
- [time.Time](https://golang.org/pkg/time/#Time)
- [url.URL](https://golang.org/pkg/net/url/#URL)
- `pointer`, `slice` and `array` of above types
- `nested` and `embedded` structs 

## TODO
Any contribution is appreciated :)

- [ ] Add support for map data type
    - [ ] Add support for map slice
- [ ] Add support for slice of structs
- [ ] Add support for [encoding.TextUnmarshaler](https://golang.org/pkg/encoding/#TextUnmarshaler) and [encoding.BinaryUnmarshaler](https://golang.org/pkg/encoding/#BinaryUnmarshaler)
- [ ] Add support for other providers
    - [ ] command line flags
    - [ ] [etcd](https://etcd.io)
    - [ ] [Consul](https://www.consul.io)
    - [ ] [Vault](https://www.vaultproject.io)
    - [ ] [Amazon SSM](https://docs.aws.amazon.com/systems-manager/latest/userguide/what-is-systems-manager.html)

## Documentation
Take a look at [docs](https://pkg.go.dev/github.com/milad-abbasi/gonfig) for more information.

## License
The library is released under the [MIT license](https://opensource.org/licenses/MIT).  
Checkout [LICENSE](https://github.com/milad-abbasi/gonfig/blob/master/LICENSE) file.
