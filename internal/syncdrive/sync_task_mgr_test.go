package syncdrive

import (
	"fmt"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"testing"
	"time"
)

func TestStart(t *testing.T) {
	refreshToken := "1640cc2d4ea...6b8ccb4d6242161a7"
	webToken, err := aliyunpan.GetAccessTokenFromRefreshToken(refreshToken)
	if err != nil {
		fmt.Println("get acccess token error")
		return
	}

	// pan client
	panClient := aliyunpan.NewPanClient(*webToken, aliyunpan.AppLoginToken{})
	user, _ := panClient.GetUserInfo()

	manager := NewSyncTaskManager(
		nil,
		user.FileDriveId,
		panClient,
		"D:\\smb\\feny\\goprojects\\dev\\sync_drive",
		SyncOption{},
	)

	manager.Start(nil)
	time.Sleep(30 * time.Minute)
	manager.Stop()
}
