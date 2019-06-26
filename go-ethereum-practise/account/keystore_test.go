package account

import "testing"

func TestCreateKeys(t *testing.T) {
	createKs("./tmp", "secret")
}

func TestImportKs(t *testing.T) {
	importKs("./tmp", "./tmp/UTC--2019-03-16T13-17-58.866652296Z--62c37c89dbe8d511c1a8c9ee6ff856e280ed334d", "secret")
}
