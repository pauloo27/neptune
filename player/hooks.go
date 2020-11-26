package player

const (
	HOOK_PLAYER_INITIALIZED = iota
	HOOK_FILE_LOAD_STARTED
	HOOK_PLAYBACK_PAUSED
	HOOK_PLAYBACK_RESUMED
	HOOK_VOLUME_CHANGED
	HOOK_FILE_LOADED
	HOOK_POSITION_CHANGED
)

type HookCallback func(err error, params ...interface{})

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

func callHooks(hookType int, err error, params ...interface{}) {
	if hooks, ok := hooks[hookType]; ok {
		for _, hook := range hooks {
			(*hook)(err, params...)
		}
	}
}
