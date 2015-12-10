package security

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
)

type AwsInternalProvider struct {
	ServerAddress string
}

// REQUEST: curl http://169.254.169.254/latest/meta-data/iam/security-credentials/ROLE
func (a *AwsInternalProvider) Apply(c *cli.Context) error {

	role := c.String("aws-role")
	if len(role) == 0 {
		return errors.New("flag aws-role not set")
	}

	resp, err := http.Get(fmt.Sprintf("%s%s", a.ServerAddress, role))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var responseJson map[string]interface{}
	fmt.Println(body)
	responseBytes := []byte(body)

	if err := json.Unmarshal(responseBytes, &responseJson); err != nil {
		return err
	}

	SetEnvironmentVariable("AWS_ACCESS_KEY_ID", responseJson["AccessKeyId"])
	SetEnvironmentVariable("AWS_SECRET_ACCESS_KEY", responseJson["SecretAccessKey"])
	SetEnvironmentVariable("AWS_SESSION_TOKEN", responseJson["Token"])

	return nil
}

func SetEnvironmentVariable(key string, value interface{}) {
	if _, ok := value.(string); ok {
		os.Setenv(key, value.(string))
	}
}
