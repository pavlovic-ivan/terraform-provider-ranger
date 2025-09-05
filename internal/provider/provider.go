// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"os"

	"github.com/g-research/ranger-go"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &rangerProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &rangerProvider{
			version: version,
		}
	}
}

// rangerProvider is the provider implementation.
type rangerProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// rangerProviderModel maps provider schema data to a Go type.
type rangerProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

// Metadata returns the provider type name.
func (p *rangerProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "ranger"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *rangerProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Interact with Ranger.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URI for Ranger. May also be provided via RANGER_HOST environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for Ranger. May also be provided via RANGER_USERNAME environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for Ranger. May also be provided via RANGER_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a Ranger client for data sources and resources.
func (p *rangerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring HashiCups client")

	// Retrieve provider data from configuration

	var config rangerProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Ranger Host",
			"The provider cannot create the Ranger client as there is an unknown configuration value for the Ranger host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the RANGER_HOST environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Ranger Username",
			"The provider cannot create the Ranger client as there is an unknown configuration value for the Ranger username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the RANGER_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Ranger API Password",
			"The provider cannot create the Ranger client as there is an unknown configuration value for the Ranger API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the RANGER_PASSWORD environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("RANGER_HOST")
	username := os.Getenv("RANGER_USERNAME")
	password := os.Getenv("RANGER_PASSWORD")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Ranger Host",
			"The provider cannot create the Ranger client as there is a missing or empty value for the Ranger host. "+
				"Set the host value in the configuration or use the RANGER_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Ranger Username",
			"The provider cannot create the Ranger client as there is a missing or empty value for the Ranger username. "+
				"Set the username value in the configuration or use the RANGER_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Ranger Password",
			"The provider cannot create the Ranger client as there is a missing or empty value for the Ranger password. "+
				"Set the password value in the configuration or use the RANGER_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Ranger client using the configuration values
	client := ranger.NewClient(host, username, password)

	ctx = tflog.SetField(ctx, "ranger_host", host)
	ctx = tflog.SetField(ctx, "ranger_username", username)
	ctx = tflog.SetField(ctx, "ranger_password", password)
	ctx = tflog.MaskAllFieldValuesStrings(ctx, "ranger_password")

	tflog.Debug(ctx, "Creating Ranger client")

	// Make the Ranger client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Ranger client", map[string]any{"success": true})

}

// DataSources defines the data sources implemented in the provider.
func (p *rangerProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewServiceDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *rangerProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPolicyResource,
	}
}
