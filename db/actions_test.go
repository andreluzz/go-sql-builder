package db

import (
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/andreluzz/go-sql-builder/builder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type User struct {
	ID            string `json:"id" sql:"id" pk:"true"`
	Username      string `json:"username" sql:"username"`
	FirstName     string `json:"firstname" sql:"first_name"`
	LastName      string `json:"lastname" sql:"last_name"`
	Email         string `json:"email" sql:"email"`
	Password      string `json:"password" sql:"password"`
	CreatedBy     string `json:"created_by" sql:"created_by"`
	UpdatedBy     string `json:"updated_by" sql:"updated_by"`
	CreatedByUser *User  `json:"created_by_user" table:"core_users" alias:"created_by_user" on:"created_by_user.id = core_users.created_by"`
	UpdatedByUser *User  `json:"updated_by_user" table:"core_users" alias:"updated_by_user" on:"updated_by_user.id = core_users.updated_by"`
}

type ActionsTestSuite struct {
	suite.Suite
	InstanceID string
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test3
func (suite *ActionsTestSuite) SetupTest() {
	var config Config
	toml.DecodeFile("config.toml", &config)
	Connect(config.Host, config.Port, config.User, config.Password, config.DBName, false)
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *ActionsTestSuite) Test001InsertStruct() {
	user := &User{
		FirstName: "Teste",
		LastName:  "ORM",
		Email:     "teste@teste.com",
		Password:  "12345",
	}

	id, err := InsertStruct("users", user)
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	assert.NoError(suite.T(), err, msg)

	suite.InstanceID = id

}

func (suite *ActionsTestSuite) Test002UpdateStruct() {
	user := &User{
		ID:    suite.InstanceID,
		Email: "user@teste.com",
	}

	err := UpdateStruct("users", user, builder.Equal("id", suite.InstanceID), "email")
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	assert.NoError(suite.T(), err, msg)
}

func (suite *ActionsTestSuite) Test003LoadStruct() {
	user := &User{}
	err := LoadStruct("users", user, builder.Equal("users.id", suite.InstanceID))
	assert.NoError(suite.T(), err, "Error loading struct")
	assert.Equal(suite.T(), "user@teste.com", user.Email)
}

func (suite *ActionsTestSuite) Test004LoadStructArray() {
	users := []User{}
	err := LoadStruct("users", &users, nil)
	assert.NoError(suite.T(), err, "Error loading array struct")
	assert.NotEmpty(suite.T(), users, "Empty array")
}

func (suite *ActionsTestSuite) Test005LoadEmbeddedStruct() {
	user := User{}
	err := LoadStruct("core_users", &user, builder.Equal("core_users.id", "57a97aaf-16da-44ef-a8be-b1caf52becd6"))
	assert.NoError(suite.T(), err, "Error loading array struct")
	assert.Equal(suite.T(), "admin", user.CreatedByUser.Username, "Invalid ceated by user first name")
	assert.Equal(suite.T(), "admin", user.UpdatedByUser.Username, "Invalid updated by user first name")
}

func (suite *ActionsTestSuite) Test006QueryStruct() {
	statement := builder.Select("id", "first_name", "email").From("users").Where(builder.Equal("id", suite.InstanceID))
	user := User{}
	err := QueryStruct(statement, &user)
	assert.NoError(suite.T(), err, "Error loading struct")
	assert.Equal(suite.T(), "user@teste.com", user.Email)
}

func (suite *ActionsTestSuite) Test007DeleteStruct() {
	err := DeleteStruct("users", builder.Equal("id", suite.InstanceID))
	assert.NoError(suite.T(), err, "Error deleting object")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestActionsSuite(t *testing.T) {
	suite.Run(t, new(ActionsTestSuite))
}
