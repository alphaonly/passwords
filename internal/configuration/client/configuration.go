package configuration

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
)

const AgentDefaultJSON = `{"POLL_INTERVAL":"2s","REPORT_INTERVAL":"10s","ADDRESS":"localhost:8080","SCHEME":"http","USE_JSON":1,"KEY":"","RATE_LIMIT":1,"CRYPTO_KEY":""}`

type AgentConfiguration struct {
	Address    string `json:"ADDRESS,omitempty"`
	EnvChanged map[string]bool
}

type AgentConfigurationOption func(*AgentConfiguration)

func UnMarshalAgentDefaults(s string) AgentConfiguration {
	ac := AgentConfiguration{}
	err := json.Unmarshal([]byte(s), &ac)
	if err != nil {
		log.Fatal("cannot unmarshal server configuration")
	}
	return ac
}

func NewAgentConfiguration() *AgentConfiguration {
	c := UnMarshalAgentDefaults(AgentDefaultJSON)
	c.EnvChanged = make(map[string]bool)

	return &c

}

func NewAgentConf(options ...AgentConfigurationOption) *AgentConfiguration {
	c := UnMarshalAgentDefaults(AgentDefaultJSON)
	c.EnvChanged = make(map[string]bool)
	for _, option := range options {
		option(&c)
	}
	return &c
}

func UpdateACFromEnvironment(c *AgentConfiguration) {

	c.Address = getEnv("ADDRESS", &StrValue{c.Address}, c.EnvChanged).(string)

}

func UpdateACFromFlags(c *AgentConfiguration) {
	dc := NewAgentConfiguration()
	var (
		a = flag.String("a", dc.Address, "Domain name and :port")
	)

	flag.Parse()

	//Если значение параметра из переменных окружения равно по умолчанию, то обновляем из флагов

	message := "variable %v  updated from flags, value %v"
	if !c.EnvChanged["ADDRESS"] {
		c.Address = *a
		log.Printf(message, "ADDRESS", c.Address)
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
func NewBoolValue(s string) VariableValue {
	changedValue, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatal("Bool Parse error")
	}
	return &BoolValue{value: changedValue}
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
