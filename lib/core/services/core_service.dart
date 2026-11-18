// lib/core/services/core_service.dart
import 'dart:async';
import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;

class HuMaNagerCoreService {
  static const String baseUrl = 'http://localhost:8081/api';

  static Future<bool> checkHealth() async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/health'));
      if (response.statusCode == 200) {
        return true;
      }
    } catch (e) {
      if (kDebugMode) {
        print('Daemon connection error: $e');
      }
    }
    return false;
  }

  static Future fetchAnatomyTree() async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/tree'));
      if (response.statusCode == 200) {
        final decoded = jsonDecode(response.body);
        if (decoded is Map) {
          return decoded;
        }
      }
    } catch (e) {
      if (kDebugMode) {
        print('Failed to fetch anatomy tree: $e');
      }
    }
    return null;
  }

  static Future fetchDirectoryContents(String path) async {
    try {
      final response = await http.get(Uri.parse('$baseUrl/dir?path=$path'));
      if (response.statusCode == 200) {
        final decoded = jsonDecode(response.body);
        if (decoded is Map) {
          return decoded;
        }
      }
    } catch (e) {
      if (kDebugMode) {
        print('Failed to fetch directory contents: $e');
      }
    }
    return null;
  }

  static Future<bool> createFolder(
    String currentPath,
    String folderName,
  ) async {
    try {
      final response = await http.post(
        Uri.parse('$baseUrl/mkdir'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'path': currentPath, 'name': folderName}),
      );
      if (response.statusCode == 200) {
        final decoded = jsonDecode(response.body);
        return decoded['success'] == true;
      }
    } catch (e) {
      if (kDebugMode) {
        print('Failed to create folder: $e');
      }
    }
    return false;
  }

  static Future<bool> uploadFile(String filePath, String fileName) async {
    try {
      var request = http.MultipartRequest('POST', Uri.parse('$baseUrl/upload'));

      request.files.add(
        await http.MultipartFile.fromPath('file', filePath, filename: fileName),
      );

      var streamedResponse = await request.send();
      var response = await http.Response.fromStream(streamedResponse);

      if (response.statusCode == 200) {
        final decoded = jsonDecode(response.body);
        return decoded['success'] == true;
      }
    } catch (e) {
      if (kDebugMode) {
        print('Failed to upload file: $e');
      }
    }
    return false;
  }
}
