package main

import (
	json "github.com/JoshuaDoes/json" //JSON wrapper to handle data more conveniently

	//std necessities
	stdjson "encoding/json"
	"strings"
)

/* Developer notes:
//  We use \x00 as the separator for arguments instead of a space character to accommodate for game arguments with spaces in them
*/

type Arguments struct {
	Game []stdjson.RawMessage `json:"game"`
	JVM  []stdjson.RawMessage `json:"jvm"`

	GameRaw   []string       `json:"-"`
	GameRules []ArgumentRule `json:"-"`

	JVMRaw   []string       `json:"-"`
	JVMRules []ArgumentRule `json:"-"`
}

func (args *Arguments) GameArgs() (arguments string) {
	arguments += strings.Join(args.GameRaw, "\x00")
	for i := 0; i < len(args.GameRules); i++ {
		valid := false
		for _, rule := range args.GameRules[i].Rules {
			valid = rule.Valid()

			if !valid {
				break
			}
		}
		if !valid {
			continue
		}

		if args.GameRules[i].RawValue != "" {
			arguments += "\x00" + args.GameRules[i].RawValue
		}

		if len(args.GameRules[i].RawValues) > 0 {
			for _, arg := range args.GameRules[i].RawValues {
				arguments += "\x00" + arg
			}
		}
	}

	return arguments
}

func (args *Arguments) JVMArgs() (arguments string) {
	arguments += strings.Join(args.JVMRaw, "\x00")
	for i := 0; i < len(args.JVMRules); i++ {
		valid := false
		for _, rule := range args.JVMRules[i].Rules {
			valid = rule.Valid()

			if !valid {
				break
			}
		}
		if !valid {
			continue
		}

		if args.JVMRules[i].RawValue != "" {
			arguments += "\x00" + args.JVMRules[i].RawValue
		}

		if len(args.JVMRules[i].RawValues) > 0 {
			for _, arg := range args.JVMRules[i].RawValues {
				if arg == "-Dos.name=Windows 10" {
					//tl;dr: Fuck you Java
					//Error: Could not find or load main class 10"
					//Caused by: java.lang.ClassNotFoundException: 10"
					continue
				}
				arguments += "\x00" + arg //As an even bigger fuck you, I'm quoteArg'ing harder
			}
		}
	}

	return arguments
}

func (args *Arguments) ParseValues() {
	args.GameRaw = make([]string, 0)
	args.GameRules = make([]ArgumentRule, 0)

	for i := 0; i < len(args.Game); i++ {
		argRule := ArgumentRule{}
		if err := json.Unmarshal(args.Game[i], &argRule); err == nil {
			argRule.ParseValues()
			args.GameRules = append(args.GameRules, argRule)
			continue
		}

		arg := ""
		if err := json.Unmarshal(args.Game[i], &arg); err == nil {
			args.GameRaw = append(args.GameRaw, arg)
			continue
		}

		log.Error("Unsupported Game argument type: ", args.Game[i])
	}

	args.JVMRaw = make([]string, 0)
	args.JVMRules = make([]ArgumentRule, 0)

	for i := 0; i < len(args.JVM); i++ {
		argRule := ArgumentRule{}
		if err := json.Unmarshal(args.JVM[i], &argRule); err == nil {
			argRule.ParseValues()
			args.JVMRules = append(args.JVMRules, argRule)
			continue
		}

		arg := ""
		if err := json.Unmarshal(args.JVM[i], &arg); err == nil {
			args.JVMRaw = append(args.JVMRaw, arg)
			continue
		}

		log.Error("Unsupported JVM argument type: ", args.JVM[i])
	}
}

type ArgumentRule struct {
	Rules []Rule             `json:"rules"`
	Value stdjson.RawMessage `json:"value"`

	RawValue  string   `json:"-"`
	RawValues []string `json:"-"`
}

func (argRule *ArgumentRule) ParseValues() {
	arg := ""
	if err := json.Unmarshal(argRule.Value, &arg); err == nil {
		argRule.RawValue = arg
		return
	}

	args := make([]string, 0)
	if err := json.Unmarshal(argRule.Value, &args); err == nil {
		argRule.RawValues = args
		return
	}

	log.Error("Unsupported argument rule type: ", argRule.Value)
}
