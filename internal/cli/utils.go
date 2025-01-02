package cli

import (
	"encoding/json"
	"fmt"
)

func prettyPrint(data any) error {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	fmt.Println(string(res))
	return nil
}
