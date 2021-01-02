package mgo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mgoClient *mongo.Client

func Connect(ctx context.Context, opts ...*options.ClientOptions) (*mongo.Client, error) {
	if mgoClient != nil {
		return mgoClient, nil
	}

	var err error
	mgoClient, err = mongo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	err = mgoClient.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return mgoClient, nil
}

func Close(ctx context.Context) error {
	return mgoClient.Disconnect(ctx)
}
