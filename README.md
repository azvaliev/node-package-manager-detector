# Node Package Manager Detector

Detect the package manager being used in a node project and output it.

## Usage

Invoke binary in a directory of an npm project or pass in an absolute path.
Invoke with `--help` for details on arguments

Will output one of:
- `yarn`
- `yarn-classic`
- `npm`
- `pnpm`

OR `ERROR: ....` if there was an error detecting the package manager
