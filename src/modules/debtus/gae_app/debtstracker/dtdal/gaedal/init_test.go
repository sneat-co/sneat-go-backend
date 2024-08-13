package gaedal

import "os"

func init() {
	os.Setenv("GAE_LONG_APP_ID", "gae-unit-tests")
	os.Setenv("GAE_PARTITION", "gae-partition")
	os.Setenv("RUN_WITH_DEVAPPSERVER", "yes")
}
