package parser_test

import (
	"testing"

	"github.com/alongubkin/cpl-compiler/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeclarationOneID(t *testing.T) {
	program, err := parser.Parse("var1 : int;")
	require.NoError(t, err)
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{
			parser.Declaration{Names: []string{"var1"}, Type: parser.Integer},
		},
	}, program)
}

func TestDeclarationMultipeIDs(t *testing.T) {
	program, err := parser.Parse("var1, var2, var3 : float;")
	require.NoError(t, err)
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{
			parser.Declaration{Names: []string{"var1", "var2", "var3"}, Type: parser.Float},
		},
	}, program)
}

func TestDeclarationInvalidType(t *testing.T) {
	_, err := parser.Parse("var1, var2, var3 : uu;")
	assert.Error(t, err)
}

func TestMultipleDeclarations(t *testing.T) {
	program, err := parser.Parse("var1, var2 : int; var3 : float; var4,var5:int;")
	require.NoError(t, err)
	assert.EqualValues(t, &parser.Program{
		Declarations: []parser.Declaration{
			parser.Declaration{Names: []string{"var1", "var2"}, Type: parser.Integer},
			parser.Declaration{Names: []string{"var3"}, Type: parser.Float},
			parser.Declaration{Names: []string{"var4", "var5"}, Type: parser.Integer},
		},
	}, program)
}
