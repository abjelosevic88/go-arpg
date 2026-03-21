package tmx

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/abjeloevic88/go-game/asset"
	"github.com/abjeloevic88/go-game/tilemap"
	"github.com/hajimehoshi/ebiten/v2"
)

// mapFile is the root <map> element.
type mapFile struct {
	XMLName  xml.Name     `xml:"map"`
	Width    int          `xml:"width,attr"`
	Height   int          `xml:"height,attr"`
	TileW    int          `xml:"tilewidth,attr"`
	TileH    int          `xml:"tileheight,attr"`
	Tilesets []tilesetRef `xml:"tileset"`
	Layers   []layerData  `xml:"layer"`
}

// layerData is a <layer> element containing tile data.
type layerData struct {
	Name string  `xml:"name,attr"`
	Data csvData `xml:"data"`
}

// csvData is the <data encoding="csv"> element.
type csvData struct {
	Encoding string `xml:"encoding,attr"`
	Content  string `xml:",chardata"`
}

// parseCSV splits CSV content into a 2D grid of integer GIDs.
func parseCSV(csv string, width, height int) ([][]int, error) {
	grid := make([][]int, 0, height)
	for _, line := range strings.Split(strings.TrimSpace(csv), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Split(line, ",")
		// Trim trailing empty field from lines ending with a comma
		if len(fields) > 0 && strings.TrimSpace(fields[len(fields)-1]) == "" {
			fields = fields[:len(fields)-1]
		}
		row := len(grid)
		if len(fields) != width {
			return nil, fmt.Errorf("tmx: row %d has %d values, want %d", row, len(fields), width)
		}
		vals := make([]int, width)
		for col, f := range fields {
			f = strings.TrimSpace(f)
			v, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("tmx: invalid GID %q at row %d col %d: %w", f, row, col, err)
			}
			vals[col] = v
		}
		grid = append(grid, vals)
	}
	if len(grid) != height {
		return nil, fmt.Errorf("tmx: got %d rows, want %d", len(grid), height)
	}
	return grid, nil
}

// applyGIDMap replaces GIDs in the grid with tile indices from gidMap.
// Unknown GIDs and GID 0 map to tile index 0.
func applyGIDMap(grid [][]int, gidMap map[int]int) [][]int {
	for r, row := range grid {
		for c, gid := range row {
			if idx, ok := gidMap[gid]; ok && gid != 0 {
				grid[r][c] = idx
			} else {
				grid[r][c] = 0
			}
		}
	}
	return grid
}


// LoadMap parses a TMX file at path, translates GIDs using gidMap,
// and returns a *tilemap.Tilemap wired to the given TileSet.
// GIDs not found in gidMap (including GID 0) default to tile index 0.
func LoadMap(path string, gidMap map[int]int, ts *tilemap.TileSet) (*tilemap.Tilemap, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("tmx: open %s: %w", path, err)
	}
	defer f.Close()

	var mf mapFile
	if err := xml.NewDecoder(f).Decode(&mf); err != nil {
		return nil, fmt.Errorf("tmx: decode %s: %w", path, err)
	}

	if len(mf.Layers) == 0 {
		return nil, fmt.Errorf("tmx: no layers in %s", path)
	}

	content := mf.Layers[0].Data.Content
	grid, err := parseCSV(content, mf.Width, mf.Height)
	if err != nil {
		return nil, err
	}

	grid = applyGIDMap(grid, gidMap)

	return &tilemap.Tilemap{
		Cols:    mf.Width,
		Rows:    mf.Height,
		TileW:   mf.TileW,
		TileH:   mf.TileH,
		Grid:    grid,
		TileSet: ts,
	}, nil
}

// TSXTileset represents a parsed TSX tileset file.
type TSXTileset struct {
	XMLName    xml.Name `xml:"tileset"`
	Name       string   `xml:"name,attr"`
	TileWidth  int      `xml:"tilewidth,attr"`
	TileHeight int      `xml:"tileheight,attr"`
	TileCount  int      `xml:"tilecount,attr"`
	Columns    int      `xml:"columns,attr"`
	Image      TSXImage `xml:"image"`
}

