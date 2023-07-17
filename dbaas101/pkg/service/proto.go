package service

import (
	"github.com/pingcap/errors"
	pingcapv1alph1 "github.com/pingcap/tidb-operator/pkg/apis/pingcap/v1alpha1"
)

type Spec struct {
	Replicas int32 `json:"replicas"`
}

type UpdateTidbClusterParam struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	PD        *Spec  `json:"pd,omitempty"`
	TiDB      *Spec  `json:"tidb,omitempty"`
	TiKV      *Spec  `json:"tikv,omitempty"`
	TiFlash   *Spec  `json:"tiflash,omitempty"`
}

func (param *UpdateTidbClusterParam) Validate() error {
	if param.Name == "" {
		return errors.New("name is empty")
	}
	if param.Namespace == "" {
		return errors.New("namespace is empty")
	}
	return nil
}

type UpdateTidbClusterResult struct{}

type CreateTidbClusterResult struct{}

type DeleteTidbClusterParam struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (param *DeleteTidbClusterParam) Validate() error {
	if param.Name == "" {
		return errors.New("name is empty")
	}
	if param.Namespace == "" {
		return errors.New("namespace is empty")
	}
	return nil
}

type DeleteTidbClusterResult struct{}

type GetTidbClusterParam struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (param *GetTidbClusterParam) Validate() error {
	if param.Name == "" {
		return errors.New("name is empty")
	}
	if param.Namespace == "" {
		return errors.New("namespace is empty")
	}
	return nil
}

// just return some fields of tidb cluster
type GetTidbClusterResult struct {
	Meta   TidbClusterMeta                  `json:"meta"`
	Spec   pingcapv1alph1.TidbClusterSpec   `json:"spec"`
	Status pingcapv1alph1.TidbClusterStatus `json:"status,omitempty"`
}

type TidbClusterMeta struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type ListTidbClusterResult struct {
	Items []TidbClusterMeta `json:"items"`
}
