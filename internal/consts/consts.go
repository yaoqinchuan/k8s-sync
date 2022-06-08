package consts

const (
	WS_PENDING  = "PENDING"
	WS_STARTING = "STARTING"
	// start timeout in seconds
	STARTING_TIMEOUT    = 300
	WS_STARTING_TIMEOUT = "STARTING_TIMEOUT"
	WS_RUNNING          = "RUNNING"
	WS_DELETING         = "DELETING"
	WS_DELETING_TIMEOUT = "DELETING_TIMEOUT"
	DELETING_TIMEOUT    = 300
	WS_STOPPING         = "STOPPING"
	WS_STOPPED          = "STOPPED"
	WS_STOPPING_TIMEOUT = "STOPPING_TIMEOUT"
	STOPPING_TIMEOUT    = 300
	WS_ERROR            = "ERROR"

	// workspace component kind
	PersistentVolumeClaim = "PersistentVolumeClaim"
	SERVICE               = "SERVICE"
	CONFIGMAP             = "CONFIGMAP"
)
