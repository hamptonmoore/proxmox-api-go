package cli_storage_test

import (
	"testing"

	_ "github.com/Telmate/proxmox-api-go/cli/command/commands"
	"github.com/Telmate/proxmox-api-go/proxmox"
	cliTest "github.com/Telmate/proxmox-api-go/test/cli"
	storagesubtests "github.com/Telmate/proxmox-api-go/test/cli/Storage/storage-sub-tests"
)

func Test_Storage_GlusterFS_1_Cleanup(t *testing.T) {
	storagesubtests.Cleanup("glusterfs-test-1", t)
}

func Test_Storage_GlusterFS_1_Create_Empty(t *testing.T) {
	s := storagesubtests.CloneJson(storagesubtests.GlusterfsEmpty)
	s.BackupRetention = &proxmox.ConfigStorageBackupRetention{}
	Test := cliTest.Test{
		InputJson: storagesubtests.InlineMarshal(s),
		Expected:  "(glusterfs-test-1)",
		Contains:  true,
		Args:      []string{"-i", "create", "storage", "glusterfs-test-1"},
	}
	Test.StandardTest(t)
}

func Test_Storage_GlusterFS_1_Get_Empty(t *testing.T) {
	storagesubtests.GlusterfsGetEmpty("glusterfs-test-1", t)
}

func Test_Storage_GlusterFS_1_Update_Full(t *testing.T) {
	s := storagesubtests.CloneJson(storagesubtests.GlusterfsFull)
	Test := cliTest.Test{
		InputJson: storagesubtests.InlineMarshal(s),
		Expected:  "(glusterfs-test-1)",
		Contains:  true,
		Args:      []string{"-i", "update", "storage", "glusterfs-test-1"},
	}
	Test.StandardTest(t)
}

func Test_Storage_GlusterFS_1_Get_Full(t *testing.T) {
	storagesubtests.GlusterfsGetFull("glusterfs-test-1", t)
}

func Test_Storage_GlusterFS_1_Delete(t *testing.T) {
	storagesubtests.Delete("glusterfs-test-1", t)
}
