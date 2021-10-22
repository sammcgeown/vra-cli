/*
Package types Copyright 2021 VMware, Inc.
SPDX-License-Identifier: BSD-2-Clause
*/
package types

// Executions - Code Stream Execution document structure
type Executions struct {
	Project            string        `json:"project"`
	ID                 string        `json:"id"`
	Name               string        `json:"name"`
	UpdatedAt          string        `json:"updatedAt"`
	Link               string        `json:"_link"`
	UpdateTimeInMicros int64         `json:"_updateTimeInMicros"`
	CreateTimeInMicros int64         `json:"_createTimeInMicros"`
	ProjectID          string        `json:"_projectId"`
	Index              int           `json:"index"`
	Notifications      []interface{} `json:"notifications"`
	Comments           string        `json:"comments"`
	Icon               string        `json:"icon"`
	Starred            struct {
	} `json:"starred"`
	Input                 interface{}   `json:"input"`
	Output                interface{}   `json:"output"`
	StageOrder            []interface{} `json:"stageOrder"`
	Stages                interface{}   `json:"stages"`
	Status                string        `json:"status"`
	StatusMessage         string        `json:"statusMessage"`
	DurationInMicros      int           `json:"_durationInMicros"`
	TotalDurationInMicros int           `json:"_totalDurationInMicros"`
	RequestTimeInMicros   int64         `json:"_requestTimeInMicros"`
	ExecutedBy            string        `json:"_executedBy"`
	PipelineLink          string        `json:"_pipelineLink"`
	Nested                bool          `json:"_nested"`
	Rollback              bool          `json:"_rollback"`
	InputMeta             interface{}   `json:"_inputMeta"`
	OutputMeta            interface{}   `json:"_outputMeta"`
	WorkspaceResults      []struct {
		Status string   `json:"status"`
		Step   string   `json:"step"`
		Logs   []string `json:"logs"`
	} `json:"workspaceResults"`
	Tags []string `json:"tags"`
}

// VariableResponse - Code Stream API Variable response
type VariableResponse struct {
	Project            string `json:"project"`
	Kind               string `json:"kind"`
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Version            string `json:"version"`
	CreatedBy          string `json:"createdBy"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	Link               string `json:"_link"`
	UpdateTimeInMicros int64  `json:"_updateTimeInMicros"`
	CreateTimeInMicros int64  `json:"_createTimeInMicros"`
	ProjectID          string `json:"_projectId"`
	Type               string `json:"type"`
	Value              string `json:"value"`
}

// VariableRequest - Code Stream API Variable Create Request
type VariableRequest struct {
	Project     string `json:"project"`
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	Value       string `json:"value"`
}

// Pipeline - Code Stream Pipeline API
type Pipeline struct {
	Project            string `json:"project"`
	Kind               string `json:"kind"`
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	CreatedBy          string `json:"createdBy"`
	UpdatedBy          string `json:"updatedBy"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	Link               string `json:"_link"`
	UpdateTimeInMicros int64  `json:"_updateTimeInMicros"`
	CreateTimeInMicros int64  `json:"_createTimeInMicros"`
	ProjectID          string `json:"_projectId"`
	Icon               string `json:"icon"`
	Enabled            bool   `json:"enabled"`
	Concurrency        int    `json:"concurrency"`
	Input              interface {
	} `json:"input"`
	Output interface {
	} `json:"output"`
	Starred struct {
	} `json:"starred"`
	StageOrder    []string               `json:"stageOrder"`
	Stages        map[string]interface{} `json:"stages"`
	Notifications struct {
		Email   []interface{} `json:"email"`
		Jira    []interface{} `json:"jira"`
		Webhook []interface{} `json:"webhook"`
	} `json:"notifications"`
	Options   []interface{} `json:"options"`
	Workspace struct {
		Image    string        `json:"image"`
		Path     string        `json:"path"`
		Endpoint string        `json:"endpoint"`
		Cache    []interface{} `json:"cache"`
		Limits   struct {
			CPU    float64 `json:"cpu"`
			Memory int     `json:"memory"`
		} `json:"limits"`
		AutoCloneForTrigger bool `json:"autoCloneForTrigger"`
	} `json:"workspace"`
	InputMeta  interface{}   `json:"_inputMeta"`
	OutputMeta interface{}   `json:"_outputMeta"`
	Warnings   []interface{} `json:"_warnings"`
	Rollbacks  []interface{} `json:"rollbacks"`
	Tags       []string      `json:"tags"`
	State      string        `json:"state"`
}

// PipelineStage - Code Stream Pipeline Stage API
type PipelineStage struct {
	Tags      []string               `json:"tags"`
	TaskOrder []string               `json:"taskOrder"`
	Tasks     map[string]interface{} `json:"tasks"`
}

// PipelineTask - Code Stream Pipeline Task API
type PipelineTask struct {
	Configured    bool              `json:"_configured"`
	Endpoints     map[string]string `json:"endpoints"`
	IgnoreFailure bool              `json:"ignoreFailure"`
	Input         struct {
		InputProperties map[string]string `json:"inputProperties"`
		Action          string            `json:"action"`
		Blueprint       string            `json:"blueprint"`
		Name            string            `json:"name"`
		Parameters      map[string]string `json:"parameters"`
		Properties      map[string]string `json:"properties"`
		Pipeline        string            `json:"pipeline"`
	} `json:"input"`
	PreCondition string   `json:"preCondition"`
	Tags         []string `json:"tags"`
	Type         string   `json:"type"`
}

