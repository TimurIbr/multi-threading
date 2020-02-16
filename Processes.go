package multi_threading

type Process struct {
	workerMessagesQueue MessageQueue
	NetworkLayer        *NetworkLayer
	node                int

	//TODO(middle): add context.h
	//TODO(hard): workerThread thread
	stopFlag bool
	workers  []workFunction
}
