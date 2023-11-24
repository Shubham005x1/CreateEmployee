package content

import (
	"context"
	"log"
	"unicode"

	"cloud.google.com/go/firestore"
)

type Employee struct {
	ID        string `firestore:"id" json:"id"`
	FirstName string `firestore:"firstname" json:"firstname"`
	LastName  string `firestore:"lastname" json:"lastname"`
	Email     string `firestore:"email" json:"email"`
	Password  string `firestore:"password" json:"password"`
	PhoneNo   string `firestore:"phoneNo" json:"phoneNo"`
	Role      string `firestore:"role" json:"role"`
}

func initializeFirestore() {
	Once.Do(func() {
		ctx := context.Background()

		// Initialize Firestore with the service account key
		var err error
		client, err = firestore.NewClient(ctx, "takeoff-task-3")
		if err != nil {
			log.Fatalf("Failed to create Firestore client: %v", err)
		}
	})
}
func validNumberEntry(name string) bool {
	for _, char := range name {
		if unicode.IsLetter(char) {
			return true
		}
	}
	return false

}
