// Copyright (c) HashiCorp, Inc.

package provider

import (
	"context"
	"fmt"

	"github.com/g-research/ranger-go"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AccessModel struct {
	Type      types.String `tfsdk:"type"`
	IsAllowed types.Bool   `tfsdk:"is_allowed"`
}

type PolicyItemModel struct {
	Accesses      []AccessModel  `tfsdk:"accesses"`
	Users         []types.String `tfsdk:"users"`
	Groups        []types.String `tfsdk:"groups"`
	DelegateAdmin types.Bool     `tfsdk:"delegate_admin"`
}

type ResourceTypeModel struct {
	Values      []types.String `tfsdk:"values"`
	IsExcludes  types.Bool     `tfsdk:"is_excludes"`
	IsRecursive types.Bool     `tfsdk:"is_recursive"`
}

type ResourcesModel struct {
	Topic       *ResourceTypeModel `tfsdk:"topic"`
	Database    *ResourceTypeModel `tfsdk:"database"`
	Table       *ResourceTypeModel `tfsdk:"table"`
	URL         *ResourceTypeModel `tfsdk:"url"`
	HiveService *ResourceTypeModel `tfsdk:"hiveservice"`
	Global      *ResourceTypeModel `tfsdk:"global"`
	UDF         *ResourceTypeModel `tfsdk:"udf"`
	Column      *ResourceTypeModel `tfsdk:"column"`
}

// orderResourceModel maps the resource schema data.
type policyResourceModel struct {
	ID             types.Int64       `tfsdk:"id"`
	GUID           types.String      `tfsdk:"guid"`
	Name           types.String      `tfsdk:"name"`
	Description    types.String      `tfsdk:"description"`
	Service        types.String      `tfsdk:"service"`
	Resources      ResourcesModel    `tfsdk:"resources"`
	IsAuditEnabled types.Bool        `tfsdk:"is_audit_enabled"`
	IsEnabled      types.Bool        `tfsdk:"is_enabled"`
	Version        types.Int64       `tfsdk:"version"`
	PolicyType     types.Int64       `tfsdk:"policy_type"`
	PolicyPriority types.Int64       `tfsdk:"policy_priority"`
	IsDenyAllElse  types.Bool        `tfsdk:"is_deny_all_else"`
	ServiceType    types.String      `tfsdk:"service_type"`
	PolicyItems    []PolicyItemModel `tfsdk:"policy_items"`
}

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &policyResource{}
	_ resource.ResourceWithConfigure = &policyResource{}
)

// NewPolicyResource is a helper function to simplify the provider implementation.
func NewPolicyResource() resource.Resource {
	return &policyResource{}
}

// policyResource is the resource implementation.
type policyResource struct {
	client *ranger.Client
}

