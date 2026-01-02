# Resources with unicode in names and titles

query "query_with_emoji_title" {
  title       = "Query with 🚀 Emoji Title"
  description = "Tests unicode emoji in title"
  sql         = "SELECT '🎉' as emoji, 'celebration' as name"
}

query "query_with_chinese" {
  title       = "中文标题 Chinese Title"
  description = "测试中文字符 Testing Chinese characters"
  sql         = "SELECT '你好' as greeting, 'chinese' as language"
}

query "query_with_arabic" {
  title       = "مرحبا Arabic Title"
  description = "Testing Arabic: مرحبا بك"
  sql         = "SELECT 'مرحبا' as greeting, 'arabic' as language"
}

query "query_with_accents" {
  title       = "Café Résumé Naïve"
  description = "Testing accented characters: é, è, ñ, ü, ö"
  sql         = "SELECT 'Ñoño' as name, 'Zürich' as city"
}

control "unicode_control" {
  title       = "🔒 Security Control with Unicode"
  description = "Проверка безопасности - Security check in Russian"
  sql         = "SELECT 'pass' as status, '资源' as resource, 'Unicode ✓' as reason"
}

dashboard "unicode_dashboard" {
  title       = "📊 Unicode Dashboard"
  description = "Dëscription with spëcial charactërs"

  card {
    title = "🌍 Global Count"
    sql   = query.query_with_emoji_title.sql
  }

  card {
    title = "日本語 Japanese"
    sql   = "SELECT 'こんにちは' as greeting"
  }
}
