package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/stackpulse/steps-sdk-go/step"
	"github.com/stackpulse/steps/steps/metabase/base"
)

const outputMessage string = "Added to group"
const addUserToGroupUrl = "api/permissions/membership"
const metabaseSessionHeader string = "X-Metabase-Session"

type MetabaseAddMemberToGroup struct {
	*base.MetabaseCommand
}

func (mb *MetabaseAddMemberToGroup) Init() error {
	var args base.MetabaseCommandArgs
	baseCommand, err := base.NewMetabaseCommand(&args)
	if err != nil {
		return err
	}
	mb.MetabaseCommand = baseCommand
	return nil
}

func (mb *MetabaseAddMemberToGroup) Run() (int, []byte, error) {
	var ret []byte
	var err error
	var token string
	token, err = mb.GetSessionToken()
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get session error: %w", err)
	}
	err = mb.GroupNameToId(token)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("group name to id error: %w", err)
	}
	err = mb.GetUserId(token)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get user id error: %w", err)
	}
	client := &http.Client{}
	userId, err := strconv.Atoi(mb.UserId)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("get user id error: %w", err)
	}
	postBody, err := json.Marshal(map[string]int{
		"user_id":  userId,
		"group_id": mb.GroupId,
	})
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("json marshal error: %w", err)
	}
	requestBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", mb.Host, addUserToGroupUrl), requestBody)
	req.Header.Add(metabaseSessionHeader, token)
	req.Header.Add("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("add member to group  error: %w", err)
	}
	defer response.Body.Close()
	output := base.Output{User: mb.Requester, Message: outputMessage}
	if ret, err = json.Marshal(output); err != nil {
		return step.ExitCodeFailure, nil, fmt.Errorf("json marshal: %w", err)
	}
	return step.ExitCodeOK, ret, nil
}

func main() {
	step.Run(&MetabaseAddMemberToGroup{})
}
