/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package store

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	extmocks "github.com/trustbloc/fabric-peer-ext/pkg/mocks"

	"github.com/trustbloc/sidetree-core-go/pkg/api/batch"
	"github.com/trustbloc/sidetree-core-go/pkg/docutil"

	"github.com/trustbloc/sidetree-fabric/pkg/context/store/mocks"
)

//go:generate counterfeiter -o ./mocks/store.gen.go --fake-name Store . store

const (
	chID = "mychannel"
	id   = "id"

	namespace = "did:sidetree"
)

func TestNew(t *testing.T) {
	c := NewClient(chID, namespace, &mocks.Store{})
	require.NotNil(t, c)
}

func TestProviderError(t *testing.T) {
	testErr := errors.New("provider error")

	s := &mocks.Store{}
	s.QueryReturns(nil, testErr)

	c := NewClient(chID, namespace, s)
	require.NotNil(t, c)

	payload, err := c.Get(id)
	require.Error(t, err)
	require.Contains(t, err.Error(), testErr.Error())
	require.Nil(t, payload)
}

func TestWriteContent(t *testing.T) {
	didID := namespace + docutil.NamespaceDelimiter + id

	vk1 := &queryresult.KV{
		Namespace: "document~diddoc",
		Key:       didID,
		Value:     []byte("{}"),
	}

	s := &mocks.Store{}
	it := extmocks.NewResultsIterator().WithResults([]*queryresult.KV{vk1})
	s.QueryReturns(it, nil)

	c := NewClient(chID, namespace, s)

	ops, err := c.Get(id)
	require.Nil(t, err)
	require.NotNil(t, ops)
	require.Equal(t, 1, len(ops))
}

func TestGetOperationsError(t *testing.T) {

	doc, err := getOperations([][]byte{[]byte("[test : 123]")})
	require.NotNil(t, err)
	require.Nil(t, doc)
	require.Contains(t, err.Error(), "invalid character")
}

func TestClient_Put(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		s := &mocks.Store{}
		c := NewClient(chID, namespace, s)

		require.NoError(t, c.Put([]*batch.Operation{{}}))
	})

	t.Run("Store error", func(t *testing.T) {
		errExpected := errors.New("injected store error")
		s := &mocks.Store{}
		s.PutReturns(errExpected)

		c := NewClient(chID, namespace, s)

		err := c.Put([]*batch.Operation{{}})
		require.EqualError(t, err, errExpected.Error())
	})
}

func TestClient_Get(t *testing.T) {
	s := &mocks.Store{}
	s.QueryReturns(&extmocks.ResultsIterator{}, nil)

	c := NewClient(chID, namespace, s)

	ops, err := c.Get("suffix")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
	require.Empty(t, ops)
}

func TestSort(t *testing.T) {
	var operations []*batch.Operation

	const testID = "id"
	delete := &batch.Operation{ID: testID, Type: "delete", TransactionTime: 2, TransactionNumber: 1}
	update := &batch.Operation{ID: testID, Type: "update", TransactionTime: 1, TransactionNumber: 7}
	create := &batch.Operation{ID: testID, Type: "create", TransactionTime: 1, TransactionNumber: 1}

	operations = append(operations, delete)
	operations = append(operations, update)
	operations = append(operations, create)

	result := sortChronologically(operations)
	require.Equal(t, create.Type, result[0].Type)
	require.Equal(t, update.Type, result[1].Type)
	require.Equal(t, delete.Type, result[2].Type)
}
