// main.dart
import 'package:flutter/material.dart';
import 'core/services/core_service.dart';

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
      home: const ExplorerScreen(),
    );
  }
}

class ExplorerScreen extends StatefulWidget {
  const ExplorerScreen({super.key});

  @override
  State<ExplorerScreen> createState() => _ExplorerScreenState();
}

class _ExplorerScreenState extends State<ExplorerScreen> {
  late Future<dynamic> _treeFuture;

  @override
  void initState() {
    super.initState();
    _treeFuture = HuMaNagerCoreService.fetchAnatomyTree();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('HuMaNager - Anatomical Explorer')),
      body: FutureBuilder<dynamic>(
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

          final data = snapshot.data as Map<String, dynamic>;
          final rootName = data['root'] ?? 'Unknown';
          final systems = List<String>.from(data['systems'] ?? []);

          return ListView(
            padding: const EdgeInsets.all(16.0),
            children: [
              Text(
                'Root Storage: $rootName',
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
              ...systems.map(
                (system) => Card(
                  child: ListTile(
                    leading: const Icon(Icons.folder, color: Colors.blueAccent),
                    title: Text(system),
                    trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                    onTap: () {},
                  ),
                ),
              ),
            ],
          );
        },
      ),
    );
  }
}
