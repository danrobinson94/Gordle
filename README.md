# Gordle

A Wordle clone built in GoLang

Type "go run main.go" in the terminal to start the app.

On startup, the game selects a random five letter word from its dictionary.
User will have 6 tries to guess this word.

The app currently displays the correct word after each attempt to help with code review.

After each attempt, the game gives you feedback on the letters from your word:

Green letters - Correct letters in the correct place.
Yellow letters - Correct letters but in the wrong place.
Grey letters - Incorrect letters, not found in the winning word.

If you've guessed a word before, you won't be allowed to guess it again, to keep you
from wasting guesses.

If you guess a word not in the game's dictionary, you get to keep trying.

After 6 unsuccessful attempts you lose and the game will close automatically.

Enjoy!
