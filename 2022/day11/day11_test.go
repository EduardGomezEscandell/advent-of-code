package day11_test

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/EduardGomezEscandell/AdventOfCode/2022/day11"
	"github.com/EduardGomezEscandell/AdventOfCode/2022/utils/testutils"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	err := testutils.CheckEnv()
	if err != nil {
		log.Printf("Setup: %v", err)
	}
	r := m.Run()
	os.Exit(r)
}

func TestReadData(t *testing.T) {
	testCases := map[string]struct {
		data []string
		want []day11.Monkey
	}{
		"single":  {data: example[:6], want: exampleMonkeys[:1]},
		"two":     {data: example[:13], want: exampleMonkeys[:2]},
		"example": {data: example, want: exampleMonkeys},
		"data":    {data: realData, want: realMonkeys},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			defer testutils.Backup(&day11.ReadDataFile)()
			day11.ReadDataFile = func() ([]byte, error) {
				return []byte(strings.Join(tc.data, "\n")), nil
			}

			got, err := day11.ParseInput()
			require.NoError(t, err)

			require.Equal(t, len(tc.want), len(got))
			for i := range got {
				requireEqualMonkeys(t, tc.want[i], got[i])
			}
		})
	}
}

func TestPart1(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		data []day11.Monkey
		want int
	}{
		"example": {data: exampleMonkeys, want: 10605},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := day11.Part1(tc.data)

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestPart2(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		data []day11.Monkey
		want int
	}{
		"empty": {data: []day11.Monkey{}, want: 0},
	}

	for name, tc := range testCases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := day11.Part2(tc.data)

			require.NoError(t, err)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestRealData(t *testing.T) {
	expected := `Result of part 1: 120056
Result of part 2: 0
`
	buff := new(bytes.Buffer)

	err := day11.Main(buff)

	require.NoError(t, err)
	require.Equal(t, expected, buff.String())
}

func requireEqualMonkeys(t *testing.T, m1, m2 day11.Monkey) {
	t.Helper()

	require.Equalf(t, m1.Id, m2.Id, "Mismatching monkey ids:\n%v\n%v", m1, m2)
	require.Equalf(t, m1.Inventory, m2.Inventory, "Mismatching monkey inventories:\n%v\n%v", m1, m2)
	require.Equalf(t, m1.TestValue, m2.TestValue, "Mismatching monkey test values:\n%v\n%v", m1, m2)
	require.Equalf(t, m1.ThrowTrue, m2.ThrowTrue, "Mismatching monkey throw true:\n%v\n%v", m1, m2)
	require.Equalf(t, m1.ThrowFalse, m2.ThrowFalse, "Mismatching monkey throw false:\n%v\n%v", m1, m2)
	require.Equalf(t, m1.Inspections, m2.Inspections, "Mismatching monkey inspectionss:\n%v\n%v", m1, m2)
	for i := uint64(0); i <= 5; i++ {
		require.Equalf(t, m1.ThrowFalse, m2.ThrowFalse, `Mismatching monkey functions:
%v
%v
Functions are not the same:
f1(%d) -> %d
f2(%d) -> %d`, m1, m2, i, m1.Inspect(i), i, m2.Inspect(i))
	}
}

var example = []string{
	"Monkey 0:",
	"  Starting items: 79, 98",
	"  Operation: new = old * 19",
	"  Test: divisible by 23",
	"    If true: throw to monkey 2",
	"    If false: throw to monkey 3",
	"",
	"Monkey 1:",
	"  Starting items: 54, 65, 75, 74",
	"  Operation: new = old + 6",
	"  Test: divisible by 19",
	"    If true: throw to monkey 2",
	"    If false: throw to monkey 0",
	"",
	"Monkey 2:",
	"  Starting items: 79, 60, 97",
	"  Operation: new = old * old",
	"  Test: divisible by 13",
	"    If true: throw to monkey 1",
	"    If false: throw to monkey 3",
	"",
	"Monkey 3:",
	"  Starting items: 74",
	"  Operation: new = old + 3",
	"  Test: divisible by 17",
	"    If true: throw to monkey 0",
	"    If false: throw to monkey 1",
}

var exampleMonkeys = []day11.Monkey{
	{
		Id:         0,
		Inventory:  []uint64{79, 98},
		Inspect:    func(x uint64) uint64 { return x * 19 },
		TestValue:  23,
		ThrowTrue:  2,
		ThrowFalse: 3,
	},
	{
		Id:         1,
		Inventory:  []uint64{54, 65, 75, 74},
		Inspect:    func(x uint64) uint64 { return x + 6 },
		TestValue:  19,
		ThrowTrue:  2,
		ThrowFalse: 0,
	},
	{
		Id:         2,
		Inventory:  []uint64{79, 60, 97},
		Inspect:    func(x uint64) uint64 { return x * x },
		TestValue:  13,
		ThrowTrue:  1,
		ThrowFalse: 3,
	},
	{
		Id:         3,
		Inventory:  []uint64{74},
		Inspect:    func(x uint64) uint64 { return x + 3 },
		TestValue:  17,
		ThrowTrue:  0,
		ThrowFalse: 1,
	},
}

var realData = []string{
	"Monkey 0:",
	"  Starting items: 89, 74",
	"  Operation: new = old * 5",
	"  Test: divisible by 17",
	"    If true: throw to monkey 4",
	"    If false: throw to monkey 7",
	"",
	"Monkey 1:",
	"  Starting items: 75, 69, 87, 57, 84, 90, 66, 50",
	"  Operation: new = old + 3",
	"  Test: divisible by 7",
	"    If true: throw to monkey 3",
	"    If false: throw to monkey 2",
	"",
	"Monkey 2:",
	"  Starting items: 55",
	"  Operation: new = old + 7",
	"  Test: divisible by 13",
	"    If true: throw to monkey 0",
	"    If false: throw to monkey 7",
	"",
	"Monkey 3:",
	"  Starting items: 69, 82, 69, 56, 68",
	"  Operation: new = old + 5",
	"  Test: divisible by 2",
	"    If true: throw to monkey 0",
	"    If false: throw to monkey 2",
	"",
	"Monkey 4:",
	"  Starting items: 72, 97, 50",
	"  Operation: new = old + 2",
	"  Test: divisible by 19",
	"    If true: throw to monkey 6",
	"    If false: throw to monkey 5",
	"",
	"Monkey 5:",
	"  Starting items: 90, 84, 56, 92, 91, 91",
	"  Operation: new = old * 19",
	"  Test: divisible by 3",
	"    If true: throw to monkey 6",
	"    If false: throw to monkey 1",
	"",
	"Monkey 6:",
	"  Starting items: 63, 93, 55, 53",
	"  Operation: new = old * old",
	"  Test: divisible by 5",
	"    If true: throw to monkey 3",
	"    If false: throw to monkey 1",
	"",
	"Monkey 7:",
	"  Starting items: 50, 61, 52, 58, 86, 68, 97",
	"  Operation: new = old + 4",
	"  Test: divisible by 11",
	"    If true: throw to monkey 5",
	"    If false: throw to monkey 4",
}

var realMonkeys = []day11.Monkey{
	{
		Id:         0,
		Inventory:  []uint64{89, 74},
		Inspect:    func(x uint64) uint64 { return x * 5 },
		TestValue:  17,
		ThrowTrue:  4,
		ThrowFalse: 7,
	},
	{
		Id:         1,
		Inventory:  []uint64{75, 69, 87, 57, 84, 90, 66, 50},
		Inspect:    func(x uint64) uint64 { return x + 3 },
		TestValue:  7,
		ThrowTrue:  3,
		ThrowFalse: 2,
	},
	{
		Id:         2,
		Inventory:  []uint64{55},
		Inspect:    func(x uint64) uint64 { return x + 7 },
		TestValue:  13,
		ThrowTrue:  0,
		ThrowFalse: 7,
	},
	{
		Id:         3,
		Inventory:  []uint64{69, 82, 69, 56, 68},
		Inspect:    func(x uint64) uint64 { return x + 5 },
		TestValue:  2,
		ThrowTrue:  0,
		ThrowFalse: 2,
	},
	{
		Id:         4,
		Inventory:  []uint64{72, 97, 50},
		Inspect:    func(x uint64) uint64 { return x + 2 },
		TestValue:  19,
		ThrowTrue:  6,
		ThrowFalse: 5,
	},
	{
		Id:         5,
		Inventory:  []uint64{90, 84, 56, 92, 91, 91},
		Inspect:    func(x uint64) uint64 { return x * 19 },
		TestValue:  3,
		ThrowTrue:  6,
		ThrowFalse: 1,
	},
	{
		Id:         6,
		Inventory:  []uint64{63, 93, 55, 53},
		Inspect:    func(x uint64) uint64 { return x * x },
		TestValue:  5,
		ThrowTrue:  3,
		ThrowFalse: 1,
	},
	{
		Id:         7,
		Inventory:  []uint64{50, 61, 52, 58, 86, 68, 97},
		Inspect:    func(x uint64) uint64 { return x + 4 },
		TestValue:  11,
		ThrowTrue:  5,
		ThrowFalse: 4,
	},
}
