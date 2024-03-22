package opcuaClient

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	// "strconv"
	"time"
	rtdebug "runtime/debug"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/debug"
	"github.com/gopcua/opcua/ua"
	// "github.com/gopcua/opcua/uasc"
)

func ConnectToDevice(ctx context.Context,endpoint string, enableDebug bool) (*opcua.Client, error) {
	debug.Enable = enableDebug
	// ctx := context.Background()
	c, err := opcua.NewClient(endpoint, opcua.SecurityMode(ua.MessageSecurityModeNone))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if err := c.Connect(ctx); err != nil {
		log.Println(err)
		return nil, err
	}

	return c, nil
	// defer c.Close(ctx)
}


func GetRootNode(c *opcua.Client) (*opcua.Node,error) {
	nodeId,err:=ua.ParseNodeID("i=84")
	if err != nil {
		return nil, err
	}
	return c.Node(nodeId), nil
}

func GetNodeByNodeId(nodeID string, c *opcua.Client) (*opcua.Node, error) {
	nodeId, err := ua.ParseNodeID(nodeID)
	if err != nil {
		return nil, err
	}
	return c.Node(nodeId), nil
}


func GetNodeChildren(node *opcua.Node,ctx context.Context) ([]*opcua.Node, error) {
	children,err := node.Children(ctx,0,ua.NodeClassFromString("Unspecified"))
	if err != nil {
		return nil, err
	}
	return children, nil
}


func GetNodeDisplayName(node *opcua.Node, ctx context.Context) (string, error) {
	name, err := node.DisplayName(ctx)
	if err != nil {
		return "", err
	}
	return name.Text, nil
}

func GetNodeDescripe(node *opcua.Node, ctx context.Context) (string, error) {
	description, err := node.Description(ctx)
	if err != nil {
		return "", err
	}
	return description.Text, nil
}

func ReadValueByNodeId(nodeID string, ctx context.Context, c *opcua.Client) (interface{}, error) {
	id, err := ua.ParseNodeID(nodeID)
	if err != nil {
		log.Printf("Read failed: %s\n", err)
		return nil, err
	}

	req := &ua.ReadRequest{
		MaxAge: 2000,
		NodesToRead: []*ua.ReadValueID{
			{NodeID: id},
		},
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	var resp *ua.ReadResponse
	for {
		resp, err = c.Read(ctx, req)
		if err == nil {
			break
		}

		// Following switch contains known errors that can be retried by the user.
		// Best practice is to do it on read operations.
		switch {
		case err == io.EOF && c.State() != opcua.Closed:
			// has to be retried unless user closed the connection
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSessionIDInvalid):
			// Session is not activated has to be retried. Session will be recreated internally.
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSessionNotActivated):
			// Session is invalid has to be retried. Session will be recreated internally.
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSecureChannelIDInvalid):
			// secure channel will be recreated internally.
			time.After(1 * time.Second)
			continue

		default:
			log.Printf("Read failed: %s", err)
			return nil, err
		}
	}

	if resp != nil && resp.Results[0].Status != ua.StatusOK {
		log.Printf("Status not OK: %v", resp.Results[0].Status)
	}

	// log.Printf("%#v", resp.Results[0].Value.Value())
	return resp.Results[0].Value.Value(), nil
}




func ReadMultiValueByNodeIds(nodeIDs []string,nodesReadValueIds []*ua.ReadValueID, ctx context.Context, c *opcua.Client) ([]interface{}, error) {
	var nodesToRead []*ua.ReadValueID
	if len(nodeIDs)==0{
		return []interface{}{},nil
	}
	if len(nodeIDs)>0{
		nodesToRead=[]*ua.ReadValueID{}
		for _,nd := range nodeIDs {
			id, err := ua.ParseNodeID(nd)
			if err != nil {
				log.Printf("Read failed: %s", err)
				return nil, err
			}
			nodesToRead = append(nodesToRead,&ua.ReadValueID{NodeID: id})
		}
	} else if len(nodesReadValueIds)>0 {
		nodesToRead = nodesReadValueIds
	} else{
		return nil,errors.New("请传参nodeIDs或nodesToRead")
	}
	
		
	req := &ua.ReadRequest{
		MaxAge: 2000,
		NodesToRead:nodesToRead,
		TimestampsToReturn: ua.TimestampsToReturnBoth,
	}

	var resp *ua.ReadResponse
	var err error
	for {
		resp, err = c.Read(ctx, req)
		if err == nil {
			break
		}

		// Following switch contains known errors that can be retried by the user.
		// Best practice is to do it on read operations.
		switch {
		case err == io.EOF && c.State() != opcua.Closed:
			// has to be retried unless user closed the connection
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSessionIDInvalid):
			// Session is not activated has to be retried. Session will be recreated internally.
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSessionNotActivated):
			// Session is invalid has to be retried. Session will be recreated internally.
			time.After(1 * time.Second)
			continue

		case errors.Is(err, ua.StatusBadSecureChannelIDInvalid):
			// secure channel will be recreated internally.
			time.After(1 * time.Second)
			continue

		default:
			log.Printf("Read failed: %s", err)
			return nil, err
		}
	}

	if resp != nil && resp.Results[0].Status != ua.StatusOK {
		log.Printf("Status not OK: %v\n", resp.Results[0].Status)
	}

	// log.Printf("%#v", resp.Results[0].Value.Value())
	results := []interface{}{}
	for _,r := range resp.Results {
		results = append(results, r.Value.Value())
	}
	return results, nil
}


func WriteMultiValueByNodeIds(nodeIDsWithValue []map[string]interface{}, ctx context.Context, c *opcua.Client) (err error) {
	defer func() {
		// except:
		if r := recover(); r != nil {
			e := fmt.Errorf("panic error:[%+v]\n%s", r, rtdebug.Stack())
			log.Println("WriteMultiValueByNodeIds 发生异常:", e)
			err = e
		} else {
			err = nil
		}
	}()

	nodesToWrite:=[]*ua.WriteValue{}
	for _,nv := range nodeIDsWithValue{
		for nd,v := range nv {

			id, err := ua.ParseNodeID(nd)
			if err != nil {
				log.Printf("Write failed: %s", err)
				return err
			}
			// vtw,_:=strconv.ParseFloat(v.(string),32)     // 后期加上类型判断，使用switch
			vtw := v
			variant, err := ua.NewVariant(vtw)
			if err != nil {
				log.Printf("Write failed: %s", err)
				return err
			}
			
			nodesToWrite = append(nodesToWrite,&ua.WriteValue{
				NodeID:      id,
				AttributeID: ua.AttributeIDValue,
				Value: &ua.DataValue{
					EncodingMask: ua.DataValueValue,
					Value:        variant,
				},
			})
		}

	}
	

	// id, err := ua.ParseNodeID(*nodeID)
	// if err != nil {
	// 	log.Fatalf("invalid node id: %v", err)
	// }

	// v, err := ua.NewVariant(*value)
	// if err != nil {
	// 	log.Fatalf("invalid value: %v", err)
	// }

	req := &ua.WriteRequest{
		NodesToWrite: nodesToWrite,
	}

	resp, err := c.Write(ctx, req)
	if err != nil {
		log.Fatalf("Write failed: %s", err)
	}
	log.Printf("%v", resp.Results[0])
	return nil
}