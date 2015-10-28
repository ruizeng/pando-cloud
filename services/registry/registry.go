package main

import (
	"github.com/PandoCloud/pando-cloud/pkg/models"
	//"github.com/PandoCloud/pando-cloud/pkg/rpcs"
	"github.com/PandoCloud/pando-cloud/pkg/server"
)

type Registry struct{}

func (r *Registry) ValidateProduct(key string, reply *models.Product) error {

}
