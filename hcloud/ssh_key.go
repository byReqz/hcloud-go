package hcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hetznercloud/hcloud-go/hcloud/schema"
)

// SSHKey represents a SSH key in the Hetzner Cloud.
type SSHKey struct {
	ID          int
	Name        string
	Fingerprint string
	PublicKey   string
}

// SSHKeyFromSchema converts a schema.SSHKey to a SSHKey.
func SSHKeyFromSchema(s schema.SSHKey) SSHKey {
	return SSHKey{
		ID:          s.ID,
		Name:        s.Name,
		Fingerprint: s.Fingerprint,
		PublicKey:   s.PublicKey,
	}
}

// SSHKeyClient is a client for the SSH keys API.
type SSHKeyClient struct {
	client *Client
}

// Get retrieves a SSH key.
func (c *SSHKeyClient) Get(ctx context.Context, id int) (*SSHKey, *Response, error) {
	req, err := c.client.NewRequest(ctx, "GET", fmt.Sprintf("/ssh_keys/%d", id), nil)
	if err != nil {
		return nil, nil, err
	}

	var body struct {
		SSHKey schema.SSHKey `json:"ssh_key"`
	}
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	sshKey := SSHKeyFromSchema(body.SSHKey)
	return &sshKey, resp, nil
}

// SSHKeyListOpts specifies options for listing SSH keys.
type SSHKeyListOpts struct {
	ListOpts
}

// List returns a list of SSH keys for a specific page.
func (c *SSHKeyClient) List(ctx context.Context, opts SSHKeyListOpts) ([]*SSHKey, *Response, error) {
	path := "/ssh_keys?" + valuesForListOpts(opts.ListOpts).Encode()
	req, err := c.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	var body struct {
		SSHKeys []schema.SSHKey `json:"ssh_keys"`
	}
	resp, err := c.client.Do(req, &body)
	if err != nil {
		return nil, nil, err
	}
	sshKeys := make([]*SSHKey, 0, len(body.SSHKeys))
	for _, s := range body.SSHKeys {
		sshKey := SSHKeyFromSchema(s)
		sshKeys = append(sshKeys, &sshKey)
	}
	return sshKeys, resp, nil
}

// All returns all SSH keys.
func (c *SSHKeyClient) All(ctx context.Context) ([]*SSHKey, error) {
	allSSHKeys := []*SSHKey{}

	opts := SSHKeyListOpts{}
	opts.PerPage = 50

	_, err := c.client.all(func(page int) (*Response, error) {
		opts.Page = page
		sshKeys, resp, err := c.List(ctx, opts)
		if err != nil {
			return resp, err
		}
		allSSHKeys = append(allSSHKeys, sshKeys...)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}

	return allSSHKeys, nil
}

// SSHKeyCreateOpts specifies parameters for creating a SSH key.
type SSHKeyCreateOpts struct {
	Name      string
	PublicKey string
}

// Validate checks if options are valid.
func (o SSHKeyCreateOpts) Validate() error {
	if o.Name == "" {
		return errors.New("missing name")
	}
	if o.PublicKey == "" {
		return errors.New("missing public key")
	}
	return nil
}

// Create creates a new SSH key with the given options.
func (c *SSHKeyClient) Create(ctx context.Context, opts SSHKeyCreateOpts) (*SSHKey, *Response, error) {
	if err := opts.Validate(); err != nil {
		return nil, nil, err
	}

	reqBody := struct {
		Name      string `json:"name"`
		PublicKey string `json:"public_key"`
	}{
		Name:      opts.Name,
		PublicKey: opts.PublicKey,
	}
	reqBodyData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, err
	}

	req, err := c.client.NewRequest(ctx, "POST", "/ssh_keys", bytes.NewReader(reqBodyData))
	if err != nil {
		return nil, nil, err
	}

	var respBody struct {
		SSHKey schema.SSHKey `json:"ssh_key"`
	}
	resp, err := c.client.Do(req, &respBody)
	if err != nil {
		return nil, resp, err
	}
	sshKey := SSHKeyFromSchema(respBody.SSHKey)
	return &sshKey, resp, nil
}

// Delete deletes a SSH key.
func (c *SSHKeyClient) Delete(ctx context.Context, id int) (*Response, error) {
	req, err := c.client.NewRequest(ctx, "DELETE", fmt.Sprintf("/ssh_keys/%d", id), nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req, nil)
}
