package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)
// Pet struct
type Pet struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
	Age  int    `json:"age"`
}

// Owner struct
type Owner struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Pets []Pet  `json:"pets"`
}

// global variables for storing owners and auto-incrementing IDs
var owners = []Owner{} // owner list
var nextOwnerID = 1    // next owner ID
var nextPetID = 1      // next pet ID

// Create a new owner (POST /owners)
func createOwner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	owner := Owner{}
	json.NewDecoder(r.Body).Decode(&owner)

	owner.ID = nextOwnerID // set ID
	nextOwnerID++          // increment nextOwnerID
	owners = append(owners, owner)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(owner)
}

// Add a new pet to an owner (POST /owners/{id}/pets)
func createPet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Wrong method", http.StatusMethodNotAllowed)
		return
	}

	ownerIDStr := r.URL.Query().Get("ownerId")
	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		http.Error(w, "Owner ID error", http.StatusBadRequest)
		return
	}

	pet := Pet{}
	json.NewDecoder(r.Body).Decode(&pet)

	for i := 0; i < len(owners); i++ { // loop to find owner
		if owners[i].ID == ownerID {
			pet.ID = nextPetID
			nextPetID++
			owners[i].Pets = append(owners[i].Pets, pet)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(pet)
			return
		}
	}

	http.Error(w, "No owner", http.StatusNotFound)
}

// Get all owners and their pets (GET /owners)
func getOwners(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "No", http.StatusMethodNotAllowed)
		return
	}

	json.NewEncoder(w).Encode(owners) // directly send data
}

// Update owner (PUT /owners/{id})
func updateOwner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT", http.StatusMethodNotAllowed)
		return
	}

	ownerIDStr := r.URL.Query().Get("id")
	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		http.Error(w, "Wrong ID", http.StatusBadRequest)
		return
	}

	var updatedOwner Owner
	json.NewDecoder(r.Body).Decode(&updatedOwner)

	for i := 0; i < len(owners); i++ {
		if owners[i].ID == ownerID {
			owners[i].Name = updatedOwner.Name
			json.NewEncoder(w).Encode(owners[i])
			return
		}
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

// Delete an owner (DELETE /owners/{id})
func deleteOwner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "No", http.StatusMethodNotAllowed)
		return
	}

	ownerIDStr := r.URL.Query().Get("id")
	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		http.Error(w, "Wrong ID", http.StatusBadRequest)
		return
	}

	for i := 0; i < len(owners); i++ {
		if owners[i].ID == ownerID {
			owners = append(owners[:i], owners[i+1:]...)
			fmt.Fprintf(w, "Deleted %d", ownerID)
			return
		}
	}

	http.Error(w, "Not found", http.StatusNotFound)
}

// Delete a pet (DELETE /owners/{ownerId}/pets/{petId})
func deletePet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "No", http.StatusMethodNotAllowed)
		return
	}

	ownerIDStr := r.URL.Query().Get("ownerId")
	ownerID, err := strconv.Atoi(ownerIDStr)
	if err != nil {
		http.Error(w, "Owner ID error", http.StatusBadRequest)
		return
	}

	petIDStr := r.URL.Query().Get("petId")
	petID, err := strconv.Atoi(petIDStr)
	if err != nil {
		http.Error(w, "Pet ID error", http.StatusBadRequest)
		return
	}

	for i := 0; i < len(owners); i++ {
		if owners[i].ID == ownerID {
			for j := 0; j < len(owners[i].Pets); j++ {
				if owners[i].Pets[j].ID == petID {
					owners[i].Pets = append(owners[i].Pets[:j], owners[i].Pets[j+1:]...)
					fmt.Fprintf(w, "Pet deleted")
					return
				}
			}
			http.Error(w, "No pet", http.StatusNotFound)
			return
		}
	}

	http.Error(w, "No owner", http.StatusNotFound)
}

// main function to handle routes
func main() {
	http.HandleFunc("/owners", createOwner)           // POST
	http.HandleFunc("/owners/all", getOwners)         // GET
	http.HandleFunc("/owners/update", updateOwner)    // PUT
	http.HandleFunc("/owners/delete", deleteOwner)    // DELETE
	http.HandleFunc("/owners/pets", createPet)        // POST
	http.HandleFunc("/owners/pets/delete", deletePet) // DELETE

	fmt.Println("Server starting :6459")
	log.Fatal(http.ListenAndServe(":6459", nil))
}
