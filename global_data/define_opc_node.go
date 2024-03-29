package globaldata

import (
	"fmt"

	"github.com/Brandon-lz/myopcua/utils"
)

type OpcNode struct {
	NodeID   string
	Name     string
	DataType string
	Value    interface{}
}

type OPCNodeVarsDFT struct {
	CurrentNodes         map[int64]*OpcNode  // 0 node1, 1 node2, 2 node3...
	CurrentValues        map[int64]interface{}  // 0 value1, 1 value2, 2 value3...
	NodeNameSets         map[string]struct{} // set of node names
	NodeIdSets           map[string]struct{} // set of node ids  golang中没有集合  nodeid unique
	NodeIdList           []string            // list of node ids
	NodeNameIndex        map[string]int64    // node name to index in CurrentValues   node name 索引
}

func NewGlobalOPCNodeVars() *OPCNodeVarsDFT {
	return &OPCNodeVarsDFT{
		CurrentNodes:  make(map[int64]*OpcNode),
		NodeNameSets:  make(map[string]struct{}),
		NodeIdSets:    make(map[string]struct{}),
		NodeIdList:    make([]string, 0),
		NodeNameIndex: make(map[string]int64),
		CurrentValues: make(map[int64]interface{}),
	}
}

func (s *OPCNodeVarsDFT) Save() error {
	return utils.Dump(s, "systemvars.obj")
}

func (s *OPCNodeVarsDFT) len() int {
	return len(s.NodeIdList)
}

func (s *OPCNodeVarsDFT) AddNode(node *OpcNode) error {
	OpcWriteLock.Lock()
	defer OpcWriteLock.Unlock()
	_, nameOk := s.NodeNameSets[node.Name]
	_, idOk := s.NodeIdSets[node.NodeID]
	if !nameOk && !idOk {
		s.CurrentNodes[int64(s.len())] = node
		s.CurrentValues[int64(s.len())] = interface{}(node.Value)
		s.NodeNameSets[node.Name] = struct{}{}
		s.NodeIdSets[node.NodeID] = struct{}{}
		s.NodeIdList = append(s.NodeIdList, node.NodeID)
		s.NodeNameIndex[node.Name] = int64(s.len())
		return nil
	}

	return fmt.Errorf("node name or node id already exists")
}

func (s *OPCNodeVarsDFT) GetNode(id int64) (*OpcNode, error) {
	if id < 0 || id >= int64(s.len()) {
		return nil, fmt.Errorf("node id out of range")
	}
	return s.CurrentNodes[id], nil
}

func (s *OPCNodeVarsDFT) GetNodeByName(name string) (*OpcNode, error) {
	index, ok := s.NodeNameIndex[name]
	if !ok {
		return nil, fmt.Errorf("node name not found")
	}
	return s.CurrentNodes[index], nil
}

func (s *OPCNodeVarsDFT) GetValueByName(name string) (interface{}, error) {
	index, ok := s.NodeNameIndex[name]
	if !ok {
		return nil, fmt.Errorf("node name not found")
	}
	return s.CurrentValues[index], nil
}


func (s *OPCNodeVarsDFT) DeleteNode(id int64) error {
	OpcWriteLock.Lock()
	defer OpcWriteLock.Unlock()
	if id < 0 || id >= int64(s.len()) {
		return fmt.Errorf("node id out of range")
	}

	node := s.CurrentNodes[id]
	delete(s.NodeNameSets, node.Name)
	delete(s.NodeIdSets, node.NodeID)
	delete(s.NodeNameIndex, node.Name)
	s.NodeIdList = append(s.NodeIdList[:id], s.NodeIdList[id+1:]...)
	for i := id; i < int64(s.len()); i++ {
		s.CurrentNodes[i] = s.CurrentNodes[i+1]
	}
	return nil
}
