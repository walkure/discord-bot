# discord-bot

Discord用botです。Goで書かれています。

現在、TwitterのURLに反応する機能が実装されています。Discordは公開アカウントかつNSFW判定されてないtweetは展開するんですが、
非公開アカウントやNSFW判定されたtweetは展開しないので、ここで展開します。

## setup

DiscordのUIはちょくちょく変わっているので、以下の説明通りにいかない場合があります。

1. Discord.devの[Application](https://discord.com/developers/applications)からアプリを追加。
   1. 「Bot」タブでTokenを生成。このトークンは一回しか表示されないのでメモっておく。これを環境変数`DISCORD_BOT_TOKEN`に設定。
      1. また、イベントを送って貰う必要があるので「Message Content Intent」をenableする。
   2. 「OAuth2」タブでインストールURLを生成。
      1. 「OAuth2 URL Generator」の「Scopes」で「bot」を選ぶと、いろんなパーミッションを選ぶ画面(Bot Permission)が表示される
      2. Bot Permissionで何かを選ぶことなく、「Integration Type」で「Guild Install」を選んでサーバへのインストールを指定する
      3. これで、Generated URLにアプリをサーバへインストールするためのリンクが生成される。
      4. botを入れたいサーバの管理権限を持った状態でURLへアクセスすると、インストールできる。
2. Twitter APIを叩くには`auth_token`が必要です。
   1. Twitterにログインした状態で、ブラウザのDeveloper Consoleを開き、Twitterのcookieから`auth_token`を取ってくる
   2. 複数のtokenを`TWITTER_AUTH_TOKENS_FROM_BROWSER`にカンマ区切りで書いておくと、ランダムに選択してAPIを呼び出します。
   3. NSFW画像をRT内でも展開するには当該アカウントでセンシティブ画像を表示する設定にする必要があります。

## Twitterの展開について

API叩く部分は[slack-unfurler](https://github.com/walkure/slack-unfurler)を流用しています。

## Execute locally with Docker

```batch
docker build -t discord-bot .
docker run --rm --env-file .env discord-bot
```

## License

MIT

## Author

walkure at 3pf.jp
