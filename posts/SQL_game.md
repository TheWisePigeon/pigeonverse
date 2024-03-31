---
title: Writing a game with SQL (and Go)
posted_at: 30 03 2024
slug: SQL_game
tldr: An article going through how I made a hangman terminal game but all the game logic is written in SQL and handled by Postgres
draft: true
---

So, last week I watched [this video](https://youtu.be/9a_O3QHLajw?si=V2QBo6MTamnDLZpw) and I thought to myself "dangg, this is crazy".
It wasn't the first time I was seeing people do crazy things with SQL and this time I decided to try and do something a little bit crazy too.
And to make it better I have been stepping away from classic backend development and exploring game dev for the past two weeks so this was a perfect
project to work on.

I decided to write a Hangman game that could be played from the terminal but the whole game logic meaning the generation of the word to be guessed, the processing
of a guess and of whether the user had won or lost the game had to be written in SQL and handled by the database. The client would just be a dummy responsible of
displaying the game state in an understandable way to the user.
In this article I will go over the code and also over my thought process while making this project, because at the time I wanted to built the project I had absolutely no idea how I could
do it (skill issues basically), but slowly one google search at a time I did it, so without further a do let's dive in :)

## The game
The concept for a Hangman game is very much simple. There is a word that the player has to guess, every time they make a wrong guess, we start drawing a hanged stickman
one stick at a time, the player loses when the drawing is complete meaning they could not guess the word.
As I said in the introduction, I wanted all the logic to be handled by my database (I used Postgres) and a client that would just be responsible of displaying the game state.
So the first thing I did was to write out all the functionalities of the game:

- Start a new game
- Accept a guess
- Process a guess
- Detect a loss
- Detect a win

With this in mind it was time to write some SQL. 

### Starting a new game
When a new game starts, we need some things both on the client and in the database. We need the word to be guessed, then we need to display something on the client that would tell
the user how many letters are contained in that word. We also need to be able to keep track of the game state, so in the client code I had this model
```go
type GameModel struct {
	ID              string //Game id
	NumberOfLetters int //Number of letters in word to guess
	WrongGuesses    int `db:"wrong_guesses"` //Amount of wrong guesses the player has made
	GuessedLetters  []string //Guessed letters that are contained in the word to be guessed
	AlreadyGuessed  []string //All guessed letters
	Finished        bool   `db:"game_state"` //Whether the game is finished or not
	WordToGuess     string `db:"word_to_guess"` //The word to be guessed
	db              *sql.DB //A connection pool to the database
}
```
If it doesn't already, this will make more sense as we go.
