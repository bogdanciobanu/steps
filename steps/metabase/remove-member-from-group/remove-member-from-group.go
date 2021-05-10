package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/steps/metabase/base"
)

const deleteUrl string = "api/permissions/membership"
const method string = "DELETE"
const metabaseSessionHeader string = "X-Metabase-Session"
const outputMessage string = "Removed member from group"

type MetabaseRemoveMemberFromGroup struct {
	*base.MetabaseCommand
}

func (mb *MetabaseRemoveMemberFromGroup) Init() error {
	var args base.MetabaseCommandArgs
	baseCommand, err := base.NewMetabaseCommand(&args)
	if err != nil {
		return err
	}
	mb.MetabaseCommand = baseCommand
	return nil
}

func (mb *MetabaseRemoveMemberFromGroup) Run() (int, []byte, error) {
	var ret []byte
	var err error
	var token string
	token, err = mb.GetSessionToken()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("Get session error: %w", err)
	}
	err = mb.GroupNameToId(token)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("Group name to id error: %w", err)
	}
	err = mb.GetUserId(token)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("Get user id error: %w", err)
	}
	err = mb.GetMembershipId(token)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("Get membership id error: %w", err)
	}
	client := &http.Client{}
	deleteBody, err := json.Marshal(map[string]int{
		"group_id": mb.GroupId,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("json marshal: %w", err)
	}
	requestBody := bytes.NewBuffer(deleteBody)
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s/%d", mb.Host, deleteUrl, mb.MembershipId), requestBody)
	req.Header.Add(metabaseSessionHeader, token)
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("Remove member from group error: %w", err)
	}
	defer response.Body.Close()
	output := base.Output{User: mb.Requester, Message: outputMessage}
	if ret, err = json.Marshal(output); err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("json marshal: %w", err)
	}
	return step.ExitCodeOK, ret, nil
}

func main() {
	step.Run(&MetabaseRemoveMemberFromGroup{})
}
