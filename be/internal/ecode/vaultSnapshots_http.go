package ecode

import (
	"github.com/go-dev-frame/sponge/pkg/errcode"
)

// vaultSnapshots business-level http error codes.
// the vaultSnapshotsNO value range is 1~999, if the same error code is used, it will cause panic.
var (
	vaultSnapshotsNO       = 53
	vaultSnapshotsName     = "vaultSnapshots"
	vaultSnapshotsBaseCode = errcode.HCode(vaultSnapshotsNO)

	ErrCreateVaultSnapshots     = errcode.NewError(vaultSnapshotsBaseCode+1, "failed to create "+vaultSnapshotsName)
	ErrDeleteByIDVaultSnapshots = errcode.NewError(vaultSnapshotsBaseCode+2, "failed to delete "+vaultSnapshotsName)
	ErrUpdateByIDVaultSnapshots = errcode.NewError(vaultSnapshotsBaseCode+3, "failed to update "+vaultSnapshotsName)
	ErrGetByIDVaultSnapshots    = errcode.NewError(vaultSnapshotsBaseCode+4, "failed to get "+vaultSnapshotsName+" details")
	ErrListVaultSnapshots       = errcode.NewError(vaultSnapshotsBaseCode+5, "failed to list of "+vaultSnapshotsName)

	// error codes are globally unique, adding 1 to the previous error code
)
