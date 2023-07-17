package service

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/pingcap/log"
	"go.uber.org/zap"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/apimachinery/pkg/types"
	_ "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	pingcapv1alph1 "github.com/pingcap/tidb-operator/pkg/apis/pingcap/v1alpha1"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

// @title		Dbaas 101 API
// @version		1.0
// @description	Dbaas 101 API.
// @BasePath		/

type API struct {
	k8sClient client.Client // access api server with cache
	k8sReader client.Reader // access api server directly
}

func NewAPI(cli client.Client, reader client.Reader) (*API, error) {
	return &API{
		k8sClient: cli,
		k8sReader: reader,
	}, nil
}

func (api *API) CreateTidbCluster(ctx context.Context, param *pingcapv1alph1.TidbCluster) (*CreateTidbClusterResult, error) {
	if err := api.k8sClient.Create(ctx, param); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil, ErrNotFound("tidb cluster already existed")
		}
		log.Error("create tidb clusters failed", zap.Error(err))
		return nil, ErrInternal(err.Error())
	}
	resp := new(CreateTidbClusterResult)
	return resp, nil
}

func (api *API) GetTidbCluster(ctx context.Context, param *GetTidbClusterParam) (*GetTidbClusterResult, error) {
	tc := new(pingcapv1alph1.TidbCluster)
	key := types.NamespacedName{
		Namespace: param.Namespace,
		Name:      param.Name,
	}
	if err := api.k8sClient.Get(ctx, key, tc); err != nil {
		if apierrors.IsNotFound(err) {
			return nil, ErrNotFound("tidb cluster not found")
		}
		log.Error("get tidb clusters failed", zap.Error(err))
		return nil, ErrInternal(err.Error())
	}

	resp := new(GetTidbClusterResult)
	resp.Meta = TidbClusterMeta{
		Name:      tc.Name,
		Namespace: tc.Namespace,
	}
	resp.Spec = tc.Spec
	resp.Status = tc.Status
	return resp, nil
}

func (api *API) ListTidbCluster(ctx context.Context) (*ListTidbClusterResult, error) {
	list := new(pingcapv1alph1.TidbClusterList)
	err := api.k8sClient.List(ctx, list)
	if client.IgnoreNotFound(err) != nil {
		log.Error("list tidb clusters failed", zap.Error(err))
		return nil, ErrInternal(err.Error())
	}
	res := &ListTidbClusterResult{
		Items: make([]TidbClusterMeta, 0, len(list.Items)),
	}
	for _, tc := range list.Items {
		res.Items = append(res.Items, TidbClusterMeta{
			Namespace: tc.Namespace,
			Name:      tc.Name,
		})
	}
	return res, nil
}

func (api *API) DeleteTidbCluster(ctx context.Context, param *DeleteTidbClusterParam) (*DeleteTidbClusterResult, error) {
	list := new(pingcapv1alph1.TidbClusterList)
	listOptions := &client.ListOptions{}
	listOptions.FieldSelector = fields.OneTermEqualSelector("metadata.name", param.Name)
	if err := api.k8sClient.List(ctx, list, listOptions); client.IgnoreNotFound(err) != nil {
		log.Error("list tidb clusters failed", zap.Error(err))
		return nil, ErrInternal(err.Error())
	}
	if len(list.Items) == 0 {
		return &DeleteTidbClusterResult{}, nil
	}
	tc := list.Items[0]
	err := api.k8sClient.Delete(ctx, &tc)
	if client.IgnoreNotFound(err) != nil {
		log.Error("delete tidb clusters failed", zap.Error(err))
		return nil, ErrInternal(err.Error())
	}
	return &DeleteTidbClusterResult{}, nil
}

// @Summary		    Crerate tidb cluster
// @Tags		    tidb cluster
// @Description	    Create tidb cluster
// @Accept			json
// @Produce			json
// @Success			200	                    {object}	CreateTidbClusterResult
// @Failure      	400  		            {object}  	ErrResp
// @Failure      	500  		            {object}  	ErrResp
// @Router			/api/v1/tidbclusters [post]
func (api *API) CreateTidbClusterHandler(ctx *gin.Context) {
	param := new(pingcapv1alph1.TidbCluster)
	if err := ctx.BindJSON(param); err != nil {
		EncodeError(ctx, ErrInvalidParameter(err.Error()))
		return
	}
	resp, err := api.CreateTidbCluster(ctx, param)
	if err != nil {
		EncodeError(ctx, err)
		return
	}
	EncodeResp(ctx, resp)
}

// @Summary		    Get tidb cluster
// @Tags		    tidb cluster
// @Description	    Get tidb cluster
// @Accept			json
// @Produce			json
// @Success			200	                    {object}	GetTidbClusterResult
// @Failure      	400  		            {object}  	ErrResp
// @Failure      	500  		            {object}  	ErrResp
// @Router			/api/v1/tidbclusters/:namespace/:name [get]
func (api *API) GetTidbClusterHandler(ctx *gin.Context) {
	param := new(GetTidbClusterParam)
	param.Namespace = ctx.Param("namespace")
	param.Name = ctx.Param("name")
	if err := param.Validate(); err != nil {
		EncodeError(ctx, ErrInvalidParameter(err.Error()))
		return
	}
	resp, err := api.GetTidbCluster(ctx, param)
	if err != nil {
		EncodeError(ctx, err)
		return
	}
	EncodeResp(ctx, resp)
}

// @Summary		    List all tidb cluster
// @Tags		    tidb cluster
// @Description	    List all tidb cluster
// @Accept			json
// @Produce			json
// @Success			200	                    {object}	ListTidbClusterResult
// @Failure      	400  		            {object}  	ErrResp
// @Failure      	500  		            {object}  	ErrResp
// @Router			/api/v1/tidbclusters [get]
func (api *API) ListTidbClusterHandler(ctx *gin.Context) {
	resp, err := api.ListTidbCluster(ctx)
	if err != nil {
		EncodeError(ctx, err)
		return
	}
	EncodeResp(ctx, resp)
}

// @Summary		Delete tidb cluster
// @Tags		tidb cluster
// @Description	Delete tidb cluster
// @Accept			json
// @Produce			json
// @Param			param			body		DeleteTidbClusterParam	    true	"param"
// @Success			200				{object}	ApplyResult
// @Failure      	400  			{object}  	ErrResp
// @Failure      	500  			{object}  	ErrResp
// @Router			/api/v1/tidbclusters/:namespace/:name [delete]
func (api *API) DeleteTidbClusterHandler(ctx *gin.Context) {
	param := &DeleteTidbClusterParam{}
	param.Namespace = ctx.Param("namespace")
	param.Name = ctx.Param("name")
	if err := param.Validate(); err != nil {
		EncodeError(ctx, ErrInvalidParameter(err.Error()))
		return
	}
	resp, err := api.DeleteTidbCluster(ctx, param)
	if err != nil {
		EncodeError(ctx, err)
		return
	}
	EncodeResp(ctx, resp)

}
