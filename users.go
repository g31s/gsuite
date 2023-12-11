package main

import (
	"fmt"

	admin "google.golang.org/api/admin/directory/v1"
    plugin "github.com/defensestation/pluginutils"
)

func (g *Gsuite) getUsersList(customerId string) ([]*admin.User, error){

	r, err := g.client.Users.List().Customer(customerId).OrderBy("email").Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve users in domain: %v", err)
	}

	if len(r.Users) > 0 {
		return r.Users, nil
	} else {
		fmt.Print("No users found.\n")
	}

	return nil, nil
}

func (g *Gsuite) getUsers(customerId string) (error) {
	users, err := g.getUsersList(customerId)
	if err != nil {
		return err
	}
	if users == nil || len(users) == 0 {
		return nil
	}

	graph, ok := g.plugin.GetGraph(employeeType)
	if !ok {
		return fmt.Errorf("unable to find %s graph", employeeType)
	}

	graph.SetReverseRelation(true)

	for _, user := range users {

		var aliasArr []map[string]interface{}
		if user.Aliases != nil {
			for _, alias := range user.Aliases {
				newAlias := map[string]interface{}{}
				newAlias["email"] = alias
				newAlias["alais"] = "alais"
				newAlias["supervisor_type"] = "is_alais_of"

				err := graph.AddNode(newAlias)
				if err != nil {
					return err
				}

				aliasArr = append(aliasArr, newAlias)
			}
		}

		userMapInterface := plugin.StructToMap(user)

		userMapInterface["relations"]    = aliasArr
		userMapInterface["personnel"]    = "personnel"
		userMapInterface["personnel_id"] = fmt.Sprintf("%s_%s", g.plugin.Name, user.PrimaryEmail)

		err := graph.AddNode(userMapInterface)
		if err != nil {
			return err
		}
	}

	return nil
}