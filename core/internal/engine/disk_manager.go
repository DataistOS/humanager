// core/internal/engine/disk_manager.go
package engine

type DiskManager struct {
	StorageRoot string
}

func NewDiskManager(root string) *DiskManager {
	return &DiskManager{StorageRoot: root}
}
