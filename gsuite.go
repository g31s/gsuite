package main

import (
	"fmt"
	"context"

	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"

	plugin "github.com/defensestation/pluginutils"
)

type Gsuite struct {
	plugin *plugin.Plugin
	client *admin.Service
}

func NewGsuitePlugin(plugin *plugin.Plugin) *Gsuite {
	return &Gsuite{
		plugin: plugin,
	}
}

func (g *Gsuite) Run(ctx context.Context) (error) {

	err := g.plugin.ValidateOptions("json_creds", "customer_id", "admin_email")
	if err != nil {
		return err
	}
	
	jsonCreds, _ := g.plugin.GetOption("json_creds");
    customerId, _ := g.plugin.GetOption("customer_id");
    adminEmail, _ := g.plugin.GetOption("admin_email");

   err = g.setGsuiteClient(ctx,jsonCreds.(string), customerId.(string), adminEmail.(string))
   if err != nil {
   		return fmt.Errorf("Unable setup gsuite client: %v", err)
   }

   // if user is given
   if _, ok := g.plugin.GetOption("users"); ok {
   		err := g.getUsers(customerId.(string))
   		if err != nil {
   			return err
   		}
   }
   return nil
}

func (g *Gsuite) setGsuiteClient(ctx context.Context, jsonCredentials, customer_id, admin_email string) (error) {
	config, err := google.JWTConfigFromJSON(
		[]byte(jsonCredentials), 
		admin.AdminDirectoryGroupScope,
		admin.AdminDirectoryUserReadonlyScope, 
		admin.AdminDirectoryOrgunitScope,
		admin.AdminDirectoryDomainScope,
		admin.AdminDirectoryRolemanagementScope,
		admin.AdminDirectoryUserschemaScope,
		admin.AdminDirectoryDeviceMobileScope,
		admin.AdminDirectoryResourceCalendarScope,
		// reports.AdminReportsAuditReadonlyScope,
		// reports.AdminReportsUsageReadonlyScope,
	)
	if err != nil {
		return fmt.Errorf("Unable to parse client secret file to config: %v", err)
	}
	// Set your admin user email
	config.Subject = admin_email

	// define service
	srv, err := admin.New(config.Client(ctx))
	if err != nil {
		return fmt.Errorf("Unable to retrieve directory Client %v", err)
	}


	g.client = srv
	return nil
}