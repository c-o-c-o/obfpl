# obfpl  
特定のファイルに特定のアプリケーションを起動するアプリケーションです。  
音声合成ソフトの出力に、様々なエフェクトや処理を自動的に行うために作成しました。  
セリフ先頭に”や＃をつけると囁き声にしたり通信機越しっぽい声にしたりします。  

# 使い方  
`obfpl.exe` を起動してください。  
音声合成ソフト出力先を `src` フォルダー、最終出力先は `dst` フォルダーに設定しています。   
これを変更したい場合はコマンドライン上で起動するか、付属している `obfpl.bat` を編集して、`obfpl.bat` から起動してください。  

```
obfpl.exe [options]
options
  -s 音声合成ソフト出力先[src]
  -d 最終出力先[dst]
  -p プロファイルのパス[profile.yml]
```

```
example
  obfpl.exe -s 音声合成ソフト出力先 -d 最終出力先
```

# 初期設定  
初期設定として  

+ ファイル名の変更、テキストファイルの字幕化
+ 字幕ファイルの改行
+ セリフ先頭に `”` が含まれている場合囁き声っぽくする
+ セリフ先頭に `（` が含まれている場合心の声っぽくする
+ セリフ先頭に `＃` が含まれている場合通信機越しの声っぽくする
+ セリフ先頭に `｜` が含まれている場合扉越しの声っぽくする
+ 出力後ファイル名に出力順をナンバリング

これらが設定されています。  
変更する場合は `profile.yml` を参照してください。  

# 再配布について  
このアプリケーションの配布パッケージには製作者が作成していないソフトウェアが付属しています。  
再配布する場合はそれらのライセンスを確認してください。  

# Licence
This software is released under the MIT License, see LICENSE.  

下記の物は配布パッケージに含まれているソフトウェアです。  

sox-14.4.2  
Sound eXchange ( SoX )  
http://sox.sourceforge.net/  

ToWhisperNet1.2  
toWhisper  
https://github.com/ksasao/toWhisper

auto-NewLine  
https://github.com/c-o-c-o/auto-NewLine

Renamer  
https://github.com/c-o-c-o/Renamer

Numbering  
https://github.com/c-o-c-o/Numbering