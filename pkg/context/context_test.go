/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package context

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	protocolApi "github.com/trustbloc/sidetree-core-go/pkg/api/protocol"
	"github.com/trustbloc/sidetree-core-go/pkg/batch/opqueue"
	"github.com/trustbloc/sidetree-fabric/pkg/mocks"
)

//go:generate counterfeiter -o ./../mocks/txnserviceprovider.gen.go --fake-name TxnServiceProvider . txnServiceProvider
//go:generate counterfeiter -o ./../mocks/txnservice.gen.go --fake-name TxnService github.com/trustbloc/fabric-peer-ext/pkg/txn/api.Service
//go:generate counterfeiter -o ./../mocks/opqueueprovider.gen.go --fake-name OperationQueueProvider . operationQueueProvider

const (
	channelID = "channel1"
	namespace = "did:sidetree"
)

func TestNew(t *testing.T) {
	txnProvider := &mocks.TxnServiceProvider{}
	dcasProvider := &mocks.DCASClientProvider{}
	opQueueProvider := &mocks.OperationQueueProvider{}
	protocolVersions := map[string]protocolApi.Protocol{}

	errExpected := errors.New("injected op queue error")
	opQueueProvider.CreateReturns(nil, errExpected)

	sctx, err := New(channelID, namespace, protocolVersions, txnProvider, dcasProvider, opQueueProvider)
	require.EqualError(t, err, errExpected.Error())
	require.Nil(t, sctx)

	opQueueProvider.CreateReturns(&opqueue.MemQueue{}, nil)

	sctx, err = New(channelID, namespace, protocolVersions, txnProvider, dcasProvider, opQueueProvider)
	require.NoError(t, err)
	require.NotNil(t, sctx)

	require.NotNil(t, sctx.Protocol())
	require.NotNil(t, sctx.CAS())
	require.NotNil(t, sctx.Blockchain())
	require.NotEmpty(t, sctx.Namespace())
	require.NotNil(t, sctx.OperationQueue())
}
