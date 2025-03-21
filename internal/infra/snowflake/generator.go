package snowflake

import (
	//"message-id-service/internal/domain"
	"id-generator/internal/domain"
	"strconv"

	"github.com/bwmarrin/snowflake"
)

type SnowflakeGenerator struct {
	node *snowflake.Node
}

func NewSnowflakeGenerator(nodeID int64) (domain.IDGenerator, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}
	return &SnowflakeGenerator{node: node}, nil
}

func (s *SnowflakeGenerator) GenerateID() (string, error) {
	id := s.node.Generate()
	return strconv.FormatInt(id.Int64(), 10), nil
}
