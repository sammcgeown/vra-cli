/*
Package types Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package types

// *** Orchestrator Types ***

// InventoryItemsList is a list of InventoryItems
type InventoryItemsList struct {
	Link []struct {
		Attributes []struct {
			Value string `json:"value"`
			Name  string `json:"name"`
		} `json:"attributes"`
		Href string `json:"href"`
		Rel  string `json:"rel"`
	} `json:"link"`
	Total int `json:"total"`
}

// WsWorkflow is a workflow
type WsWorkflow struct {
	Href             string             `json:"href"`
	Relations        Relations          `json:"relations"`
	ID               string             `json:"id"`
	CustomizedIcon   bool               `json:"customized-icon"`
	Name             string             `json:"name"`
	Version          string             `json:"version"`
	Description      string             `json:"description"`
	CategoryID       string             `json:"category-id"`
	InputParameters  []InputParameters  `json:"input-parameters"`
	OutputParameters []OutputParameters `json:"output-parameters"`
}

// OutputParameters defines the Output Parameters
type OutputParameters struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

// CategoryContext is a category context
type CategoryContext struct {
	CategoryType     string `json:"categoryType"`
	Description      string `json:"description"`
	Name             string `json:"name"`
	ParentCategoryID string `json:"parentCategoryId"`
}

// WsCategory is a category
type WsCategory struct {
	Href      string    `json:"href"`
	Relations Relations `json:"relations"`
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	PathIds   []string  `json:"path-ids"`
	Type      string    `json:"type"`
}

// Attributes is the name/value pair attributes of a resource
type Attributes struct {
	Value string `json:"value,omitempty"`
	Name  string `json:"name"`
}

// WsCategoryRequest is a category Request
type WsCategoryRequest struct {
	Name             string `json:"name"`
	Type             string `json:"type"`
	ParentCategoryID string `json:"parent-category-id"`
}

// WsAction is an Orchestrator Action
type WsAction struct {
	Href               string            `json:"href"`
	Relations          Relations         `json:"relations"`
	ID                 string            `json:"id"`
	OutputType         string            `json:"output-type"`
	Name               string            `json:"name"`
	Module             string            `json:"module"`
	Description        string            `json:"description"`
	Version            string            `json:"version"`
	Fqn                string            `json:"fqn"`
	Script             string            `json:"script"`
	BundleHasContent   bool              `json:"bundleHasContent"`
	RuntimeMemoryLimit int               `json:"runtimeMemoryLimit"`
	RuntimeTimeout     int               `json:"runtimeTimeout"`
	InputParameters    []InputParameters `json:"input-parameters"`
}

// Link is a relational link
type Link struct {
	Attributes []Attributes `json:"attributes,omitempty"`
	Href       string       `json:"href"`
	Rel        string       `json:"rel"`
}

// Relations is an array of links
type Relations struct {
	Link []Link `json:"link"`
}

// InputParameters is the input parameters for an action
type InputParameters struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	Name        string `json:"name"`
}

// WsPackages is a list of references to packages
type WsPackages struct {
	Link  []PackageLink `json:"link"`
	Start int           `json:"start"`
	Total int           `json:"total"`
}

// PackageLink is a reference to a package
type PackageLink struct {
	Attribute []Attribute `json:"attributes"`
	Href      string      `json:"href"`
	Type      string      `json:"type"`
	Rel       string      `json:"rel"`
}

// WsPackage is a package
type WsPackage struct {
	Workflows       []Workflows       `json:"workflows"`
	Actions         []Actions         `json:"actions"`
	Configurations  []Configurations  `json:"configurations"`
	Resources       []Resources       `json:"resources"`
	PolicyTemplates []PolicyTemplates `json:"policyTemplates"`
	UsedPlugins     []UsedPlugins     `json:"usedPlugins"`
	ID              string            `json:"id"`
	Href            string            `json:"href"`
	Name            string            `json:"name"`
	Description     string            `json:"description"`
}

// ExportPackageOptions represents the options for exporting a package
type ExportPackageOptions struct {
	ExportConfigurationAttributeValues      bool
	ExportConfigSecureStringAttributeValues bool
	ExportGlobalTags                        bool
	ViewContents                            bool
	AddToPackage                            bool
	EditContents                            bool
}

// Attribute is a name/value pair attribute
type Attribute struct {
	DisplayValue string `json:"displayValue"`
	Value        string `json:"value"`
	Name         string `json:"name"`
}

// Workflows is a list of workflows
type Workflows struct {
	Attribute []Attribute `json:"attribute"`
	Href      string      `json:"href"`
	Type      string      `json:"type"`
	Rel       string      `json:"rel"`
}

// Actions is a list of actions
type Actions struct {
	Attribute []Attribute `json:"attribute"`
	Href      string      `json:"href"`
	Type      string      `json:"type"`
	Rel       string      `json:"rel"`
}

// Configurations is a list of configurations
type Configurations struct {
	Attribute []Attribute `json:"attribute"`
	Href      string      `json:"href"`
	Type      string      `json:"type"`
	Rel       string      `json:"rel"`
}

// Resources is a list of resources
type Resources struct {
	Attribute []Attribute `json:"attribute"`
	Href      string      `json:"href"`
	Type      string      `json:"type"`
	Rel       string      `json:"rel"`
}

// PolicyTemplates is a list of policy templates
type PolicyTemplates struct {
	Attribute []Attribute `json:"attribute"`
	Href      string      `json:"href"`
	Type      string      `json:"type"`
	Rel       string      `json:"rel"`
}

// UsedPlugins is a used plugin definition
type UsedPlugins struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	ServerVersion string `json:"server-version"`
}
