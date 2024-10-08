package fly

import (
	"context"
	"fmt"
)

func (c *Client) ListMachines(ctx context.Context, appName string, includeDeleted *bool, region, summary *string) ([]Machine, error) {
	queryParams := map[string]string{
		"region": "",
	}

	if includeDeleted != nil {
		queryParams["include_deleted"] = fmt.Sprintf("%v", *includeDeleted)
	}
	if region != nil {
		queryParams["region"] = *region
	}
	if summary != nil {
		queryParams["summary"] = *summary
	}

	var res []Machine
	endpoint := fmt.Sprintf("/apps/%s/machines", appName)
	if err := c.doRequest(ctx, "GET", endpoint, nil, queryParams, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) CreateMachine(ctx context.Context, appName string, reqBody CreateMachineRequest) (*Machine, error) {
	var res Machine
	endpoint := fmt.Sprintf("/apps/%s/machines", appName)
	if err := c.doRequest(ctx, "POST", endpoint, reqBody, nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetMachine(ctx context.Context, appName string, machineID string) (*Machine, error) {
	var res Machine
	endpoint := fmt.Sprintf("/apps/%s/machines/%s", appName, machineID)
	if err := c.doRequest(ctx, "GET", endpoint, nil, nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) UpdateMachine(ctx context.Context, appName string, machineID string, reqBody UpdateMachineRequest) (*Machine, error) {
	var res Machine
	endpoint := fmt.Sprintf("/apps/%s/machines/%s", appName, machineID)
	if err := c.doRequest(ctx, "POST", endpoint, reqBody, nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) DeleteMachine(ctx context.Context, appName string, machineID string, force *bool) error {
	endpoint := fmt.Sprintf("/apps/%s/machines/%s", appName, machineID)
	queryParams := map[string]string{}

	if force != nil {
		queryParams["force"] = fmt.Sprintf("%v", *force)
	}

	return c.doRequest(ctx, "DELETE", endpoint, nil, queryParams, nil, nil)
}

func (c *Client) ListMachineEvents(ctx context.Context, appName string, machineID string) (*[]MachineEvent, error) {
	var res []MachineEvent
	endpoint := fmt.Sprintf("/apps/%s/machines/%s/events", appName, machineID)
	if err := c.doRequest(ctx, "GET", endpoint, nil, nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) ListProcesses(ctx context.Context, appName string, machineID string) (*[]ProcessStat, error) {
	var res []ProcessStat
	endpoint := fmt.Sprintf("/apps/%s/machines/%s/ps", appName, machineID)
	if err := c.doRequest(ctx, "GET", endpoint, nil, nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) CreateLease(ctx context.Context, appName string, machineID string, reqBody CreateLeaseRequest, leaseNonce string) (*Lease, error) {
	headers := map[string]string{}
	if leaseNonce != "" {
		headers["fly-machine-lease-nonce"] = leaseNonce
	}

	var res Lease
	endpoint := fmt.Sprintf("/apps/%s/machines/%s/lease", appName, machineID)
	if err := c.doRequest(ctx, "POST", endpoint, reqBody, nil, headers, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetLease(ctx context.Context, appName string, machineID string) (*Lease, error) {
	var res Lease
	endpoint := fmt.Sprintf("/apps/%s/machines/%s/lease", appName, machineID)
	if err := c.doRequest(ctx, "GET", endpoint, nil, nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) ExecuteCommand(ctx context.Context, appName, machineID string, reqBody MachineExecRequest) (*ExecResponse, error) {
	var res ExecResponse
	endpoint := fmt.Sprintf("/apps/%s/machines/%s/exec", appName, machineID)
	if err := c.doRequest(ctx, "POST", endpoint, reqBody, nil, nil, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetMachineMetadata(ctx context.Context, appName string, machineID string) (map[string]string, error) {
	var res map[string]string
	endpoint := fmt.Sprintf("/apps/%s/machines/%s/metadata", appName, machineID)
	if err := c.doRequest(ctx, "GET", endpoint, nil, nil, nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *Client) UpdateMachineMetadata(ctx context.Context, appName, machineID, key string, value string) error {
	endpoint := fmt.Sprintf("/apps/%s/machines/%s/metadata/%s", appName, machineID, key)
	if err := c.doRequest(ctx, "POST", endpoint, map[string]string{"value": value}, nil, nil, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteMachineMetadata(ctx context.Context, appName, machineID, key string) error {
	endpoint := fmt.Sprintf("/apps/%s/machines/%s/metadata/%s", appName, machineID, key)
	if err := c.doRequest(ctx, "DELETE", endpoint, nil, nil, nil, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) RestartMachine(ctx context.Context, appName, machineID string, timeout, signal *string) error {
	endpoint := fmt.Sprintf("/apps/%s/machines/%s/restart", appName, machineID)
	queryParams := map[string]string{}

	if timeout != nil {
		queryParams["timeout"] = *timeout
	}
	if signal != nil {
		queryParams["signal"] = *signal
	}

	return c.doRequest(ctx, "POST", endpoint, nil, queryParams, nil, nil)
}

func (c *Client) SendSignal(ctx context.Context, appName, machineID string, reqBody SignalRequest) error {
	endpoint := fmt.Sprintf("/apps/%s/machines/%s/signal", appName, machineID)
	if err := c.doRequest(ctx, "POST", endpoint, reqBody, nil, nil, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) UncordonMachine(ctx context.Context, appName, machineID string) error {
	endpoint := fmt.Sprintf("/apps/%s/machines/%s/uncordon", appName, machineID)
	if err := c.doRequest(ctx, "POST", endpoint, nil, nil, nil, nil); err != nil {
		return err
	}
	return nil
}