// Metadata returns the resource type name.
func (r *policyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

// Schema defines the schema for the resource.
func (r *policyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resourceSchema := schema.SingleNestedAttribute{
		Attributes: map[string]schema.Attribute{
			"values": schema.ListAttribute{
				Description: "List of resource values.",
				ElementType: types.StringType,
				Required:    true,
			},
			"is_excludes": schema.BoolAttribute{
				Description: "If true, the policy applies to all topics except those specified in the values list.",
				Default:     booldefault.StaticBool(false),
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_recursive": schema.BoolAttribute{
				Description: "If true, the policy applies recursively to all sub-topics.",
				Default:     booldefault.StaticBool(false),
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
		Optional: true,
	}

	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The ID of the policy.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"guid": schema.StringAttribute{
				Description: "The GUID of the policy.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the policy.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description of the policy.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"service": schema.StringAttribute{
				Description: "The name of the service this policy applies to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_audit_enabled": schema.BoolAttribute{
				Description: "Enable or disable audit logging for this policy.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"is_enabled": schema.BoolAttribute{
				Description: "Enable or disable this policy.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"version": schema.Int64Attribute{
				Description: "The version of the policy.",
				Computed:    true,
			},
			"policy_type": schema.Int64Attribute{
				Description: "The type of the policy. This is typically used to differentiate between different policy types in Ranger.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"policy_priority": schema.Int64Attribute{
				Description: "The priority of the policy. Policies with lower numbers are evaluated first.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"is_deny_all_else": schema.BoolAttribute{
				Description: "If true, this policy denies all other access not explicitly allowed by other policies.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"service_type": schema.StringAttribute{
				Description: "The type of service this policy applies to, such as 'kafka'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"resources": schema.SingleNestedAttribute{
				Description: "Resources to which the policy applies.",
				Attributes: map[string]schema.Attribute{
					// Kafka Resources
					"topic": resourceSchema,

					// Hive Resources
					"database":    resourceSchema,
					"table":       resourceSchema,
					"url":         resourceSchema,
					"hiveservice": resourceSchema,
					"global":      resourceSchema,
					"udf":         resourceSchema,
					"column":      resourceSchema,
				},
				Required: true,
			},
			"policy_items": schema.ListNestedAttribute{
				Description: "List of policy items that define the access controls for this policy.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"accesses": schema.ListNestedAttribute{
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										Required:    true,
										Description: "The type of access such as 'publish', 'consume', etc.",
									},
									"is_allowed": schema.BoolAttribute{
										Optional:    true,
										Computed:    true,
										Description: "If true, the access is allowed; if false, it is denied.",
										Default:     booldefault.StaticBool(true),
										PlanModifiers: []planmodifier.Bool{
											boolplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
							Description: "List of accesses that define the permissions granted by this policy item.",
							Required:    true,
							Validators: []validator.List{
								listvalidator.SizeAtLeast(1),
							},
						},
						"users": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "List of users to which this policy applies.",
							PlanModifiers: []planmodifier.List{
								listplanmodifier.UseStateForUnknown(),
							},
						},
						"groups": schema.ListAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "List of groups to which this policy applies.",
							PlanModifiers: []planmodifier.List{
								listplanmodifier.UseStateForUnknown(),
							},
						},
						"delegate_admin": schema.BoolAttribute{
							Optional:    true,
							Computed:    true,
							Description: "If true, allows the user to delegate admin privileges.",
							Default:     booldefault.StaticBool(false),
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
		},
	}
}

func ConvertModelToPolicy(plan *policyResourceModel) *ranger.Policy {
	// Generate API request body from plan
	policy := ranger.Policy{
		Name:           plan.Name.ValueString(),
		Service:        plan.Service.ValueString(),
		Resources:      ranger.Resources{},
		IsAuditEnabled: plan.IsAuditEnabled.ValueBool(),
		IsEnabled:      plan.IsEnabled.ValueBool(),
		PolicyType:     int(plan.PolicyType.ValueInt64()),
		PolicyPriority: int(plan.PolicyPriority.ValueInt64()),
		IsDenyAllElse:  plan.IsDenyAllElse.ValueBool(),
	}

	// Set optional fields
	if plan.Description.ValueString() != "" {
		policy.Description = plan.Description.ValueString()
	} else {
		policy.Description = ""
	}

	if plan.ServiceType.ValueString() != "" {
		policy.ServiceType = plan.ServiceType.ValueString()
	}

	// Populate Resources directly from the plan
	if plan.Resources.Topic != nil {
		policy.Resources.Topic = &ranger.ResourceType{
			Values:      make([]string, len(plan.Resources.Topic.Values)),
			IsExcludes:  plan.Resources.Topic.IsExcludes.ValueBool(),
			IsRecursive: plan.Resources.Topic.IsRecursive.ValueBool(),
		}
		for i, value := range plan.Resources.Topic.Values {
			policy.Resources.Topic.Values[i] = value.ValueString()
		}
	}

	if plan.Resources.Database != nil {
		policy.Resources.Database = &ranger.ResourceType{
			Values:      make([]string, len(plan.Resources.Database.Values)),
			IsExcludes:  plan.Resources.Database.IsExcludes.ValueBool(),
			IsRecursive: plan.Resources.Database.IsRecursive.ValueBool(),
		}
		for i, value := range plan.Resources.Database.Values {
			policy.Resources.Database.Values[i] = value.ValueString()
		}
	}

	if plan.Resources.Table != nil {
		policy.Resources.Table = &ranger.ResourceType{
			Values:      make([]string, len(plan.Resources.Table.Values)),
			IsExcludes:  plan.Resources.Table.IsExcludes.ValueBool(),
			IsRecursive: plan.Resources.Table.IsRecursive.ValueBool(),
		}
		for i, value := range plan.Resources.Table.Values {
			policy.Resources.Table.Values[i] = value.ValueString()
		}
	}

	if plan.Resources.URL != nil {
		policy.Resources.URL = &ranger.ResourceType{
			Values:      make([]string, len(plan.Resources.URL.Values)),
			IsExcludes:  plan.Resources.URL.IsExcludes.ValueBool(),
			IsRecursive: plan.Resources.URL.IsRecursive.ValueBool(),
		}
		for i, value := range plan.Resources.URL.Values {
			policy.Resources.URL.Values[i] = value.ValueString()
		}
	}

	if plan.Resources.HiveService != nil {
		policy.Resources.HiveService = &ranger.ResourceType{
			Values:      make([]string, len(plan.Resources.HiveService.Values)),
			IsExcludes:  plan.Resources.HiveService.IsExcludes.ValueBool(),
			IsRecursive: plan.Resources.HiveService.IsRecursive.ValueBool(),
		}
		for i, value := range plan.Resources.HiveService.Values {
			policy.Resources.HiveService.Values[i] = value.ValueString()
		}
	}

	if plan.Resources.Global != nil {
		policy.Resources.Global = &ranger.ResourceType{
			Values:      make([]string, len(plan.Resources.Global.Values)),
			IsExcludes:  plan.Resources.Global.IsExcludes.ValueBool(),
			IsRecursive: plan.Resources.Global.IsRecursive.ValueBool(),
		}
		for i, value := range plan.Resources.Global.Values {
			policy.Resources.Global.Values[i] = value.ValueString()
		}
	}

	if plan.Resources.UDF != nil {
		policy.Resources.UDF = &ranger.ResourceType{
			Values:      make([]string, len(plan.Resources.UDF.Values)),
			IsExcludes:  plan.Resources.UDF.IsExcludes.ValueBool(),
			IsRecursive: plan.Resources.UDF.IsRecursive.ValueBool(),
		}
		for i, value := range plan.Resources.UDF.Values {
			policy.Resources.UDF.Values[i] = value.ValueString()
		}
	}

	if plan.Resources.Column != nil {
		policy.Resources.Column = &ranger.ResourceType{
			Values:      make([]string, len(plan.Resources.Column.Values)),
			IsExcludes:  plan.Resources.Column.IsExcludes.ValueBool(),
			IsRecursive: plan.Resources.Column.IsRecursive.ValueBool(),
		}
		for i, value := range plan.Resources.Column.Values {
			policy.Resources.Column.Values[i] = value.ValueString()
		}
	}

	// Populate PolicyItems
	if plan.PolicyItems != nil {
		policy.PolicyItems = make([]ranger.PolicyItem, len(plan.PolicyItems))
		for i, item := range plan.PolicyItems {
			policyItem := ranger.PolicyItem{
				DelegateAdmin: item.DelegateAdmin.ValueBool(),
			}
			if item.Users != nil {
				policyItem.Users = make([]string, len(item.Users))
				for j, user := range item.Users {
					policyItem.Users[j] = user.ValueString()
				}
			}
			if item.Groups != nil {
				policyItem.Groups = make([]string, len(item.Groups))
				for j, group := range item.Groups {
					policyItem.Groups[j] = group.ValueString()
				}
			} else {
				policyItem.Groups = nil // Ensure Groups is nil if not provided
			}
			if item.Accesses != nil {
				policyItem.Accesses = make([]ranger.Access, len(item.Accesses))
				for j, access := range item.Accesses {
					policyItem.Accesses[j] = ranger.Access{
						Type:      access.Type.ValueString(),
						IsAllowed: access.IsAllowed.ValueBool(),
					}
				}
			}

			policy.PolicyItems[i] = policyItem
		}
	} else {
		policy.PolicyItems = nil // Ensure it's set to nil if no items
	}

	return &policy
}

func ConvertPolicyToModel(policy *ranger.Policy) *policyResourceModel {
	// Overwrite items with refreshed state
	model := policyResourceModel{
		ID:             types.Int64Value(int64(policy.ID)),
		Name:           types.StringValue(policy.Name),
		IsEnabled:      types.BoolValue(policy.IsEnabled),
		Version:        types.Int64Value(int64(policy.Version)),
		Service:        types.StringValue(policy.Service),
		GUID:           types.StringValue(policy.GUID),
		IsAuditEnabled: types.BoolValue(policy.IsAuditEnabled),
		PolicyType:     types.Int64Value(int64(policy.PolicyType)),
		PolicyPriority: types.Int64Value(int64(policy.PolicyPriority)),
		IsDenyAllElse:  types.BoolValue(policy.IsDenyAllElse),
		ServiceType:    types.StringValue(policy.ServiceType),
		Resources:      ResourcesModel{},
		PolicyItems:    make([]PolicyItemModel, len(policy.PolicyItems)),
	}

	if policy.Description != "" {
		model.Description = types.StringValue(policy.Description)
	} else {
		model.Description = types.StringNull()
	}

	if policy.Resources.Topic != nil {
		model.Resources.Topic = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.Topic.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.Topic.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.Topic.IsRecursive),
		}
		for i, value := range policy.Resources.Topic.Values {
			model.Resources.Topic.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.Database != nil {
		model.Resources.Database = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.Database.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.Database.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.Database.IsRecursive),
		}
		for i, value := range policy.Resources.Database.Values {
			model.Resources.Database.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.Table != nil {
		model.Resources.Table = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.Table.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.Table.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.Table.IsRecursive),
		}
		for i, value := range policy.Resources.Table.Values {
			model.Resources.Table.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.URL != nil {
		model.Resources.URL = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.URL.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.URL.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.URL.IsRecursive),
		}
		for i, value := range policy.Resources.URL.Values {
			model.Resources.URL.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.HiveService != nil {
		model.Resources.HiveService = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.HiveService.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.HiveService.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.HiveService.IsRecursive),
		}
		for i, value := range policy.Resources.HiveService.Values {
			model.Resources.HiveService.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.Global != nil {
		model.Resources.Global = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.Global.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.Global.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.Global.IsRecursive),
		}
		for i, value := range policy.Resources.Global.Values {
			model.Resources.Global.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.UDF != nil {
		model.Resources.UDF = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.UDF.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.UDF.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.UDF.IsRecursive),
		}
		for i, value := range policy.Resources.UDF.Values {
			model.Resources.UDF.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.Column != nil {
		model.Resources.Column = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.Column.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.Column.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.Column.IsRecursive),
		}
		for i, value := range policy.Resources.Column.Values {
			model.Resources.Column.Values[i] = types.StringValue(value)
		}
	}

	// Populate PolicyItems
	if policy.PolicyItems != nil {
		model.PolicyItems = make([]PolicyItemModel, len(policy.PolicyItems))
		for i, item := range policy.PolicyItems {
			model.PolicyItems[i] = PolicyItemModel{
				DelegateAdmin: types.BoolValue(item.DelegateAdmin),
			}
			if item.Users != nil {
				model.PolicyItems[i].Users = make([]types.String, len(item.Users))
				for j, user := range item.Users {
					model.PolicyItems[i].Users[j] = types.StringValue(user)
				}
			}

			if item.Groups != nil {
				model.PolicyItems[i].Groups = make([]types.String, len(item.Groups))
				for j, group := range item.Groups {
					model.PolicyItems[i].Groups[j] = types.StringValue(group)
				}
			}
			if item.Accesses != nil {
				model.PolicyItems[i].Accesses = make([]AccessModel, len(item.Accesses))
				for j, access := range item.Accesses {
					model.PolicyItems[i].Accesses[j] = AccessModel{
						Type:      types.StringValue(access.Type),
						IsAllowed: types.BoolValue(access.IsAllowed),
					}
				}
			}
		}
	} else {
		model.PolicyItems = nil // Ensure it's set to nil if no items
	}

	return &model
}

