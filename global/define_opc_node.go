package globaldata

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Brandon-lz/myopcua/utils"
)

type OpcNode struct {
	NodeID   string
	Name     string
	DataType string
	Value    interface{}
}

type OPCNodeVarsDFT struct {
	TimeStamp     time.Time
	CurrentNodes  map[int64]*OpcNode    // 0 node1, 1 node2, 2 node3...
	CurrentValues map[int64]interface{} // 0 value1, 1 value2, 2 value3...
	NodeNameSets  map[string]struct{}   // set of node names
	NodeIdSets    map[string]struct{}   // set of node ids  golang中没有集合  nodeid unique
	NodeIdList    []string              // list of node ids
	NodeNameIndex map[string]int64      // node name to index in CurrentValues   node name 索引
}

func NewGlobalOPCNodeVars() *OPCNodeVarsDFT {
	return &OPCNodeVarsDFT{
		CurrentNodes:  make(map[int64]*OpcNode),
		CurrentValues: make(map[int64]interface{}),
		NodeNameSets:  make(map[string]struct{}),
		NodeIdSets:    make(map[string]struct{}),
		NodeIdList:    make([]string, 0),
		NodeNameIndex: make(map[string]int64),
	}
}

func (s *OPCNodeVarsDFT) Save() error {
	return utils.DumpObj2Local(s, "systemvars.obj")
}

func (s *OPCNodeVarsDFT) len() int {
	return len(s.NodeIdList)
}

func (s *OPCNodeVarsDFT) AddNode(node *OpcNode) error {
	OpcReadLock.Lock()
	defer OpcReadLock.Unlock()
	_, nameOk := s.NodeNameSets[node.Name]
	_, idOk := s.NodeIdSets[node.NodeID]
	if !nameOk && !idOk {
		s.CurrentNodes[int64(len(s.NodeIdList))] = node
		s.CurrentValues[int64(len(s.NodeIdList))] = interface{}(nil)
		s.NodeNameSets[node.Name] = struct{}{}
		s.NodeIdSets[node.NodeID] = struct{}{}
		s.NodeNameIndex[node.Name] = int64(len(s.NodeIdList))
		s.NodeIdList = append(s.NodeIdList, node.NodeID) // must be last add in
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

func (s *OPCNodeVarsDFT) GetNodeByNodeId(nodeId string) (*OpcNode, error) {
	_, ok := s.NodeIdSets[nodeId]
	if !ok {
		return nil, fmt.Errorf("node name not found")
	}
	for i, ni := range s.NodeIdList {
		if ni == nodeId {
			return s.CurrentNodes[int64(i)], nil
		}
	}
	return nil, fmt.Errorf("node name not found")
}

func (s *OPCNodeVarsDFT) GetValueByName(name string) (interface{}, error) {
	index, ok := s.NodeNameIndex[name]
	if !ok {
		return nil, fmt.Errorf("node name not found")
	}
	slog.Debug(fmt.Sprintf("GetValueByName name:%s index:%d", name, index))
	slog.Debug(fmt.Sprintf("CurrentValues:%v", s.CurrentValues))
	return s.CurrentValues[index], nil
}

func (s *OPCNodeVarsDFT) DeleteNode(id int64) error {
	OpcReadLock.Lock()
	defer OpcReadLock.Unlock()
	if id < 0 || id >= int64(s.len()) {
		return fmt.Errorf("node id out of range")
	}

	node := s.CurrentNodes[id]
	delete(s.NodeNameSets, node.Name)
	delete(s.NodeIdSets, node.NodeID)
	delete(s.NodeNameIndex, node.Name)
	s.NodeIdList = append(s.NodeIdList[:id], s.NodeIdList[id+1:]...)
	delete(s.CurrentNodes, id)
	delete(s.CurrentValues, id)
	return nil
}


func (s *OPCNodeVarsDFT) DeleteNodeByNodeId(nodeId string) error {
	for i,id := range s.NodeIdList {
		if id == nodeId {
			return s.DeleteNode(int64(i))
		}
	}
	return fmt.Errorf("node id not found")
}
