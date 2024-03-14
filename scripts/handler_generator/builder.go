package main

import (
	"cmp"
	"fmt"
	"go/ast"
	"reflect"
	"slices"
	"strconv"
)

type switchCase struct {
	EventsubMessageType    string
	EventsubMessageVersion string
	EventType              string
	HandlerFieldName       string
}

func getSubHandler(node ast.Node) (*ast.StructType, bool) {
	var subHandler *ast.StructType
	ast.Inspect(node, func(n ast.Node) bool {
		ts, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		if ts.Name.Name == "SubHandler" {
			subHandler, _ = ts.Type.(*ast.StructType)
			return false
		}
		return true
	})
	return subHandler, subHandler != nil
}

func getTagValue(tag string, key string) (string, bool) {
	unquoted, err := strconv.Unquote(tag)
	if err != nil {
		return "", false
	}
	structTag := reflect.StructTag(unquoted)
	return structTag.Lookup(key)
}

func buildHandlerCases(subHandler *ast.StructType) []switchCase {
	var cases []switchCase

	for _, field := range subHandler.Fields.List {
		if field.Tag == nil {
			continue
		}

		eventsubMessageType, ok := getTagValue(field.Tag.Value, "eventsub-type")
		if !ok {
			continue
		}

		eventsubMessageVersion, ok := getTagValue(field.Tag.Value, "eventsub-version")
		if !ok {
			fmt.Printf("Warning: skipping message %q because no eventsub-version tag", eventsubMessageType)
			continue
		}

		typ := field.Type.(*ast.IndexExpr)
		eventTyp := typ.Index.(*ast.SelectorExpr)
		eventTypeName := fmt.Sprintf("%s.%s", eventTyp.X.(*ast.Ident), eventTyp.Sel.Name)

		cases = append(cases, switchCase{
			EventsubMessageType:    eventsubMessageType,
			EventsubMessageVersion: eventsubMessageVersion,
			EventType:              eventTypeName,
			HandlerFieldName:       field.Names[0].Name,
		})
	}

	// Sort by eventsub message type
	slices.SortFunc(cases, func(a, b switchCase) int {
		return cmp.Compare(a.EventsubMessageType, b.EventsubMessageType)
	})

	return cases
}