func UpdatePlanFromPolicy(plan *policyResourceModel, policy *ranger.Policy, resp *resource.UpdateResponse) {
	// Map response body to schema and populate Computed attribute values
	plan.ID = types.Int64Value(int64(policy.ID))
	plan.GUID = types.StringValue(policy.GUID)
	plan.Name = types.StringValue(policy.Name)
	plan.Service = types.StringValue(policy.Service)
	plan.Version = types.Int64Value(int64(policy.Version))
	plan.IsEnabled = types.BoolValue(policy.IsEnabled)
	plan.IsAuditEnabled = types.BoolValue(policy.IsAuditEnabled)
	plan.PolicyType = types.Int64Value(int64(policy.PolicyType))
	plan.PolicyPriority = types.Int64Value(int64(policy.PolicyPriority))
	plan.IsDenyAllElse = types.BoolValue(policy.IsDenyAllElse)
	plan.Resources = ResourcesModel{}

	// Populate optional fields
	if policy.Description != "" {
		plan.Description = types.StringValue(policy.Description)
	} else {
		plan.Description = types.StringNull()
	}

	if policy.ServiceType != "" {
		plan.ServiceType = types.StringValue(policy.ServiceType)
	} else {
		plan.ServiceType = types.StringNull()
	}

	// Populate Resources directly from the plan
	if policy.Resources.Topic != nil {
		plan.Resources.Topic = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.Topic.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.Topic.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.Topic.IsRecursive),
		}
		for i, value := range policy.Resources.Topic.Values {
			plan.Resources.Topic.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.Database != nil {
		plan.Resources.Database = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.Database.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.Database.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.Database.IsRecursive),
		}
		for i, value := range policy.Resources.Database.Values {
			plan.Resources.Database.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.Table != nil {
		plan.Resources.Table = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.Table.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.Table.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.Table.IsRecursive),
		}
		for i, value := range policy.Resources.Table.Values {
			plan.Resources.Table.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.URL != nil {
		plan.Resources.URL = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.URL.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.URL.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.URL.IsRecursive),
		}
		for i, value := range policy.Resources.URL.Values {
			plan.Resources.URL.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.HiveService != nil {
		plan.Resources.HiveService = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.HiveService.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.HiveService.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.HiveService.IsRecursive),
		}
		for i, value := range policy.Resources.HiveService.Values {
			plan.Resources.HiveService.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.Global != nil {
		plan.Resources.Global = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.Global.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.Global.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.Global.IsRecursive),
		}
		for i, value := range policy.Resources.Global.Values {
			plan.Resources.Global.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.UDF != nil {
		plan.Resources.UDF = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.UDF.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.UDF.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.UDF.IsRecursive),
		}
		for i, value := range policy.Resources.UDF.Values {
			plan.Resources.UDF.Values[i] = types.StringValue(value)
		}
	}

	if policy.Resources.Column != nil {
		plan.Resources.Column = &ResourceTypeModel{
			Values:      make([]types.String, len(policy.Resources.Column.Values)),
			IsExcludes:  types.BoolValue(policy.Resources.Column.IsExcludes),
			IsRecursive: types.BoolValue(policy.Resources.Column.IsRecursive),
		}
		for i, value := range policy.Resources.Column.Values {
			plan.Resources.Column.Values[i] = types.StringValue(value)
		}
	}

	// Populate PolicyItems
	if policy.PolicyItems != nil {
		plan.PolicyItems = make([]PolicyItemModel, len(policy.PolicyItems))
		for i, item := range policy.PolicyItems {
			plan.PolicyItems[i] = PolicyItemModel{
				DelegateAdmin: types.BoolValue(item.DelegateAdmin),
				Accesses:      make([]AccessModel, len(item.Accesses)),
			}
			if item.Users != nil {
				plan.PolicyItems[i].Users = make([]types.String, len(item.Users))
				for j, user := range item.Users {
					plan.PolicyItems[i].Users[j] = types.StringValue(user)
				}
			}
			if item.Groups != nil {
				plan.PolicyItems[i].Groups = make([]types.String, len(item.Groups))
				for j, group := range item.Groups {
					plan.PolicyItems[i].Groups[j] = types.StringValue(group)
				}
			}
			for j, access := range item.Accesses {
				plan.PolicyItems[i].Accesses[j] = AccessModel{
					Type:      types.StringValue(access.Type),
					IsAllowed: types.BoolValue(access.IsAllowed),
				}
			}
		}
	} else {
		plan.PolicyItems = nil // Ensure it's set to nil if no items
	}
}

