package audit

import "github.com/techjanitor/pram-libs/db"

// Audit adds an action to the audit log
type Audit struct {
	User   uint
	Ib     uint
	Ip     string
	Action string
	Info   string
}

var (
	AuditNewThread       = "Thread Created"
	AuditReply           = "Replied"
	AuditNewTag          = "Tag Created"
	AuditAddTag          = "Tag Added"
	AuditCloseThread     = "Thread Closed"
	AuditOpenThread      = "Thread Opened"
	AuditStickyThread    = "Thread Stickied"
	AuditUnstickyThread  = "Thread Unstickied"
	AuditDeletePost      = "Post Deleted"
	AuditDeleteThread    = "Thread Deleted"
	AuditPurgePost       = "Post Purged"
	AuditPurgeThread     = "Thread Purged"
	AuditDeleteTag       = "Tag Deleted"
	AuditDeleteImageTag  = "Image Tag Deleted"
	AuditPurge           = "Deleted Items Purged"
	AuditFlushCache      = "Cache Flushed"
	AuditBanIp           = "IP Banned"
	AuditBanImage        = "Image Banned"
	AuditRegister        = "Account Registered"
	AuditChangePassword  = "Password Changed"
	AuditEmailUpdate     = "Email Updated"
	AuditFavoriteRemoved = "Favorite Removed"
	AuditFavoriteAdded   = "Favorite Added"
)

// Submit will insert audit info into the audit log
func (a *Audit) Submit() (err error) {

	// Get Database handle
	dbase, err := db.GetDb()
	if err != nil {
		return
	}

	// Insert data into audit table
	ps, err := dbase.Prepare("INSERT INTO audit (user_id,ib_id,audit_ip,audit_time,audit_action,audit_info) VALUES (?,?,?,NOW(),?,?)")
	if err != nil {
		return
	}
	defer ps.Close()

	_, err = ps.Exec(a.User, a.Ib, a.Ip, a.Action, a.Info)
	if err != nil {
		return
	}

	return
}