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
	ID          string `json:"id" sql:"id" pk:"true"`
	FirstName   string `json:"firstName" sql:"first_name"`
	LastName    string `json:"lastName" sql:"last_name"`
	Email       string `json:"email" sql:"email"`
	Description string `json:"description" sql:"value" alias:"description" table:"translations" on:"description.structure_id = users.id and description.structure_field = 'description'"`
	Profile     string `json:"profile" sql:"value" alias:"prf" table:"translations" on:"prf.structure_id = users.id and prf.structure_field = 'profile'"`
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
		FirstName:   "Teste",
		LastName:    "ORM",
		Email:       "teste@teste.com",
		Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		Profile:     "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
	}

	id, err := InsertStruct("users", user)
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	assert.NoError(suite.T(), err, msg)

	statement := builder.Insert("translations", "structure_type", "structure_field", "structure_id", "value", "language_code")
	statement.Values("user", "description", id, user.Description, "pt-br")
	//statement.Values("user", "profile", id, user.Profile, "pt-br")

	err = Exec(statement)
	if err != nil {
		msg = err.Error()
	}
	assert.NoError(suite.T(), err, msg)

	values := []interface{}{
		id,
		user.Profile,
	}
	rawQuery := builder.Raw("insert into translations (structure_type, structure_field, structure_id, value, language_code) values ('user', 'profile', $1, $2, 'pt-br')", values...)
	err = Exec(rawQuery)
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

	err := UpdateStruct("users", user, "email")
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	assert.NoError(suite.T(), err, msg)
}

func (suite *ActionsTestSuite) Test003LoadStruct() {
	user := &User{
		ID: suite.InstanceID,
	}

	err := LoadStruct("users", user)
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	assert.NoError(suite.T(), err, msg)
}

func (suite *ActionsTestSuite) Test004DeleteStruct() {
	user := &User{
		ID: suite.InstanceID,
	}

	err := DeleteStruct("users", user)
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	assert.NoError(suite.T(), err, msg)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestActionsSuite(t *testing.T) {
	suite.Run(t, new(ActionsTestSuite))
}
