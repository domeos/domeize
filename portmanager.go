package main

import (
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var PORT_SCAN_FILE = [...]string{"/proc/net/raw", "/proc/net/raw6", "/proc/net/tcp", "/proc/net/tcp6",
	"/proc/net/udp", "/proc/net/udp6", "/proc/net/udplite", "/proc/net/udplite6"}

//var SPECIAL_AVAILABLE_PORTS = [...]int {80, 443, 53, 8080, 2379, 4001}
var AVAILABLE_PORT_START int = 20000
var AVAILABLE_PORT_END int = 20999

type PortSet struct {
	m map[int]bool
	sync.RWMutex
}

func NewSet() *PortSet {
	return &PortSet{
		m: map[int]bool{},
	}
}

func (s *PortSet) Add(item int) {
	s.Lock()
	defer s.Unlock()
	s.m[item] = true
}

func (s *PortSet) Remove(item int) {
	s.Lock()
	defer s.Unlock()
	delete(s.m, item)
}

func (s *PortSet) Has(item int) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[item]
	return ok
}

func (s *PortSet) Len() int {
	return len(s.List())
}

func (s *PortSet) Clear() {
	s.m = map[int]bool{}
}

func (s *PortSet) IsEmpty() bool {
	if s.Len() == 0 {
		return true
	}
	return false
}

func (s *PortSet) List() []interface{} {
	s.RLock()
	defer s.RUnlock()
	list := []interface{}{}
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

func GetAllAvailablePorts() (availablePorts []int, err error) {
	var portSet *PortSet = NewSet()
	for _, file := range PORT_SCAN_FILE {
		portInfoStr, readFileError := ioutil.ReadFile(file)
		if os.IsNotExist(readFileError) {
			continue
		}
		if readFileError != nil {
			log.Println("[Warn] Get Available Ports: ", readFileError.Error())
			return nil, errors.New("can't open proc file to check port:" + readFileError.Error())
		}
		portInfo := strings.Split(string(portInfoStr), "\n")
		for idx, line := range portInfo {
			if idx > 0 && len(line) > 3 {
				portHex := strings.Split(strings.Split(line, ":")[2], " ")[0]
				portTmp, convErr := strconv.ParseInt(portHex, 16, 32)
				port := int(portTmp)
				if convErr != nil {
					log.Println("[Warn] Get Available Ports: ", convErr.Error())
					err = errors.New("parse port error")
				} else if !portSet.Has(port) {
					portSet.Add(port)
				}
			}
		}
	}
	for aport := AVAILABLE_PORT_START; aport <= AVAILABLE_PORT_END; aport++ {
		if !portSet.Has(aport) {
			availablePorts = append(availablePorts, aport)
		}
	}
	return availablePorts, err
}

func GetAvailablePorts(portNum int) ([]int, error) {
	var ports []int
	if portNum < 0 {
		return nil, errors.New("portNum not positive:" + string(portNum))
	} else if portNum == 0 {
		return ports, nil
	}
	availablePorts, err := GetAllAvailablePorts()
	if err != nil {
		return nil, err
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < portNum; i++ {
		pos := r.Intn(len(availablePorts))
		ports = append(ports, availablePorts[pos])
	}
	return ports, nil
}
