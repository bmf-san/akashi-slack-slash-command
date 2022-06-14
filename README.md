# akashi-slack-slash-command
SlackのスラッシュコマンドでAkashiのAPIを利用するアプリケーション。

SlackのスラッシュコマンドでCloud Functionsの関数を呼び出しで処理を実行します。
関数はSheets APIとAkashi APIをコールしてレスポンスを返します。

# 手順
## 前提
- gcloudコマンドが利用できること
  - アプリケーションをデプロイするために、gcloudコマンドが利用できる必要があります。
- Cloud Build APIが有効化されていること
  - Cloud Functionsに関数をデプロイするために必要です。
- Google Sheet APIが有効化されていること
  - Google Cloud Platform上でGoogle Sheet APIを有効化しておく必要があります。
- AKASHI APIが有効化されていること

## 準備
### Cloud Functionに関数を作成
Cloud Functionsにて関数を作成しておく。トリガータイプは`HTTP`、認証は`未認証の呼び出しを許可`を選択。`HTTPSが必須`はチェック。
関数デプロイ後、Cloud Functionsの関数の詳細>トリガーに記載されているトリガーURLをコピーしておく。

### Slack App
#### Slack Appを作成
##### 
[Create an app](https://api.slack.com/apps?new_app=1)にて、`From scratch`を
押下。

![スクリーンショット 2022-06-15 21 24 12](https://user-images.githubusercontent.com/13291041/173827906-4e45e346-1b8e-4f42-bd93-e0a61efa5684.png)


`App Name`を入力、`Pick a workspace to develop your app in:`にてワークスペースを選択して、`Create App`を押下。

![スクリーンショット 2022-06-15 21 24 51](https://user-images.githubusercontent.com/13291041/173827918-eac89b94-3238-4fda-974c-1b5ffdc6c9fe.png)

#### Slash Commandを設定
設定画面（ex. `https://api.slack.com/apps/*******`）にて、Slash Commandsを選択。

![スクリーンショット 2022-06-15 21 36 37](https://user-images.githubusercontent.com/13291041/173828344-539e181c-f033-4545-b0fd-b2ddcd3c90a3.png)

`Create New Command`を押下し、`Command`、`Short Description`、`Usage Hint`、`Escape channels, users, and links sent to your app`は任意で設定。

`Request URL`はCloud Functionsの関数の詳細>トリガーに記載されているトリガーURLを入力。`https://REGION名-PROJECTのID.cloudfunctions.net/関数名`の形式。

![スクリーンショット 2022-06-15 21 30 52](https://user-images.githubusercontent.com/13291041/173827928-a0033277-0e07-4a26-bff8-a0cdbd19c30c.png)

`Save`を押下。

#### Slack Appをインストール
設定画面（ex. `https://api.slack.com/apps/*******`）にて、`Install App`を押下。

![スクリーンショット 2022-06-15 21 37 09](https://user-images.githubusercontent.com/13291041/173828510-6ee1e474-c5c8-4cf8-b683-992f3d2e5790.png)

`Install to workspace`を押下して、任意のワークスペースにAppをインストール。

#### Signing Secretを取得
設定画面（ex. `https://api.slack.com/apps/*******`）にて、`Basic Infomation`を押下。App Credentialsという項目に、`Signing Secret`があるので、値をコピーしておく。

### サービスアカウントの発行
GCPでサービスアカウントを発行。

任意のサービスアカウント名を入力したら`作成して続行`を押下。ロールの調整はせずに、`完了`を押下。
![スクリーンショット 2022-06-15 21 44 21](https://user-images.githubusercontent.com/13291041/173829849-2e886e02-8dc8-4a7d-952b-ac9624acf32a.png)

サービスアカウント一覧から作成したサービスアカウントを選択、新しい鍵をJSON形式で作成。

![スクリーンショット 2022-06-15 21 47 14](https://user-images.githubusercontent.com/13291041/173830444-3161edde-daf0-44c7-ba89-a1785efc1edc.png)

作成された秘密鍵（jsonファイル）をbase64でエンコードした値をコピーしておく。
`base64 service_account.json`

### Spread Sheetsを作成
Spread Sheetを作成。

以下シートに貼り付けるためのコピペ用データ。

```
employee_id	slack_user_name	akashi_company_id	akashi_api_token
dummy_id	dummy_name	dummy_com_id	dummy_token
dummy_id	dummy_name	dummy_com_id	dummy_token
dummy_id	dummy_name	dummy_com_id	dummy_token
dummy_id	dummy_name	dummy_com_id	dummy_token
dummy_id	dummy_name	dummy_com_id	dummy_token
```

`employee_id`は社員番号。
`slack_user_name`はslackのユーザ名。
`akashi_company_id`はAkashi APIのコールで利用する企業ID。
`akashi_api_token`はAkashiでユーザーごとに発行するAPIトークン。

作成したSpread Sheetの共有設定を開き、サービスアカウントのメールを追加する。
サービスアカウントのメールは、サービスアカウントの詳細で確認できる。

### .env.yamlの作成
`.env.yaml`を作成。
```sh
cp .env.yaml.example .env.yaml
```

`SLACK_SIGINING_SECRET`はSlack Appのsigning secret。
`SERVICE_ACCOUNT`はbase64エンコードしたサービスアカウントの値。
`SPREAD_SHEET_ID`はSpread SheetのID。IDはSpread Sheetのリンクから確認できる。`https://docs.google.com/spreadsheets/d/<SPREAD SHEEET ID>/edit#gid=0`

## 関数を再デプロイ
```sh
export FUNC=関数名
make deploy
```

## 使い方

|   コマンド   |                                                              内容                                                              |
| ------------ | ------------------------------------------------------------------------------------------------------------------------------ |
| /akashi      | 出勤または退勤します。出勤が打刻されていない場合は出勤。出勤が打刻されているまたは退勤が打刻されている場合は退勤を打刻します。 |
| /akashi 出勤 | 出勤を打刻します。                                                                                                             |
| /akashi 退勤 | 退勤を打刻します。                                                                                                             |
※Slash CommandのCommandを`akashi`とした場合を想定。

# コスト
[Cloud Functions - Cloudfunctionsの料金](https://cloud.google.com/functions/pricing?hl=ja)

# クォータ
[Sheets API - Usage limits](https://developers.google.com/sheets/api/limits)

# References
- [Google Cloud - Go on Google App Engine](https://cloud.google.com/appengine/docs/go)
- [Google Cloud - App Engine pricing](https://cloud.google.com/appengine/pricing)
- [www.serversus.work - Google App Engine(GAE)を無料枠で収めるための勘所](https://www.serversus.work/topics/p1uaj4jrv8b5x70hwe6p/)
- [github.com - slack-go/slack](https://github.com/slack-go/slack)