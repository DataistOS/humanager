// lib/main.dart
import 'package:flutter/material.dart';
import 'core/services/core_service.dart';
import 'features/explorer/presentation/screens/explorer_screen.dart'
    as feature_explorer;

void main() {
  runApp(const HuMaNagerApp());
}

class HuMaNagerApp extends StatelessWidget {
  const HuMaNagerApp({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'HuMaNager',
      theme: ThemeData.dark(),
      home: const AnatomyRootScreen(),
    );
  }
}

class AnatomyRootScreen extends StatefulWidget {
  const AnatomyRootScreen({super.key});

  @override
  State<AnatomyRootScreen> createState() => _AnatomyRootScreenState();
}

class _AnatomyRootScreenState extends State<AnatomyRootScreen> {
  late Future _treeFuture;

  @override
  void initState() {
    super.initState();
    _treeFuture = HuMaNagerCoreService.fetchAnatomyTree();
  }

  String _resolveInitialPath(String itemName) {
    if (itemName == 'biological_system' ||
        itemName.startsWith('biological_system/')) {
      return itemName;
    }
    return 'biological_system/$itemName';
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('HuMaNager - Anatomical Explorer'),
        actions: [
          IconButton(
            icon: const Icon(Icons.folder_open),
            tooltip: 'باز کردن ریشه مخزن (~/5ec0ndbr41n)',
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) =>
                      const feature_explorer.ExplorerScreen(currentPath: ''),
                ),
              );
            },
          ),
        ],
      ),
      body: FutureBuilder(
        future: _treeFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError || snapshot.data == null) {
            return const Center(
              child: Text('⚠️ در حال اتصال به موتور گو (دیمن خاموش است)...'),
            );
          }

          final data = snapshot.data as Map;
          final rootName = (data['root'] ?? '~/5ec0ndbr41n').toString();
          final rawSystems = data['systems'];
          final List<String> systems = rawSystems is List
              ? rawSystems.map((e) => e.toString()).toList()
              : <String>[];

          return ListView(
            padding: const EdgeInsets.all(16.0),
            children: [
              Text(
                'Root Storage: $rootName  (~/5ec0ndbr41n)',
                style: const TextStyle(
                  fontSize: 18,
                  fontWeight: FontWeight.bold,
                ),
              ),
              const Divider(),
              const Text(
                'Biological Systems:',
                style: TextStyle(fontSize: 14, color: Colors.grey),
              ),
              const SizedBox(height: 8),
              ...systems.map((system) {
                final targetPath = _resolveInitialPath(system);
                return Card(
                  child: ListTile(
                    leading: const Icon(Icons.folder, color: Colors.blueAccent),
                    title: Text(system),
                    subtitle: Text(
                      targetPath,
                      style: const TextStyle(fontSize: 12, color: Colors.grey),
                    ),
                    trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                    onTap: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => feature_explorer.ExplorerScreen(
                            currentPath: targetPath,
                          ),
                        ),
                      );
                    },
                  ),
                );
              }),
            ],
          );
        },
      ),
    );
  }
}
