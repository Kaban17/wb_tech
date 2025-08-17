# Go Cut Utility

A simple command-line utility, analogous to `cut`, for splitting lines by a delimiter and printing selected fields.

## Features

-   Select specific fields (`-f`) by number, range, or a combination.
-   Specify a custom delimiter (`-d`).
-   Optionally, only print lines that contain the delimiter (`-s`).

## Usage

The utility reads from STDIN.

### Flags

| Flag | Description                                     | Default    |
| ---- | ----------------------------------------------- | ---------- |
| `-f` | Select only these fields (e.g., "1,3", "2-4")   | (required) |
| `-d` | Use a custom delimiter instead of TAB           | `	`       |
| `-s` | Only print lines containing the delimiter       | `false`    |

### Examples

#### Basic Field Selection

Given the input `a,b,c`:

```sh
echo "a,b,c" | go run main.go -f "1,3" -d ","
```

**Output:**

```
a,c
```

#### Using a Different Delimiter

Given the input `one two three`:

```sh
echo "one two three" | go run main.go -f "1,3" -d " "
```

**Output:**

```
one three
```

#### Filtering Lines without the Delimiter

Given the input:
```
a,b
c d
e,f
```

```sh
printf "a,b\nc d\ne,f" | go run main.go -f "1,2" -d "," -s
```

**Output:**

```
a,b
e,f
```

#### Complex Field Selection with Ranges

Given the input `1:2:3:4:5:6`:

```sh
echo "1:2:3:4:5:6" | go run main.go -f "1,3-4,6" -d ":"
```

**Output:**

```
1:3:4:6
```

```