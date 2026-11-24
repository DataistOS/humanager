// lib/core/services/upupa_bridge_service.dart
import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:upupa_dart/upupa_dart.dart';

class UpupaBridgeService {
  static const String defaultBaseUrl = 'http://127.0.0.1:8090';

  final UpupaClient client;
  final String baseUrl;

  UpupaBridgeService({String baseUrl = defaultBaseUrl, AuthStore? authStore})
    : baseUrl = baseUrl,
      client = UpupaClient(baseUrl, authStore: authStore);

  Future<List<Record>> fetchRecords(
    String collectionName, {
    int page = 1,
    int perPage = 30,
    String? filter,
    String? sort,
  }) {
    return client.records.getList(
      collectionName,
      page: page,
      perPage: perPage,
      filter: filter,
      sort: sort,
    );
  }

  Future<Record> fetchRecord(String collectionName, String recordId) {
    return client.records.getOne(collectionName, recordId);
  }

  Future<Record> updateAvatarPath(
    String collectionName,
    String recordId,
    String relativePath, {
    String fieldName = 'avatarpath',
  }) {
    return client.records.update(collectionName, recordId, {
      fieldName: relativePath,
    });
  }

  Future syncFile({
    required String filePath,
    required String fileName,
    required String collection,
    String? collectionId,
    String? recordId,
    String? slug,
    String? targetField,
    String? targetFilename,
  }) async {
    try {
      final uri = Uri.parse('$baseUrl/api/v1/files/sync');
      final request = http.MultipartRequest('POST', uri);

      request.files.add(
        await http.MultipartFile.fromPath('file', filePath, filename: fileName),
      );

      request.fields['collection'] = collection;
      if (collectionId != null && collectionId.isNotEmpty) {
        request.fields['collection_id'] = collectionId;
      }
      if (recordId != null && recordId.isNotEmpty) {
        request.fields['record_id'] = recordId;
      }
      if (slug != null && slug.isNotEmpty) {
        request.fields['slug'] = slug;
      }
      if (targetField != null && targetField.isNotEmpty) {
        request.fields['target_field'] = targetField;
      }
      if (targetFilename != null && targetFilename.isNotEmpty) {
        request.fields['target_filename'] = targetFilename;
      }

      final streamedResponse = await request.send();
      final response = await http.Response.fromStream(streamedResponse);

      if (response.statusCode >= 200 && response.statusCode < 300) {
        final decoded = jsonDecode(response.body);
        if (decoded is Map) {
          return decoded;
        }
      }
    } catch (_) {}
    return null;
  }

  void dispose() {
    client.dispose();
  }
}
