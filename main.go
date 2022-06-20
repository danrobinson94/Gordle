package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

/* Setting a limit for number of letters in the word. When we grab words from our dictionary we
can limit our winning word to those with a length equal to the lettersLimit const */

const (
	lettersLimit int = 5
)

// Grabs words from a txt file, puts them in array of strings, then returns.
func chooseWord() []string {
	// Get txt file at https:
	res, err := http.Get("https://raw.githubusercontent.com/dwyl/english-words/master/words_alpha.txt")
	if err != nil {
		log.Fatalln(err)
	}
	// Reads all words until the end of the file, assigns them to "body"
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// Assigns the var created by ioutil.ReadAll to a string and splits it every time that there is a line break
	words := strings.Split(string(body), "\r\n")
	possibleWords := []string{}
	// For the entire length of the array of words that we have
	for _, word := range words {
		// If the word's length is below the lettersLimit we assigned above
		if len(word) == lettersLimit {
			// Append it to the list of possibleWords
			possibleWords = append(possibleWords, strings.ToUpper(word))
		}
	}
	// rand.Seed(time.Now().Unix())
	return possibleWords
}

/* Takes a string array and string as parameters. Checks if the exact string is contained in the array */
func contains(s []string, str string) bool {
	// For the length of the string array sent
	for _, v := range s {
		// if one of the strings in the array we sent is equal to the single string sent, return true
		if v == strings.ToUpper(str) {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("/--------------/")
	fmt.Println("WAR GAMES")
	fmt.Println("/--------------/")
	fmt.Println("Only joking, I am a humble wordle clone.")
	fmt.Println("Do you respond well to good-cop/bad-cop? Asking for a friend.")
	fmt.Println("/--------------/")
	words := chooseWord()	

	/* Select from the array of words, choosing a random number (rand.Intn) from 0 through
	 the length (len) of the words array. */
	rand.Seed(time.Now().Unix())
	chosenWord := words[rand.Intn(len(words))] 
	// Display the winning word for Jon and Andy's convenience when trying out the app
	fmt.Println("CORRECT WORD --", chosenWord)

	reader := bufio.NewReader(os.Stdin)
	// Initialize the count of guesses made
	tryCount := 0
	/* Keep an array of the words attempted so that the user doesn't have to waste an attempt
	on words they previously tried */
	triesArray := make([]string, 6, 6)
	// Let the user guess 6 different words from our dictionary of five letter words
	for tryCount < 6 {
			fmt.Print("ATTEMPT ", tryCount+1, "/6"," --> ")
			// Set the user's input to the variable 'text'
			text, _ := reader.ReadString('\n')
			// Clean up string (text). every time there is a line break in user input ("\n"), replace it 
			// with "" an infinite number of times (-1), then make it uppercase
			text = strings.Replace(text, "\n", "", -1)
			text = strings.ToUpper(text)
			// Check if the text input is in our dictionary of five letter words
			if contains(words, text) {
				// Check if the user won, and exit if they did
				if chosenWord == text {
					fmt.Println("------")	
					fmt.Println("BY THE BEARD OF ZEUS, YOU'VE DONE IT YOU LITTLE SHIT. WELL DONE, I SUPPOSE. UNTIL TOMORROW.")
					fmt.Println("------")
					os.Exit(0)
				// If the user did not win, check first if they've already guessed that word
				// Give them another try if they have
				} else if contains(triesArray, text) {
					fmt.Println("------")	
					fmt.Println("DON'T BE SILLY, YOU'VE ALREADY TRIED THIS WORD")	
					fmt.Println("------")	
				// If the user did not win and they have not attempted this word yet, check what
				// letters they guessed correctly
				} else {
					/* Break the attempted word and chosenWord into a list of substrings and see if any of the 
					letters are in the chosenWord */
					attemptLetters := strings.Split(text, "")
					correctLetters := strings.Split(chosenWord, "")
					var duplicatesFlag bool
					duplicateLetters := make([]string, 5)
					// Check if there are duplicate letters to decide what logic to use
					for i := 0; i < 5; i++ {
						duplicates := strings.Count(text, attemptLetters[i])
						if duplicates > 1 {
							duplicatesFlag = true
							duplicateLetters = append(duplicateLetters, attemptLetters[i])
						}
					}
					var greenLetters string
					// If there are duplicates, use more complicated logic
					if duplicatesFlag == true {
						// Create struct for how we keep track of what colors letters receive
						type letterColor struct {
							color string
							letter string
							yellowCount int 
							greenCount int
						}
						var letterColors = []letterColor{}
						duplicateGreenLetters := make([]string, 5)
						// Loop through letters to give grey or green values as needed
						for i := 0; i < 5; i++ {

							// If the letter from user input at index i matches the chosenWord letter at i, 
							// assign it "GREEN"
							if attemptLetters[i] == correctLetters[i] {
								newLetterColor := letterColor{
									color: "GREEN",
									letter: attemptLetters[i],
								}
								letterColors = append(letterColors, newLetterColor)
								greenLetters += attemptLetters[i]
								// If the letter isn't in the word, assign it "GREY"
							} else if !contains(correctLetters, attemptLetters[i]) || 
							contains(correctLetters, attemptLetters[i]) {
								newLetterColor := letterColor{
									color: "GREY",
									letter: attemptLetters[i],
								}
								letterColors = append(letterColors, newLetterColor)
							} 
							// See how many duplicate green letters we have and append to list if necessary
							duplicateGreens := strings.Count(greenLetters, attemptLetters[i])
							if duplicateGreens > 1 && 
							!contains(duplicateGreenLetters, attemptLetters[i]) {
								duplicateGreenLetters = append(duplicateGreenLetters, attemptLetters[i])
							}
						}
							yellowsAdded := make([]string, 5)
							// Loop through letters, and if they haven't been made "GREEN", decide
							// If they should be made yellow. 
							for i:= 0; i < 5; i++ {
								// See how many of the duplicate letters in user input are actually duplicates
								// In the chosen word
								correctDuplicates := strings.Count(chosenWord, attemptLetters[i])
								// Checks to see if any attempted letters should
								// made yellow, making sure to not change any green letters to yellow.
								// This will need to be simplified in the future
								if contains(correctLetters, attemptLetters[i]) && 
									letterColors[i].color != "GREEN" &&
									((!contains(duplicateLetters, attemptLetters[i]) && 
									!contains(yellowsAdded, attemptLetters[i])) ||
									(contains(duplicateGreenLetters, attemptLetters[i]) &&
									contains(duplicateLetters, attemptLetters[i])) || 
									(!contains(duplicateGreenLetters, attemptLetters[i]) && 
									!contains(yellowsAdded, attemptLetters[i])) &&
									strings.Count(greenLetters, attemptLetters[i]) < correctDuplicates &&
									letterColors[i].color != "GREEN") ||
									((contains(correctLetters, attemptLetters[i]) && 
									!contains(yellowsAdded, attemptLetters[i])) &&
									strings.Count(greenLetters, attemptLetters[i]) < correctDuplicates &&
									letterColors[i].color != "GREEN") {
										// Now we can make the letter yellow
										letterColors[i].color = "YELLOW"
										letterColors[i].yellowCount += 1
										yellowsAdded = append(yellowsAdded, attemptLetters[i])
								}
							}
							fmt.Println("CORRECT WORD --", chosenWord)
							// Output results 
							for i:=0; i < len(letterColors); i++ {
								
								if letterColors[i].color == "GREEN" {
									fmt.Println(attemptLetters[i], " - ", "GREEN  - CORRECT LETTER AND PLACEMENT")
								}
								if letterColors[i].color == "YELLOW" {
									fmt.Println(attemptLetters[i], " - ", "YELLOW - RIGHT LETTER, WRONG PLACE")
								}
								if letterColors[i].color == "GREY" {
									fmt.Println(attemptLetters[i], " - ", "GREY   - WRONG LETTER. CHIN UP, YOU'LL GET IT NEXT TIME!")
								}
							}
					// If the word the user entered doesn't have duplicates, user more performant logic
					} else {
						fmt.Println("CORRECT WORD --", chosenWord) 
							yellowLetters := make([]string, 5)
							greenLetters := make([]string, 5)
					// Loop through 5 times, the amount of letters in chosenWord
					for i := range chosenWord {

						// attemptedLetters = append(attemptedLetters, attemptLetters[i])
						// currentAttemptLetters = append(currentAttemptLetters, attemptLetters[i])

						// Check if the index letter of the attempted word exactly matches the chosen word
						if attemptLetters[i] == correctLetters[i] {
							fmt.Println(attemptLetters[i], " - ", "GREEN  - CORRECT LETTER AND PLACEMENT")
							greenLetters = append(greenLetters, attemptLetters[i])
						// Else, check if it is not an exact match but the letter is still in the word
						// using the above contains function
						} else if (attemptLetters[i] != correctLetters[i]) && 
							contains(correctLetters, attemptLetters[i]) {
								fmt.Println(attemptLetters[i], " - ", "YELLOW - RIGHT LETTER, WRONG PLACE")
								yellowLetters = append(yellowLetters, attemptLetters[i])
						// Else, we will tell the user that their letter is not in the word
						} else {
							fmt.Println(attemptLetters[i], " - ", "GREY   - WRONG LETTER. I CANNOT UNDERSTATE HOW DISSAPPOINTED I AM.")
						}
					}
				}
					// The following actions are done regardless of any correct letters
					fmt.Println("------")	
					fmt.Println("GOOD TRY")
					fmt.Println("------")	
					// Add the attempted word to our triesArray and increase the tryCount by 1
					triesArray[tryCount] = text
					tryCount++
				}
			// If the word the user entered is not five letters long, or is a five letter word not
			// contained in our dictionary, display this message
			} else {
				fmt.Println("------")	
				fmt.Println("'", text, "'", "IS NOT IN MY DICTIONARY OF FIVE LETTER WORDS. TRY AGAIN.")
				fmt.Println("------")
				fmt.Println("OR, FEEL FREE TO BRIBE ME AND I MAY CONSIDER ADDING YOUR WORD.")
				fmt.Println("------")
			}
	}
	// If the user attempts 6 times without winning and triggering the exit code, display the following
	fmt.Println("WELL... ALMOST")
	fmt.Println("------")	
	for i := 0; i < 75; i++ {
		fmt.Println("YOU LOSE I WIN")
	}
	fmt.Println("-------")
	fmt.Println("I SINCERELY APOLOGIZE, IT WAS RUDE OF ME TO GLOAT. SLIP OF THE TONGUE. PLEASE DO COME BACK ANY TIME YOU WOULD LIKE TO ENDURE ANOTHER HORRIFYING LOSS.")

}