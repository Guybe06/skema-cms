package main

const (
	APIVersion = "1.0.0"

	FeatureAuth          = "authentification"
	FeatureOrganizations = "organisations"
	FeatureMemberships   = "membres"
	FeatureConnections   = "connexions"
	FeatureCollections   = "collections"
	FeatureContent       = "contenu"
	FeatureStorage       = "stockage"
	FeaturePublicAPI     = "api-publique"
	FeatureAPIKeys       = "cles-api"
	FeatureDashboard     = "tableau-de-bord"
	FeatureMCP           = "mcp"
	FeatureNotifications = "notifications"
)

var availableFeatures = []string{
	FeatureAuth,
	FeatureOrganizations,
	FeatureMemberships,
	FeatureConnections,
	FeatureCollections,
	FeatureContent,
	FeatureStorage,
	FeaturePublicAPI,
	FeatureAPIKeys,
	FeatureDashboard,
	FeatureMCP,
	FeatureNotifications,
}