// Create creates the resource and sets the initial Terraform state.
// Create a new resource.
func (r *policyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan policyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyToCreate := ConvertModelToPolicy(&plan)

	// Create new policy
	policy, err := r.client.CreatePolicy(policyToCreate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating policy",
			"Could not create policy, unexpected error: "+err.Error(),
		)
		return
	}

	UpdatePlanFromPolicy(&plan, policy, nil)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *policyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state policyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed policy from Ranger
	policy, err := r.client.GetPolicy(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Ranger policy",
			"Could not read policy ID "+state.ID.String()+": "+err.Error(),
		)
		return
	}

	state = *ConvertPolicyToModel(policy)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *policyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan policyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	policyToUpdate := ConvertModelToPolicy(&plan)

	var state policyResourceModel
	diags = req.State.Get(ctx, &state)

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	policyToUpdate.ID = int(state.ID.ValueInt64())

	// Update existing policy
	policy, err := r.client.UpdatePolicy(policyToUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating Ranger policy",
			"Could not update policy ID "+plan.ID.String()+": "+err.Error(),
		)
		return
	}

	UpdatePlanFromPolicy(&plan, policy, resp)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *policyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state policyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing policy
	err := r.client.DeletePolicy(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Ranger policy",
			"Could not delete policy, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *policyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*ranger.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *ranger.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}
