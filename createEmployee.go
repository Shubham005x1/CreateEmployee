package content

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/Shubham005x1/MyValidations/validations"
	"golang.org/x/crypto/bcrypt"
)

var (
	client     *firestore.Client
	logClient  *logging.Client
	onceClient sync.Once
)

func init() {
	functions.HTTP("CreateEmployee", CreateEmployee)
}

// CreateEmployee handles the creation of an employee record.
func CreateEmployee(w http.ResponseWriter, r *http.Request) {
	// Create a new background context.
	ctx := context.Background()

	// Initialize Firestore connection.
	initializeFirestore()
	initializePubsub()
	// Initialize the Logging client for logging events.
	logClient, _ = logging.NewClient(ctx, "takeoff-task-3")

	// Ensure the logClient is closed after the function completes.
	defer logClient.Close()

	// Create a logger for this function.
	logger := logClient.Logger("my-log")

	// Log an entry indicating that the CreateEmployee function has started.
	logger.Log(logging.Entry{
		Payload:  "CreateEmployee function started",
		Severity: logging.Info,
	})

	// Parse the request body to retrieve employee data.
	var emp Employee // Assuming there's a struct named Employee.
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		http.Error(w, "Invalid request body!", http.StatusBadRequest)
		logger.Log(logging.Entry{
			Payload:  "Invalid request body!",
			Severity: logging.Error,
		})
		return
	}

	// Check if the employee ID contains characters (non-numeric).
	if validNumberEntry(emp.ID) {
		http.Error(w, "Id cannot be character", http.StatusBadRequest)
		return
	}

	// Query Firestore to check if an employee with the same ID already exists.
	collectionRef := client.Collection("employees")
	query := collectionRef.Where("id", "==", emp.ID)
	data, err := query.Documents(ctx).GetAll()
	if err != nil {
		// Handle error querying Firestore.
		http.Error(w, fmt.Sprintf("Error querying Firestore: %v", err), http.StatusInternalServerError)
		return
	}

	if len(data) > 0 {
		// At least one document was found, indicating an existing employee with the same ID.
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Employee with the same ID already exists"))
		return
	}

	// Check if any required fields are empty.
	if emp.ID == "" || emp.FirstName == "" || emp.LastName == "" || emp.Email == "" || emp.Password == "" || emp.PhoneNo == "" || emp.Role == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	// Validate FirstName and LastName to ensure they don't contain numbers.
	if validations.ValidNameEntry(emp.FirstName) {
		http.Error(w, "Name Cannot contain Numbers please enter valid Name", http.StatusBadRequest)
		return
	}
	if validations.ValidNameEntry(emp.LastName) {
		http.Error(w, "LastName Cannot contain Numbers please enter valid Name", http.StatusBadRequest)
		return
	}
	// Assuming `emp` has a `Password` field.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(emp.Password), bcrypt.DefaultCost)
	if err != nil {
		// Handle error
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	emp.Password = string(hashedPassword)
	logger.Log(logging.Entry{
		Payload:  "Password Hashed and Saved in Database",
		Severity: logging.Info,
	})

	// Validate the email format.
	err = validations.IsValidEmail(emp.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = validations.IsNumberValid(emp.PhoneNo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log that the request body has been parsed successfully.
	logger.Log(logging.Entry{
		Payload:  "Request body parsed successfully",
		Severity: logging.Info,
	})
	employeeJSON, err := json.Marshal(emp)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error serializing employee to JSON: %v", err), http.StatusInternalServerError)
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error serializing employee to JSON: %v", err),
			Severity: logging.Error,
		})
		return
	}

	msg := &pubsub.Message{
		Data: employeeJSON,
	}
	// result, _ := topic.Publish(ctx, msg).Get(ctx)
	// serverID, err := result.Get(ctx)
	result := topic.Publish(ctx, msg)
	serverID, err := result.Get(ctx)
	if err != nil {
		// Handle error publishing to Pub/Sub.
		http.Error(w, fmt.Sprintf("Error publishing to Pub/Sub: %v", err), http.StatusInternalServerError)
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Error publishing to Pub/Sub: %v", err),
			Severity: logging.Error,
		})
		return
	}
	if serverID == "" {
		// The message was not successfully published.
		http.Error(w, "Failed to publish message to Pub/Sub", http.StatusInternalServerError)
		logger.Log(logging.Entry{
			Payload:  "Failed to publish message to Pub/Sub",
			Severity: logging.Error,
		})
		return
	}
	logger.Log(logging.Entry{
		Payload:  "Message published to Pub/Sub",
		Severity: logging.Info,
	})

	// Add the employee data to Firestore.
	_, _, err = client.Collection("employees").Add(ctx, emp)
	if err != nil {
		// Handle error adding employee to Firestore.
		http.Error(w, fmt.Sprintf("Failed to add employee: %v", err), http.StatusInternalServerError)
		logger.Log(logging.Entry{
			Payload:  fmt.Sprintf("Failed to add employee: %v", err),
			Severity: logging.Error,
		})
		return
	}

	// Log that the employee was added successfully.
	logger.Log(logging.Entry{
		Payload:  "Employee added successfully",
		Severity: logging.Info,
	})
	defer pubsubClient.Close()
	// Respond with a status code indicating success and a message.
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Employee created successfully"))
}
