package sneatgaeapp

import (
	"log"
	"os"
)

func initFirebase() {
	if gCloudProject := os.Getenv("GCLOUD_PROJECT"); gCloudProject != "" {
		log.Println("Sneat: GCLOUD_PROJECT:", gCloudProject)
	} else {
		log.Println("Sneat: GCLOUD_PROJECT is not set")
	}
	if firebaseAuthEmulatorHost := os.Getenv("FIREBASE_AUTH_EMULATOR_HOST"); firebaseAuthEmulatorHost != "" {
		log.Println("Sneat: FIREBASE_AUTH_EMULATOR_HOST:", firebaseAuthEmulatorHost)
	} else {
		log.Println("Sneat: FIREBASE_AUTH_EMULATOR_HOST is not set")
	}
	if firestoreEmulatorHost := os.Getenv("FIRESTORE_EMULATOR_HOST"); firestoreEmulatorHost != "" {
		log.Println("Sneat: FIRESTORE_EMULATOR_HOST:", firestoreEmulatorHost)
	} else {
		log.Println("Sneat: FIRESTORE_EMULATOR_HOST is not set")
	}
}
