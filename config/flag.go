package config

import (
	"os"

	"github.com/mylxsw/glacier/starter/app"
)

func initCmdFlags(ins *app.App) {

	ins.AddStringFlag("listen", ":8080", "Web 服务监听地址")
	ins.AddStringFlag("api-listen", ":8081", "API 服务监听地址")
	ins.AddBoolFlag("enable-migrate", "是否启用迁移，启用后，当数据结构有更新时，会自动更新数据库")
	ins.AddDurationFlag("start-delay", 0, "服务启动延迟时间，用于在服务启动前做一些初始化工作，例如 Docker 环境下等待初始化数据库等")

	ins.AddStringFlag("base-url", "", "Web 服务的基础 URL，例如 https://web.aicode.cc")
	ins.AddBoolFlag("production", "是否为生产环境，生产环境下只有正式的支付渠道可用")
	ins.AddStringFlag("socks5-proxy", "", "socks5 proxy")
	ins.AddStringFlag("proxy-url", "", "HTTP 代理放置，支持 http、https、socks5，代理类型由 URL schema 决定，如果 scheme 为空，则默认为 http")
	ins.AddStringFlag("db-uri", "root:12345@tcp(127.0.0.1:3306)/aiserver?charset=utf8mb4&parseTime=True&loc=Local", "database url")
	ins.AddStringFlag("session-secret", "aidea-secret", "用户会话加密密钥")
	ins.AddBoolFlag("enable-recordchat", "是否记录聊天历史记录（目前只做记录，没有实际作用，只是为后期增加多端聊天记录同步做准备）")
	ins.AddBoolFlag("enable-cors", "是否启用跨域请求支持")
	ins.AddBoolFlag("enable-websocket", "是否启用 WebSocket 支持")
	ins.AddBoolFlag("debug-with-sql", "是否在日志中输出 SQL 语句")
	ins.AddBoolFlag("enable-api-keys", "是否启用 API Keys 功能")
	ins.AddBoolFlag("enable-model-rate-limit", "是否启用模型请求频率限制，当前限制只支持每分钟 5 次/用户")
	ins.AddStringFlag("universal-link-config", "", "universal link 配置文件路径，留空则使用默认的 universal link，配置文件格式参考 https://developer.apple.com/documentation/xcode/supporting-associated-domains")
	ins.AddBoolFlag("should-bind-phone", "是否需要绑定手机号码")

	ins.AddStringFlag("temp-dir", "/tmp", "临时文件存储目录")

	ins.AddStringFlag("redis-host", "127.0.0.1", "redis host")
	ins.AddIntFlag("redis-port", 6379, "redis port")
	ins.AddStringFlag("redis-password", "", "redis password")

	ins.AddIntFlag("queue-workers", 0, "任务队列工作线程（Goroutine）数量，设置为 0 则不启用任务队列")
	ins.AddBoolFlag("enable-scheduler", "是否启用定时任务")

	ins.AddBoolFlag("enable-custom-home-models", "是否启用自定义首页模型，启用后注意执行 2023101701-ddl.sql 数据迁移")

	ins.AddBoolFlag("enable-openai", "是否启用 OpenAI")
	ins.AddBoolFlag("openai-azure", "使用 Azure 的 OpenAI 服务")
	ins.AddStringFlag("openai-apiversion", "2023-05-15", "required when openai-azure is true")
	ins.AddBoolFlag("openai-autoproxy", "使用 Socks5 代理访问 OpenAI 服务")
	ins.AddStringFlag("openai-organization", "", "openai organization")
	ins.AddStringSliceFlag("openai-servers", []string{"https://api.openai.com/v1"}, "OpenAI 服务地址，配置多个时会自动在多个服务之间平衡负载，不要忘记在在 URL 后面添加 /v1")
	ins.AddStringSliceFlag("openai-keys", []string{}, "OpenAI Keys，如果指定多个，会在多个服务之间平衡负载")

	ins.AddBoolFlag("enable-openai-dalle", "是否启用 OpenAI DALL-E")
	ins.AddBoolFlag("dalle-using-openai-setting", "是否使用 OpenAI 的配置")
	ins.AddBoolFlag("openai-dalle-azure", "使用 Azure 的 OpenAI 服务")
	ins.AddStringFlag("openai-dalle-apiversion", "2023-05-15", "required when openai-dalle-azure is true")
	ins.AddBoolFlag("openai-dalle-autoproxy", "使用 Socks5 代理访问 OpenAI 服务")
	ins.AddStringFlag("openai-dalle-organization", "", "openai organization")
	ins.AddStringSliceFlag("openai-dalle-servers", []string{"https://api.openai.com/v1"}, "OpenAI 服务地址，配置多个时会自动在多个服务之间平衡负载，不要忘记在在 URL 后面添加 /v1")
	ins.AddStringSliceFlag("openai-dalle-keys", []string{}, "OpenAI Keys，如果指定多个，会在多个服务之间平衡负载")

	ins.AddBoolFlag("enable-fallback-openai", "是否启用备用 OpenAI 服务")
	ins.AddBoolFlag("fallback-openai-azure", "使用 Azure 的 OpenAI 服务")
	ins.AddStringFlag("fallback-openai-apiversion", "2023-05-15", "required when fallback-openai-azure is true")
	ins.AddBoolFlag("fallback-openai-autoproxy", "使用 Socks5 代理访问 OpenAI 服务")
	ins.AddStringFlag("fallback-openai-organization", "", "openai organization")
	ins.AddStringSliceFlag("fallback-openai-servers", []string{"https://api.openai.com/v1"}, "OpenAI 服务地址，配置多个时会自动在多个服务之间平衡负载，不要忘记在在 URL 后面添加 /v1")
	ins.AddStringSliceFlag("fallback-openai-keys", []string{}, "OpenAI Keys，如果指定多个，会在多个服务之间平衡负载")

	ins.AddBoolFlag("enable-anthropic", "是否启用 Anthropic")
	ins.AddBoolFlag("anthropic-autoproxy", "使用 socks5 代理访问 Anthropic 服务")
	ins.AddStringFlag("anthropic-server", "https://api.anthropic.com", "anthropic server")
	ins.AddStringFlag("anthropic-apikey", "", "anthropic server key")

	ins.AddBoolFlag("enable-googleai", "是否启用 GoogleAI")
	ins.AddBoolFlag("googleai-autoproxy", "使用 socks5 代理访问 GoogleAI 服务")
	ins.AddStringFlag("googleai-server", "https://generativelanguage.googleapis.com", "googleai server")
	ins.AddStringFlag("googleai-key", "", "googleai server key")

	ins.AddBoolFlag("enable-baiduwxai", "是否启用百度文心千帆大模型")
	ins.AddStringFlag("baiduwx-key", "", "百度文心大模型 Key")
	ins.AddStringFlag("baiduwx-secret", "", "百度文心大模型 Secret")

	ins.AddBoolFlag("enable-dashscopeai", "是否启用阿里灵积平台(通义千问)")
	ins.AddStringFlag("dashscope-key", "", "阿里灵积平台密钥")
	ins.AddStringSliceFlag("dashscope-keys", []string{}, "阿里灵积平台密钥，这里所有的 Keys 会和 dashscope-key 合并到一起，随机均摊请求负载")

	ins.AddBoolFlag("enable-xfyunai", "是否启用讯飞 星火 AI")
	ins.AddStringFlag("xfyun-appid", "", "讯飞星火 APP ID")
	ins.AddStringFlag("xfyun-apikey", "", "讯飞星火 API Key")
	ins.AddStringFlag("xfyun-apisecret", "", "讯飞星火 API Secret")

	ins.AddBoolFlag("enable-sensenovaai", "是否启用商汤日日新 AI")
	ins.AddStringFlag("sensenova-keyid", "", "商汤日日新 Key ID")
	ins.AddStringFlag("sensenova-keysecret", "", "商汤日日新 Key Secret")

	ins.AddBoolFlag("enable-baichuan", "是否启用百川大模型")
	ins.AddStringFlag("baichuan-apikey", "", "百川大模型 API Key")
	ins.AddStringFlag("baichuan-secret", "", "百川大模型 API Secret")

	ins.AddBoolFlag("enable-sky", "是否启用天工 AI")
	ins.AddStringFlag("sky-appkey", "", "天工 AI APP Key")
	ins.AddStringFlag("sky-appsecret", "", "天工 AI APP Secret")

	ins.AddBoolFlag("enable-gpt360", "是否启用 360 智脑大模型")
	ins.AddStringFlag("gpt360-apikey", "", "360 智脑大模型 API Key")

	ins.AddBoolFlag("enable-zhipuai", "是否启用智谱 AI")
	ins.AddStringFlag("zhipuai-key", "", "智谱 AI API Key")

	ins.AddBoolFlag("enable-moonshot", "是否启用月之暗面")
	ins.AddStringFlag("moonshot-apikey", "", "月之暗面 API Key")

	ins.AddBoolFlag("enable-oneapi", "是否启用 OneAPI")
	ins.AddStringFlag("oneapi-server", "", "one-server server")
	ins.AddStringFlag("oneapi-key", "", "one-server key")

	ins.AddStringSliceFlag("openrouter-support-models", []string{"01-ai/yi-34b-chat"}, "openrouter 支持的模型")
	ins.AddBoolFlag("enable-openrouter", "是否启用 OpenRouter")
	ins.AddBoolFlag("openrouter-autoproxy", "使用 socks5 代理访问 OpenRouter 服务")
	ins.AddStringFlag("openrouter-server", "https://openrouter.ai/api/v1", "openrouter server")
	ins.AddStringFlag("openrouter-key", "", "openrouter key")

	ins.AddBoolFlag("enable-stabilityai", "是否启用 StabilityAI 文生图、图生图服务")
	ins.AddBoolFlag("stabilityai-autoproxy", "使用 socks5 代理访问 StabilityAI 服务")
	ins.AddStringFlag("stabilityai-organization", "", "stabilityai organization")
	ins.AddStringSliceFlag("stabilityai-servers", []string{"https://api.stability.ai"}, "stabilityai servers")
	ins.AddFlags(app.StringEnvFlag("stabilityai-key", "", "stabilityai key", "STABILITYAI_KEY"))

	ins.AddBoolFlag("enable-leapai", "是否启用 LeapAI 文生图、图生图服务")
	ins.AddBoolFlag("leapai-autoproxy", "使用 socks5 代理访问 Leap 服务")
	ins.AddStringSliceFlag("leapai-servers", []string{"https://api.tryleap.ai"}, "leapai servers")
	ins.AddFlags(app.StringEnvFlag("leapai-key", "", "stabilityai key", "LEAPAI_API_KEY"))

	ins.AddBoolFlag("enable-deepai", "是否启用 DeepAI 超分辨率、上色服务")
	ins.AddBoolFlag("deepai-autoproxy", "deepai auto proxy")
	ins.AddStringSliceFlag("deepai-servers", []string{"https://api.deepai.org"}, "deepai servers")
	ins.AddFlags(app.StringEnvFlag("deepai-key", "", "deepai key", "DEEPAI_KEY"))

	ins.AddBoolFlag("enable-getimgai", "是否启用 getimg.ai 文生图、图生图服务")
	ins.AddBoolFlag("getimgai-autoproxy", "使用 socks5 代理访问 getimg.ai 服务")
	ins.AddStringFlag("getimgai-server", "https://api.getimg.ai", "getimgai server")
	ins.AddFlags(app.StringEnvFlag("getimgai-key", "", "getimgai key", "GETIMGAI_KEY"))

	ins.AddBoolFlag("enable-leptonai", "是否启用 lepton.ai 的模型服务")
	ins.AddBoolFlag("leptonai-autoproxy", "使用 socks5 代理访问 lepton.ai 服务")
	ins.AddStringSliceFlag("leptonai-qr-servers", []string{"https://aiqr.lepton.run"}, "lepton.ai QR servers")
	ins.AddStringSliceFlag("leptonai-keys", []string{os.Getenv("LEPTONAI_KEY")}, "lepton.ai keys")

	ins.AddBoolFlag("enable-fromstonai", "是否启用 6pen 的文生图、图生图服务")
	ins.AddStringFlag("fromston-server", "https://ston.6pen.art", "fromston server")
	ins.AddStringFlag("fromston-key", "", "fromston key")

	ins.AddBoolFlag("enable-translate", "是否启用翻译服务")
	ins.AddStringFlag("translate-server", "https://openapi.youdao.com/api", " 有道翻译 API 地址")
	ins.AddStringFlag("translate-appid", "", "有道翻译 APPID")
	ins.AddStringFlag("translate-appkey", "", "有道翻译 APPKEY")

	ins.AddStringFlag("storage-appkey", "", "七牛云 APP KEY")
	ins.AddStringFlag("storage-secret", "", "七牛云 APP SECRET")
	ins.AddStringFlag("storage-bucket", "aicode", "七牛云存储 Bucket 名称")
	ins.AddStringFlag("storage-callback", "https://YOUR_SERVER_HOST/v1/callback/storage/qiniu", "七牛云存储上传回调接口")
	ins.AddStringSliceFlag("storage-callbacks", []string{}, "七牛云存储上传回调接口，支持多个回调地址")
	ins.AddStringFlag("storage-callback-host", "", "七牛云存储上传回调 Host")
	ins.AddStringFlag("storage-domain", "", "七牛云存储资源访问域名（也可以用 CDN 域名），例如 https://cdn.example.com")
	ins.AddStringFlag("storage-region", "z0", "七牛云存储区域，可选值：z0, z1, z2, na0, as0, cn-east-2, ap-northeast-1")

	ins.AddStringFlag("apple-keyid", "", "apple sign in key id")
	ins.AddStringFlag("apple-teamid", "", "apple sign in team id")
	ins.AddStringFlag("apple-secret", "", "apple sign in secret")

	ins.AddBoolFlag("enable-mail", "是否启用邮件发送服务")
	ins.AddStringFlag("mail-from", "", "mail from")
	ins.AddStringFlag("mail-host", "", "mail host")
	ins.AddIntFlag("mail-port", 465, "mail port")
	ins.AddStringFlag("mail-username", "", "mail username")
	ins.AddStringFlag("mail-password", "", "mail password")
	ins.AddBoolFlag("mail-ssl", "mail ssl")

	ins.AddStringFlag("tencent-id", "", "tencent app id")
	ins.AddStringFlag("tencent-key", "", "tencent app key")
	ins.AddStringFlag("tencent-smssdkappid", "", "tencent sms sdk app id")
	ins.AddStringFlag("tencent-smstemplateid", "", "腾讯短信验证码模板 ID")
	ins.AddStringFlag("tencent-smssign", "AIdea", "腾讯短信签名")
	ins.AddBoolFlag("tencent-voice", "是否使用腾讯的语音转文本服务，不启用则使用 OpenAI 的 Whisper 模型")
	ins.AddBoolFlag("enable-tencentai", "是否启用腾讯混元大模型 AI 服务")

	ins.AddStringFlag("aliyun-key", "", "aliyun app key")
	ins.AddStringFlag("aliyun-secret", "", "aliyun app secret")
	ins.AddStringFlag("aliyun-smstemplateid", "", "阿里云短信验证码模板 ID")
	ins.AddStringFlag("aliyun-smssign", "AIdea", "阿里云短信签名")
	ins.AddBoolFlag("enable-contentdetect", "是否启用内容安全检测（使用阿里云的内容安全服务）")

	ins.AddBoolFlag("enable-applepay", "启用 Apple 应用内支付")

	ins.AddBoolFlag("enable-alipay", "启用支付宝支付支持，需要指定 alipay-xxx 的所有配置项")
	ins.AddStringFlag("alipay-appid", "", "支付宝 APP ID")
	ins.AddStringFlag("alipay-app-private-key", "path/to/alipay-app-private-key.txt", "支付宝 APP 私钥存储路径")
	ins.AddStringFlag("alipay-app-public-key", "path/to/appCertPublicKey_2021004100000000.crt", "支付宝 APP 公钥证书存储路径")
	ins.AddStringFlag("alipay-root-cert", "path/to/alipayRootCert.crt", "支付宝根证书路径")
	ins.AddStringFlag("alipay-public-key", "path/to/alipayCertPublicKey_RSA2.crt", "支付宝公钥证书路径")
	ins.AddStringFlag("alipay-notify-url", "https://ai-api.aicode.cc/v1/payment/callback/alipay-notify", "支付宝支付回调地址")
	ins.AddStringFlag("alipay-return-url", "https://ai-api.aicode.cc/public/payment/alipay-return", "支付宝支付 return url")
	ins.AddBoolFlag("alipay-sandbox", "是否使用支付宝沙箱环境")

	ins.AddStringSliceFlag("sms-channels", []string{}, "启用的短信通道，支持腾讯云和阿里云: tencent, aliyun，多个值时随机每次发送随机选择")

	ins.AddStringFlag("log-path", "", "日志文件存储目录，留空则写入到标准输出")
	ins.AddBoolFlag("log-colorful", "是否启用彩色日志")

	ins.AddBoolFlag("dingding-slack-mode", "启用后，钉钉通知将发送至 Slack，而不是钉钉群")
	ins.AddStringFlag("dingding-token", "", "钉钉群通知 Token，留空则不通知")
	ins.AddStringFlag("dingding-secret", "", "钉钉群通知 Secret")

	ins.AddBoolFlag("cnlocal-mode", "是否启用国产化模式，启用后，将使用 cnlocal-vendor/cnlocal-model 指定的模型替代数字人默认的 GPT 模型")
	ins.AddBoolFlag("cnlocal-onlyios", "国产化模式只对 IOS 系统有效，客户端版本 > 1.0.4")
	ins.AddStringFlag("cnlocal-vendor", "讯飞星火", "国产化模型服务商，目前支持讯飞星火、灵积、文心千帆、商汤日日新")
	ins.AddStringFlag("cnlocal-model", "generalv2", "国产化模型名称，讯飞星火支持 generalv2, 灵积支持 qwen-v1, 商汤日日新支持 nova-ptc-xl-v1，文心千帆支持 model_ernie_bot_turbo、model_badiu_llama2_70b、model_baidu_llama2_7b_cn、model_baidu_chatglm2_6b_32k、model_baidu_aquila_chat7b、model_baidu_bloomz_7b")

	ins.AddStringFlag("default-img2img-model", "lb-realistic-versionv4.0", "默认的图生图模型，值取自数据表 image_model.model_id")
	ins.AddStringFlag("default-txt2img-model", "sb-stable-diffusion-xl-1024-v1-0", "默认的文生图模型，值取自数据表 image_model.model_id")

	ins.AddStringFlag("img2img-recognition-provider", "", "图生图图像识别处理模型，用于识别图像内容，生成图生图的提示语，目前支持 xfyun，留空则表示不启用")

	ins.AddStringFlag("price-table-file", "", "价格表文件路径，留空则使用默认价格表")

	ins.AddStringFlag("font-path", "", "字体文件路径")
	ins.AddStringFlag("service-status-page", "", "服务状态页面，留空则不启用服务状态页面")

	ins.AddBoolFlag("free-chat-enabled", "是否启用免费聊天功能，启用后，未登录可以免费使用部分模型")
	ins.AddIntFlag("free-chat-daily-limit", 5, "每日免费次数，基于客户端 IP 限制")
	ins.AddIntFlag("free-chat-daily-global-limit", 1000, "每日全局免费次数，总的限制次数，不管是哪个 IP")

	ins.AddStringFlag("wechat-appid", "", "微信开放平台 APP ID")
	ins.AddStringFlag("wechat-secret", "", "微信开放平台 APP Secret")
	ins.AddBoolFlag("wechatpay-enabled", "是否启用微信支付")
	ins.AddStringFlag("wechatpay-mchid", "", "微信支付商户号")
	ins.AddStringFlag("wechatpay-cert-serial-number", "", "微信支付 API 证书序列号")
	ins.AddStringFlag("wechatpay-cert-private-key-path", "", "微信支付 API 证书私钥路径")
	ins.AddStringFlag("wechatpay-api-v3-key", "", "微信支付 APIv3 密钥")
	ins.AddStringFlag("wechatpay-notify-url", "", "微信支付回调地址，格式 https://YOUR_SERVER_HOST/v1/payment/callback/wechat-pay/notify")

	ins.AddBoolFlag("enable-stripe", "是否启用 stripe 支付")
	ins.AddStringFlag("stripe-publishable-key", "", "stripe publishable key")
	ins.AddStringFlag("stripe-secret-key", "", "stripe secret key")
	ins.AddStringFlag("stripe-webhook-secret", "", "stripe webhook secret")

	ins.AddStringSliceFlag("default-home-models", []string{"gpt-3.5-turbo", "gpt-4"}, "默认的首页模型，值取自数据表 models.model_id")
	ins.AddStringSliceFlag("default-home-models-ios", []string{"chat-3.5", "chat-4"}, "默认的首页模型（IOS），值取自数据表 models.model_id")
	ins.AddStringFlag("default-role-model", "gpt-4o", "默认的角色模型，值取自数据表 models.model_id")

	ins.AddStringFlag("text-to-voice-engine", "openai", "文本转语音引擎，支持 minimax/openai/azure")
	ins.AddStringFlag("text-to-voice-azure-region", "eastus", "Azure 语音服务区域")
	ins.AddStringFlag("text-to-voice-azure-key", "", "Azure 语音服务 Key")

	ins.AddBoolFlag("enable-voice-to-text", "是否启用语音转文本功能")
	ins.AddBoolFlag("enable-text-to-voice", "是否启用文本转语音功能")

	ins.AddStringFlag("minimax-api-key", "", "Minimax API Key")
	ins.AddStringFlag("minimax-group-id", "", "Minimax Group ID")

	ins.AddBoolFlag("enable-summarizer", "是否启用聊天记录总结功能")
	ins.AddStringFlag("summarizer-model", "gpt-4o-mini", "总结模型名称")

	ins.AddStringFlag("flux-api-server", "https://api.bfl.ml", "flux api server")
	ins.AddStringFlag("flux-api-key", "", "flux api key")

	ins.AddStringFlag("bigmodel-search-api-key", "", "BigModel 搜索 API Key")
	ins.AddStringFlag("bochaai-search-api-key", "", "Bochaai 搜索 API Key")
	ins.AddStringFlag("search-engine", "bigmodel", "搜索引擎，支持 bigmodel/bocha-web/bocha-ai")
	ins.AddStringSliceFlag("available-search-engines", []string{}, "可用的搜索引擎")
	ins.AddStringFlag("search-assistant-model", "gpt-4o-mini", "搜索助手模型")
	ins.AddStringFlag("search-assistant-api-base", "https://api.openai.com/v1", "搜索助手 API Base")
	ins.AddStringFlag("search-assistant-api-key", "", "搜索助手 API Key")
}