// CreateExecutionRequest - Code Stream Create Execution Request
type CreateExecutionRequest struct {
	Comments string      `json:"comments"`
	Input    interface{} `json:"input"`
}

// CreateExecutionResponse - Code Stream Create Execution Response
type CreateExecutionResponse struct {
	Comments      string      `json:"comments"`
	Source        string      `json:"source"`
	Input         interface{} `json:"input"`
	ExecutionLink string      `json:"executionLink"`
	Tags          []string    `json:"tags"`
}

// Endpoint - Code Stream Create Endpoint
type Endpoint struct {
	Project            string      `json:"project"`
	Kind               string      `json:"kind"`
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	CreatedBy          string      `json:"createdBy"`
	UpdatedBy          string      `json:"updatedBy"`
	CreatedAt          string      `json:"createdAt"`
	UpdatedAt          string      `json:"updatedAt"`
	Link               string      `json:"_link"`
	UpdateTimeInMicros int64       `json:"_updateTimeInMicros"`
	CreateTimeInMicros int64       `json:"_createTimeInMicros"`
	ProjectID          string      `json:"_projectId"`
	Type               string      `json:"type"`
	IsRestricted       bool        `json:"isRestricted"`
	Properties         interface{} `json:"properties"`
	IsLocked           bool        `json:"isLocked"`
	ValidationOutput   string      `json:"validationOutput"`
}

// CustomIntegration - Code Stream Custom Integration
type CustomIntegration struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	Version            string `json:"version"`
	CreatedBy          string `json:"createdBy"`
	UpdatedBy          string `json:"updatedBy"`
	CreatedAt          string `json:"createdAt"`
	UpdatedAt          string `json:"updatedAt"`
	Link               string `json:"_link"`
	UpdateTimeInMicros int64  `json:"_updateTimeInMicros"`
	CreateTimeInMicros int64  `json:"_createTimeInMicros"`
	Status             string `json:"status"`
	Yaml               string `json:"yaml"`
}

// PipelineImportResponse - Code Stream Pipeline Import Response
type PipelineImportResponse struct {
	Name          string `yaml:"name"`
	Status        string `yaml:"status"`
	StatusMessage string `yaml:"statusMessage"`
}

// Project - Project-Service struct
type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	OrgID          string `json:"orgId"`
	Administrators []struct {
		Email string `json:"email"`
		Type  string `json:"type"`
	} `json:"administrators"`
	Members     []interface{} `json:"members"`
	Viewers     []interface{} `json:"viewers"`
	Constraints struct {
	} `json:"constraints"`
	Properties struct {
		NamingTemplate string `json:"__namingTemplate"`
	} `json:"properties"`
	OperationTimeout int  `json:"operationTimeout"`
	SharedResources  bool `json:"sharedResources"`
}

// ProjectList - Project List
type ProjectList struct {
	Content  []Project `json:"content"`
	Pageable struct {
		Offset int `json:"offset"`
		Sort   struct {
			Sorted   bool `json:"sorted"`
			Unsorted bool `json:"unsorted"`
			Empty    bool `json:"empty"`
		} `json:"sort"`
		QueryInfo struct {
			CustomOptions struct {
			} `json:"customOptions"`
			Expand []interface{} `json:"expand"`
			Select []interface{} `json:"select"`
			Sort   struct {
				Sorted   bool `json:"sorted"`
				Unsorted bool `json:"unsorted"`
				Empty    bool `json:"empty"`
			} `json:"sort"`
		} `json:"queryInfo"`
		PageNumber int  `json:"pageNumber"`
		PageSize   int  `json:"pageSize"`
		Paged      bool `json:"paged"`
		Unpaged    bool `json:"unpaged"`
	} `json:"pageable"`
	Last          bool `json:"last"`
	TotalPages    int  `json:"totalPages"`
	TotalElements int  `json:"totalElements"`
	Sort          struct {
		Sorted   bool `json:"sorted"`
		Unsorted bool `json:"unsorted"`
		Empty    bool `json:"empty"`
	} `json:"sort"`
	Number           int  `json:"number"`
	First            bool `json:"first"`
	NumberOfElements int  `json:"numberOfElements"`
	Size             int  `json:"size"`
	Empty            bool `json:"empty"`
}

// PipelineYaml - Code Stream Pipeline YAML
type PipelineYaml struct {
	Project     string      `yaml:"project"`
	Kind        string      `yaml:"kind"`
	Name        string      `yaml:"name"`
	Icon        string      `yaml:"icon"`
	Enabled     bool        `yaml:"enabled"`
	Description string      `yaml:"description"`
	Concurrency int         `yaml:"concurrency"`
	Input       interface{} `yaml:"input"`
	InputMeta   interface{} `yaml:"_inputMeta"`
	Workspace   interface{} `yaml:"workspace"`
	StageOrder  []string    `yaml:"stageOrder"`
	Stages      interface{} `yaml:"stages"`
}

// EndpointYaml - Code Stream Endpoint YAML
type EndpointYaml struct {
	Project     string            `yaml:"project"`
	Kind        string            `yaml:"kind"`
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Type        string            `yaml:"type"`
	Properties  map[string]string `yaml:"properties"`
}
