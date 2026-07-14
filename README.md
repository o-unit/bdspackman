# bdspackman

Minecraft Bedrock Editionのサーバでビヘイビアパック / リソースパックの管理をするための簡易ツール。

### どんなツール？
TUIでminecraft bedrock版の各種パックを管理するためのツール(になるように頑張り中)。  
- .../[behavior / resource]_packs と .../worlds/[WorldDir]/[behavior / resource]_packs の4ディレクトリにあるパックの一覧表示  
- world_[behavior / resource]_packs.jsonの出力  
- 一覧にあるパックの有効/無効切り替え  
- パックの適用順の入れ替え  
ができます。  

近い将来の目標はパックの追加 / 移動 / 削除の実装。  
少し遠い目標はパックの更新。

### 動作確認済み環境
- Ubuntu 24.04 (amd64)

### 現在の実装状況
✅パックの一覧表示  
✅world_behavior_packs.jsonの更新  
✅world_resource_packs.jsonの更新  
✅パックの有効化 / 無効化  
❌パックの移動  
❌パックの削除  
❌パックの追加  
❌パックの更新  

### 使い方
リリースページからzipファイルをダウンロード。  
解凍したら実行ファイルができるので、PATHが通っている場所に移動させて  
minecraft bedrockサーバを実行しているユーザが実行できるようにしてあげてください。  

```$ bdspackman --help
Usage of bdspackman:
  -dirname
        Show directory name column
  -export-dir string
        Directory to export world pack json files
  -export-prefix string
        Prefix added to exported filenames
  -lang string
        Language code (default "ja_JP")
  -serverdir string
        Bedrock Dedicated Server directory (default ".")
  -systempack
        Show system packs
  -uuid
        Show UUID column
  -world string
        World name (required)```

### FAQ
これから書きます  