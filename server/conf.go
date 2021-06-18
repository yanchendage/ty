package server


var Conf = map[string]interface{}{
	"debug"				: true,

	"name" 				: "TY Server",
	"ip" 				: "127.0.0.1",
	"ipVersion" 		: "tcp4",
	"port" 				: 7726,
	// max of connection
	"maxConn" 			: 10000,
	// max of packet size
	"maxPacketSize" 	: 4096,
	// number of workerPool
	"workerPoolNumber" 	: uint32(4),

	"maxTaskQueueLen"	: 1024,

	"buffMessageQueueSize"	: 1000,

}


