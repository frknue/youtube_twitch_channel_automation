package video

import (
	"fmt"
	"testing"
)

// // TestConcatenateVideos function tests the ConcatenateVideos function to ensure it works as expected.
// func TestConcatenateVideos(t *testing.T) {
// 	videoFiles := []string{
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/LivelyRichCourgetteKevinTurtle-vBUCnKWIJjiCnGmp.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/ToughResilientSangGivePLZ-Xdq6YwG6pS9sG7Dq.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/IncredulousChillyRaccoonTriHard-WzAWPRV4rqF_aJ0s.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/SweetJollyBottleDxAbomb-yjgTv-mwfWwh7ep-.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/AntsyQuaintSoybeanNerfRedBlaster-iDZ5a5tE9Ge1sus8.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/MiniatureElegantBobaOSfrog-W0LPhgYJdnGgC-zN.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/DreamyThankfulAlpacaMrDestructoid-abJX18fBdEKNjB9D.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/UnsightlyAntediluvianButterYouDontSay-f8RrCE2BIcJZGqWB.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/NastyCrispyToothTebowing-w2Tq7dcfezLYgmvJ.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/CoweringNimbleSnoodOneHand-r6h5b4OGdo4GsLm_.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/BitterDifficultBarracudaShazBotstix-TS-DpzM7rzplQp_f.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/CovertFlirtyDotterelOSkomodo-cCdqfO2QxKzhXz-q.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/BovinePoliteQueleaMcaT-X7VYxXc41R33CXlx.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/HeartlessFaintNeanderthalOptimizePrime-fgc5Vtpyf2Q-V70F.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/HorribleFrailAsparagusCurseLit-i7YiwnLFkK3f57Kj.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/PoliteAltruisticBunnyFloof-OcgcIbqDMb102buf.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/SmallSquareMuleWow-aN7HRNcuHoZtilUu.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/TriumphantVictoriousCiderTBTacoLeft-8uAYZPeRc75xQhIF.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/BlindingRacyMallardPermaSmug-bNaB1KrGy2S37bZb.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/ScaryFaithfulGrassTebowing-POEy6w2C2YTtyLqb.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/ConsiderateNimbleClintFreakinStinkin-OUa004NbsUL27etw.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/SpineyBoringGoblinSaltBae-4AS3pWFwWAD4BsQz.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/SneakyEncouragingSwanBuddhaBar-fDNcOGnprdzn4FCb.mp4",
// 		"/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/ChillyResourcefulPonyBloodTrail-cE-7FAhEtt62_5yT.mp4",
// 	}

// 	outputFile := "/Users/furkanulker/git/private/youtube_twitch_channel_automation/output/d719ab16-4653-4179-95ac-1e5f92cfb10c/output.mp4"

// 	err := ConcatenateVideos(videoFiles, outputFile)

// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

func TestCreateVideoTitle(t *testing.T) {
	gameID := "538054672"

	videoTitle, videoEpisode, err := CreateVideoTitle(gameID)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Video title:", videoTitle)
	fmt.Println("Video episode:", videoEpisode)
}

func TestGetGameTitleByID(t *testing.T) {
	gameID := "538054672"
	gameTitle, err := GetGameTitleByID(gameID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Game title:", gameTitle)
}
