package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	userSlice = userSlice[:0]
	userInfoData, _ := os.ReadFile("./cuda_on_time/user_info")
	lines := strings.Split(string(userInfoData), "\n")

	for _, line := range lines {
		fragments := strings.Split(line, ":")
		uName := fragments[0]
		uid := fragments[2]
		userSlice = append(userSlice, User{UserName: uName, Uid: uid})
	}
}

func GetUser() []User {
	userSliceSync.Do(initUserSlice)
	return userSlice
}

var (
	uidUnameMap  map[string]string
	uidUnameSync sync.Once
)

func initUidUserNameMap() {
	userSlice = GetUser()
	for _, user := range userSlice {
		uidUnameMap[user.Uid] = user.UserName
	}
}

func GetUserNameViaUid(uid string) string {
	uidUnameSync.Do(initUidUserNameMap)
	return uidUnameMap[uid]
}

func GetUidViaPid(pid string) string {
	cmd := "ps -p " + pid + " -o uid"
	res := ExecAndGetRes(cmd)
	lines := strings.Split(res, "\n")
	if len(lines) == 1 {
		return "none"
	}
	res = lines[1]
	uid := strings.Fields(res)[0]
	return uid
}

type CudaInfo struct {
	cudaNum     string
	pid         []string
	gpuUsage    int
	memoryUsage int
}

//func GetCudaInfoSlice() []CudaInfo {
//	cudaInfoSlice := make([]CudaInfo, 6)
//	result := ExecAndGetRes("nvidia-smi")
//	regex, _ := regexp.Compile(`[0-9]{1,2}ab%\s+Default`)
//	cudaUsages := regex.FindAllString(result, -1)
//	for cudaNum, cudaUsage := range cudaUsages {
//		cudaInfoSlice[cudaNum].cudaNum = string(rune(cudaNum))
//		cudaUsage = strings.Fields(cudaUsage)[0]
//		cudaUsageInt, _ := strconv.Atoi(cudaUsage[:len(cudaUsages)-1])
//		cudaInfoSlice[cudaNum].gpuUsage = cudaUsageInt
//	}
//
//	regex, _ = regexp.Compile(`^\d{0,5} / 24576MiB$`)
//	gRamUsages := regex.FindAllString(result, -1)
//	for cudaNum, gRamUsage := range gRamUsages {
//		gRamUsage = strings.Fields(gRamUsage)[0]
//		cudaInfoSlice[cudaNum].memoryUsage, _ = strconv.Atoi(gRamUsage[:len(gRamUsages)-3])
//	}
//
//	infoLines := strings.Split(result, "\n")
//	regex, _ = regexp.Compile(`ID\s+ID\s+Usage`)
//	startIndex := -1
//	for i, str := range infoLines {
//		if regex.MatchString(str) {
//			startIndex = i + 1
//			break
//		}
//	}
//
//	infoLines = infoLines[startIndex : len(infoLines)-1]
//}
