package tmx

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/abjeloevic88/go-game/tilemap"
	"pgregory.net/rapid"
)

func TestParseCSV_HappyPath(t *testing.T) {
	csv := "1,2,3\n4,5,6\n"
	grid, err := parseCSV(csv, 3, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := [][]int{{1, 2, 3}, {4, 5, 6}}
	for r := range want {
		for c := range want[r] {
			if grid[r][c] != want[r][c] {
				t.Errorf("grid[%d][%d] = %d, want %d", r, c, grid[r][c], want[r][c])
			}
		}
	}
}

func TestParseCSV_TrailingComma(t *testing.T) {
	csv := "1,2,3,\n4,5,6,\n"
	grid, err := parseCSV(csv, 3, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if grid[0][0] != 1 || grid[1][2] != 6 {
		t.Errorf("unexpected grid values: %v", grid)
	}
}

func TestParseCSV_WhitespaceAndEmptyLines(t *testing.T) {
	csv := "\n  1, 2, 3 \n  4, 5, 6  \n\n"
	grid, err := parseCSV(csv, 3, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := [][]int{{1, 2, 3}, {4, 5, 6}}
	for r := range want {
		for c := range want[r] {
			if grid[r][c] != want[r][c] {
				t.Errorf("grid[%d][%d] = %d, want %d", r, c, grid[r][c], want[r][c])
			}
		}
	}
}

func TestParseCSV_WrongRowLength(t *testing.T) {
	csv := "1,2,3\n4,5\n"
	_, err := parseCSV(csv, 3, 2)
	if err == nil {
		t.Fatal("expected error for wrong row length")
	}
	want := "tmx: row 1 has 2 values, want 3"
	if err.Error() != want {
		t.Errorf("error = %q, want %q", err.Error(), want)
	}
}

func TestParseCSV_InvalidInteger(t *testing.T) {
	csv := "1,abc,3\n4,5,6\n"
	_, err := parseCSV(csv, 3, 2)
	if err == nil {
		t.Fatal("expected error for invalid integer")
	}
	if got := err.Error(); len(got) == 0 {
		t.Error("expected non-empty error message")
	}
}

func TestParseCSV_WrongRowCount(t *testing.T) {
	csv := "1,2,3\n"
	_, err := parseCSV(csv, 3, 2)
	if err == nil {
		t.Fatal("expected error for wrong row count")
	}
}

func TestApplyGIDMap_BasicMapping(t *testing.T) {
	grid := [][]int{{1, 38, 0}, {38, 1, 1}}
	gidMap := map[int]int{1: 0, 38: 1}
	result := applyGIDMap(grid, gidMap)
	want := [][]int{{0, 1, 0}, {1, 0, 0}}
	for r := range want {
		for c := range want[r] {
			if result[r][c] != want[r][c] {
				t.Errorf("result[%d][%d] = %d, want %d", r, c, result[r][c], want[r][c])
			}
		}
	}
}

func TestApplyGIDMap_UnknownGIDDefaultsToZero(t *testing.T) {
	grid := [][]int{{99, 100}}
	gidMap := map[int]int{1: 0}
	result := applyGIDMap(grid, gidMap)
	if result[0][0] != 0 || result[0][1] != 0 {
		t.Errorf("unknown GIDs should map to 0, got %v", result)
	}
}

func TestApplyGIDMap_GIDZeroDefaultsToZero(t *testing.T) {
	grid := [][]int{{0}}
	gidMap := map[int]int{0: 5}
	result := applyGIDMap(grid, gidMap)
	if result[0][0] != 0 {
		t.Errorf("GID 0 should always map to 0, got %d", result[0][0])
	}
}

func TestApplyGIDMap_EmptyGrid(t *testing.T) {
	grid := [][]int{}
	gidMap := map[int]int{1: 0}
	result := applyGIDMap(grid, gidMap)
	if len(result) != 0 {
		t.Errorf("expected empty grid, got %v", result)
	}
}

// Feature: tiled-map-loading, Property 2: GID Mapping Correctness
// Validates: Requirements 2.2, 2.3, 2.4
func TestProperty_GIDMappingCorrectness(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random grid dimensions
		width := rapid.IntRange(1, 64).Draw(t, "width")
		height := rapid.IntRange(1, 48).Draw(t, "height")

		// Generate a random 2D grid with values 0-100
		grid := make([][]int, height)
		for r := 0; r < height; r++ {
			grid[r] = make([]int, width)
			for c := 0; c < width; c++ {
				grid[r][c] = rapid.IntRange(0, 100).Draw(t, fmt.Sprintf("grid[%d][%d]", r, c))
			}
		}

		// Generate a random GID map with 0-20 entries, keys 0-100, values 0-10
		numEntries := rapid.IntRange(0, 20).Draw(t, "numEntries")
		gidMap := make(map[int]int, numEntries)
		for i := 0; i < numEntries; i++ {
			key := rapid.IntRange(0, 100).Draw(t, fmt.Sprintf("gidKey[%d]", i))
			val := rapid.IntRange(0, 10).Draw(t, fmt.Sprintf("gidVal[%d]", i))
			gidMap[key] = val
		}

		// Deep copy the grid before applying (applyGIDMap modifies in-place)
		original := make([][]int, height)
		for r := 0; r < height; r++ {
			original[r] = make([]int, width)
			copy(original[r], grid[r])
		}

		// Apply GID mapping
		result := applyGIDMap(grid, gidMap)

		// Assert each cell matches expected mapping
		for r := 0; r < height; r++ {
			for c := 0; c < width; c++ {
				gid := original[r][c]
				var expected int
				if gid == 0 {
					expected = 0
				} else if idx, ok := gidMap[gid]; ok {
					expected = idx
				} else {
					expected = 0
				}
				if result[r][c] != expected {
					t.Fatalf("grid[%d][%d]: gid=%d, got %d, want %d (gidMap=%v)",
						r, c, gid, result[r][c], expected, gidMap)
				}
			}
		}
	})
}

// Feature: tiled-map-loading, Property 1: CSV Parse Round-Trip
// Validates: Requirements 1.2
func TestProperty_CSVParseRoundTrip(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		width := rapid.IntRange(1, 64).Draw(t, "width")
		height := rapid.IntRange(1, 48).Draw(t, "height")

		// Generate a random 2D grid with values 0-999
		grid := make([][]int, height)
		for r := 0; r < height; r++ {
			grid[r] = make([]int, width)
			for c := 0; c < width; c++ {
				grid[r][c] = rapid.IntRange(0, 999).Draw(t, fmt.Sprintf("grid[%d][%d]", r, c))
			}
		}

		// Format the grid as CSV
		var sb strings.Builder
		for _, row := range grid {
			for c, v := range row {
				if c > 0 {
					sb.WriteByte(',')
				}
				fmt.Fprintf(&sb, "%d", v)
			}
			sb.WriteByte('\n')
		}

		// Parse the CSV and assert equality
		parsed, err := parseCSV(sb.String(), width, height)
		if err != nil {
			t.Fatalf("parseCSV failed: %v", err)
		}

		if len(parsed) != height {
			t.Fatalf("row count: got %d, want %d", len(parsed), height)
		}
		for r := 0; r < height; r++ {
			if len(parsed[r]) != width {
				t.Fatalf("row %d col count: got %d, want %d", r, len(parsed[r]), width)
			}
			for c := 0; c < width; c++ {
				if parsed[r][c] != grid[r][c] {
					t.Fatalf("grid[%d][%d]: got %d, want %d", r, c, parsed[r][c], grid[r][c])
				}
			}
		}
	})
}

