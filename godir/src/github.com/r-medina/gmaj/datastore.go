//
//  contains the API and interal functions to interact with the Key-Value store
//  that the Chord ring is providing
//

package gmaj

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/r-medina/gmaj/gmajpb"
)

var errNoDatastore = errors.New("gmaj: node does not have a datastore")

//
// External API Into Datastore
//

// Get a value in the datastore, provided an abitrary node in the ring
func Get(node *Node, key string) ([]byte, error) {
	if node == nil {
		return nil, errors.New("Node cannot be nil")
	}

	return node.get(key)
}

// Put a key/value in the datastore, provided an abitrary node in the ring.
// This is useful for testing.
func Put(node *Node, key string, val []byte, ver int32) (bool, error) {
	if node == nil {
		return false, errors.New("Node cannot be nil")
	}

	return node.put(key, val, ver)
}

// locate helps find the appropriate node in the ring.
func (node *Node) locate(key string) (*gmajpb.Node, error) {
	hashed, err := hashKey(key)
	if err != nil {
		return nil, err
	}

	return node.findSuccessor(hashed)
}

// obtainNewKeys is called when a node joins a ring and wants to request keys
// from its successor.
func (node *Node) obtainNewKeys() error {
	node.succMtx.RLock()
	succ := node.successor
	node.succMtx.RUnlock()

	// TODO(asubiotto): Test the case where there are two nodes floating around
	// that need keys.
	// Assume new predecessor has been set.
	prevPredecessor, err := node.getPredecessorRPC(succ)
	if err != nil {
		return err
	}

	return node.transferKeysRPC(
		succ, node.Id, prevPredecessor,
	) // implicitly correct even when prevPredecessor.ID == nil
}

//
// RPCs to assist with interfacing with the datastore ring
//

func (node *Node) getKey(key string) ([]byte, error) {
	node.dsMtx.RLock()
	val, ok := node.datastore[key]
	node.dsMtx.RUnlock()
	if !ok {
		return nil, errors.New("key does not exist")
	}

	node.tsMtx.RLock()
	_, ok1 := node.tempstore[key]
	node.tsMtx.RUnlock()
	if ok1 {
		return nil, nil
	}

	return val, nil
}

func (node *Node) putKeyVal(keyVal *gmajpb.KeyVal) (bool, error) {
	key := keyVal.Key
	val := keyVal.Val
	ver := keyVal.Ver

	// node.dsMtx.RLock()
	// _, exists := node.datastore[key]
	// node.dsMtx.RUnlock()
	// if exists {
	// 	return errors.New("cannot modify an existing value")
	// }

	if ver == 1 {

		node.tsMtx.Lock()
		delete(node.tempstore, key)
		node.tsMtx.Unlock()

		node.dsMtx.Lock()
		node.datastore[key] = val
		node.dsMtx.Unlock()

		return true, nil
	}

	if ver == 2 {

		node.tsMtx.Lock()
		delete(node.tempstore, key)
		node.tsMtx.Unlock()

		return true, nil
	}

	if ver == 0 {
		ret := false
		node.tsMtx.Lock()
		_, ok := node.tempstore[key]
		if ok == false {
			node.tempstore[key] = val
			ret = true
		} else {
			ret = false
		}
		node.tsMtx.Unlock()

		return ret, nil
	}

	return false, nil
}

func (node *Node) get(key string) ([]byte, error) {
	remoteNode, err := node.locate(key)
	if err != nil {
		return nil, err
	}

	// TODO(asubiotto): Smart retries on not found error. Implement channel
	// that notifies when stabilize has been called.

	// Retry on error because it might be due to temporary unavailability
	// (e.g. write happened while transferring nodes).
	val, err := node.getKeyRPC(remoteNode, key)
	if err != nil {
		<-time.After(config.RetryInterval)
		remoteNode, err = node.locate(key)
		if err != nil {
			return nil, err
		}

		val, err = node.getKeyRPC(remoteNode, key)
		if err != nil {
			return nil, err
		}
	}

	return val, nil
}

func (node *Node) put(key string, val []byte, ver int32) (bool, error) {
	remoteNode, err := node.locate(key)
	if err != nil {
		return false, err
	}

	return node.putKeyValRPC(remoteNode, key, val, ver)
}

func (node *Node) transferKeys(fromID []byte, toNode *gmajpb.Node) error {
	if idsEqual(toNode.Id, node.Id) {
		return nil
	}

	node.dsMtx.Lock()
	defer node.dsMtx.Unlock()

	toDelete := []string{}
	for key, val := range node.datastore {
		hashedKey, err := hashKey(key)
		if err != nil {
			return err
		}

		// Check that the hashed_key lies in the correct range before putting
		// the value in our predecessor.
		if betweenRightIncl(hashedKey, fromID, toNode.Id) {
			if _, err := node.putKeyValRPC(toNode, key, val, 1); err != nil {
				return err
			}

			toDelete = append(toDelete, key)
		}
	}

	for _, key := range toDelete {
		delete(node.datastore, key)
	}

	return nil
}

// DatastoreString write the contents of a node's data store to stdout.
func (node *Node) DatastoreString() (str string) {
	buf := bytes.Buffer{}

	defer func() { str = buf.String() }()

	buf.WriteString(fmt.Sprintf(
		"Node-%v datastore:", IDToString(node.Id),
	))

	const maxLen = 64

	node.dsMtx.RLock()
	if len(node.datastore) == 0 {
		defer node.dsMtx.RUnlock()
		return
	}

	for key, val := range node.datastore {
		buf.WriteString("\n")
		buf.WriteString(key)
		buf.WriteString(": ")
		if len(val) >= maxLen {
			buf.WriteString(fmt.Sprintf("%s... (truncated)", val[:maxLen]))
		} else {
			buf.WriteString(fmt.Sprintf("%s", val))
		}
	}
	node.dsMtx.RUnlock()

	return
}
