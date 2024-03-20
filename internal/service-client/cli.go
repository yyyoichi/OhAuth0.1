package serviceclient

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/yyyoichi/OhAuth0.1/internal/database"
)

type (
	Brawser struct {
		codeReceiver      CodeReceiver
		accessTokenClient AccessTokenClient
		resourceClient    ResourceClient

		mu                     sync.Mutex
		currentServiceClientId *string
		accessTokens           map[string]string
		refreshTokens          map[string]string
	}
	BrawserConfig struct {
		RedirectURI       string
		AuthServerURI     string
		ResourceServerURI string
	}
)

func NewBrawser(config BrawserConfig) *Brawser {
	var b Brawser
	b.codeReceiver = NewCodeReceiver(config.RedirectURI)
	b.accessTokenClient = NewAccessTokenClient(config.AuthServerURI)
	b.resourceClient = NewResourceClient(config.ResourceServerURI)

	b.mu = sync.Mutex{}
	b.currentServiceClientId = nil
	b.accessTokens = map[string]string{}
	b.refreshTokens = map[string]string{}
	return &b
}

func (b *Brawser) Brawse(ctx context.Context, input string) (*output, error) {
	command := ParseCommand(input)
	switchedAnySite := b.currentServiceClientId == nil

	switch command.command {
	case help:
		return helpOutput, nil
	case status:
		return newStatusOutput(b.accessTokens, *b.currentServiceClientId), nil
	case showSites:
		return showSitesOutput, nil
	case switchsite:
		id := command.args[0]
		if id != "500" && id != "501" {
			return nil, errors.New("unknown site id")
		}
		b.moveToServiceClient(id)
		return newSwitchSiteOutput(id), nil
	case login:
		if ok := switchedAnySite; !ok {
			return nil, errors.New("please switch site")
		}
		if err := b.login(ctx); err != nil && errors.Is(err, ErrAlreadyLogin) {
			return nil, fmt.Errorf("cannot login: %w", err)
		}
		return newLoginSuccededOutput(*b.currentServiceClientId), nil
	case logout:
		if ok := switchedAnySite; !ok {
			return nil, errors.New("you have not logged in to any site")
		}
		id := *b.currentServiceClientId
		b.logout()
		return newLogoutOutput(id), nil
	case viewProfile:
		if ok := switchedAnySite; !ok {
			return nil, errors.New("you have not logged in to any site")
		}
		profile, err := b.viewProfile(ctx)
		if err != nil {
			return nil, fmt.Errorf("canno get profile: %w", err)
		}
		return newViewProfileOutput(profile), nil
	}

	return nil, errors.New("unknown command")
}

func (b *Brawser) moveToServiceClient(id string) {
	b.currentServiceClientId = &id
}
func (b *Brawser) logout() {
	delete(b.accessTokens, *b.currentServiceClientId)
}

var (
	ErrAlreadyLogin = errors.New("already login")
)

func (b *Brawser) login(ctx context.Context) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if _, found := b.accessTokens[*b.currentServiceClientId]; found {
		return ErrAlreadyLogin
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, time.Duration(3)*time.Minute)
	defer cancel()
	b.codeReceiver.Init()
	b.codeReceiver.Start(timeoutCtx)
	codeResult := b.codeReceiver.Receive()
	if codeResult.err != nil {
		return fmt.Errorf("cannot receive authorization code: %w", codeResult.err)
	}
	slog.InfoContext(ctx, "receive Authorization Code!")

	// get accesstoken
	token, err := b.accessTokenClient.GetByCode(ctx, codeResult.code, AccessTokenRequestParam{
		ClientId:     *b.currentServiceClientId,
		ClientSecret: database.CLIENT_SECRET,
	})
	if err != nil {
		return fmt.Errorf("cannot get accesstoken: %w", err)
	}
	slog.InfoContext(ctx, "get Access Token!")
	b.accessTokens[*b.currentServiceClientId] = token.AccessToken
	b.refreshTokens[*b.currentServiceClientId] = token.RefreshToken
	return nil
}

func (b *Brawser) viewProfile(ctx context.Context) (map[string]any, error) {
	token, found := b.accessTokens[*b.currentServiceClientId]
	if !found {
		return nil, fmt.Errorf("access token is not found")
	}
	p, err := b.resourceClient.ViewProfile(ctx, token)
	if err != nil {
		if err := b.refreshToken(ctx); err != nil {
			return nil, fmt.Errorf("cannot get refresh token: %w", err)
		}
		if p, err = b.resourceClient.ViewProfile(ctx, token); err != nil {
			return nil, fmt.Errorf("cannot get profile: %w", err)
		}
	}
	var profile = map[string]any{}
	profile["id"] = p.UserId
	profile["age"] = p.Age
	profile["profile"] = p.Profile
	profile["name"] = p.Name
	return profile, nil
}

