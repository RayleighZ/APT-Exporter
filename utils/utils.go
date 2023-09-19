package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type User struct {
	Uid      string `json:"uid"`
	UserName string `json:"user_name"`
}

func ExecAndGetRes(cmd string) string {
	command := exec.Command("bash", "-c", cmd)
	stdout, err := command.Output()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	return string(stdout)
}

var (
	userSlice     []User
	userSliceSync sync.Once
)

func initUserSlice() {
	content, err := os.ReadFile("./cuda_on_time/user.json")
	if err != nil {
		fmt.Println("Error opening file:", err)
		panic(1)
	}
	jsonString := string(content)
	err = json.Unmarshal([]byte(jsonString), &userSlice)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		panic(2)
	}
}

func GetUser() []User {
	userSliceSync.Do(initUserSlice)
	return userSlice
}
