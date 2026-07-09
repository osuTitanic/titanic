package cdn

import (
	"context"
	"net/http"
)

func (c *Client) AdminDeleteFile(ctx context.Context, objectKey string) error {
	objectPath, err := objectPath(objectKey)
	if err != nil {
		return err
	}
	return c.performRequestJson(ctx, http.MethodDelete, "/admin/files/"+objectPath, nil, nil, nil, nil)
}
