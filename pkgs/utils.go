package pkgs

import "net/http"

func GetServerHealth() (string, error) {
	resp, err := http.Get("http://api.spacetraders.io/game/status")

	if err != nil {
		return "", err
	}

	return resp.Status, nil
}
