# RemotePlayOther
Steam Remote Play Togetherを非対応ゲームでできるようにする

## 事前設定
1. main.goをコンパイルする(windows:`go build -o ./main.exe -ldflags -H=windowsgui`)
2. `example_games.json`を`games.json`に書き換え 例を参考に中身に"ゲーム名"と"実行ファイルまでのフルパス"を追加していく

## How To Use
1. SteamでRemote Play Togetherに対応したゲームを開く
2. 歯車=>管理=>ローカルファイルを閲覧 を開く
3. ゲームの実行ファイルをほかの名前に変える
4. コンパイルしたファイルをコピーし元のゲームの実行ファイルの名前にする
5. `games.json`をコピーする