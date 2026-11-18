// lib/features/explorer/presentation/screens/explorer_screen.dart
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:file_picker/file_picker.dart';
import '../../../../core/services/core_service.dart';

class ExplorerScreen extends StatefulWidget {
  final String currentPath;
  const ExplorerScreen({super.key, this.currentPath = ''});

  @override
  State<ExplorerScreen> createState() => _ExplorerScreenState();
}

class _ExplorerScreenState extends State<ExplorerScreen> {
  late Future _dirFuture;
  bool _isUploading = false;

  @override
  void initState() {
    super.initState();
    _loadDirectory();
  }

  void _loadDirectory() {
    _dirFuture = HuMaNagerCoreService.fetchDirectoryContents(
      widget.currentPath,
    );
  }

  Future<void> _pickAndUploadFile() async {
    try {
      FilePickerResult? result = await FilePicker.pickFiles();

      if (result != null && result.files.single.path != null) {
        String filePath = result.files.single.path!;
        String fileName = result.files.single.name;

        setState(() {
          _isUploading = true;
        });

        bool success = await HuMaNagerCoreService.uploadFile(
          filePath,
          fileName,
        );

        if (!mounted) return;

        setState(() {
          _isUploading = false;
        });

        if (success) {
          if (!mounted) return;
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('فایل "$fileName" با موفقیت آپلود و طبقه‌بندی شد!'),
            ),
          );
          _loadDirectory();
        } else {
          if (!mounted) return;
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(content: Text('خطا در آپلود یا طبقه‌بندی فایل')),
          );
        }
      }
    } catch (e) {
      if (!mounted) return;
      setState(() {
        _isUploading = false;
      });
      if (kDebugMode) {
        print('Pick file error: $e');
      }
    }
  }

  void _showCreateFolderDialog(BuildContext context) {
    final TextEditingController controller = TextEditingController();

    showDialog(
      context: context,
      builder: (dialogContext) => AlertDialog(
        title: const Text('ساخت پوشه جدید در کالبد'),
        content: TextField(
          controller: controller,
          decoration: const InputDecoration(
            hintText: 'نام پوشه (مثلاً memory_cache)',
          ),
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(dialogContext),
            child: const Text('لغو'),
          ),
          ElevatedButton(
            onPressed: () async {
              if (controller.text.isNotEmpty) {
                Navigator.pop(dialogContext);
                bool success = await HuMaNagerCoreService.createFolder(
                  widget.currentPath,
                  controller.text,
                );
                if (success) {
                  if (!mounted) return;
                  setState(() {
                    _loadDirectory();
                  });
                }
              }
            },
            child: const Text('ساخت'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          widget.currentPath.isEmpty ? 'HuMaNager - Root' : widget.currentPath,
        ),
      ),
      body: _isUploading
          ? const Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  CircularProgressIndicator(),
                  SizedBox(height: 16),
                  Text('در حال انتقال و دسته‌بندی هوشمند فایل به کالبد...'),
                ],
              ),
            )
          : FutureBuilder(
              future: _dirFuture,
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const Center(child: CircularProgressIndicator());
                }
                if (snapshot.hasError || snapshot.data == null) {
                  return const Center(
                    child: Text('⚠️ خطا در دریافت اطلاعات از دیمن گو'),
                  );
                }

                final data = snapshot.data as Map;
                final items = data['items'] ?? [];

                if (items.isEmpty) {
                  return const Center(child: Text('📂 این پوشه خالی است'));
                }

                return ListView.builder(
                  itemCount: items.length,
                  itemBuilder: (context, index) {
                    final item = items[index];
                    final bool isDir = item['is_dir'];
                    final String name = item['name'];
                    final String path = item['path'];

                    return Card(
                      child: ListTile(
                        leading: Icon(
                          isDir ? Icons.folder : Icons.insert_drive_file,
                          color: isDir ? Colors.amber : Colors.blueAccent,
                        ),
                        title: Text(name),
                        trailing: isDir
                            ? const Icon(Icons.arrow_forward_ios, size: 16)
                            : null,
                        onTap: () {
                          if (isDir) {
                            Navigator.push(
                              context,
                              MaterialPageRoute(
                                builder: (context) =>
                                    ExplorerScreen(currentPath: path),
                              ),
                            );
                          } else {}
                        },
                      ),
                    );
                  },
                );
              },
            ),
      floatingActionButton: Column(
        mainAxisAlignment: MainAxisAlignment.end,
        children: [
          FloatingActionButton(
            heroTag: 'upload_file',
            onPressed: _isUploading ? null : _pickAndUploadFile,
            backgroundColor: Colors.green,
            child: const Icon(Icons.cloud_upload),
          ),
          const SizedBox(height: 12),
          FloatingActionButton(
            heroTag: 'create_folder',
            onPressed: () => _showCreateFolderDialog(context),
            backgroundColor: Colors.blueAccent,
            child: const Icon(Icons.create_new_folder),
          ),
        ],
      ),
    );
  }
}
