# Changelog

## [0.2.0] - 2026-11-24

### Added
- Added `UpupaBridgeService` (`lib/core/services/upupa_bridge_service.dart`) to forward files and record metadata to Upupa (`POST /api/v1/files/sync`) and update target record fields (`avatarpath`).
- Added `HuMaNagerCoreService.uploadViaUpupaSync` (`lib/core/services/core_service.dart`) to support synchronized uploads through Upupa.
- Added `HandleGetTree` (`core/internal/api/handler.go`) and config path resolution in the daemon (`core/cmd/daemon/main.go`).
- Added file path inspection and clipboard copy dialog in `ExplorerScreen` (`lib/features/explorer/presentation/screens/explorer_screen.dart`).

### Changed
- Aligned storage root (`~/5ec0ndbr41n`), anatomy tree generation (`core/internal/engine/tree_builder.go`), and `core/configs/anatomy_tree.json` with the on-disk anatomical hierarchy.
- Updated `HandleListDirectory` and `SetupRouter` (`core/internal/api/router.go`) to return normalized slash paths, non-null item lists, and live tree data.
- Updated `lib/main.dart` and `ExplorerScreen` navigation to resolve paths under `biological_system` and refresh directory state after uploads and folder creation.
- Updated collection/file routing in `core/internal/routing/classifier.go`, Flutter dependencies (`pubspec.yaml`, `pubspec.lock`), and `VERSION` to `0.2.0`.

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