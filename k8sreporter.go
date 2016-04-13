package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"io/ioutil"
)

//type podSpec struct {
//	Metadata metadatast
//}
//
//type metadatast struct {
//	Annotations map[string]string
//}

type response struct {
	ResultCode int
	ResultMsg string
}

func getPatchUrl(server string, clusterName string) (string, error) {
	ns := GetEnv("MY_POD_NAMESPACE")
	name := GetEnv("MY_POD_NAME")
	if len(ns) == 0 || len(name) == 0 || len(server) == 0 {
		log.Fatalf("failed to get url:namespace:%s, podname:%s, k8sserver:%s", ns, name, server)
		return "", errors.New("failed to get k8s api server")
	}
	server = strings.Replace(server, "http://", "", 1)
	return fmt.Sprintf("http://%s/api/agent/pod/%s/%s/%s/annotation", server, clusterName, ns, name), nil
}

func patchInfo(apiAddr string, body io.Reader) error {
	req, err := http.NewRequest("PUT", apiAddr, body)
	if err != nil {
		return err
	}
//	req.Header.Add("Content-Type", "application/strategic-merge-patch+json")
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		bodystr := string(body);
		errInfo := fmt.Sprintf("report port not ok, code:%d, ret:%s", resp.StatusCode, bodystr)
		return errors.New(errInfo)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var respcontent response
	err = json.Unmarshal(data, &respcontent)
	if err != nil {
		return err
	}
	if respcontent.ResultCode != 200 {
		errInfo := fmt.Sprintf("report port not ok, code:%d, ret:%s", respcontent.ResultCode, respcontent.ResultMsg)
		return errors.New(errInfo)
	}
	return nil
}

func ReportInfos(server string, clusterName string, portEnvs map[string]string) error {

//	st := podSpec{metadatast{portEnvs}}
	body, err := json.Marshal(portEnvs)
	if err != nil {
		return err
	}

	url, err := getPatchUrl(server, clusterName)
	if err != nil {
		return err
		err.Error()
	}
	s := string(body)

	return patchInfo(url, strings.NewReader(s))
}

//func main() {
//	envs := make(map[string]string)
//	envs["conan"] = "tesdkf"
//	ReportInfos(envs)
//}

//PATCH /api/v1/namespaces/{namespace}/pods/{name}
