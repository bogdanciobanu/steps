package base

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	gabs "github.com/Jeffail/gabs/v2"
	"github.com/stackpulse/steps-sdk-go/env"
	"github.com/stackpulse/steps-sdk-go/log"
)

const getSessionUrl string = "api/user"
const getMembershipUrl string = "api/permissions/membership"
const getGroupsUrl string = "api/permissions/group"
const metabaseSessionHeader string = "X-Metabase-Session"

type MetabaseCommandArgs struct {
	Username     string `env:"USERNAME,required"`
	Password     string `env:"PASSWORD,required"`
	Host         string `env:"HOST,required" envDefault:"metabase.data.svc.cluster.local"`
	Requester    string `env:"REQUESTER"`
	GroupName    string `env:"GROUP_NAME,required"`
	UserId       string
	MembershipId int
	GroupId      int
}

type Output struct {
	User    string
	Message string
}

type MetabaseCommand struct {
	MetabaseCommandArgs
}

func NewMetabaseCommand(args *MetabaseCommandArgs) (*MetabaseCommand, error) {
	if err := env.Parse(args); err != nil {
		return nil, fmt.Errorf("parse args: %w", err)
	}
	return &MetabaseCommand{MetabaseCommandArgs: *args}, nil
}

func (mb *MetabaseCommand) GroupNameToId(token string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", mb.Host, getGroupsUrl), nil)
	req.Header.Add(metabaseSessionHeader, token)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Http error: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Http error: %w", err)
	}
	groups, err := gabs.ParseJSON([]byte(body))
	if err != nil {
		return fmt.Errorf("json parse error: %w", err)
	}
	log.Debug("groups list %w", groups)
	for _, group := range groups.Children() {
		group_name, ok := group.S("name").Data().(string)
		if !ok {
			return fmt.Errorf("group_name couldn't be converted to string")
		}
		if strings.ToLower(group_name) == strings.ToLower(mb.GroupName) {
			group_id, ok := group.S("id").Data().(float64)
			if !ok {
				return fmt.Errorf("id couldn't be converted to float")
			}
			mb.GroupId = int(group_id)
			break
		}

	}
	return nil
}

func (mb *MetabaseCommand) GetMembershipId(token string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", mb.Host, getMembershipUrl), nil)
	req.Header.Add(metabaseSessionHeader, token)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Http error: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Http error: %w", err)
	}
	memberships, err := gabs.ParseJSON([]byte(body))
	if err != nil {
		return fmt.Errorf("json parse error: %w", err)
	}
	log.Debug("Membership list %w", memberships)
	for _, membership := range memberships.S(mb.UserId).Children() {
		group_id, ok := membership.S("group_id").Data().(float64)
		if !ok {
			return fmt.Errorf("group_id couldn't be converted to float")
		}
		if int(group_id) == mb.GroupId {
			membership_id, ok := membership.S("membership_id").Data().(float64)
			if !ok {
				return fmt.Errorf("membership_id couldn't be converted to float")
			}
			mb.MembershipId = int(membership_id)
		}

	}
	return nil
}

func (mb *MetabaseCommand) GetUserId(token string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/", mb.Host, getSessionUrl), nil)
	req.Header.Add(metabaseSessionHeader, token)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Http error: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Http error: %w", err)
	}
	users, err := gabs.ParseJSON([]byte(body))
	if err != nil {
		return fmt.Errorf("json parse error: %w", err)
	}
	log.Debug("User list %w", users)
	for _, user := range users.Children() {
		email, ok := user.S("email").Data().(string)
		if !ok {
			return fmt.Errorf("email couldn't be converted to string")
		}
		if email == mb.Requester {
			userId, ok := user.S("id").Data().(float64)
			if !ok {
				return fmt.Errorf("userID couldn't be converted to float64")
			}
			mb.UserId = fmt.Sprintf("%g", userId)
			break
		}
	}
	return nil
}

func (mb *MetabaseCommand) GetSessionToken() (string, error) {
	postBody, err := json.Marshal(map[string]string{
		"username": mb.Username,
		"password": mb.Password,
	})
	if err != nil {
		return "", fmt.Errorf("json marshal error: %w", err)
	}
	requestBody := bytes.NewBuffer(postBody)
	resp, err := http.Post(fmt.Sprintf("%s/api/session", mb.Host), "application/json", requestBody)
	if err != nil {
		return "", fmt.Errorf("Http error: %w", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Http error: %w", err)
	}
	token, err := gabs.ParseJSON([]byte(body))
	if err != nil {
		return "", fmt.Errorf("json parse error: %w", err)
	}
	result, ok := token.S("id").Data().(string)
	if !ok {
		return "", fmt.Errorf("couldn't convert token to string")
	}
	return result, nil
}
