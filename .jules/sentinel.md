## 2025-05-18 - Cross-Platform Path Separators in Filename Sanitization
**Vulnerability:** Using `filepath.Separator` with `strings.ReplaceAll` to sanitize file names received from remote APIs only replaces the host OS path separator (e.g. `\` on Windows, `/` on POSIX).
**Learning:** Remote API inputs often send POSIX paths (`/`) regardless of the user's OS. On Windows, replacing `filepath.Separator` (`\`) leaves `/` untouched, enabling path traversal attacks.
**Prevention:** Always replace both `/` and `\` explicit separators when sanitizing untrusted file names or path components.
