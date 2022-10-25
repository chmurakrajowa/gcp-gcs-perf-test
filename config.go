package main

import "time"

type config struct {
	ApplicationName    string `env:"NAME" default:"gcp-gcs-perf-stat"`
	ApplicationVersion string `env:"VERSION" default:"v0.1.2"`
	InstanceGUID       string `env:"CF_INSTANCE_GUID"`

	LogLevel   string `env:"LOGLEVEL" default:"TRACE"` // "PANIC"|"FATAL"|"ERROR"|"WARNING"|"INFO"|"DEBUG"|"TRACE"
	LogAs      string `env:"LOGAS" default:"text"`     // "text"|"json"
	DebugLevel int    `env:"DEBUGLEVEL" default:"0"`
	LogDetails bool   `env:"LOGDETAILS" default:"false"` // "true"|"false"

	ShutdownWaitTime  time.Duration `env:"SHUTDOWNWAIT_TIME" default:"480s"`
	ConnectionTimeout time.Duration `env:"CONNECTION_TIMEOUT" default:"8s"`
	Timezone          string        `env:"TIMEZONE" default:"Local"`

	BucketName string `env:"BUCKET_NAME"`
	NumObjects int    `env:"NUM_OBJECTS" default:"256"`
}
