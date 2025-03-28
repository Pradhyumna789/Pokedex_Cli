package main

import (
  "fmt"
  "strings"
  "bufio" 
  "os"
  "encoding/json" 
  "net/http"
  "time"
  "github.com/Pradhyumna789/Pokedex_Cli/internal/pokecache"
  "bytes"
  "io"
  "math/rand"
)

var cache *pokecache.Cache

var mapOfCaughtPokemon = make(map[string]Pokemon)

type cliCommand struct {
  name string
  description string
  callback func([]string) error
}

func commandExit(args []string) error {
  fmt.Println("Closing the Pokedex... Goodbye!")
  os.Exit(0)
  return nil
}

func commandHelp(commands map[string]cliCommand) func([]string) error {
  return func(args []string) error {
    fmt.Println("Welcome to the Pokedex!") 
    fmt.Println("Usage:")
    fmt.Println("")

    for _, v := range commands {
      fmt.Printf("%s: %s\n", v.name, v.description) 
    }

      return nil 
    }
}

func cleanInput(text string) []string { 
  lowered_sliced_text := strings.Fields(strings.ToLower(text))
  return lowered_sliced_text
}


type pokeApiResponse struct {
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type exploreCommandJson struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height    int `json:"height"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order        int `json:"order"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []any  `json:"past_abilities"`
	PastTypes     []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	} `json:"past_types"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       any    `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  any    `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      any    `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale any    `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

func fetchLocations(args []string) error {

    url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/")

    if cachedData, found := cache.Get(url); found {
      var pokeapiRes pokeApiResponse
      decoder := json.NewDecoder(bytes.NewReader(cachedData))
      err_decode := decoder.Decode(&pokeapiRes)

      if err_decode != nil {
        return fmt.Errorf("cachedData not successfully decoded %w", err_decode)
      }

      for _, location := range pokeapiRes.Results {
        fmt.Println(location.Name) 
      }

      return nil

    }

    client := &http.Client{
      Timeout: time.Second * 20,
    }

    req, err := http.NewRequest("GET", url, nil) 
    
    if err != nil {
      return fmt.Errorf("error creating a GET request %w", err)
    }

    res, err := client.Do(req)
    if err != nil {
      return fmt.Errorf("error getting a response %w", err)
    }

    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
      return fmt.Errorf("Error in converting response's body to a slice of bytes %w", err)
    }

    cache.Add(url, body)

    var pokeApiRes pokeApiResponse
    decoder := json.NewDecoder(bytes.NewReader(body)) 
    err_decode := decoder.Decode(&pokeApiRes)

    if err_decode != nil {
      return fmt.Errorf("error decoding json %w", err_decode)
    }
   
    for _, location := range pokeApiRes.Results {
      fmt.Println(location.Name)
    }

    return nil
}

func fetchLocationsBackwards(args []string) error {
    url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/")

    if cachedData, found := cache.Get(url); found {
      var pokeapiRes pokeApiResponse
      decoder := json.NewDecoder(bytes.NewReader(cachedData))
      err_decode := decoder.Decode(&pokeapiRes)

      if err_decode != nil {
        return fmt.Errorf("cachedData not successfully decoded %w", err_decode)
      }

      for _, location := range pokeapiRes.Results {
        fmt.Println(location.Name) 
      }

      return nil

    }

    client := &http.Client{
      Timeout: time.Second * 20,
    }

    req, err := http.NewRequest("GET", url, nil) 
    
    if err != nil {
      return fmt.Errorf("error creating a GET request %w", err)
    }

    res, err := client.Do(req)
    if err != nil {
      return fmt.Errorf("error getting a response %w", err)
    }

    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
      return fmt.Errorf("Error in converting response's body to a slice of bytes %w", err)
    }

    cache.Add(url, body)

    var pokeApiRes pokeApiResponse
    decoder := json.NewDecoder(bytes.NewReader(body)) 
    err_decode := decoder.Decode(&pokeApiRes)

    if err_decode != nil {
      return fmt.Errorf("error decoding json %w", err_decode)
    }
   
    for i := len(pokeApiRes.Results) - 1; i >= 0; i-- {
      fmt.Println(pokeApiRes.Results[i].Name)
    }

    return nil
}

func commandExplore() func([]string) error  {
  return func(args []string) error {
     if len(args) == 0 {
        return fmt.Errorf("Location are name is required")
     }

     locationName := args[0]
     url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", locationName)

     if cachedData, found := cache.Get(url); found {
      var exploreJson exploreCommandJson
      decoder := json.NewDecoder(bytes.NewReader(cachedData))
      err_decode := decoder.Decode(&exploreJson)

      if err_decode != nil {
        return fmt.Errorf("cachedData not successfully decoded %w", err_decode)
      }

      fmt.Println("Found Pokemon:")
      for _, pokemonEncounter := range exploreJson.PokemonEncounters {
        fmt.Printf(" - %s\n", pokemonEncounter.Pokemon.Name) 
      }

      return nil

    }

    client := &http.Client{
      Timeout: time.Second * 20,
    }

    req, err := http.NewRequest("GET", url, nil) 
    
    if err != nil {
      return fmt.Errorf("error creating a GET request %w", err)
    }

    res, err := client.Do(req)
    if err != nil {
      return fmt.Errorf("error getting a response %w", err)
    }

    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
      return fmt.Errorf("Error in converting response's body to a slice of bytes %w", err)
    }

    cache.Add(url, body)

    var exploreJson exploreCommandJson
    decoder := json.NewDecoder(bytes.NewReader(body)) 
    err_decode := decoder.Decode(&exploreJson)

    if err_decode != nil {
      return fmt.Errorf("error decoding json %w", err_decode)
    }
   
    fmt.Println("Found Pokemon:")
    for _, pokemonEncounter := range exploreJson.PokemonEncounters {
      fmt.Printf(" - %s\n", pokemonEncounter.Pokemon.Name) 
    }

    return nil
  }
} 

func commandCatch() func([]string) error {
  return func (args []string) error {
    if len(args) == 0 {
      return fmt.Errorf("Requires pokemon name to catch")
    } 

    pokemonName := args[0]
    url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemonName)

    client := &http.Client{
      Timeout: time.Second * 20,
    }

    req, err := http.NewRequest("GET", url, nil) 
    if err != nil {
      return fmt.Errorf("error creating a GET request %w", err)
    }

    res, err := client.Do(req)
    if err != nil {
      return fmt.Errorf("error getting a response %w", err)
    }

    defer res.Body.Close()

    body, err := io.ReadAll(res.Body)
    if err != nil {
      return fmt.Errorf("Error in converting response's body to a slice of bytes %w", err)
    }

    var pokemonNameJson Pokemon 
    decoder := json.NewDecoder(bytes.NewReader(body)) 
    err_decode := decoder.Decode(&pokemonNameJson)

    if err_decode != nil {
      return fmt.Errorf("error decoding json %w", err_decode)
    }

    fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

    // Seed uses the provided seed value to initialize the generator to a deterministic state.
    // Seed should not be called concurrently with any other [Rand] method.
    // difficulty of catching a pokemon is decided on the pokemon's base experience
    
    catchChance := 100 - (pokemonNameJson.BaseExperience / 10)  

    if rand.Intn(100) < catchChance {
      fmt.Printf("%s escaped!\n", pokemonName)
    } else {
      mapOfCaughtPokemon[pokemonName] = pokemonNameJson
      fmt.Printf("%s was caught!\n", pokemonName)
    }

    return nil
  }
}

func commandInspect() func([]string) error {
  return func(args []string) error {
    if len(args) == 0 {
      return fmt.Errorf("Error - requires pokemon name to inspect it")
    }

    pokemonName := args[0]

    caughtPokemon, found := mapOfCaughtPokemon[pokemonName]
    if !found {
      return fmt.Errorf("you have not caught that pokemon")
    } else {
      fmt.Printf("Name: %s\n", caughtPokemon.Name)
      fmt.Printf("Height: %d\n", caughtPokemon.Height)
      fmt.Printf("Weight: %d\n", caughtPokemon.Weight)
      fmt.Println("Stats: ")

      for _, stat := range caughtPokemon.Stats {
        fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.BaseStat)
      } 

      fmt.Println("Types:")

      for _, pokemonType := range caughtPokemon.Types { 
        fmt.Printf(" - %s\n", pokemonType.Type.Name)
      }
    }

    return nil
  }
}

func main() {

  cache = pokecache.NewCache(10 * time.Second)

  commandsRegistry := make(map[string]cliCommand)

  commandsRegistry["help"] = cliCommand{
        name: "help",
        description: "Displays a help message",
        callback: commandHelp(commandsRegistry), // parentheses after commandHelp because it's returning a higher order funciton (closure)
  }

  commandsRegistry["exit"] = cliCommand{
        name: "exit",
        description: "Exit the Pokedex",
        callback: commandExit,
  }

  commandsRegistry["map"] = cliCommand {
      name: "map",
      description: "shows next 20 locations of the pokemon",
      callback: fetchLocations,
  }

  commandsRegistry["mapb"] = cliCommand {
      name: "mapb",
      description: "shows previous 20 locations of the pokemon",
      callback: fetchLocationsBackwards,
  }

  commandsRegistry["explore"] = cliCommand {
      name: "explore",
      description: "explore pokemons in a particular location by it's name",
      callback: commandExplore(), // parentheses after commandExplore because this is also returning a higher order function like commandHelp (closure)
  }

  commandsRegistry["catch"] = cliCommand {
      name: "catch",
      description: "catch some pokemon",
      callback: commandCatch(),
  }

  commandsRegistry["inspect"] = cliCommand {
      name: "inspect",
      description: "inspect details of the caught pokemon",
      callback: commandInspect(),
  }

  scanner := bufio.NewScanner(os.Stdin)

  fmt.Println("Welcome to the Pokedex!")
  for {
      fmt.Print("Pokedex > ")

      scanner.Scan()
      scannedText := scanner.Text() 

      displaySlice := cleanInput(scannedText) 
      displayWord := displaySlice[0]

      commandEntered, found := commandsRegistry[displayWord]

      if !found {
        fmt.Println("Unknown command")
      } else {
        args := []string{}
        if len(displaySlice) > 1 {
          args = displaySlice[1:]
        }

        err := commandEntered.callback(args)
        if err != nil {
          fmt.Println("Error executing command: ", err)
        }
      }       

  }  

}