func (b *Brawser) refreshToken(ctx context.Context) error {
	refreshToken := b.refreshTokens[*b.currentServiceClientId]
	token, err := b.accessTokenClient.GetByRefreshToken(ctx, refreshToken, AccessTokenRequestParam{
		ClientId:     *b.currentServiceClientId,
		ClientSecret: database.CLIENT_SECRET,
	})
	if err != nil {
		return err
	}
	b.accessTokens[*b.currentServiceClientId] = token.AccessToken
	b.refreshTokens[*b.currentServiceClientId] = token.RefreshToken
	return nil
}

type (
	Command struct {
		command command
		args    []string
	}
	command string

	output struct {
		messageId messageId
		message   string
	}
	messageId uint
)

const (
	unknown     command = "unknown"
	help        command = "help"
	status      command = "status"
	showSites   command = "show-sites"
	switchsite  command = "switch-site"
	login       command = "login"
	logout      command = "logout"
	viewProfile command = "view-profile"

	unknownMsgId messageId = iota
	helpMsgId
	statusMsgId
	showsiteMsgId
	switchsiteMsgId
	loginSucceededMsgId
	logoutMsgId
	viewProfileMsgId
)

func ParseCommand(input string) Command {
	cmds := strings.Split(input, " ")
	switch command(cmds[0]) {
	case help:
		return Command{command: help}
	case status:
		return Command{command: status}
	case showSites:
		return Command{command: showSites}
	case switchsite:
		return Command{command: switchsite, args: cmds[1:]}
	case login:
		return Command{command: login}
	case logout:
		return Command{command: logout}
	case viewProfile:
		return Command{command: viewProfile}
	default:
		return Command{command: unknown}
	}
}

// Brawser outputs

var (
	helpOutput = &output{
		messageId: helpMsgId,
		message: `
- status
- show-sites
- switch-site [id]
- login
- logout
- help
- view-profile
`,
	}

	newStatusOutput = func(tokens map[string]string, currentId string) *output {
		if len(tokens) == 0 {
			return &output{
				messageId: statusMsgId,
				message: `
You have not logged in to any site.
				`,
			}
		}
		ids := make([]string, 0, len(tokens))
		for id := range tokens {
			ids = append(ids, id)
		}
		return &output{
			messageId: statusMsgId,
			message: fmt.Sprintf(`
You're logged in to site %v 
and now you're browsing %s.
			`, ids, currentId),
		}
	}

	showSitesOutput = &output{
		messageId: showsiteMsgId,
		message: fmt.Sprintf(`
Available Sites
- %s: Id[ %s ] 
- %s: Id[ %s ]
`,
			database.MockServiceClient500.Name, database.MockServiceClient500.Id,
			database.MockServiceClient501.Name, database.MockServiceClient501.Id,
		),
	}

	newSwitchSiteOutput = func(id string) *output {
		var message string
		switch id {
		case database.MockServiceClient500.Id:
			message = fmt.Sprintf(`
🎆🎆🎆🎆🎆🎆🎆🎆🎆🎆🎆🎆🎆
///////////////////////////////////
【 %s 】
///////////////////////////////////
`,
				database.MockServiceClient500.Name)
		case database.MockServiceClient501.Name:
			message = fmt.Sprintf(`
🦭🦭🦭🦭🦭🦭🦭🦭🦭🦭🦭🦭🦭
///////////////////////////////////
【 %s 】
///////////////////////////////////
`,
				database.MockServiceClient500.Name)
		}
		return &output{
			messageId: switchsiteMsgId,
			message:   message,
		}
	}

	newLoginSuccededOutput = func(id string) *output {
		return &output{
			messageId: loginSucceededMsgId,
			message: fmt.Sprintf(`
🚀Login to %s!!
`,
				id),
		}
	}

	newLogoutOutput = func(id string) *output {
		return &output{
			messageId: logoutMsgId,
			message: fmt.Sprintf(`
logout from %s
`,
				id),
		}
	}

	newViewProfileOutput = func(profile map[string]any) *output {
		var ps = make([]string, 0, len(profile))
		for key, val := range profile {
			ps = append(ps, fmt.Sprintf(`- %s: %v`, key, val))
		}
		return &output{
			messageId: viewProfileMsgId,
			message: fmt.Sprintf(`
%s
`,
				strings.Join(ps, "\n")),
		}
	}
)
