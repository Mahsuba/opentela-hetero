# OpenTela Baseline Build

## 1. Objective

The objective of this experiment is to build the original OpenTela project without modifying its source code.

This establishes a working baseline for the OpenTela-Hetero research project.

Before introducing hardware capability-aware placement, we need to understand and verify that the original OpenTela system builds and runs successfully.

## 2. Experimental Environment

- Operating system: Ubuntu (WSL2)
- Architecture: AMD64 (x86_64)
- Go version: Go 1.27.1
- Repository: OpenTela
- Build branch: main
- Build commit: fc9875d

## 3. Build Procedure

The original OpenTela repository was cloned and built in the Ubuntu environment.

The official release build command was:

```bash
make build-release
```

The command successfully compiled the OpenTela project and generated two architecture-specific binaries.

## 4. Build Results

The release build completed successfully.

Generated binaries:

```text
build/release/otela-amd64
build/release/otela-arm64
```

The AMD64 binary is appropriate for the x86_64 development machine used in this experiment.

The ARM64 binary was also generated, but it was not tested on ARM hardware.

## 5. CLI Verification

The AMD64 binary was tested using the following command:

```bash
./build/release/otela-amd64 --help
```

The command executed successfully and displayed the available OpenTela commands, including:

- init
- start
- peer-id
- probe
- wallet
- instance
- update
- version

This confirms that the compiled binary can execute and expose its command-line interface.

## 6. Version Verification

The following command was executed:

```bash
./build/release/otela-amd64 version
```

Output:

```text
OpenTela version main (commit: fc9875d)
(built at: 2026-11-13T14:45:25-0800)
```

This confirms the version and commit information embedded in the generated binary.

## 7. Repository Integrity

After building OpenTela, the following command was executed from the original repository root:

```bash
git status --short
```

No modified or untracked files were reported.

This confirms that the baseline build did not leave visible source-code changes in the original repository's working tree.

## 8. Issues Encountered

The initial attempt to run the release build command from the repository root failed because the command was executed from the wrong directory.

This was corrected by navigating to the appropriate source directory before running the build.

The corrected build completed successfully.

This was a command-location issue, not a source-code compilation failure.

## 9. Conclusion

The original OpenTela project was successfully compiled in the Ubuntu WSL2 environment.

The generated AMD64 binary passed basic CLI and version verification.

The original repository's working tree remained clean after the experiment.

Therefore, we have established a working baseline for subsequent OpenTela-Hetero development and evaluation.

Future experiments can build on this baseline to investigate hardware capability discovery, metadata propagation, and compatibility-aware placement.

