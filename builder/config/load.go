package config

import (
	"github.com/spf13/viper"

	"github.com/uor-framework/uor-client-go/builder/api/v1alpha1"
)

// ReadCollectionConfig read the specified config into a CollectionConfiguration type.
func ReadCollectionConfig(configName string) (v1alpha1.DataSetConfiguration, error) {
	var configuration v1alpha1.DataSetConfiguration
	cfg, err := readInConfig(configName, configuration)
	if err != nil {
		return configuration, err
	}
	return cfg.(v1alpha1.DataSetConfiguration), nil
}

// ReadSchemaConfig read the specified config into a SchemaConfiguration type.
func ReadSchemaConfig(configName string) (v1alpha1.SchemaConfiguration, error) {
	var configuration v1alpha1.SchemaConfiguration
	cfg, err := readInConfig(configName, configuration)
	if err != nil {
		return configuration, err
	}
	return cfg.(v1alpha1.SchemaConfiguration), nil
}

// ReadAttributeQuery read the specified config into a AttributeQuery type.
func ReadAttributeQuery(configName string) (v1alpha1.AttributeQuery, error) {
	var configuration v1alpha1.AttributeQuery
	cfg, err := readInConfig(configName, configuration)
	if err != nil {
		return configuration, err
	}
	return cfg.(v1alpha1.AttributeQuery), nil
}

func readInConfig(configName string, object interface{}) (interface{}, error) {
	viper.SetConfigName(configName)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&object)
	if err != nil {
		return nil, err
	}
	return object, nil
}
