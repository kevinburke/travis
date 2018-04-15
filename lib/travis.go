package travis

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	types "github.com/kevinburke/go-types"
	"github.com/kevinburke/rest"
)

// GetToken looks in a file for the Travis API token. organization is your
// Github username ("kevinburke") or organization ("golang").
func GetToken(organization string) (string, error) {
	var filename string
	var f io.ReadCloser
	var err error
	checkedLocations := make([]string, 1)
	if cfg, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		filename = filepath.Join(cfg, "travis")
		f, err = os.Open(filename)
		checkedLocations[0] = filename
	} else {
		var homeDir string
		user, userErr := user.Current()
		if userErr == nil {
			homeDir = user.HomeDir
		} else {
			homeDir = os.Getenv("HOME")
		}
		filename = filepath.Join(homeDir, "cfg", "travis")
		f, err = os.Open(filename)
		checkedLocations[0] = filename
		if err != nil { //fallback
			rcFilename := filepath.Join(homeDir, ".travis")
			f, err = os.Open(rcFilename)
			checkedLocations = append(checkedLocations, rcFilename)
		}
	}
	if err != nil {
		err = fmt.Errorf(`Couldn't find a config file in %s.

Add a configuration file with your Travis token, like this:

[organizations]

    [organizations.kevinburke]
    token = "aabbccddeeff00"

Go to https://travis-ci.org/profile/<your-username> if you need to find your token.
`, strings.Join(checkedLocations, " or "))
		return "", err
	}
	defer f.Close()
	var c FileConfig
	_, err = toml.DecodeReader(bufio.NewReader(f), &c)
	if err != nil {
		return "", err
	}
	org, err := getCaseInsensitiveOrg(organization, c.Organizations)
	if err != nil {
		return "", err
	}
	return org.Token, nil
}

const Version = "0.1"

const Host = "https://api.travis-ci.org"
const WebHost = "https://travis-ci.org"

type Client struct {
	*rest.Client
	token string
}

func NewClient(token string) *Client {
	return &Client{
		Client: rest.NewClient("", "", Host),
		token:  token,
	}
}

func (c *Client) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	req, err := c.Client.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+c.token)
	req.Header.Set("Travis-API-Version", "3")
	return req, nil
}

// getCaseInsensitiveOrg finds the key in the list of orgs. This is a case
// insensitive match, so if key is "ShyP" and orgs has a key named "sHYp",
// that will count as a match.
func getCaseInsensitiveOrg(key string, orgs map[string]organization) (organization, error) {
	for k, _ := range orgs {
		lower := strings.ToLower(k)
		if _, ok := orgs[lower]; !ok {
			orgs[lower] = orgs[k]
			delete(orgs, k)
		}
	}
	lowerKey := strings.ToLower(key)
	if o, ok := orgs[lowerKey]; ok {
		return o, nil
	} else {
		return organization{}, fmt.Errorf(`Couldn't find organization %s in the config.

Go to https://travis-ci.org/profile/<your-username> if you need to create a token.
`, key)
	}
}

type FileConfig struct {
	Organizations map[string]organization
}

type organization struct {
	Token string
}

type Pagination struct {
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	Count   int  `json:"count"`
	IsFirst bool `json:"is_first"`
	IsLast  bool `json:"is_last"`
}

type Response struct {
	Type           string     `json:"@type"`
	HREF           string     `json:"@href"`
	Representation string     `json:"@representation"`
	Pagination     Pagination `json:"@pagination"`
	Data           interface{}
}

type Build struct {
	Type           string         `json:"@type"`
	HREF           string         `json:"@href"`
	Representation string         `json:"@representation"`
	ID             int64          `json:"id"`
	Number         string         `json:"number"`
	State          string         `json:"state"`
	PreviousState  string         `json:"previous_state"`
	Duration       int64          `json:"duration"`
	StartedAt      time.Time      `json:"started_at"`
	FinishedAt     types.NullTime `json:"finished_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	Branch         Branch         `json:"branch"`
	Repository     Repository     `json:"repository"`
}

type Branch struct {
	Type           string `json:"@type"`
	HREF           string `json:"@href"`
	Representation string `json:"@representation"`
	Name           string `json:"name"`
}

type Repository struct {
	Type           string `json:"@type"`
	HREF           string `json:"@href"`
	Representation string `json:"@representation"`
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
}

func (r *Response) UnmarshalJSON(b []byte) error {
	r2 := make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &r2); err != nil {
		return err
	}
	if err := json.Unmarshal(r2["@type"], &r.Type); err != nil {
		return err
	}
	if err := json.Unmarshal(r2["@href"], &r.HREF); err != nil {
		return err
	}
	if err := json.Unmarshal(r2["@representation"], &r.Representation); err != nil {
		return err
	}
	if err := json.Unmarshal(r2["@pagination"], &r.Pagination); err != nil {
		return err
	}
	if err := json.Unmarshal(r2[r.Type], &r.Data); err != nil {
		return err
	}
	return nil
}
