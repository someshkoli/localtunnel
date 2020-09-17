package tunnel

// Nat - mapping of connection id's with its connection details
type Nat map[string]*Connection

// Insert - adds a new tunnel connection
func (nat Nat) Insert(id string, conn *Connection) (string, *Connection) {
	nat[id] = conn
	return id, conn
}

// Lookup - looks up for available tunnel
func (nat Nat) Lookup(id string) (*Connection, bool) {
	conn, ok := nat[id]
	return conn, ok
}

// Remove - Removes nat entry from the storage
func (nat Nat) Remove(id string) {
	delete(nat, id)
}
