package main

import (
	"fmt"
	"log"

	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func main() {

	apiURL, err := url.Parse("https://YourUrl")
	if err != nil {
		log.Println(err)
	}
	APIKey := "API-YourAPIKey"
	spaceName := "Default"
	environmentNames := []string{"Environment1", "Environment2"}

	// Get reference to space
	space := GetSpace(apiURL, APIKey, spaceName)

	// Get references to environments
	for i := 0; i < len(environmentNames); i++ {
		environment := GetEnvironment(apiURL, APIKey, space, environmentNames[i])

		// Check to see if environment already exists
		if environment == nil {
			environment := octopusdeploy.NewEnvironment(environmentNames[i])
			client := octopusAuth(apiURL, APIKey, space.ID)
			client.Environments.Add(environment)
		} else {
			fmt.Println("Environment " + environment.Name + " already exists.")
		}
	}
}

func octopusAuth(octopusURL *url.URL, APIKey, space string) *octopusdeploy.Client {
	client, err := octopusdeploy.NewClient(nil, octopusURL, APIKey, space)
	if err != nil {
		log.Println(err)
	}

	return client
}

func GetSpace(octopusURL *url.URL, APIKey string, spaceName string) *octopusdeploy.Space {
	client := octopusAuth(octopusURL, APIKey, "")

	// Get specific space object
	space, err := client.Spaces.GetByName(spaceName)

	if err != nil {
		log.Println(err)
	} else {
		fmt.Println("Retrieved space " + space.Name)
	}

	return space
}

func GetEnvironment(octopusURL *url.URL, APIKey string, space *octopusdeploy.Space, EnvironmentName string) *octopusdeploy.Environment {
	client := octopusAuth(octopusURL, APIKey, space.ID)

	environment, err := client.Environments.GetByName(EnvironmentName)

	if err != nil {
		log.Println(err)
	}

	if len(environment) > 0 {
		return environment[0]
	}

	return nil
}