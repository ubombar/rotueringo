package pinpong

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Node struct {
	wg               *sync.WaitGroup
	startWg          *sync.WaitGroup
	Id               int
	IncomingChannels []chan interface{}
	OutgoingChannels []chan interface{}
	ChannelIds       []int
}

func createResponse(reponses ...interface{}) []interface{} {
	return reponses
}

func (n *Node) onStart(toId int) interface{} {
	payload := fmt.Sprintf("ping: %v", n.Id)
	return payload
}

func (n *Node) onIncomingMessage(fromId int, payload interface{}) []interface{} {
	if strMessage, ok := payload.(string); ok {
		splittedMessage := strings.Split(strMessage, ":")
		protocolName := splittedMessage[0]
		packetContent := splittedMessage[1]

		switch protocolName {
		case "ping":
			// Response with pong
			fmt.Printf("Node%v: Received ping from %v\n", n.Id, packetContent)

			return createResponse(fmt.Sprintf("pong: %v", n.Id))
		case "pong":
			// Response with ping (yes infinite loop)
			return createResponse(fmt.Sprintf("ping: %v", n.Id))
		}

	}

	return nil
}

func (n *Node) run() {
	defer n.wg.Done()

	n.startWg.Done()
	fmt.Printf("Node %v Started!\n", n.Id)
	n.startWg.Wait()
	// Every node has started here!

	for i := 0; i < len(n.ChannelIds); i++ {
		adjacentChannelId := n.ChannelIds[i]
		n.OutgoingChannels[i] <- n.onStart(adjacentChannelId)
	}

	for {
		for i := 0; i < len(n.ChannelIds); i++ {
			adjacentChannelId := n.ChannelIds[i]

			select {
			case payload := <-n.IncomingChannels[i]:
				responsePayloads := n.onIncomingMessage(adjacentChannelId, payload)

				if responsePayloads != nil {
					for j := 0; j < len(responsePayloads); j++ {
						n.OutgoingChannels[i] <- responsePayloads[j]
					}
				}
			default:
				time.Sleep(time.Millisecond * 200)
			}
		}
	}
}

func (n *Node) Start() {
	go n.run()
}

func CreateNode(id int, wg, startWg *sync.WaitGroup) *Node {
	wg.Add(1)
	numberOfLinks := 8
	startWg.Add(1)

	node := &Node{
		wg:               wg,
		startWg:          startWg,
		Id:               id,
		IncomingChannels: make([]chan interface{}, 0, numberOfLinks),
		OutgoingChannels: make([]chan interface{}, 0, numberOfLinks),
		ChannelIds:       make([]int, 0),
	}
	return node
}

func ConnectNodes(node1, node2 *Node) {
	bufferSize := 10
	incomingToNode1 := make(chan interface{}, bufferSize)
	incomingToNode2 := make(chan interface{}, bufferSize)

	node1.IncomingChannels = append(node1.IncomingChannels, incomingToNode1)
	node2.IncomingChannels = append(node2.IncomingChannels, incomingToNode2)

	node1.OutgoingChannels = append(node1.OutgoingChannels, incomingToNode2)
	node2.OutgoingChannels = append(node2.OutgoingChannels, incomingToNode1)

	node1.ChannelIds = append(node1.ChannelIds, node2.Id)
	node2.ChannelIds = append(node2.ChannelIds, node1.Id)

}
