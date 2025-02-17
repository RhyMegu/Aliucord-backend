package modules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Aliucord/Aliucord-backend/common"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func init() {
	modules = append(modules, initAutoReplies)
}

func r(regex string) *regexp.Regexp {
	return regexp.MustCompile("(?i)" + regex)
}

var lastReplyCache = map[discord.UserID]string{}

func reply(msg *gateway.MessageCreateEvent, reply string) {
	if lastReplyCache[msg.Author.ID] == reply {
		return
	}
	_, err := s.SendTextReply(msg.ChannelID, reply, msg.ID)
	logger.LogIfErr(err)
	lastReplyCache[msg.Author.ID] = reply
}

const (
	JustAsk          = "https://dontasktoask.com/"
	CheckThePins     = "<a:checkpins:859804429536198676>"
	MentionHelp      = "Rule 9: Don't dm or mention for support"
	ElaborateHelp    = "We can't help you if you don't tell us your issue. "
	InstallPlugins   = "https://cdn.discordapp.com/attachments/811261298997460992/875552420363636766/21-08-13-03-31-20.mp4" //people would rather watch a video than opening the docu so
	InstallThemes    = "Themer plugin <https://discord.com/channels/811255666990907402/811261298997460992/845243103757467658>"
	CreateThemes     = "Read this documentation: https://github.com/Aliucord/documentation/tree/main/theme-dev"
	SearchThemes     = "Check <#824357609778708580> and search on there, maybe there's the theme you want"
	AliuCrash        = "Send crashlogs (check Crashes in Settings, and copy the most recent), if there aren't any crashlogs, then we can't do nothing about it. (If you want, you can try send a logcat, check <https://pastebin.com/pNhXwhrd>)"
	FreeNitro        = "Not possible. Nitrospoof exists for \"free\" emotes, for anything else buy nitro."
	Usage            = "Go to the plugin's repository and read the readme. Chances are the dev added a description."
	BetterInternet   = "This happens when you have an old/misbehaving router. Use mobile data (~120mb usage) or maybe a VPN (*or just get better internet*)."
	PluginDownloader = "PluginDownloader is now a part of Aliucord. (It won't be present in the plugin list) If the option to download plugins is still missing, update Aliucord."
)

func initAutoReplies() {
	cfg := config.AutoReplyConfig
	if !cfg.Enabled {
		return
	}

	PRD := fmt.Sprintf("%s 👉 <#%s>", CheckThePins, cfg.PRD)
	FindPlugin := fmt.Sprintf("Search in <#%s> and <#%s>. If it doesn't exist, then %s in <#%s>",
		cfg.PluginsList, cfg.NewPlugins, CheckThePins, cfg.PRD)

	autoRepliesString := map[string]string{
		"a plugin to":           PRD,
		"can you make":          PRD,
		"how do i use":          Usage,
		"free nitro":            FreeNitro,
		"handshake exception":   BetterInternet,
		"connection terminated": BetterInternet,
	}

	autoRepliesRegex := map[*regexp.Regexp]string{
		r("^(?:i need )?help(?: me)?$"):               ElaborateHelp,
		r("<@!?\\d{2,19}> help"):                      MentionHelp,
		r("help <@!?\\d{2,19}>"):                      MentionHelp,
		r("animated (profile|avatar|pfp)"):            FreeNitro,
		r("^is there a plugin .+"):                    FindPlugin,
		r("^where(?: i)?s(?: the )?.+ plugin$"):       FindPlugin,
		r("^can (?:anyone|you) help(?: me)?\\??$"):    JustAsk,
		r("can'?t download plugin ?downloader"):       PluginDownloader,
		r("where(?: i)s(?: the)? plugin ?downloader"): PluginDownloader,
		r("how (to|do I|do you) install plugins"):     InstallPlugins,
		r("how (to|do I|do you) install themes"):      InstallThemes,
		r("how (to|do I|do you|can i|) create themes"):CreateThemes,
		r("(does anyone know|is there) a theme that"): SearchThemes,
		r("my aliucord (crashed|keeps crashing|crash|crashes)"): AliuCrash,
	}

	s.AddHandler(func(msg *gateway.MessageCreateEvent) {
		if msg.Member == nil || len(msg.Attachments) > 0 || (msg.ReferencedMessage != nil && msg.ReferencedMessage.Author.ID == msg.Author.ID) || msg.Author.Bot || strings.HasPrefix(msg.Content, "Quick Aliucord ") {
			return
		}

		c, err := s.Channel(msg.ChannelID)
		if err == nil {
			if c.ID != cfg.PRD && c.ParentID != cfg.SupportCategory {
				return
			}
		} else {
			logger.Println(err)
		}

		for _, role := range msg.Member.RoleIDs {
			if common.HasRole(cfg.IgnoredRoles, role) {
				return
			}
		}

		for regex, value := range autoRepliesRegex {
			if regex.MatchString(msg.Content) {
				reply(msg, value)
				return
			}
		}

		content := strings.ToLower(msg.Content)
		for trigger, value := range autoRepliesString {
			if strings.Contains(content, trigger) {
				reply(msg, value)
				return
			}
		}
	})
}
