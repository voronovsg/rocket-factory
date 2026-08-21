package testcontainers

const (
	MongoContainerName = "mongo-inventory-test"
	MongoPort          = "27017"

	IAMStubContainerName = "iam-stub"
	IAMStubPort          = "50053"
	IAMStubDockerfile    = "platform/pkg/testcontainers/stubiam/Dockerfile"

	MongoImageNameKey = "MONGO_IMAGE_NAME"
	MongoHostKey      = "MONGO_HOST"
	MongoPortKey      = "MONGO_PORT"
	MongoDatabaseKey  = "MONGO_DATABASE"
	MongoUsernameKey  = "MONGO_INITDB_ROOT_USERNAME"
	MongoPasswordKey  = "MONGO_INITDB_ROOT_PASSWORD" //nolint:gosec
	MongoAuthDBKey    = "MONGO_AUTH_DB"
)
