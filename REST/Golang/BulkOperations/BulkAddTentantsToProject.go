package main

import (
	"fmt"
	"log"

	"net/url"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
)

func main() {

	apiURL, err := url.Parse("https://YourURL")
	if err != nil {
		log.Println(err)
	}
	APIKey := "API-YourAPIKey"
	spaceName := "Default"
	projectName := "MyProject"
	environmentNameList := []string{"Environment", "List"}
	tenantTag := "TENANT TAG TO FILTER ON"  //Format = [Tenant Tag Set Name]/[Tenant Tag] "Tenant Type/Customer"
	whatIf := false
	maxNumberOfTenants := 1
	tenantsUpdated := 0

	// Get reference to space
	space := GetSpace(apiURL, APIKey, spaceName)

	// Get project reference
	project := GetProject(apiURL, APIKey, space, projectName)

	// Get envrionment ids
	environments := []string{}
	for i := 0; i < len(environmentNameList); i++ {
		environment := GetEnvironment(apiURL, APIKey, space, environmentNameList[i])

		if nil != environment {
			environments = append(environments, environment.ID)
		}
	}

	// Get tenants
	tenants := GetTenantsByTag(apiURL, APIKey, space, tenantTag)

	// Loop through teneants
	for i := 0; i < len(tenants); i++ {
		tenantUpdated := false
		if len(tenants[i].ProjectEnvironments) == 0 {
			// Add everything
			projectEnvironments := make(map[string][]string)
			projectEnvironments[project.ID] = environments
			tenants[i].ProjectEnvironments = projectEnvironments
			tenantUpdated = true
		} else {
			projectEnvironments := tenants[i].ProjectEnvironments

			for e := 0; e < len(environments); e++ {
				if !contains(projectEnvironments[project.ID], environments[e]) {
					// Add
					existingEntries := []string{}
					existingEntries = projectEnvironments[project.ID]
					existingEntries = append(existingEntries, environments[e])
					projectEnvironments[project.ID] = existingEntries
					tenantUpdated = true
				}
			}
		}

		if tenantUpdated {
			if whatIf {
				fmt.Println(tenants[i])
			} else {
				client := octopusAuth(apiURL, APIKey, space.ID)
				client.Tenants.Update(tenants[i])
			}

			tenantsUpdated++
		}

		if maxNumberOfTenants == tenantsUpdated {
			break
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

func GetProject(octopusURL *url.URL, APIKey string, space *octopusdeploy.Space, projectName string) *octopusdeploy.Project {
	// Create client
	client := octopusAuth(octopusURL, APIKey, space.ID)

	// Get project
	project, err := client.Projects.GetByName(projectName)

	if err != nil {
		log.Println(err)
	}

	if project != nil {
		fmt.Println("Retrieved project " + project.Name)
	} else {
		fmt.Println("Project " + projectName + " not found!")
	}

	return project
}

func GetEnvironment(octopusURL *url.URL, APIKey string, space *octopusdeploy.Space, EnvironmentName string) *octopusdeploy.Environment {
	client := octopusAuth(octopusURL, APIKey, space.ID)

	environments, err := client.Environments.GetByName(EnvironmentName)

	if err != nil {
		log.Println(err)
	}

	for i := 0; i < len(environments); i++ {
		if environments[i].Name == EnvironmentName {
			return environments[i]
		}
	}

	return nil
}

func GetTenantsByTag(octopusURL *url.URL, APIKey string, space *octopusdeploy.Space, tagName string) []*octopusdeploy.Tenant {
	// Create client
	client := octopusAuth(octopusURL, APIKey, space.ID)

	tenants, err := client.Tenants.GetAll()

	if err != nil {
		log.Println(err)
	}

	tenantsWithTag := []*octopusdeploy.Tenant{}

	for i := 0; i < len(tenants); i++ {
		if contains(tenants[i].TenantTags, tagName) {
			tenantsWithTag = append(tenantsWithTag, tenants[i])
		}
	}

	return tenantsWithTag
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}