package snowflakeId

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	EPOCH      int64 = 1577836800000
	NODE_BITS  int64 = 10
	STEP_BITS  int64 = 12
	MAX_NODE   int64 = -1 ^ (-1 << NODE_BITS)
	MAX_STEP   int64 = -1 ^ (-1 << STEP_BITS)
	TIMESHIFT        = NODE_BITS + STEP_BITS
	NODE_SHIFT       = STEP_BITS
)

type Node struct {
	mu        sync.Mutex
	epoch     time.Time
	time      int64
	nodeID    int64
	step      int64
	nodeMax   int64
	nodeMask  int64
	timeShift int64
	nodeShift int64
}

func NewNode(nodeID int64) (*Node, error) {
	if nodeID < 0 && nodeID > MAX_NODE {
		return nil, fmt.Errorf("Node ID must be between 0 and %d", MAX_NODE)
	}

	return &Node{
		epoch:     time.Unix(EPOCH/1000, (EPOCH%1000)*1000000),
		nodeID:    nodeID,
		nodeMax:   MAX_NODE,
		nodeMask:  MAX_NODE << STEP_BITS,
		timeShift: TIMESHIFT,
		nodeShift: NODE_SHIFT,
	}, nil
}

func (n *Node) GenerateID() int64 {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Since(n.epoch).Milliseconds()

	if now == n.time {
		n.step = (n.step + 1) & MAX_STEP
		if n.step == 0 {
			for now <= n.time {
				now = time.Since(n.epoch).Milliseconds()
			}
		}
	} else {
		n.step = 0
	}

	id := (now << n.timeShift) | (n.nodeID << n.nodeShift) | (n.step)
	return id
}

type ID int64

func (id ID) Int64() int64 {
	return int64(id)
}
func (id ID) String(base int) (string, error) {
	if base < 2 && base > 37 {
		return "", fmt.Errorf("please select the base value between 2 - 36")
	}
	return strconv.FormatInt(int64(id), base), nil
}
