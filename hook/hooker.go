package hook

type HookCallback func(params ...interface{})

var hooks = make(map[int][]*HookCallback)

func RegisterHook(hookType int, cb HookCallback) {
	if currentHooks, ok := hooks[hookType]; ok {
		hooks[hookType] = append(currentHooks, &cb)
	} else {
		hooks[hookType] = []*HookCallback{&cb}
	}
}

func RegisterHooks(hookTypes []int, cb HookCallback) {
	for _, hookType := range hookTypes {
		RegisterHook(hookType, cb)
	}
}

func CallHooks(hookType int, params ...interface{}) {
	if hooks, ok := hooks[hookType]; ok {
		for _, hook := range hooks {
			(*hook)(params...)
		}
	}
}
