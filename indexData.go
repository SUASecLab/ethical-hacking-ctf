package main

import (
	"context"
	"log"
	"math"
	"slices"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

func createReturnDataStructure(token, successMessage, errorMessage string) IndexData {
	bashFlags := []Flag{}
	suasploitableFlags := []Flag{}
	bonusFlags := []Flag{}
	examFlags := []Flag{}
	var bashProgress float64 = 0.0
	var suasploitableProgress float64 = 0.0
	var bonusProgress float64 = 0.0
	var examProgress = 0.0

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
							numberOfBonusFlags := 0
							numberOfExamFlags := 0
							numberOfCollectedBashFlags := 0
							numberOfCollectedSUASploitableFlags := 0
							numberOfCollectedBonusFlags := 0
							numberOfCollectedExamFlags := 0
							knownFlags := []Flag{}
							for _, flag := range allFlags {
								// Check if we already know this flag
								if slices.Contains(knownFlags, flag) {
									continue
								}

								// Add flags to list of known flags
								knownFlags = append(knownFlags, flag)

								// Create union of flags and available flags for that user
								if user.AvailableFlags != nil && slices.Contains(user.AvailableFlags, flag.Flag) {
									// Count flags
									switch flag.Type {
									case "bash":
										numberOfBashFlags += 1
									case "suasploitable":
										numberOfSUASploitableFlags += 1
									case "bonus":
										numberOfBonusFlags += 1
									case "exam":
										numberOfExamFlags += 1
									}

									// Check for each flag if the flag was collected
									if slices.Contains(user.CollectedFlags, flag.Flag) {
										// Truncate flag prefix
										// The second line ensures we don't run into an error if there is no hypen (which never should be the case)
										flagParts := strings.SplitAfterN(flag.Flag, "-", 2)
										flag.Flag = flagParts[len(flagParts)-1]

										switch flag.Type {
										case "bash":
											bashFlags = append(bashFlags, flag)
											numberOfCollectedBashFlags += 1
										case "suasploitable":
											suasploitableFlags = append(suasploitableFlags, flag)
											numberOfCollectedSUASploitableFlags += 1
										case "bonus":
											bonusFlags = append(bonusFlags, flag)
											numberOfCollectedBonusFlags += 1
										case "exam":
											examFlags = append(examFlags, flag)
											numberOfCollectedExamFlags += 1
										}
									}
								}
							}

							// Calculate progress by dividing collected flags by available flags per CtF
							if numberOfBashFlags != 0 {
								bashProgress = (float64(numberOfCollectedBashFlags) / float64(numberOfBashFlags) * 100)
							}
							if numberOfSUASploitableFlags != 0 {
								suasploitableProgress = (float64(numberOfCollectedSUASploitableFlags) / float64(numberOfSUASploitableFlags) * 100)
							}
							if numberOfBonusFlags != 0 {
								bonusProgress = (float64(numberOfCollectedBonusFlags) / float64(numberOfBonusFlags) * 100)
							}
							if numberOfExamFlags != 0 {
								examProgress = (float64(numberOfCollectedExamFlags) / float64(numberOfExamFlags) * 100)
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
		BonusProgress:         int(math.RoundToEven(bonusProgress)),
		ExamProgress:          int(math.RoundToEven(examProgress)),
		BashFlags:             bashFlags,
		SUASploitableFlags:    suasploitableFlags,
		BonusFlags:            bonusFlags,
		ExamFlags:             examFlags,
	}
	return data
}
