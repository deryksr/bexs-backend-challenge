package service

import (
	"errors"
	"reflect"
	"strconv"
	"testing"
)

func TestParseCitiesToString(test *testing.T) {
	testGroup := []struct {
		name     string
		input    []*City
		expected string
	}{
		{
			"Get empty string from empty input",
			[]*City{},
			"",
		},
		{
			"Get two cities as string",
			[]*City{&City{"A", false, nil}, &City{"B", false, nil}},
			"A - B",
		},
		{
			"Get two cities as string",
			[]*City{
				&City{"A", false, nil},
				&City{"B", false, nil},
				&City{"C", false, nil},
				&City{"D", false, nil},
			},
			"A - B - C - D",
		},
	}

	for _, testCase := range testGroup {
		test.Run(testCase.name, func(test *testing.T) {
			result := parseCitiesToString(testCase.input)

			if result != testCase.expected {
				test.Errorf("Expected <%v> but got <%v>", testCase.expected, result)
			}
		})
	}
}

func TestGetBestRoute(test *testing.T) {
	testGroup := []struct {
		name         string
		source       string
		target       string
		graphInput   [][]string
		resultWanted Route
		errorWanted  error
	}{
		{
			"The best route is the only route",
			"A",
			"B",
			[][]string{
				{"A", "B", "5"},
			},
			Route{
				Paths: []string{"A - B"},
				Cost:  5,
			},
			nil,
		},
		{
			"The cities are not connected",
			"B",
			"A",
			[][]string{
				{"A", "C", "5"},
				{"B", "C", "7"},
				{"C", "B", "1"},
			},
			Route{},
			errors.New("None route has been found between B - A"),
		},
		{
			"The source are not on the graph",
			"H",
			"A",
			[][]string{
				{"A", "C", "5"},
				{"B", "C", "7"},
				{"C", "B", "1"},
			},
			Route{},
			errors.New("source <H> has not found"),
		},
		{
			"The target are not on the graph",
			"A",
			"X",
			[][]string{
				{"A", "C", "5"},
				{"B", "C", "7"},
				{"C", "B", "1"},
			},
			Route{},
			errors.New("target <X> has not found"),
		},
		{
			"Has only one best route",
			"A",
			"D",
			[][]string{
				{"A", "B", "5"},
				{"B", "C", "9"},
				{"B", "D", "1"},
				{"C", "A", "2"},
				{"A", "D", "9"},
				{"A", "C", "3"},
			},
			Route{
				Paths: []string{"A - B - D"},
				Cost:  6,
			},
			nil,
		},
	}

	for _, testCase := range testGroup {
		test.Run(testCase.name, func(test *testing.T) {
			graph := GetGraphInstance()
			for _, line := range testCase.graphInput {
				origin := City{line[0], false, nil}
				destination := City{line[1], false, nil}
				cost, _ := strconv.Atoi(line[2])
				graph.AddRoad(&origin, &destination, cost)
			}

			bestRoute, err := GetBestRoute(testCase.source, testCase.target)

			if err != nil && err.Error() != testCase.errorWanted.Error() {
				test.Errorf(
					"TestGetBestRoute() got an unexpected error: - want: <%v> but got: <%v>",
					testCase.errorWanted,
					err,
				)
			}

			if !reflect.DeepEqual(bestRoute, testCase.resultWanted) {
				test.Errorf(
					"TestGetBestRoute() Expected: <%v> but got <%v>",
					testCase.resultWanted,
					bestRoute,
				)
			}
			CleanGraph()
		})
	}
}

func TestGetAllRoutes(test *testing.T) {
	var testGroup = []struct {
		name         string
		source       string
		target       string
		graphInput   [][]string
		resultWanted RouteList
		errorWanted  error
	}{
		{
			"The cities are not connected",
			"B",
			"A",
			[][]string{
				{"A", "C", "5"},
				{"B", "C", "7"},
				{"C", "B", "1"},
			},
			RouteList{},
			errors.New("None route has been found between B - A"),
		},
		{
			"The source are not on the graph",
			"H",
			"A",
			[][]string{
				{"A", "C", "5"},
				{"B", "C", "7"},
				{"C", "B", "1"},
			},
			RouteList{},
			errors.New("source <H> has not found"),
		},
		{
			"The target are not on the graph",
			"A",
			"X",
			[][]string{
				{"A", "C", "5"},
				{"B", "C", "7"},
				{"C", "B", "1"},
			},
			RouteList{},
			errors.New("target <X> has not found"),
		},
		{
			"Has two paths to get the target",
			"A",
			"C",
			[][]string{
				{"A", "B", "5"},
				{"B", "C", "7"},
				{"A", "C", "1"},
			},
			RouteList{
				Route{
					Paths: []string{"A - C"},
					Cost:  1,
				},
				Route{
					Paths: []string{"A - B - C"},
					Cost:  12,
				},
			},
			nil,
		},
		{
			"Get all paths to target",
			"A",
			"D",
			[][]string{
				{"A", "B", "5"},
				{"B", "C", "9"},
				{"B", "D", "1"},
				{"C", "A", "2"},
				{"A", "D", "9"},
				{"A", "C", "3"},
				{"C", "D", "8"},
			},
			[]Route{
				Route{
					Paths: []string{"A - B - D"},
					Cost:  6,
				},
				Route{
					Paths: []string{"A - D"},
					Cost:  9,
				},
				Route{
					Paths: []string{"A - C - D"},
					Cost:  11,
				},
				Route{
					Paths: []string{"A - B - C - D"},
					Cost:  22,
				},
			},
			errors.New(""),
		},
		{
			"Has two paths with the same cost",
			"A",
			"D",
			[][]string{
				{"A", "B", "5"},
				{"B", "C", "9"},
				{"B", "D", "1"},
				{"C", "A", "2"},
				{"A", "D", "6"},
				{"A", "C", "3"},
				{"C", "D", "8"},
			},
			[]Route{
				Route{
					Paths: []string{"A - B - D", "A - D"},
					Cost:  6,
				},
				Route{
					Paths: []string{"A - C - D"},
					Cost:  11,
				},
				Route{
					Paths: []string{"A - B - C - D"},
					Cost:  22,
				},
			},
			errors.New(""),
		},
	}

	for _, testCase := range testGroup {
		test.Run(testCase.name, func(test *testing.T) {
			graph := GetGraphInstance()
			for _, line := range testCase.graphInput {
				origin := City{line[0], false, nil}
				destination := City{line[1], false, nil}
				cost, _ := strconv.Atoi(line[2])
				graph.AddRoad(&origin, &destination, cost)
			}

			allRoutes, err := GetAllRoutes(testCase.source, testCase.target)
			if err != nil && err.Error() != testCase.errorWanted.Error() {
				test.Errorf(
					"TestGetAllRoutes() got an unexpected error: - want: <%v> but got: <%v>",
					testCase.errorWanted,
					err,
				)
			}

			if !reflect.DeepEqual(allRoutes, testCase.resultWanted) {
				test.Errorf(
					"GetAllRoutes() Expected: <%v> but got <%v>",
					testCase.resultWanted,
					allRoutes,
				)
			}
			CleanGraph()
		})
	}
}
