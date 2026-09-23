## Quick start

The erosion command generates a heightmap, runs the droplet simulation, and
writes the original, eroded, and difference images:

```sh
go run ./cmd/erosion \
	-width 512 \
	-height 512 \
	-depth 8 \
	-seed 42 \
	-droplets 200000 \
	-steps 30 \
	-beforeOutput before.png \
	-afterOutput after.png \
	-differenceOutput difference.png
```

The command produces:

- `before.png`: the generated 16-bit grayscale heightmap.
- `after.png`: the heightmap after erosion.
- `difference.png`: a diagnostic image with red for material removed and blue
  for material deposited.

| Before | After | Difference |
| --- | --- | --- |
| ![The before heightmap](https://github.com/Lukefl/ErosionSim/blob/main/before.png) | ![The after heightmap](https://github.com/Lukefl/ErosionSim/blob/main/after.png) | ![The difference](https://github.com/Lukefl/ErosionSim/blob/main/difference.png) |

For a quick default run:

```sh
go run ./cmd/erosion
```

## Generate only

Generate a 512x512 PNG:

```sh
go run ./cmd/heightmap
```

Choose dimensions, a seed, and an output path:

```sh
go run ./cmd/heightmap -width 1024 -height 768 -depth 8 -seed 7 -output terrain.png
```

The heightmap generator accepts these flags:

```sh
go run ./cmd/heightmap -help
```

| Flag | Default | Description |
| --- | ---: | --- |
| `-width` | `512` | Heightmap width in cells. |
| `-height` | `512` | Heightmap height in cells. |
| `-depth` | `4` | Number of ridged sine-wave detail layers. |
| `-seed` | `42` | Random seed for repeatable terrain. |
| `-output` | `heightmap.png` | Output PNG path. |

The erosion command uses the same `-width`, `-height`, `-depth`, and `-seed`
flags, plus:

| Flag | Default | Description |
| --- | ---: | --- |
| `-droplets` | `10000` | Number of simulated water droplets. |
| `-steps` | `30` | Maximum movement steps per droplet. |
| `-beforeOutput` | `heightmap.png` | Original terrain PNG path. |
| `-afterOutput` | `eroded.png` | Eroded terrain PNG path. |
| `-differenceOutput` | `difference.png` | Erosion diagnostic PNG path. |

`-depth` controls how many ridged sine-wave layers are summed. Lower values
create broader terrain; higher values add progressively finer peaks and
foothills. A value of `4` is used by default.

Run the tests with:

```sh
go test ./...
```