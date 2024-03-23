package globaldata

import (
	"earth/utils"
	"fmt"
)

type OpcNode struct {
	NodeID   string
	Name     string
	DataType string
	Value    interface{}
}

type SystemVarsDFT struct {
	CurrentValues map[int64]*OpcNode  // 0 node1, 1 node2, 2 node3...
	NodeNameSets  map[string]struct{} // set of node names
	NodeIdSets    map[string]struct{} // set of node ids  golang中没有集合  nodeid unique
	NodeIdList    []string            // list of node ids
	NodeNameIndex map[string]int64    // node name to index in CurrentValues   node name 索引
}

func NewSystemVarsDFT() *SystemVarsDFT {
	return &SystemVarsDFT{
		CurrentValues: make(map[int64]*OpcNode),
		NodeNameSets:  make(map[string]struct{}),
		NodeIdSets:    make(map[string]struct{}),
		NodeIdList:    make([]string, 0),
		NodeNameIndex: make(map[string]int64),
	}
}

func (s *SystemVarsDFT) Save() error {
	return utils.Dump(s,"systemvars.obj")
}


func (s *SystemVarsDFT) len() int {
	return len(s.NodeIdList)
}

func (s *SystemVarsDFT) AddNode(node *OpcNode) error {
	OpcWriteLock.Lock()
	defer OpcWriteLock.Unlock()
	_, nameOk := s.NodeNameSets[node.Name]
	_, idOk := s.NodeIdSets[node.NodeID]
	if !nameOk && !idOk {
		s.CurrentValues[int64(s.len())] = node
		s.NodeNameSets[node.Name] = struct{}{}
		s.NodeIdSets[node.NodeID] = struct{}{}
		s.NodeIdList = append(s.NodeIdList, node.NodeID)
		s.NodeNameIndex[node.Name] = int64(s.len())
		return nil
	}

	return fmt.Errorf("node name or node id already exists")
}

func (s *SystemVarsDFT) GetNode(id int64) (*OpcNode, error) {
	if id < 0 || id >= int64(s.len()) {
		return nil, fmt.Errorf("node id out of range")
	}
	return s.CurrentValues[id], nil
}

func (s *SystemVarsDFT) GetNodeByName(name string) (*OpcNode, error) {
	index, ok := s.NodeNameIndex[name]
	if !ok {
		return nil, fmt.Errorf("node name not found")
	}
	return s.CurrentValues[index], nil
}

func (s *SystemVarsDFT) DeleteNode(id int64) error {
	OpcWriteLock.Lock()
	defer OpcWriteLock.Unlock()
	if id < 0 || id >= int64(s.len()) {
		return fmt.Errorf("node id out of range")
	}

	node := s.CurrentValues[id]
	delete(s.NodeNameSets, node.Name)
	delete(s.NodeIdSets, node.NodeID)
	delete(s.NodeNameIndex, node.Name)
	s.NodeIdList = append(s.NodeIdList[:id], s.NodeIdList[id+1:]...)
	for i := id; i < int64(s.len()); i++ {
		s.CurrentValues[i] = s.CurrentValues[i+1]
	}
	return nil
}
