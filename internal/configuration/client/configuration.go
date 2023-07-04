package configuration

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"strconv"
)

const ClientDefaultJSON = `{"ADDRESS":"localhost:8080"}`

const COMMAND_HELP = `Commands: 	
	For user's data:
	NEW 	- registers new given user's data if the user does not exist,
            	format: NEW  <name> <surname> <phone>  (for previously given login and password) 
	
	UPDATE	- updates user's data user's data is authorized by login and password
				format: UPDATE <name> <surname> <phone>
				there may be  sign "/" means the fields value is left unedited
				and empty sign "*" means the field should be empty
				Examples UPDATE Paul Anderson /    - phone remains untouched 
	For user's accounts data:
	
	ADD     - adds new account data for authorized user, 
				format: ADD <account> <password> <description>
	
	GET     - prints given user's account data to stdout if exists		   
				format: GET <account>
	
	EDIT    - updates user's account data 
				format: EDIT <account> <password> <description>
	if the command is not given, it prints user's all accounts data to stdout		    	   
`

type ClientConfiguration struct {
	Address    string `json:"ADDRESS,omitempty"`
	Login      string `json:"LOGIN,omitempty"`
	Password   string `json:"PASSWORD,omitempty"`
	Command    string `json:"COMMAND,omitempty"`
	EnvChanged map[string]bool
}

type ClientConfigurationOption func(*ClientConfiguration)

func UnMarshalAgentDefaults(s string) ClientConfiguration {
	ac := ClientConfiguration{}
	err := json.Unmarshal([]byte(s), &ac)
	if err != nil {
		log.Fatal("cannot unmarshal server configuration")
	}
	return ac
}

func NewClientConfiguration() *ClientConfiguration {
	c := UnMarshalAgentDefaults(ClientDefaultJSON)
	c.EnvChanged = make(map[string]bool)

	return &c

}

func NewClientConf(options ...ClientConfigurationOption) *ClientConfiguration {
	c := UnMarshalAgentDefaults(ClientDefaultJSON)
	c.EnvChanged = make(map[string]bool)
	for _, option := range options {
		option(&c)
	}
	return &c
}

func UpdateCCFromEnvironment(c *ClientConfiguration) {

	c.Address = getEnv("ADDRESS", &StrValue{c.Address}, c.EnvChanged).(string)

}

func UpdateCCFromFlags(cc *ClientConfiguration) {
	// cc = NewClientConfiguration()
	var (
		a = flag.String("a", cc.Address, "Domain name and :port")
		l = flag.String("l", cc.Login, "User login")
		p = flag.String("p", cc.Password, "User password")
		c = flag.String("c", cc.Command, COMMAND_HELP)

		messageVariableUpdated  = "variable %v  has been overwritten from flags, value %v"
		messageVariableReceived = "variable %v  received from flags, value %v"
	)

	flag.Parse()

	if !cc.EnvChanged["ADDRESS"] {
		cc.Address = *a
		log.Printf(messageVariableUpdated, "ADDRESS", cc.Address)
	}

	cc.Login = *l
	cc.EnvChanged["LOGIN"] = true
	log.Printf(messageVariableReceived, "LOGIN", cc.Login)

	cc.Password = *p
	cc.EnvChanged["PASSWORD"] = true
	log.Printf(messageVariableReceived, "PASSWORD", cc.Password)

	cc.Command = *c
	cc.EnvChanged["COMMAND"] = true
	log.Printf(messageVariableReceived, "COMMAND", cc.Command)
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
