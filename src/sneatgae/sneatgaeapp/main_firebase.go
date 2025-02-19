package sneatgaeapp

import (
	"log"
	"os"
	"strings"
)

func logFirebaseEmulatorVars() {
	gCloudProject := os.Getenv("GOOGLE_CLOUD_PROJECT")

	if gCloudProject == "" {
		log.Println("WARNING: Sneat: GCLOUD_PROJECT is not set")
	} else if strings.HasPrefix(gCloudProject, "demo-") {
		logHost := func(name string) {
			if value := os.Getenv(name); value != "" {
				log.Printf("SNEAT: %s: %s://%s", name, "http", value)
			} else {
				log.Printf("Sneat: %s is not set", name)
			}
		}
		logHost("FIREBASE_AUTH_EMULATOR_HOST")
		logHost("FIRESTORE_EMULATOR_HOST")
	}
}
