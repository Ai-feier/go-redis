package config

import (
	"bufio"
	"fmt"
	"go-redis/lib/logger"
	"io"
	"os"
	"reflect"
	"strconv"
	"strings"
)

// ServerProperties defines global config properties
type ServerProperties struct {
	Bind           string `cfg:"bind"`
	Port           int    `cfg:"port"`
	AppendOnly     bool   `cfg:"appendOnly"`
	AppendFilename string `cfg:"appendFilename"`
	MaxClients     int    `cfg:"maxclients"`
	RequirePass    string `cfg:"requirepass"`
	Databases      int    `cfg:"databases"`

	Peers []string `cfg:"peers"`
	Self  string   `cfg:"self"`
}

// Properties holds global config properties
var Properties *ServerProperties

func init() {
	// default config
	Properties = &ServerProperties{
		Bind:       "127.0.0.1",
		Port:       6379,
		AppendOnly: false,
	}
}

func parse(src io.Reader) *ServerProperties {
	config := &ServerProperties{}

	// read config file
	rawMap := make(map[string]string)
	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 && line[0] == '#' {
			continue
		}
		pivot := strings.IndexAny(line, " ")
		if pivot > 0 && pivot < len(line)-1 { // separator found
			key := line[0:pivot]
			value := strings.Trim(line[pivot+1:], " ")
			rawMap[strings.ToLower(key)] = value
		}
	}
	if err := scanner.Err(); err != nil {
		logger.Fatal(err)
	}

	// parse format
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)
	n := t.Elem().NumField()
	for i := 0; i < n; i++ {
		field := t.Elem().Field(i)
		fieldVal := v.Elem().Field(i)
		key, ok := field.Tag.Lookup("cfg")
		if !ok {
			key = field.Name
		}
		value, ok := rawMap[strings.ToLower(key)]
		if ok {
			// fill config
			switch field.Type.Kind() {
			case reflect.String:
				fieldVal.SetString(value)
			case reflect.Int:
				intValue, err := strconv.ParseInt(value, 10, 64)
				if err == nil {
					fieldVal.SetInt(intValue)
				}
			case reflect.Bool:
				boolValue := "yes" == value
				fieldVal.SetBool(boolValue)
			case reflect.Slice:
				if field.Type.Elem().Kind() == reflect.String {
					slice := strings.Split(value, ",")
					fieldVal.Set(reflect.ValueOf(slice))
				}
			}
		}
	}
	return config
}

// SetupConfig read config file and store properties into Properties
func SetupConfig(configFilename string) {
	file, err := os.Open(configFilename)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	Properties = parse(file)
}

// String returns a formatted string representation of ServerProperties.
func (s *ServerProperties) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Bind: %s\n", s.Bind))
	sb.WriteString(fmt.Sprintf("Port: %d\n", s.Port))
	sb.WriteString(fmt.Sprintf("AppendOnly: %v\n", s.AppendOnly))
	sb.WriteString(fmt.Sprintf("AppendFilename: %s\n", s.AppendFilename))
	sb.WriteString(fmt.Sprintf("MaxClients: %d\n", s.MaxClients))
	sb.WriteString(fmt.Sprintf("RequirePass: %s\n", s.RequirePass))
	sb.WriteString(fmt.Sprintf("Databases: %d\n", s.Databases))

	sb.WriteString("Peers: [")
	for i, peer := range s.Peers {
		sb.WriteString(peer)
		if i < len(s.Peers)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString("]\n")

	sb.WriteString(fmt.Sprintf("Self: %s\n", s.Self))

	return sb.String()
}