// Unit tests for LoadMap
// Validates: Requirements 1.1, 1.2, 1.3, 1.4, 3.1, 3.2, 3.3, 3.4

func TestLoadMap_HappyPath(t *testing.T) {
	gidMap := map[int]int{1: 0, 38: 1}
	ts := &tilemap.TileSet{}
	tm, err := LoadMap("../maps/map.tmx", gidMap, ts)
	if err != nil {
		t.Fatalf("LoadMap failed: %v", err)
	}
	// Verify dimensions
	if tm.Cols != 32 {
		t.Errorf("Cols = %d, want 32", tm.Cols)
	}
	if tm.Rows != 24 {
		t.Errorf("Rows = %d, want 24", tm.Rows)
	}
	if tm.TileW != 32 {
		t.Errorf("TileW = %d, want 32", tm.TileW)
	}
	if tm.TileH != 32 {
		t.Errorf("TileH = %d, want 32", tm.TileH)
	}
	// Spot-check tile positions from actual TMX data:
	// Row 0: 1,38,38,... → GID 1=grass(0), GID 38=path(1)
	if tm.Grid[0][0] != 0 {
		t.Errorf("Grid[0][0] = %d, want 0 (grass for GID 1)", tm.Grid[0][0])
	}
	if tm.Grid[0][1] != 1 {
		t.Errorf("Grid[0][1] = %d, want 1 (path for GID 38)", tm.Grid[0][1])
	}
	// Row 1: 1,1,38,... → GID 1=grass(0), GID 1=grass(0)
	if tm.Grid[1][0] != 0 {
		t.Errorf("Grid[1][0] = %d, want 0 (grass for GID 1)", tm.Grid[1][0])
	}
	if tm.Grid[1][1] != 0 {
		t.Errorf("Grid[1][1] = %d, want 0 (grass for GID 1)", tm.Grid[1][1])
	}
	// Row 2, col 3: GID 19 → not in gidMap → defaults to 0
	if tm.Grid[2][3] != 0 {
		t.Errorf("Grid[2][3] = %d, want 0 (unknown GID 19 defaults to 0)", tm.Grid[2][3])
	}
}

