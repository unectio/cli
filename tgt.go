package main

type Target struct {
	Do	func(cmd int, namep *string)
}

func getTarget(t string) *Target {
	switch t {
	case "fn", "func", "function":
		return &Target { Do: doFunction }
	case "code":
		return &Target { Do: doCode }
	case "rt", "router":
		return &Target { Do: doRouter }
	case "repo", "repository":
		return &Target { Do: doRepo }
	case "sec", "secret":
		return &Target { Do: doSecret }
	case "tg", "trig", "trigger":
		return &Target { Do: doTrigger }
	case "am", "auth_method":
		return &Target { Do: doAuth }
	}

	return &Target { Do: func(_ int, _ *string) { usage() } }
}

func listTargets() []string {
	return []string {
		"fn | func | function",
		"code",
		"rt | router",
		"repo | repository",
		"sec | secret",
		"tg | trig | trigger",
		"am | auth_method",
	}
}

func doTargetCmd(cmd int, namep *string, actions map[int]func(namep *string)) {
	fn, ok := actions[cmd]
	if !ok {
		fn = func(_ *string){ usage() }
	}

	fn(namep)
}
