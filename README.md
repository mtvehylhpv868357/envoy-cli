# envoy-cli

A lightweight CLI for managing and switching between environment variable profiles across projects.

---

## Installation

```bash
go install github.com/yourusername/envoy-cli@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envoy-cli/releases).

---

## Usage

Create and manage named environment profiles, then load them into your current shell session.

```bash
# Create a new profile
envoy create myproject --set DB_URL=postgres://localhost/dev --set DEBUG=true

# List all profiles
envoy list

# Switch to a profile
envoy use myproject

# Show variables in a profile
envoy show myproject

# Remove a profile
envoy delete myproject
```

Profiles are stored locally in `~/.envoy/profiles` and can be scoped per directory using the `--local` flag.

```bash
# Create a project-local profile
envoy create staging --local --set API_KEY=abc123
```

---

## Why envoy-cli?

- No dependencies or daemons — just a single binary
- Works with any shell (bash, zsh, fish)
- Safe secret handling — values are stored with restricted file permissions
- Easy to integrate into CI/CD pipelines

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any major changes.

---

## License

MIT © [yourusername](https://github.com/yourusername)