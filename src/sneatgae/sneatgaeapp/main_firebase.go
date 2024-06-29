package sneatgaeapp

import (
	"log"
	"os"
)

func initFirebase() {
	if gCloudProject := os.Getenv("GCLOUD_PROJECT"); gCloudProject != "" {
		//logus.Println("Sneat: GCLOUD_PROJECT:", gCloudProject)
	} else {
		log.Println("WARNING: Sneat: GCLOUD_PROJECT is not set")
	}

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
