package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	pixela "github.com/ebc-2in2crc/pixela4go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ErrNeglect is error that had reported by pa.
var ErrNeglect = errors.New("neglect")

type pixelaClientFactory struct {
	user    pixelaUser
	profile pixelaUserProfile
	graph   pixelaGraph
	pixel   pixelaPixel
	webhook pixelaWebhook
}

type pixelaUser interface {
	Create(input *pixela.UserCreateInput) (*pixela.Result, error)
	Update(input *pixela.UserUpdateInput) (*pixela.Result, error)
	Delete() (*pixela.Result, error)
}

type pixelaUserProfile interface {
	Update(input *pixela.UserProfileUpdateInput) (*pixela.Result, error)
	URL() string
}

type pixelaGraph interface {
	Create(input *pixela.GraphCreateInput) (*pixela.Result, error)
	GetAll() (*pixela.GraphDefinitions, error)
	Get(input *pixela.GraphGetInput) (*pixela.GraphDefinition, error)
	GetSVG(input *pixela.GraphGetSVGInput) (string, error)
	URL(input *pixela.GraphURLInput) string
	Stats(input *pixela.GraphStatsInput) (*pixela.Stats, error)
	Update(input *pixela.GraphUpdateInput) (*pixela.Result, error)
	Delete(input *pixela.GraphDeleteInput) (*pixela.Result, error)
	GetPixelDates(input *pixela.GraphGetPixelDatesInput) (*pixela.Pixels, error)
	Stopwatch(input *pixela.GraphStopwatchInput) (*pixela.Result, error)
	Add(input *pixela.GraphAddInput) (*pixela.Result, error)
}

type pixelaPixel interface {
	Create(input *pixela.PixelCreateInput) (*pixela.Result, error)
	Increment(input *pixela.PixelIncrementInput) (*pixela.Result, error)
	Decrement(input *pixela.PixelDecrementInput) (*pixela.Result, error)
	Get(input *pixela.PixelGetInput) (*pixela.Quantity, error)
	Update(input *pixela.PixelUpdateInput) (*pixela.Result, error)
	Delete(input *pixela.PixelDeleteInput) (*pixela.Result, error)
}

type pixelaWebhook interface {
	Create(input *pixela.WebhookCreateInput) (*pixela.WebhookCreateResult, error)
	GetAll() (*pixela.WebhookDefinitions, error)
	Invoke(input *pixela.WebhookInvokeInput) (*pixela.Result, error)
	Delete(input *pixela.WebhookDeleteInput) (*pixela.Result, error)
}

func (p *pixelaClientFactory) User() pixelaUser {
	if p.user != nil {
		return p.user
	}
	return pixela.New(getUsername(), getToken()).User()
}

func (p *pixelaClientFactory) UserProfile() pixelaUserProfile {
	if p.profile != nil {
		return p.profile
	}
	return pixela.New(getUsername(), getToken()).UserProfile()
}

func (p *pixelaClientFactory) Graph() pixelaGraph {
	if p.graph != nil {
		return p.graph
	}
	return pixela.New(getUsername(), getToken()).Graph()
}

func (p *pixelaClientFactory) Pixel() pixelaPixel {
	if p.pixel != nil {
		return p.pixel
	}
	return pixela.New(getUsername(), getToken()).Pixel()
}

func (p *pixelaClientFactory) Webhook() pixelaWebhook {
	if p.webhook != nil {
		return p.webhook
	}
	return pixela.New(getUsername(), getToken()).Webhook()
}

var pixelaClient = pixelaClientFactory{}

func getBoolFlag(name string) *bool {
	return getBoolPtr(viper.GetBool(name))
}

func getBoolPtr(v bool) *bool {
	if !v {
		return nil
	}
	return &v
}

func getStringFlag(name string) *string {
	return getStringPtr(viper.GetString(name))
}

func getStringPtr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func showHelp(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func marshalResult(result *pixela.Result) (string, error) {
	b, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("marshal object failed: %w", err)
	}
	return string(b), nil
}
