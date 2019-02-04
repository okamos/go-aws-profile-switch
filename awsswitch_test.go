package awsswitch

import (
	"path/filepath"
	"testing"
)

var loadtests = []map[string]interface{}{
	map[string]interface{}{
		"index": 0, "key": "aws_access_key_id", "expect": "profile1_id",
	},
	map[string]interface{}{
		"index": 2, "key": "aws_secret_access_key", "expect": "profile2_secret",
	},
	map[string]interface{}{
		"index": 1, "key": "region", "expect": "us-east-1",
	},
	map[string]interface{}{
		"index": 0, "key": "output", "expect": "json",
	},
}

func TestLoadCredentialsNormalFile(t *testing.T) {
	credentialsPath = filepath.Join("testdata", "normal")
	credentials, _ := LoadCredentials()
	for _, v := range loadtests {
		got := credentials[v["index"].(int)][v["key"].(string)]
		want := v["expect"].(string)
		if got != want {
			t.Errorf("got: %s, want: %s", got, want)
		}
	}
}

func TestLoadCredentialsSpaceFile(t *testing.T) {
	credentialsPath = filepath.Join("testdata", "space")
	credentials, _ := LoadCredentials()
	for _, v := range loadtests {
		got := credentials[v["index"].(int)][v["key"].(string)]
		want := v["expect"].(string)
		if got != want {
			t.Errorf("got: %s, want: %s", got, want)
		}
	}
}

func TestLoadCredentialsCommentFile(t *testing.T) {
	credentialsPath = filepath.Join("testdata", "comment")
	credentials, comments := LoadCredentials()
	for _, v := range loadtests {
		got := credentials[v["index"].(int)][v["key"].(string)]
		want := v["expect"].(string)
		if got != want {
			t.Errorf("got: %s, want: %s", got, want)
		}
	}
	got := string(comments)
	want := `# region=us-east-1
# [profile3]
# aws_access_key_id=profile3_id
# aws_secret_access_key=profile3_secret
# region=ap-northeast-1
# output=text
`
	if got != want {
		t.Errorf("got: %s, want: %s", got, want)
	}
}
