#TODO:今後自分で記載していく(Golandではプラグインが豊富な為いらないかも)
#(静的解析ツールは今後導入)lint:
#lint:
		#golangci-lint run ./...

# goのtestコードの実行例です。テスト前にlintを指定する、というMakeflieの記述例です
#test: lint
		#go test ./...

#debug:
	#open http://localhost:8080/
	#go run ./cmd/...

# DBを本番と別の環境を用意したいとき
#deploy-test:
	#gcloud config set project {テスト用のGCP project-id}
	#gcloud app deploy app_dev.yaml

# 本番と同様のDBで、アプリケーションだけ変えたい場合
#deploy-dev:
	#gcloud config set project {本番/開発用のGCP projectId}
	#gcloud app deploy app_dev.yaml # service名を変えたymlを用意しておく

#deploy-prod:
	#gcloud config set project {本番のGCP projectId}
	#gcloud app deploy app_production.yaml


