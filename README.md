# bdspackman

Minecraft Bedrock Dedicated Serverでビヘイビアパック / リソースパックを管理するためのTUIツールです。

## どんなツール？

TUIでMinecraft Bedrock版の各種パックを管理できます。

- サーバ全体の`behavior_packs` / `resource_packs`と、ワールド配下の`behavior_packs` / `resource_packs`の4ディレクトリにあるパックを一覧表示
- `world_behavior_packs.json` / `world_resource_packs.json`の更新
- `--export-prefix`や`--export-dir`を使った`world_*_packs.json`の別名・別ディレクトリ出力
- 一覧にあるパックの有効 / 無効切り替え
- パックの適用順の入れ替え
- パックをサーバ全体のパックディレクトリとワールド配下のパックディレクトリの間で移動
- パックの削除（通常は`delpacks/`へ退避、`--force-delete`指定時は完全削除）
- パックディレクトリ名のリネーム
- UUID列、ディレクトリ名列、システムパック表示の切り替え

今後の目標は、パックの追加と更新の実装です。

## 動作確認済み環境

- Ubuntu 24.04 (amd64)

## 現在の実装状況

- ✅ パックの一覧表示
- ✅ `world_behavior_packs.json`の更新
- ✅ `world_resource_packs.json`の更新
- ✅ パックの有効化 / 無効化
- ✅ パックの適用順の入れ替え
- ✅ パックの移動
- ✅ パックの削除
- ✅ パックディレクトリ名のリネーム
- ✅ `world_*_packs.json`のエクスポート
- ❌ パックの追加
- ❌ パックの更新

## 使い方

リリースページからzipファイルをダウンロードしてください。

解凍すると実行ファイルができるので、PATHが通っている場所に移動し、Minecraft Bedrock Dedicated Serverを実行しているユーザが実行できるようにしてください。

```console
$ bdspackman --help
Usage of bdspackman:
  -dirname
        Show directory name column
  -export-dir string
        Directory to export world pack json files
  -export-prefix string
        Prefix added to exported filenames
  -force-delete
        Delete packs permanently instead of moving them to delpacks.
  -lang string
        Language code (default "ja_JP")
  -serverdir string
        Bedrock Dedicated Server directory (default ".")
  -systempack
        Show system packs
  -uuid
        Show UUID column
  -version
        show version
  -world string
        World name. If omitted, level-name in server.properties is used.
```

### 起動例

```console
# カレントディレクトリをサーバディレクトリとして起動します。
# --worldを省略した場合はserver.propertiesのlevel-nameを使用します。
bdspackman

# サーバディレクトリとワールド名を明示して起動します。
bdspackman -serverdir /opt/minecraft-bedrock-server -world "Bedrock level"

# UUID列とディレクトリ名列を表示します。
bdspackman -uuid -dirname

# 保存時に現在のworld_*_packs.jsonを直接更新せず、指定ディレクトリへprefix付きで出力します。
bdspackman -export-dir ./exports -export-prefix test-
```

### TUIの主なキー操作

| キー | 操作 |
| --- | --- |
| `↑` / `↓` | カーソル移動 |
| `Space` | 選択中のパックの有効 / 無効切り替え |
| `Ctrl+↑` / `Ctrl+↓` | 選択中のパックの適用順を入れ替え |
| `Enter` | 変更内容の保存確認へ進む |
| `M` | 選択中のパックをサーバ全体 / ワールド配下のパックディレクトリ間で移動 |
| `D` | 選択中のパックの削除確認へ進む |
| `R` | 選択中のパックディレクトリ名をリネーム |
| `Esc` / `Ctrl+C` | 終了 |

削除時は、通常はサーバディレクトリ配下の`delpacks/`へパックを移動します。`--force-delete`を指定した場合は、`delpacks/`へ退避せず完全に削除します。

## FAQ

これから書きます。
