package server

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
	"strings"
)

const ServerDefaultJSON = `{
"RUN_ADDRESS":"localhost:8080",
"DATABASE_URI": "postgres://postgres:mypassword@localhost:5432/yandex",
"AUTH_TIMEOUT": 100
}`

type ServerConfiguration struct {
	RunAddress  string `json:"RUN_ADDRESS,omitempty"`
	Port        string `json:"PORT,omitempty"`
	DatabaseURI string `json:"DATABASE_URI,omitempty"`
	AuthTimeout int    `json:"AUTH_TIMEOUT,omitempty"`
	EnvChanged  map[string]bool
}

type ServerConfigurationOption func(*ServerConfiguration)

func UnMarshalServerDefaults(s string) ServerConfiguration {
	sc := ServerConfiguration{}
	err := json.Unmarshal([]byte(s), &sc)
	if err != nil {
		log.Fatal("cannot unmarshal server configuration")
	}
	return sc

}

func NewServerConfiguration() *ServerConfiguration {
	c := UnMarshalServerDefaults(ServerDefaultJSON)
	c.Port = ":" + strings.Split(c.RunAddress, ":")[1]
	c.EnvChanged = make(map[string]bool)
	return &c

}

func NewServerConf(options ...ServerConfigurationOption) *ServerConfiguration {
	c := UnMarshalServerDefaults(ServerDefaultJSON)
	c.EnvChanged = make(map[string]bool)
	for _, option := range options {
		option(&c)
	}
	return &c
}

func UpdateSCFromEnvironment(c *ServerConfiguration) {
	c.RunAddress = getEnv("RUN_ADDRESS", &StrValue{c.RunAddress}, c.EnvChanged).(string)
	//PORT is derived from ADDRESS
	c.Port = ":" + strings.Split(c.RunAddress, ":")[1]
	c.DatabaseURI = getEnv("DATABASE_URI", &StrValue{c.DatabaseURI}, c.EnvChanged).(string)
	c.AuthTimeout = getEnv("AUTH_TIMEOUT", &IntValue{c.AuthTimeout}, c.EnvChanged).(int)
}

func UpdateSCFromFlags(c *ServerConfiguration) {

	dc := NewServerConfiguration()

	var (
		a = flag.String("a", dc.RunAddress, "Domain name and :port")
		t = flag.Int("t", dc.AuthTimeout, "The timeout, user authorized, in seconds")
		d = flag.String("d", dc.DatabaseURI, "Database destination string")
	)
	flag.Parse()

	message := "variable %v  updated from flags, value %v"
	//Если значение из переменных равно значению по умолчанию, тогда берем из flagS
	if !c.EnvChanged["RUN_ADDRESS"] {
		c.RunAddress = *a
		c.Port = ":" + strings.Split(c.RunAddress, ":")[1]
		log.Printf(message, "RUN_ADDRESS", c.RunAddress)
		log.Printf(message, "PORT", c.Port)
	}
	if !c.EnvChanged["AUTH_TIMEOUT"] {
		c.AuthTimeout = *t
		log.Printf(message, "AUTH_TIMEOUT", c.AuthTimeout)
	}
	if !c.EnvChanged["DATABASE_URI"] {
		c.DatabaseURI = *d
		log.Printf(message, "DATABASE_URI", c.DatabaseURI)
	}
}

type VariableValue interface {
	Get() interface{}
	Set(string)
}
type StrValue struct {
	value string
}

func (v *StrValue) Get() interface{} {
	return v.value
}
func NewStrValue(s string) VariableValue {
	return &StrValue{value: s}
}
func (v *StrValue) Set(s string) {
	v.value = s
}

type IntValue struct {
	value int
}

func (v IntValue) Get() interface{} {
	return v.value
}
func (v *IntValue) Set(s string) {
	var err error
	v.value, err = strconv.Atoi(s)
	if err != nil {
		log.Fatal("Int Parse error")
	}
}

func NewIntValue(s string) VariableValue {
	changedValue, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("Int64 Parse error")
	}
	return &IntValue{value: changedValue}
}

type BoolValue struct {
	value bool
}

func (v BoolValue) Get() interface{} {
	return v.value
}
func (v *BoolValue) Set(s string) {
	var err error
	v.value, err = strconv.ParseBool(s)
	if err != nil {
		log.Fatal("Bool Parse error")
	}
}

func getEnv(variableName string, variableValue VariableValue, changed map[string]bool) (changedValue interface{}) {
	var stringVal string

	if variableValue == nil {
		log.Fatal("nil pointer in getEnv")
	}
	var exists bool
	stringVal, exists = os.LookupEnv(variableName)
	if !exists {
		log.Printf("variable "+variableName+" not presented in environment, remains default:%v", variableValue.Get())
		changed[variableName] = false
		return variableValue.Get()
	}
	variableValue.Set(stringVal)
	changed[variableName] = true
	log.Println("variable " + variableName + " presented in environment, value: " + stringVal)

	return variableValue.Get()
}
