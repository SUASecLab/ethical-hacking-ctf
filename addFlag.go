package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"

	"go.mongodb.org/mongo-driver/bson"
)

func addFlagToUser(w http.ResponseWriter, token, flagIdentifier, flagInput string) {
	// check if such a flag exists
	_, _, ctfClient, ctfCollection, success := connect("ctf")
	if !success {
		log.Println("Can not connect to ctf collection")
		fmt.Fprint(w, "Error, can not connect to ctf collection")
		return
	}
	var flag Flag
	err := ctfCollection.FindOne(context.TODO(), bson.D{{Key: "flag", Value: flagIdentifier}}).Decode(&flag)
	defer func() {
		if err := ctfClient.Disconnect(context.TODO()); err != nil {
			log.Println("Can not disconnect from database:", err)
		}
	}()
	log.Println(flagIdentifier)
	if err != nil {
		// flag does not exist
		// show error message
		log.Println("Could not find flag: ", err)
		runTemplate(w, createReturnDataStructure(token, "", "Entered flag does not exist"))
	} else {
		// flag exists
		usersCollectionContext, _, usersClient, usersCollection, success := connect("users")
		if !success {
			log.Println("Can not connect to users collection")
			fmt.Fprint(w, "Error, can not connect to users collection")
			return
		}

		//get uuid from token (is already validated earlier)
		uuid, success := getUuidFromToken(token)
		if !success {
			log.Println("Can not receive UUID: ", err)
			fmt.Fprint(w, "Error, can not identify user")
			return
		}
		user, err := getUserByUuid(usersCollection, uuid)
		if err != nil {
			// user does not exist
			// show error message
			log.Println("Can not find user: ", err)
			log.Println("Requested uuid: ", uuid)
			runTemplate(w, createReturnDataStructure(token, "", "User account does not exist"))
		} else {
			// user exists
			if user.CollectedFlags == nil {
				// create array and claim flag if array does not exist yet
				user.CollectedFlags = []string{flagIdentifier}
			} else {
				if slices.Contains(user.CollectedFlags, flagIdentifier) {
					// flag is already collected
					runTemplate(w, createReturnDataStructure(token, "", "Flag \""+flagInput+"\" already collected"))
					defer func() {
						if err := usersClient.Disconnect(context.TODO()); err != nil {
							log.Println("Can not disconnect from database:", err)
						}
					}()
					return
				}
				// append flag
				user.CollectedFlags = append(user.CollectedFlags, flagIdentifier)
			}

			//claim flag for the user in the database
			filter := bson.D{{Key: "uuid", Value: uuid}}
			update := bson.D{{Key: "$set",
				Value: bson.D{{
					Key:   "collectedFlags",
					Value: user.CollectedFlags,
				}}}}
			_, err := usersCollection.UpdateOne(usersCollectionContext, filter, update)

			// claiming not possible
			if err != nil {
				log.Println("Could not claim token: ", err)
				runTemplate(w, createReturnDataStructure(token, "", "Could not claim token"))
			} else {
				// flag claimed successfully
				runTemplate(w, createReturnDataStructure(token, "Claimed flag \""+flagInput+"\" successfully", ""))
			}
		}

		defer func() {
			if err := usersClient.Disconnect(context.TODO()); err != nil {
				log.Println("Can not disconnect from database:", err)
			}
		}()
	}
}
