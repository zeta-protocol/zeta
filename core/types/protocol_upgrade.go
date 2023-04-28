package types

type UpgradeStatus struct {
	AcceptedReleaseInfo *ReleaseInfo
	ReadyToUpgrade      bool
}

type ReleaseInfo struct {
	ZetaReleaseTag     string
	UpgradeBlockHeight uint64
}
