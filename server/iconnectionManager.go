package server


type IConnectionManager interface {
	Add(connection IConnection)
	Remove(connectionID uint32)
	Get(connectionID uint32)(IConnection, error)
	Count() int
	Clean()
}
