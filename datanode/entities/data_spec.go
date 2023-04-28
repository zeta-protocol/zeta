// Copyright (c) 2022 Gobalsky Labs Limited
//
// Use of this software is governed by the Business Source License included
// in the LICENSE.DATANODE file and at https://www.mariadb.com/bsl11.
//
// Change Date: 18 months from the later of the date of the first publicly
// available Distribution of this version of the repository, and 25 June 2022.
//
// On the date above, in accordance with the Business Source License, use
// of this software will be governed by version 3 or later of the GNU General
// Public License.

package entities

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zeta-protocol/zeta/core/types"
	v2 "github.com/zeta-protocol/zeta/protos/data-node/api/v2"
	datapb "github.com/zeta-protocol/zeta/protos/zeta/data/v1"

	zetapb "github.com/zeta-protocol/zeta/protos/zeta"
)

type ExternalDataSourceSpec struct {
	Spec *DataSourceSpec
}

func (s *ExternalDataSourceSpec) ToProto() *zetapb.ExternalDataSourceSpec {
	return &zetapb.ExternalDataSourceSpec{
		Spec: s.Spec.ToProto(),
	}
}

func ExternalDataSourceSpecFromProto(spec *zetapb.ExternalDataSourceSpec, txHash TxHash, zetaTime time.Time) (*ExternalDataSourceSpec, error) {
	if spec.Spec != nil {
		ds, err := DataSourceSpecFromProto(spec.Spec, txHash, zetaTime)
		if err != nil {
			return nil, err
		}

		return &ExternalDataSourceSpec{
			Spec: ds,
		}, nil
	}

	return &ExternalDataSourceSpec{
		Spec: &DataSourceSpec{},
	}, nil
}

type (
	_Spec  struct{}
	SpecID = ID[_Spec]
)

type (
	Signer  []byte
	Signers = []Signer
)

type DataSourceSpecConfiguration struct {
	Signers Signers
	Filters []Filter
}

type DataSourceSpec struct {
	ID        SpecID
	CreatedAt time.Time
	UpdatedAt time.Time
	Data      *DataSourceDefinition
	Status    DataSourceSpecStatus
	TxHash    TxHash
	ZetaTime  time.Time
}

type DataSourceSpecRaw struct {
	ID        SpecID
	CreatedAt time.Time
	UpdatedAt time.Time
	Signers   Signers
	Filters   []Filter
	Status    DataSourceSpecStatus
	TxHash    TxHash
	ZetaTime  time.Time
}

func DataSourceSpecFromProto(spec *zetapb.DataSourceSpec, txHash TxHash, zetaTime time.Time) (*DataSourceSpec, error) {
	if spec != nil {
		id := SpecID(spec.Id)

		data := &DataSourceDefinition{}
		if spec.Data != nil {
			filters := FiltersFromProto(spec.Data.GetFilters())
			signers, err := SerializeSigners(types.SignersFromProto(spec.Data.GetSigners()))
			if err != nil {
				return nil, err
			}

			data.External = &DataSourceDefinitionExternal{
				Signers: signers,
				Filters: filters,
			}
		}

		return &DataSourceSpec{
			ID:        id,
			CreatedAt: time.Unix(0, spec.CreatedAt),
			UpdatedAt: time.Unix(0, spec.UpdatedAt),
			Data:      data,
			Status:    DataSourceSpecStatus(spec.Status),
			TxHash:    txHash,
			ZetaTime:  zetaTime,
		}, nil
	}
	return nil, nil
}

func (ds *DataSourceSpec) ToProto() *zetapb.DataSourceSpec {
	filters := []*datapb.Filter{}
	signers := []*datapb.Signer{}

	if ds.Data != nil {
		desSigners := DeserializeSigners(ds.Data.External.Signers)
		signers = types.SignersIntoProto(desSigners)
		filters = filtersToProto(ds.Data.External.Filters)
	}

	return &zetapb.DataSourceSpec{
		Id:        ds.ID.String(),
		CreatedAt: ds.CreatedAt.UnixNano(),
		UpdatedAt: ds.UpdatedAt.UnixNano(),
		Data: zetapb.NewDataSourceDefinition(
			zetapb.DataSourceDefinitionTypeExt,
		).SetOracleConfig(
			&zetapb.DataSourceSpecConfiguration{
				Signers: signers,
				Filters: filters,
			},
		),
		Status: zetapb.DataSourceSpec_Status(ds.Status),
	}
}

func (ds *DataSourceSpec) ToOracleProto() *zetapb.OracleSpec {
	return &zetapb.OracleSpec{
		ExternalDataSourceSpec: &zetapb.ExternalDataSourceSpec{
			Spec: ds.ToProto(),
		},
	}
}

func (ds DataSourceSpec) Cursor() *Cursor {
	return NewCursor(DataSourceSpecCursor{ds.ZetaTime, ds.ID}.String())
}

func (ds DataSourceSpec) ToOracleProtoEdge(_ ...any) (*v2.OracleSpecEdge, error) {
	return &v2.OracleSpecEdge{
		Node:   ds.ToOracleProto(),
		Cursor: ds.Cursor().Encode(),
	}, nil
}

func SerializeSigners(signers []*types.Signer) (Signers, error) {
	if len(signers) > 0 {
		sigList := Signers{}

		for _, signer := range signers {
			data, err := signer.Serialize()
			if err != nil {
				return nil, err
			}
			sigList = append(sigList, data)
		}

		return sigList, nil
	}

	return Signers{}, nil
}

func DeserializeSigners(data Signers) []*types.Signer {
	if len(data) > 0 {
		signers := []*types.Signer{}
		for _, s := range data {
			signer := types.DeserializeSigner(s)
			signers = append(signers, signer)
		}

		return signers
	}

	return nil
}

type DataSourceSpecCursor struct {
	ZetaTime time.Time `json:"zetaTime"`
	ID       SpecID    `json:"id"`
}

func (ds DataSourceSpecCursor) String() string {
	bs, err := json.Marshal(ds)
	if err != nil {
		panic(fmt.Errorf("could not marshal oracle spec cursor: %w", err))
	}
	return string(bs)
}

func (ds *DataSourceSpecCursor) Parse(cursorString string) error {
	if cursorString == "" {
		return nil
	}
	return json.Unmarshal([]byte(cursorString), ds)
}
