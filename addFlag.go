package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func addFlagToUser(w http.ResponseWriter, token, flagIdentifier, flagInput string) {
	// check if the requested flag exists

	// 1. connect to database
	_, _, ctfClient, ctfCollection, success := connect("ctf")
	if !success {
		log.Println("Can not connect to ctf collection")
		fmt.Fprint(w, "Error, can not connect to ctf collection")
		return
	}

	// 2. search flag
	var flag Flag
	var err error

	// If the flag is prefixed with suasploitable,
	// it could be a bonus flag
	if strings.HasPrefix(flagIdentifier, "suasploitable") {
		// Create the identifier the flag would have if it was a bonus flag
		bonusFlagIdentifier := "bonus-" + flagInput

		// Search for matching suasploitable and bonus flags
		err = ctfCollection.FindOne(context.TODO(),
			bson.D{{
				Key: "flag",
				Value: bson.D{{
					Key:   "$in",
					Value: bson.A{flagIdentifier, bonusFlagIdentifier},
				}},
			}}).Decode(&flag)

		// Check if flag is bonus flag
		if flag.Type == "bonus" {
			// Use bonus flag identifier from now on
			flagIdentifier = bonusFlagIdentifier
		}
		// Otherwise, it is another type of flag
		// Here we can continue with a normal search
	} else {
		err = ctfCollection.FindOne(context.TODO(), bson.D{{Key: "flag", Value: flagIdentifier}}).Decode(&flag)
	}
	defer disconnect(ctfClient)

	// 3. evaluate result
	if err != nil {
		// flag does not exist
		// show error message
		log.Println("Could not find flag: ", err)
		runTemplate(w, createReturnDataStructure(token, "", "Entered flag does not exist"))
		return
	}

	// if we are here, the flag surely exists
	// next we evaluate the user
	// first we have to get the UUID from the token (already validated, so there is no need for another validation)

	// 1. connect to database
	usersCollectionContext, _, usersClient, usersCollection, success := connect("users")
	if !success {
		log.Println("Can not connect to users collection")
		fmt.Fprint(w, "Error, can not connect to users collection")
		return
	}

	// 2. obtain uuid
	uuid, success := getUuidFromToken(token)
	if !success {
		log.Println("Can not receive UUID: ", err)
		fmt.Fprint(w, "Error, can not identify user")
		defer disconnect(usersClient)
		return
	}

	// 3. get user information -> also checks if user does exist
	user, err := getUserByUuid(usersCollection, uuid)
	if err != nil {
		// user does not exist
		// show error message
		log.Println("Can not find user: ", err)
		log.Println("Requested uuid: ", uuid)
		runTemplate(w, createReturnDataStructure(token, "", "User account does not exist"))
		defer disconnect(usersClient)
		return
	}

	// user exists
	// check if the given flag is available for the user
	if user.AvailableFlags == nil || !slices.Contains(user.AvailableFlags, flagIdentifier) {
		// send error message, flag exists but is not available for that user
		runTemplate(w, createReturnDataStructure(token, "", "Flag \""+flagInput+"\" is not available"))
		defer disconnect(usersClient)
		return
	}

	// user exists and flag is avialble -> add flag
	if user.CollectedFlags == nil {
		// create array and claim flag if array does not exist yet
		user.CollectedFlags = []string{flagIdentifier}
	} else {
		if slices.Contains(user.CollectedFlags, flagIdentifier) {
			// flag is already collected
			runTemplate(w, createReturnDataStructure(token, "", "Flag \""+flagInput+"\" already collected"))
			defer disconnect(usersClient)
			return
		}
		// append flag
		user.CollectedFlags = append(user.CollectedFlags, flagIdentifier)
	}

	// update flag in database
	filter := bson.D{{Key: "uuid", Value: uuid}}
	update := bson.D{{Key: "$set",
		Value: bson.D{{
			Key:   "collectedFlags",
			Value: user.CollectedFlags,
		}}}}

	_, err = usersCollection.UpdateOne(usersCollectionContext, filter, update)

	// claiming not possible
	if err != nil {
		log.Println("Could not claim token: ", err)
		runTemplate(w, createReturnDataStructure(token, "", "Could not claim token"))
		defer disconnect(usersClient)
		return
	}

	// flag claimed successfully
	successMessage := "Claimed flag \"" + flagInput + "\" successfully!"
	if flag.Type == "bonus" {
		successMessage = "Claimed flag \"" + flagInput + "\" successfully! Congratulations! You found a bonus flag!"
	}
	runTemplate(w, createReturnDataStructure(token, successMessage, ""))
	defer disconnect(usersClient)
}
