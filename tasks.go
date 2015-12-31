package golph

import (
	"errors"
)

const tasksQueryPath = "api/maniphest.query"
const tasksFetchPath = "api/maniphest.info"
const tasksCreatePath = "api/maniphest.createtask"
const tasksUpdatePath = "api/maniphest.update"

// TasksService is an interface for interfacing with tasks (Maniphest)
// See: https://secure.phabricator.com/conduit/ (and search for maniphest)
type TasksService interface {
	List(*ListOptions) ([]Task, *Response, error)
	Search(*TaskSearchRequest) ([]Task, *Response, error)
	Get(string) (*Task, *Response, error)
	Create(*TaskCreateRequest) (*Task, *Response, error)
	Update(*TaskUpdateRequest) (*Response, error)
	Delete(string) (*Response, error)
}

// TasksServiceOp handles communication with the conduit methods
type TasksServiceOp struct {
	client *Client
}

var _ TasksService = &TasksServiceOp{}

// Task represents a Phabricator task.
/*
{
  "id": "5000",
  "phid": "PHID-TASK-1234",
  "authorPHID": "PHID-USER-1234",
  "ownerPHID": null,
  "ccPHIDs": [
    "PHID-USER-5000"
  ],
  "status": "resolved",
  "statusName": "Resolved",
  "isClosed": true,
  "priority": "Needs Triage",
  "priorityColor": "violet",
  "title": "Task title here",
  "description": "Task description here",
  "projectPHIDs": [
    "PHID-PROJ-1",
    "PHID-PROJ-2"
  ],
  "uri": "https://phabricator.example.com/T5000",
  "auxiliary": { },
  "objectName": "T5000",
  "dateCreated": "1415646583",
  "dateModified": "1442371259",
  "dependsOnTaskPHIDs": []
}
*/
type Task struct {
	PHID          string   `json:"phid"`
	Author        string   `json:"authorPHID"`
	Owner         string   `json:"ownerPHID"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Projects      []string `json:"projectPHIDs"`
	CCs           []string `json:"ccPHIDs"`
	Status        string   `json:"status"`
	URI           string   `json:"uri"`
	ObjectName    string   `json:"objectName"`
	StatusName    string   `json:"statusName"`
	IsClosed      bool     `json:"isClosed"`
	Priority      string   `json:"priority"`
	PriorityColor string   `json:"priorityColor"`
}

func (f Task) String() string {
	return Stringify(f)
}

type TaskGetRequest struct {
	TaskId string `form:"task_id"`
}

type TaskSearchRequest struct {
	IDs          string   `form:"ids"`
	PHIDs        string   `form:"phids"`
	OwnerPHIDs   []string `form:"ownerPHIDs"`
	AuthorPHIDs  []string `form:"authorPHIDs"`
	ProjectPHIDs []string `form:"projectPHIDs"`
	FullText     string   `form:"fullText"`
	Status       string   `form:"status"`
	Order        string   `form:"order"`
	Limit        string   `form:"limit"`
	Offset       string   `form:"offset"`
}

// TaskCreateRequest represents a request to create a Task.
type TaskCreateRequest struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	Projects    string `form:"projectPHIDs"`
	OwnerPHID   string `form:"ownerPHID"`
	CCs         string `form:"ccPHIDs"`
	Priority    string `form:"priority"`
}

// TaskUpdateRequest represents a request to create a Task.
type TaskUpdateRequest struct {
	Id          string `form:"id"`
	PHID        string `form:"phid"`
	Title       string `form:"title"`
	Description string `form:"description"`
	Projects    string `form:"projectPHIDs"`
	OwnerPHIDs  string `form:"ownerPHID"`
	CCPHIDs     string `form:"ccPHIDs"`
	Priority    string `form:"priority"`
	Comment     string `form:"comments"`
}

type SingleTaskResponse struct {
	Task      Task   `json:"result"`
	ErrorCode string `json:"error_code,omitempty"`
	ErrorInfo string `json:"error_info,omitempty"`
}

type TaskResponse struct {
	Tasks     map[string]Task `json:"result"`
	ErrorCode string          `json:"error_code,omitempty"`
	ErrorInfo string          `json:"error_info,omitempty"`
}

// Search for tasks
func (f *TasksServiceOp) Search(searchRequest *TaskSearchRequest) ([]Task, *Response, error) {
	path := tasksQueryPath

	req, err := f.client.NewRequest("POST", path, searchRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(TaskResponse)
	resp, err := f.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	var list []Task
	for _, task := range root.Tasks {
		list = append(list, task)
	}
	return list, resp, err
}

// List all tasks.
func (f *TasksServiceOp) List(opt *ListOptions) ([]Task, *Response, error) {
	path := tasksQueryPath
	path, err := addOptions(path, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := f.client.NewRequest("POST", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(TaskResponse)
	resp, err := f.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	var list []Task
	for _, task := range root.Tasks {
		list = append(list, task)
	}
	return list, resp, err
}

// Get an individual task (through the tasksFetchPath, maniphest.info).
func (f *TasksServiceOp) Get(task_id string) (*Task, *Response, error) {
	searchRequest := &TaskGetRequest{
		TaskId: task_id, // This is the "5000" part of "T5000"
	}

	req, err := f.client.NewRequest("POST", tasksFetchPath, searchRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(SingleTaskResponse)
	resp, err := f.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	if root.ErrorCode != "" {
		return nil, resp, errors.New(root.ErrorInfo)
	}

	return &root.Task, resp, err
}

// Create a task
func (f *TasksServiceOp) Create(createRequest *TaskCreateRequest) (*Task, *Response, error) {
	path := tasksCreatePath

	req, err := f.client.NewRequest("POST", path, createRequest)
	if err != nil {
		return nil, nil, err
	}

	root := new(SingleTaskResponse)
	resp, err := f.client.Do(req, root)
	if err != nil {
		return nil, resp, err
	}

	return &root.Task, resp, err
}

// Update a Task
func (f *TasksServiceOp) Update(updateRequest *TaskUpdateRequest) (*Response, error) {
	path := tasksUpdatePath

	req, err := f.client.NewRequest("POST", path, updateRequest)
	if err != nil {
		return nil, err
	}

	root := new(SingleTaskResponse)
	resp, err := f.client.Do(req, root)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// Delete a Task (Not Available Yet!)
func (f *TasksServiceOp) Delete(task_id string) (*Response, error) {
	return nil, errors.New("Delete is not available for Tasks")
}
