package main

import (
	"algorithm-go/pkg/pinpong"
	"sync"
)

func main() {

	wg := sync.WaitGroup{}
	startWg := sync.WaitGroup{}

	n0 := pinpong.CreateNode(0, &wg, &startWg)
	n1 := pinpong.CreateNode(1, &wg, &startWg)
	n2 := pinpong.CreateNode(2, &wg, &startWg)
	n3 := pinpong.CreateNode(3, &wg, &startWg)
	n4 := pinpong.CreateNode(4, &wg, &startWg)

	pinpong.ConnectNodes(n0, n1)
	pinpong.ConnectNodes(n4, n3)
	pinpong.ConnectNodes(n3, n1)
	pinpong.ConnectNodes(n3, n2)

	n0.Start()
	n1.Start()
	n2.Start()
	n3.Start()
	n4.Start()

	wg.Wait()
}
