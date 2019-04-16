package db

import (
	"encoding/json"
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
	//Groups      []Group `json:"groups" readonly:"true" embedded:"slice" alias:"grp" table:"groups" on:"grp.id = grp_usr.group_id" relation_alias:"grp_usr" relation_table:"groups_users" relation_on:"users.id = grp_usr.user_id"`
}

type SimpleUser struct {
	ID        string  `json:"id" sql:"id" pk:"true"`
	FirstName string  `json:"firstName" sql:"first_name"`
	LastName  string  `json:"lastName" sql:"last_name"`
	Email     string  `json:"email" sql:"email"`
	Groups    []Group `json:"groups" readonly:"true" embedded:"slice" alias:"grp" table:"groups" on:"grp.id = grp_usr.group_id" relation_alias:"grp_usr" relation_table:"groups_users" relation_on:"users.id = grp_usr.user_id"`
}

type Group struct {
	ID     string `json:"id" sql:"id" pk:"true"`
	Code   string `json:"code" sql:"code"`
	Active bool   `json:"active" sql:"active"`
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

	err := UpdateStruct("users", user, builder.Equal("id", suite.InstanceID), "email")
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	assert.NoError(suite.T(), err, msg)
}

func (suite *ActionsTestSuite) Test003LoadStruct() {
	user := &SimpleUser{}
	jsonByte, err := LoadStruct("users", user, builder.Equal("users.id", "059fa339-025c-4104-ab55-c764d3028f63"))
	json.Unmarshal(jsonByte, user)
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	assert.NoError(suite.T(), err, msg)
}

func (suite *ActionsTestSuite) Test004LoadStructArray() {
	users := []User{}

	jsonByte, err := LoadStruct("users", users, nil)
	json.Unmarshal(jsonByte, &users)
	msg := ""
	if err != nil {
		msg = err.Error()
	}
	assert.NoError(suite.T(), err, msg)
}

func (suite *ActionsTestSuite) Test005DeleteStruct() {
	err := DeleteStruct("users", builder.Equal("id", suite.InstanceID))
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