func TestLoadMap_FileNotFound(t *testing.T) {
	gidMap := map[int]int{1: 0}
	ts := &tilemap.TileSet{}
	_, err := LoadMap("nonexistent.tmx", gidMap, ts)
	if err == nil {
		t.Fatal("expected error for non-existent file, got nil")
	}
}

func TestLoadMap_MalformedXML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.tmx")
	if err := os.WriteFile(path, []byte("<map><broken"), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	gidMap := map[int]int{1: 0}
	ts := &tilemap.TileSet{}
	_, err := LoadMap(path, gidMap, ts)
	if err == nil {
		t.Fatal("expected error for malformed XML, got nil")
	}
}

func TestLoadMap_TileSetReference(t *testing.T) {
	gidMap := map[int]int{1: 0, 38: 1}
	ts := &tilemap.TileSet{}
	tm, err := LoadMap("../maps/map.tmx", gidMap, ts)
	if err != nil {
		t.Fatalf("LoadMap failed: %v", err)
	}
	if tm.TileSet != ts {
		t.Error("returned Tilemap's TileSet pointer does not match the one passed in")
	}
}

// Feature: tiled-map-loading, Property 3: LoadMap Preserves TMX Dimensions and Grid Layout
// Validates: Requirements 1.1, 3.1, 3.2, 3.3
func TestProperty_LoadMapPreservesDimensionsAndGrid(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random map dimensions and tile sizes
		width := rapid.IntRange(1, 64).Draw(t, "width")
		height := rapid.IntRange(1, 48).Draw(t, "height")
		tileW := rapid.IntRange(8, 64).Draw(t, "tilewidth")
		tileH := rapid.IntRange(8, 64).Draw(t, "tileheight")

		// Generate random CSV GID data (values 1-100, GID 0 means empty)
		gids := make([][]int, height)
		for r := 0; r < height; r++ {
			gids[r] = make([]int, width)
			for c := 0; c < width; c++ {
				gids[r][c] = rapid.IntRange(1, 100).Draw(t, fmt.Sprintf("gid[%d][%d]", r, c))
			}
		}

		// Build CSV string
		var csvBuf strings.Builder
		for _, row := range gids {
			for c, v := range row {
				if c > 0 {
					csvBuf.WriteByte(',')
				}
				fmt.Fprintf(&csvBuf, "%d", v)
			}
			csvBuf.WriteByte('\n')
		}

		// Build TMX XML
		tmxXML := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<map version="1.10" orientation="orthogonal" width="%d" height="%d" tilewidth="%d" tileheight="%d">
 <layer name="test" width="%d" height="%d">
  <data encoding="csv">
%s
  </data>
 </layer>
</map>`, width, height, tileW, tileH, width, height, csvBuf.String())

		// Write to temp file
		dir, err := os.MkdirTemp("", "tmx-test-*")
		if err != nil {
			t.Fatalf("create temp dir: %v", err)
		}
		defer os.RemoveAll(dir)
		path := filepath.Join(dir, "test.tmx")
		if err := os.WriteFile(path, []byte(tmxXML), 0644); err != nil {
			t.Fatalf("write temp file: %v", err)
		}

		// Generate a random GID map (0-20 entries)
		numEntries := rapid.IntRange(0, 20).Draw(t, "numEntries")
		gidMap := make(map[int]int, numEntries)
		for i := 0; i < numEntries; i++ {
			key := rapid.IntRange(1, 100).Draw(t, fmt.Sprintf("gidKey[%d]", i))
			val := rapid.IntRange(0, 10).Draw(t, fmt.Sprintf("gidVal[%d]", i))
			gidMap[key] = val
		}

		// Call LoadMap
		tm, err := LoadMap(path, gidMap, &tilemap.TileSet{})
		if err != nil {
			t.Fatalf("LoadMap failed: %v", err)
		}

		// Assert dimensions match
		if tm.Cols != width {
			t.Fatalf("Cols: got %d, want %d", tm.Cols, width)
		}
		if tm.Rows != height {
			t.Fatalf("Rows: got %d, want %d", tm.Rows, height)
		}
		if tm.TileW != tileW {
			t.Fatalf("TileW: got %d, want %d", tm.TileW, tileW)
		}
		if tm.TileH != tileH {
			t.Fatalf("TileH: got %d, want %d", tm.TileH, tileH)
		}

		// Assert grid cells match expected GID-mapped values
		for r := 0; r < height; r++ {
			for c := 0; c < width; c++ {
				gid := gids[r][c]
				var expected int
				if gid == 0 {
					expected = 0
				} else if idx, ok := gidMap[gid]; ok {
					expected = idx
				} else {
					expected = 0
				}
				if tm.Grid[r][c] != expected {
					t.Fatalf("Grid[%d][%d]: gid=%d, got %d, want %d (gidMap=%v)",
						r, c, gid, tm.Grid[r][c], expected, gidMap)
				}
			}
		}
	})
}
