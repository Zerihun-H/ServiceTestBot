package main

import "fmt"

var MainMessage string = `[ㅤ](https://hexaddis.com/tut.gif)◍◍◌◌◌◎◈  መግለጫ  ◈◎◌◌◌◍◍  

◈ ብዙም እርብሻ የሌለበት ቦታ ቢሆን ይመረጣል \!\
◈ ከሚሰጥዎት የድምፅ አብነት ጋር ለማመሳሰል ይሞክሩ \!\
◈ ከ 1 ሰኮንድ ያልብለጥ ቢሆን ተመራጭ ነው \!\


◦ ለመጀመር ከስር ያለውን አዝራር ይጫኑ ◦`

var VoiceRequestMessage string = `[ㅤ](https://hexaddis.com/voice/%d.ogg) ከታች ያለውን ቃል ደግመው ይላኩልኝ
ㅤㅤㅤ◌◎◍ \#%s\ ◍◎◌

ㅤ%s
ㅤ
`

var ThanksMessage string = `┌˚❀̥──◌─ ላምባ ──◌─❀̥˚┐
ㅤㅤ	ለትብብሮዎ እናመሰግናለን 
ㅤㅤㅤㅤ	ለመቀጠል ↴
`

var profile string = `❁✼✼✭✤✥✤✬❉❈❋✷❊✵❉
❋
✼     ጠቅላላ የድምፅ ቅጂ 
✾ㅤㅤㅤㅤㅤ↳%d
✥  ቀሪ  ያልተቀዳ ድምፅ ብዛት 
❊ㅤㅤㅤㅤㅤ↳%d
❈ㅤㅤየተጣራ የድምፅ ብዛት
❋ㅤㅤㅤㅤㅤ↳%d
✼ㅤㅤᴥ የተጋባዥ ብዛት ᴥ 
❉ㅤㅤㅤㅤㅤ↳%d
❊ㅤㅤㅤㅤ𐃫 ደረጃ 𐃫
✵ㅤㅤㅤㅤㅤ↳%d
❉`

var AlertMessage string = `[ㅤ](https://hexaddis.com/manual.jpg)እባክዎ ድምፅ ብቻ ይላኩልኝ 

ደምፅ ለመቅዳት የ ማይክራፎን መልክቱን ይጫኑ`

var BlockNotice string = `ለግዜው ስለታገዱ የ ቦቱን አስተዳደር ያናግሩ \!\

➥ [[ዋና አስተዳደር]((tg://user?id=395490182)](https://t.me/Tom201513)
➥ [[ምክትል አስተዳደር](tg://user?id=1279237180)](https://t.me/LambaSupport)`

var BlockedNotice string = `User %s Blocked By`

var WordList = []string{
	"ሰላም", "ላምባ", "ከ", "የ", "በ", "ለ", "ኛ", "ደህና", "አደርሽ", "ዋልሽ", "አመሸሽ", "ነሽ",
	"አይ", "ተዪው", "ይቅር", "አልፈልግም", "በቃኝ", "ይበቃኛል", "ቻው", "እሺ", "አዎ", "እፈልጋለሁ",
	"ቀጥይ", "ይደገም", "ድገሚ", "ድገሚልኝ", "ድገሚው", "ነው", "ላኪልኝ", "ላኪ", "ላኪው", "አንብቢልኝ",
	"አዲስ", "ያዢልኝ", "አለ", "አለኝ", "ንገሪኝ", "ሆነ", "ቀኑ", "ላይ", "ተልኮልኛል", "ነበረኝ", "ዘርዝሪሊኝ", "ዘርዝሪ",
	"አንብቢ", "ስንት", "መልዕክት", "ያዢ", "አስቀምጪ", "የዛሬው", "የዛሬ", "ሲደመር", "ሲቀነስ", "ሲካፈል", "ሲባዛ", "ደምሪ",
	"ቀንሺ", "አካፍዪ", "አባዢ", "ሰኞ", "ማክሰኞ", "ዕሮብ", "ሀሙስ", "አርብ", "ቅዳሜ", "እሁድ", "ነገ", "ትናት", "ዛሬ", "ድምጽ", "ቀንሺ",
	"አጥፊ", "ዝቅ", "አርጊ", "ጨምሪልኝ", "ስልክ", "ቁጥሩ", "ቁጥር", "ማንቂያ", "ደውል", "ሙይልኝ", "ሙይ", "ሰአት", "ንገሪኝ", "ቀን", "ቀኑን",
	"ዜሮ", "አንድ", "ሁለት", "ሶስት", "አራት", "አምስት", "ስድስት", "ሰባት", "ስምንት", "ዘጠኝ", "አስር", "አስራ", "ሃያ", "ሰላሳ", "አርባ",
	"ሃምሳ", "ስልሳ", "ሰባ", "ሰማንያ", "ዘጠና", "መቶ", "ሺ", "ሚሊየን", "የት", "ምን", "ማን", "እንዴት", "ለምን", "አይነቶች",
	"ክፍሎች", "ዝርዝሮች", "ዋጋ", "ያህል", "መች", "ትርጉም", "ትርጉሙ", "መቼ", "ማለት", "ምንድን", "ነኝ", "ያለሁበት",
	"ቦታ", "ያለሁት", "ያህል", "ይርቃል", "እዚህ", "ምን", "አዲስ", "ነገር", "ውሎ", "ዜና", "ነበር", "አየሩ", "ጃኬት",
	"ጃንጥላ", "መያዝ", "አለብኝ", "የአየር", "ሁኔታ", "መልበስ", "ደውይ", "ደውይልኝ", "ደውይለት",
	"ደውይላት", "ማስታወሻ", "ቴሌግራም", "ኢሜል", "መጽሐፍ", "ክፍል", "ታሪክ", "ምዕራፍ", "ሚስኮል", "ትዕዛዝ"}

func (s *Service) ProfileMsgBuilder(userID int64, msgID int) string {
	if _, found := s.Users[userID]; !found {
		s.CreateUser(userID, 0, msgID)
		return fmt.Sprintf(profile, 0, 0, 0, 0)
	}
	var totalVoice, totalconfirmed int
	totalVoice = len(s.Users[userID].Datasets)
	remainVoice := len(s.WordList) - totalVoice
	Invition := len(s.Users[userID].Invited)

	for _, data := range s.Users[userID].Datasets {
		if data.Confirmed {
			totalconfirmed = +1
		}
	}

	return fmt.Sprintf(profile, totalVoice, remainVoice, totalconfirmed, Invition, 21)
}