// TSXImage represents the <image> element inside a TSX tileset.
type TSXImage struct {
	Source string `xml:"source,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

// tilesetRef represents a <tileset> element in a TMX file.
type tilesetRef struct {
	FirstGID int    `xml:"firstgid,attr"`
	Source   string `xml:"source,attr"`
}

// LoadTSX parses a TSX file at the given path and returns the tileset metadata.
func LoadTSX(path string) (*TSXTileset, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("tmx: open %s: %w", path, err)
	}
	defer f.Close()

	var ts TSXTileset
	if err := xml.NewDecoder(f).Decode(&ts); err != nil {
		return nil, fmt.Errorf("tmx: decode TSX %s: %w", path, err)
	}
	return &ts, nil
}

// MarshalTSX serializes a TSXTileset back to XML bytes.
func MarshalTSX(ts *TSXTileset) ([]byte, error) {
	return xml.MarshalIndent(ts, "", "  ")
}

// SliceSpritesheet slices a spritesheet image into individual tile sub-images.
// It returns a slice of length tileCount where index i is the sub-image at
// pixel position ((i % columns) * tileW, (i / columns) * tileH).
func SliceSpritesheet(sheet *ebiten.Image, tileW, tileH, columns, tileCount int) ([]*ebiten.Image, error) {
	tiles := make([]*ebiten.Image, tileCount)
	for i := 0; i < tileCount; i++ {
		x := (i % columns) * tileW
		y := (i / columns) * tileH
		sub, err := asset.SubImage(sheet, x, y, tileW, tileH)
		if err != nil {
			return nil, fmt.Errorf("tmx: slice tile %d: %w", i, err)
		}
		tiles[i] = sub
	}
	return tiles, nil
}

// BuildTileSet creates a TileSet mapping firstGID + i to tiles[i] for each sub-image.
func BuildTileSet(tiles []*ebiten.Image, firstGID int) *tilemap.TileSet {
	m := make(map[int]*ebiten.Image, len(tiles))
	for i, img := range tiles {
		m[firstGID+i] = img
	}
	return &tilemap.TileSet{Tiles: m}
}

// LoadMapFromTMX parses a TMX file, resolves and loads the referenced TSX
// tileset and spritesheet, and returns a fully assembled Tilemap with raw GIDs
// stored in the grid (no GID translation).
func LoadMapFromTMX(tmxPath string) (*tilemap.Tilemap, error) {
	f, err := os.Open(tmxPath)
	if err != nil {
		return nil, fmt.Errorf("tmx: open %s: %w", tmxPath, err)
	}
	defer f.Close()

	var mf mapFile
	if err := xml.NewDecoder(f).Decode(&mf); err != nil {
		return nil, fmt.Errorf("tmx: decode %s: %w", tmxPath, err)
	}

	if len(mf.Tilesets) == 0 {
		return nil, fmt.Errorf("tmx: no tileset in %s", tmxPath)
	}

	ref := mf.Tilesets[0]
	tsxPath := filepath.Join(filepath.Dir(tmxPath), ref.Source)

	tsx, err := LoadTSX(tsxPath)
	if err != nil {
		return nil, err
	}

	imgPath := filepath.Join(filepath.Dir(tsxPath), tsx.Image.Source)
	sheet, err := asset.LoadImage(imgPath)
	if err != nil {
		return nil, err
	}

	tiles, err := SliceSpritesheet(sheet, tsx.TileWidth, tsx.TileHeight, tsx.Columns, tsx.TileCount)
	if err != nil {
		return nil, err
	}

	ts := BuildTileSet(tiles, ref.FirstGID)

	if len(mf.Layers) == 0 {
		return nil, fmt.Errorf("tmx: no layers in %s", tmxPath)
	}

	grid, err := parseCSV(mf.Layers[0].Data.Content, mf.Width, mf.Height)
	if err != nil {
		return nil, err
	}

	return &tilemap.Tilemap{
		Cols:    mf.Width,
		Rows:    mf.Height,
		TileW:   mf.TileW,
		TileH:   mf.TileH,
		Grid:    grid,
		TileSet: ts,
	}, nil
}
