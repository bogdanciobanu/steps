package get

import (
	"fmt"
	"time"

	"github.com/Jeffail/gabs/v2"
	base2 "github.com/stackpulse/steps/kubectl/base"
)

const ContainerStatusesPath = "status.containerStatuses"
const InitContainerStatusesPath = "status.initContainerStatuses"

const PodNameJsonKey = "name"

var parsingConfiguration = map[string]*base2.JsonParseConfig{
	PodNameJsonKey:           {ParseFunc: base2.JsonPathStringParser, Args: []string{"metadata.name"}},
	"status":                 {ParseFunc: base2.JsonPathStringParser, Args: []string{"status.phase"}},
	"age":                    {ParseFunc: base2.JsonPathDurationFromDate, Args: []string{"metadata.creationTimestamp", time.RFC3339}},
	"IP":                     {ParseFunc: base2.JsonPathStringParser, Args: []string{"status.podIP"}},
	"IPs":                    {ParseFunc: base2.JsonPathObjectArrayKeyValue, Args: []string{"status.podIPs", "ip"}},
	"node":                   {ParseFunc: base2.JsonPathStringParser, Args: []string{"spec.nodeName"}},
	"initContainersCount":    {ParseFunc: base2.JsonPathArrayLength, Args: []string{"spec.initContainers"}},
	"containersCount":        {ParseFunc: base2.JsonPathArrayLength, Args: []string{"spec.containers"}},
	"containersRestarts":     {ParseFunc: ContainersRestartCount},
	"initContainersRestarts": {ParseFunc: InitContainersRestartCount},
	"initFinished":           {ParseFunc: InitFinished},
	"containerStatuses":      {ParseFunc: ContainerStatuses},
	"initContainerStatuses":  {ParseFunc: InitContainerStatuses},
}

type Pod struct {
	Name                   string
	Status                 string
	Age                    string
	IP                     string   `json:"IP"`
	IPs                    []string `json:"IPs"`
	Node                   string
	InitContainersCount    int
	ContainersCount        int
	ContainersRestarts     int
	InitContainersRestarts int
	InitFinished           bool
	ContainerStatuses      map[string]ContainerStatus
	InitContainerStatuses  map[string]ContainerStatus
}

type ContainerStatus struct {
	ContainerID  string `json:"containerID"`
	FinishedAt   time.Time
	Image        string
	LastState    map[string]interface{}
	Ready        bool
	RestartCount int
	StartedAt    time.Time
	State        string
}

type Pods struct {
	Items []Pod
}

// Aggregating the restarts from all the containers
func getRestartCount(path string, item *gabs.Container) int {
	containerStatuses := item.Path(path).Children()
	totalRestartCount := 0

	for _, containerStatus := range containerStatuses {
		restartCount, _ := containerStatus.S("restartCount").Data().(int)
		totalRestartCount += restartCount
	}

	return totalRestartCount
}

func ContainersRestartCount(args []string, item *gabs.Container) (interface{}, error) {
	return getRestartCount(ContainerStatusesPath, item), nil
}

func InitContainersRestartCount(args []string, item *gabs.Container) (interface{}, error) {
	return getRestartCount(InitContainerStatusesPath, item), nil
}

func isAllContainersReady(path string, item *gabs.Container) bool {
	for _, status := range item.Path(path).Children() {
		if ready, _ := status.S("ready").Data().(bool); !ready {
			return false
		}
	}
	return true
}

func InitFinished(args []string, item *gabs.Container) (interface{}, error) {
	return isAllContainersReady(InitContainerStatusesPath, item), nil
}

func getStringValInObject(item *gabs.Container, val string) string {
	for _, subItem := range item.ChildrenMap() {
		if subItem.Exists(val) {
			objVal, _ := subItem.S(val).Data().(string)
			return objVal
		}
	}
	return ""
}

func getFirstChildInObject(item *gabs.Container) string {
	for key, _ := range item.ChildrenMap() {
		return key
	}
	return ""
}

func getContainersStatuses(path string, item *gabs.Container) (map[string]map[string]interface{}, error) {
	statuses := item.Path(path).Children()
	ret := make(map[string]map[string]interface{}, len(statuses))

	for _, status := range statuses {
		if !status.Exists("name") {
			return nil, fmt.Errorf("can't find container name for status: %v", status)
		}
		containerName, ok := status.S("name").Data().(string)
		if !ok {
			return nil, fmt.Errorf("can't convert container name to string for status: %v", status)
		}
		ret[containerName] = make(map[string]interface{})
		ret[containerName]["containerID"] = status.S("containerID").Data()
		ret[containerName]["image"] = status.S("image").Data()
		ret[containerName]["lastState"] = status.S("lastState").Data()
		ret[containerName]["ready"] = status.S("ready").Data()
		ret[containerName]["restartCount"] = status.S("restartCount").Data()
		ret[containerName]["startedAt"] = getStringValInObject(status.S("state"), "startedAt")
		ret[containerName]["finishedAt"] = getStringValInObject(status.S("state"), "finishedAt")
		ret[containerName]["state"] = getFirstChildInObject(status.S("state"))
	}

	return ret, nil
}

func ContainerStatuses(args []string, item *gabs.Container) (interface{}, error) {
	return getContainersStatuses(ContainerStatusesPath, item)
}

func InitContainerStatuses(args []string, item *gabs.Container) (interface{}, error) {
	return getContainersStatuses(InitContainerStatusesPath, item)
}
