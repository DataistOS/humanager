# Changelog

## [0.1.0] - 2026-11-18

### Added
- **Backend (Go Daemon):**
    - Initialized Go backend architecture (`humanager/internal`) supporting modular anatomy tree configuration.
    - Implemented core REST API router and handlers (`/api/health`, `/api/tree`, `/api/mkdir`, `/api/upload`) running on port `8081`.
    - Added hierarchical anatomical structure mapping via `anatomy_tree.json` under `core/configs/`.
    - Implemented smart disk manager and file classification engine (`Classifier`) targeting structured storage (`~/5ec0ndbr41n`).
    - Resolved Go module compatibility (updated to Go 1.20+ for `any` type support).

- **Frontend (Flutter App):**
    - Developed cross-platform Flutter frontend ("Anatomical Explorer") integrated with the Go backend API.
    - Added biological systems explorer interface (`digestive_system`, `genital_system`, `nervous_system`, etc.).
    - Integrated `file_picker` and HTTP service layers for seamless file uploading and browsing.

- **DevOps & Tooling:**
    - Configured comprehensive unified `.gitignore` for both Go and Flutter/Dart environments.
    - Organized project repository structure and pubspec dependencies.