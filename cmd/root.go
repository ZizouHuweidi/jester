package cmd

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"

	"github.com/spf13/cobra"
)

type Joke struct {
	Type      string `json:"type"`
	Setup     string `json:"setup"`
	Punchline string `json:"punchline"`
}

var (
	content  embed.FS
	category string
)

var rootCmd = &cobra.Command{
	Use:   "jester",
	Short: "your own cli jester",
	Long: `jester is a CLI that provides you with a random joke or a joke from a specific category if you specify. For example:

   jester
   jester dad
   jester -l
`,
	Run: func(cmd *cobra.Command, args []string) {
		if listCategories, _ := cmd.Flags().GetBool("list"); listCategories {
			getAllCategories()
		} else if len(args) > 0 {
			getRandomJoke(args[0])
		} else {
			getRandomJoke("")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("list", "l", false, "Get all available joke categories")
}

func readJokesFromFile() ([]Joke, error) {
	file, err := content.ReadFile("../jokes.json")
	if err != nil {
		return nil, err
	}

	var jokes []Joke
	err = json.Unmarshal(file, &jokes)
	if err != nil {
		return nil, err
	}

	return jokes, nil
}

func getRandomJoke(category string) {
	jokes, err := readJokesFromFile()
	if err != nil {
		fmt.Println("Error reading jokes:", err)
		os.Exit(1)
	}

	if category != "" {
		categoryJokes := getCategoryJokes(jokes, category)
		if len(categoryJokes) == 0 {
			fmt.Printf("No jokes found for category: %s\n\nUse 'joker -l' to get all available joke categories", category)
		} else {
			randomJoke := getRandomJokeFromCategory(categoryJokes)
			fmt.Printf("%s\n\n%s\n\n", randomJoke.Setup, randomJoke.Punchline)
		}
	} else {
		randomJoke := getRandomJokeFromAllCategories(jokes)
		fmt.Printf("%s\n\n%s\n\n", randomJoke.Setup, randomJoke.Punchline)
	}
}

func getCategoryJokes(jokes []Joke, category string) []Joke {
	var categoryJokes []Joke
	for _, joke := range jokes {
		if joke.Type == category {
			categoryJokes = append(categoryJokes, joke)
		}
	}
	return categoryJokes
}

func getAllCategories() {
	jokes, err := readJokesFromFile()
	if err != nil {
		fmt.Println("Error reading jokes:", err)
		os.Exit(1)
	}

	uniqueCategories := make(map[string]struct{})
	for _, joke := range jokes {
		uniqueCategories[joke.Type] = struct{}{}
	}

	for category := range uniqueCategories {
		fmt.Println(category)
	}
}

func getRandomJokeFromAllCategories(jokes []Joke) *Joke {
	return &jokes[rand.Intn(len(jokes))]
}

func getRandomJokeFromCategory(jokes []Joke) *Joke {
	return &jokes[rand.Intn(len(jokes))]
}
