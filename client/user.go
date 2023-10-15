package client

import (
	"a21hc3NpZ25tZW50/config"
	"a21hc3NpZ25tZW50/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UserClient interface {
	Login(email, password string) (respCode int, err error)
	Register(fullname, email, password, idCardPath string) (respCode int, err error)

	GetUserTaskCategory(token string) (*[]model.UserTaskCategory, error)
	// GetUserProfile(email string) (*model.User, error) 
}

type userClient struct {
}

func NewUserClient() *userClient {
	return &userClient{}
}

func (u *userClient) Login(email, password string) (respCode int, err error) {
	datajson := map[string]string{
		"email":    email,
		"password": password,
	}

	data, err := json.Marshal(datajson)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/v1/user/login"), bytes.NewBuffer(data))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	fmt.Printf("RESP: %+v\n", resp)

	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	if err != nil {
		return -1, err
	} else {
		return resp.StatusCode, nil
	}
}

func (u *userClient) Register(fullname, email, password, idCardPath string) (respCode int, err error) {
	datajson := map[string]string{
		"fullname": fullname,
		"email":    email,
		"password": password,
		"id_card":  idCardPath,
	}

	data, err := json.Marshal(datajson)
	if err != nil {
		return -1, err
	}

	req, err := http.NewRequest("POST", config.SetUrl("/api/v1/user/register"), bytes.NewBuffer(data))
	if err != nil {
		return -1, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	if err != nil {
		return -1, err
	} else {
		return resp.StatusCode, nil
	}
}

func (u *userClient) GetUserTaskCategory(token string) (*[]model.UserTaskCategory, error) {
	client, err := GetClientWithCookie(token)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", config.SetUrl("/api/v1/user/tasks"), nil)
	// fmt.Printf("REQ: %+v\n", req)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	// fmt.Printf("RESP: %+v\n", resp)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	// fmt.Printf("BODY: %+v\n", string(b))
	if err != nil {
		return nil, err
	}
	fmt.Printf("STATUS CODE: %+v\n", resp.StatusCode)
	if resp.StatusCode != 200 {
		return nil, errors.New("status code not 200")
	}

	var userTasks []model.UserTaskCategory
	err = json.Unmarshal(b, &userTasks)
	if err != nil {
		return nil, err
	}

	return &userTasks, nil
}

// func (u *userClient) GetUserProfile(email string) (*model.User, error) {
// 	client, err := GetClientWithCookie(token)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req, err := http.NewRequest("GET", config.SetUrl("/api/v1/user/profile"), nil)
// 	fmt.Printf("REQ: %+v\n", req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Set("Content-Type", "application/json")

// 	resp, err := client.Do(req)
// 	fmt.Printf("RESP: %+v\n", resp)

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer resp.Body.Close()

// 	b, err := ioutil.ReadAll(resp.Body)
// 	fmt.Printf("BODY: %+v\n", string(b))
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Printf("STATUS CODE: %+v\n", resp.StatusCode)
// 	if resp.StatusCode != 200 {
// 		return nil, errors.New("status code not 200")
// 	}

// 	var user model.User
// 	err = json.Unmarshal(b, &user)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &user, nil
// }
