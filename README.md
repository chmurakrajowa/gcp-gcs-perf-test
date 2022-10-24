# gcp-gcs-perf-test

## Configuration environment variables
- `BUCKET_NAME` (required): GCS Bucket name to write objects to
- `NUM_OBJECTS` (default: `256`): Number of objects to create and write to

- `NAME` (default: `gcp-gcs-perf-test`): Application name used to identify app in the platform and logs
- `VERSION` (default: `v0.1.1`): Application version
- `CF_INSTANCE_GUID` (default: auto generated): Instance GUID

- `LOGLEVEL` (default: `TRACE`): Log Level. Possible values: `PANIC`, `FATAL`, `ERROR`, `WARNING`, `INFO`, `DEBUG` or `TRACE`
- `LOGAS` (default: `text`): Log As `text` or `json`
- `DEBUGLEVEL` (default: `0`): Debug Level
- `LOGDETAILS` (default: `false`): Log additional details

- `TIMEZONE` (default: `Local`): Timezone

- `CONNECTION_TIMEOUT` (default: `8s`): Connection timeout       
- `SHUTDOWNWAIT_TIME` (default: `480s`): Shutdown wait time 
