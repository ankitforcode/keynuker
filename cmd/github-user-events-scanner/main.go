// Licensed to the Apache Software Foundation (ASF) under one or more contributor license agreements;
// and to You under the Apache License, Version 2.0.  See LICENSE in project root for full license + copyright.

package main

import (
	"encoding/json"

	"github.com/tleyden/keynuker"
	"github.com/tleyden/keynuker/keynuker-go-common"
)

// Scan Github user events for AWS keys

func main() {
	keynuker_go_common.LogMemoryUsageLoop()
	keynuker_go_common.RegistorOrInvokeActionStdIo(OpenWhiskCallback)
}

func OpenWhiskCallback(value json.RawMessage) (interface{}, error) {

	var params keynuker.ParamsScanGithubUserEventsForAwsKeys

	err := json.Unmarshal(value, &params)
	if err != nil {
		return nil, err
	}

	if err := params.Validate(); err != nil {
		return params, err
	}

	params = params.WithDefaultKeynukerOrg()

	// If any checkpoints are missing (null), set a default checkpoint of 12 hours to
	// prevent excessive unwanted historical scanning
	params = params.SetDefaultCheckpointsForMissing(keynuker_go_common.DefaultCheckpointEventTimeWindow)

	fetcher := keynuker.NewGoGithubUserEventFetcher(params.GithubAccessToken, params.GithubApiUrl)

	scanner := keynuker.NewGithubUserEventsScanner(fetcher)

	docWrapper, err := scanner.ScanAwsKeys(params)
	if err != nil {
		return nil, err
	}

	return docWrapper, nil
}
