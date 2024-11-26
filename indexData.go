package main

import (
	"context"
	"log"
	"math"
	"slices"

	"go.mongodb.org/mongo-driver/bson"
)

func createReturnDataStructure(token, successMessage, errorMessage string) IndexData {
	bashFlags := []Flag{}
	suasploitableFlags := []Flag{}
	var bashProgress float64 = 0.0
	var suasploitableProgress float64 = 0.0

	// get uuid
	uuid, success := getUuidFromToken(token)
	if !success {
		log.Println("Can not receive UUID")
		errorMessage = "Can not fetch user account from database"
	} else {
		// get collected flags from database
		// get user account first -> connect to db
		_, _, usersClient, usersCollection, success := connect("users")
		if !success {
			log.Println("Can not connect to users collection")
			errorMessage = "Can not fetch user data"
		} else {
			// get user by uuid
			user, err := getUserByUuid(usersCollection, uuid)
			if err != nil {
				log.Println("Can not get user: ", err)
				errorMessage = "Can not fetch user information"
			} else {
				// flags collected by user are stored in user.CollectedFlags
				// next, fetch all flags -> connect to db first
				_, _, ctfClient, ctfCollection, success := connect("ctf")
				if !success {
					log.Println("Can not connect to ctf collection")
					errorMessage = "Can not connect to ctf collection"
				} else {
					// fetch all flags
					cursor, err := ctfCollection.Find(context.TODO(), bson.D{})
					if err != nil {
						log.Println("Can not fetch flags from collection")
						errorMessage = "Can not fetch stored flags"
					} else {
						var allFlags []Flag
						err = cursor.All(context.TODO(), &allFlags)
						if err != nil {
							log.Println("Can not get cursor to iterate over flags")
							errorMessage = "Can not fetch flag information"
						} else {
							// figure out which flags the user collected and append them to the according list
							// furthermore count the number of total and collected flags
							numberOfBashFlags := 0
							numberOfSUASploitableFlags := 0
							numberOfCollectedBashFlags := 0
							numberOfCollectedSUASploitableFlags := 0
							for _, flag := range allFlags {
								// only count available flags
								if slices.Contains(user.AvailableFlags, flag.Flag) {
									switch flag.Type {
									case "bash":
										numberOfBashFlags += 1
									case "suasploitable":
										numberOfSUASploitableFlags += 1
									}
									if slices.Contains(user.CollectedFlags, flag.Flag) {
										switch flag.Type {
										case "bash":
											bashFlags = append(bashFlags, flag)
											numberOfCollectedBashFlags += 1
										case "suasploitable":
											suasploitableFlags = append(suasploitableFlags, flag)
											numberOfCollectedSUASploitableFlags += 1
										}
									}
								}
							}

							// calculate progress
							if numberOfBashFlags != 0 {
								bashProgress = (float64(numberOfCollectedBashFlags) / float64(numberOfBashFlags) * 100)
							}
							if numberOfSUASploitableFlags != 0 {
								suasploitableProgress = (float64(numberOfCollectedSUASploitableFlags) / float64(numberOfSUASploitableFlags) * 100)
							}
						}
					}
				}
				defer disconnect(ctfClient)
			}
		}

		defer disconnect(usersClient)
	}

	data := IndexData{
		Token:                 token,
		SuccessMessage:        successMessage,
		ErrorMessage:          errorMessage,
		BashProgress:          int(math.RoundToEven(bashProgress)),
		SUASploitableProgress: int(math.RoundToEven(suasploitableProgress)),
		BashFlags:             bashFlags,
		SUASploitableFlags:    suasploitableFlags,
	}
	return data
}
