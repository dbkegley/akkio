### Creating a new identicon

This project uses [Bubbletea](https://github.com/charmbracelet/bubbletea)
to run interactively in the terminal.

To run the program natively with Go: `go mod download && go run main.go`

Or with docker: `docker run --rm -it ghcr.io/dbkegley/akkio:latest`

This program uses a 2d matrix (15 x 15) to hold the inputs. The dimensions are completely arbitrary
and the size of the matrix is controlled by the `const MAX = 15` at the top of `main.go`.
The input string is loaded into the matrix in
[diagonal slices](https://stackoverflow.com/questions/1779199/traverse-matrix-in-diagonal-strips)

Some things of note:
- Unicode characters are allowed.
- Running natively in a terminal that supports truecolor will yield
  different results than running in docker where we are
  [limited to 256 colors](https://forums.docker.com/t/docker-run-with-colorful-output/24542).
- Each background and foreground color is deterministic. After the base colors are initialized,
  each new input adds the `rune` ([int32](https://go.dev/blog/strings))
  value to the previous color value to produce a new color.

Some possible improvements:
- The colors produced aren't always very nice to look at.
  Could exclude specific color ranges
- Add toggles to change how the matrix is loaded.
  - by row
  - by column
- Add toggles for different types of matrix transformations
  - invert background/foreground colors of each cell
  - rotate the matrix
  - transpose the matrix

