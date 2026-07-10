# 涓汉灏忕珯鍚庣寮€鍙戣鍒?
## 1. 鐩爣

杩欎唤瑙勫垝鍩轰簬浣犲綋鍓嶅伐浣滃尯閲岀殑 `tiktok_demo` 缁撴瀯鏉ヨ璁★紝鐩爣涓嶆槸閲嶆柊鍙戞槑涓€濂楁灦鏋勶紝鑰屾槸鎶婂畠鏀舵暃鎴愭洿閫傚悎涓汉灏忕珯鐨勫崟浣撳悗绔€?
鏍稿績鍘熷垯锛?
1. 鍏堝仛鍗曚綋锛屼笉鎬ョ潃鎷嗗井鏈嶅姟
2. 鐩綍杈圭晫瑕佹竻鏅帮紝鍚庣画鏂逛究杩佺Щ鍒?`Kitex`
3. 绗竴鐗堝厛璺戦€氬唴瀹归棴鐜細鐧诲綍銆佸啓鏂囩珷銆佸彂鏂囩珷銆佺湅鏂囩珷
4. 绗簩鐗堝啀鑰冭檻缂撳瓨銆佸璞″瓨鍌ㄣ€丄I 鑳藉姏

## 2. 浠?tiktok_demo 鍙互鐩存帴鍊熼壌浠€涔?
浣犺繖涓?`tiktok_demo` 宸茬粡缁欎簡涓€涓緢濂界殑 Hertz 椤圭洰楠ㄦ灦锛?
- `main.go` 璐熻矗鍒濆鍖栧拰鍚姩鏈嶅姟
- `biz/dal` 璐熻矗鏁版嵁灞傚垵濮嬪寲
- `biz/router` 璐熻矗璺敱娉ㄥ唽
- `biz/handler` 璐熻矗璇锋眰缁戝畾鍜屽搷搴?- `biz/service` 璐熻矗涓氬姟閫昏緫
- `biz/mw` 璐熻矗涓棿浠舵垨鍩虹缁勪欢
- `pkg` 璐熻矗閰嶇疆銆侀敊璇爜銆佸伐鍏峰嚱鏁般€佸父閲?
浠庝綘鍒氭墠鐪嬬殑浠ｇ爜閾捐矾閲岋紝鍙互鎬荤粨鎴愯繖鏉′富绾匡細

```text
main -> Init -> dal.Init -> router -> handler -> service -> db
```

杩欐潯涓荤嚎闈炲父閫傚悎鐩存帴澶嶇敤鍒颁釜浜哄皬绔欓噷銆?
## 3. 鍚庣鎬讳綋瀹氫綅

杩欎釜鍚庣鐨勮亴璐ｆ湁涓夌被锛?
1. 椤甸潰鏀拺
   - 棣栭〉
   - Blog 鍒楄〃椤?   - Blog 璇︽儏椤?   - About 椤?
2. 鍐呭绠＄悊
   - 鐢ㄦ埛鐧诲綍
   - 鏂囩珷澧炲垹鏀规煡
   - 鑽夌/鍙戝竷
   - 鏍囩鍜屽垎绫荤鐞?   - 鍥剧墖涓婁紶

3. 宸ョ▼鑳藉姏
   - 閰嶇疆鍔犺浇
   - 閴存潈
   - 缁熶竴鍝嶅簲
   - 閿欒鐮?   - 鏃ュ織
   - 鍒嗛〉鍜屾煡璇㈠皝瑁?
## 4. 鎺ㄨ崘椤圭洰缁撴瀯

寤鸿浣犳柊椤圭洰鏈€缁堟敹鏁涙垚杩欐牱锛?
```text
personal_site/
鈹溾攢鈹€ biz/
鈹?  鈹溾攢鈹€ handler/
鈹?  鈹?  鈹溾攢鈹€ home/
鈹?  鈹?  鈹溾攢鈹€ blog/
鈹?  鈹?  鈹溾攢鈹€ about/
鈹?  鈹?  鈹溾攢鈹€ admin/
鈹?  鈹?  鈹溾攢鈹€ auth/
鈹?  鈹?  鈹斺攢鈹€ upload/
鈹?  鈹溾攢鈹€ service/
鈹?  鈹?  鈹溾攢鈹€ home/
鈹?  鈹?  鈹溾攢鈹€ article/
鈹?  鈹?  鈹溾攢鈹€ tag/
鈹?  鈹?  鈹溾攢鈹€ category/
鈹?  鈹?  鈹溾攢鈹€ auth/
鈹?  鈹?  鈹斺攢鈹€ upload/
鈹?  鈹溾攢鈹€ dal/
鈹?  鈹?  鈹溾攢鈹€ db/
鈹?  鈹?  鈹斺攢鈹€ init.go
鈹?  鈹溾攢鈹€ model/
鈹?  鈹?  鈹溾攢鈹€ entity/
鈹?  鈹?  鈹溾攢鈹€ dto/
鈹?  鈹?  鈹斺攢鈹€ view/
鈹?  鈹溾攢鈹€ router/
鈹?  鈹?  鈹溾攢鈹€ site/
鈹?  鈹?  鈹溾攢鈹€ admin/
鈹?  鈹?  鈹斺攢鈹€ register.go
鈹?  鈹斺攢鈹€ mw/
鈹?      鈹溾攢鈹€ session/
鈹?      鈹溾攢鈹€ logger/
鈹?      鈹斺攢鈹€ recover/
鈹溾攢鈹€ pkg/
鈹?  鈹溾攢鈹€ configs/
鈹?  鈹溾攢鈹€ constants/
鈹?  鈹溾攢鈹€ errno/
鈹?  鈹溾攢鈹€ utils/
鈹?  鈹斺攢鈹€ response/
鈹溾攢鈹€ templates/
鈹溾攢鈹€ static/
鈹溾攢鈹€ migrations/
鈹溾攢鈹€ main.go
鈹斺攢鈹€ go.mod
```

杩欎釜缁撴瀯鍜?`tiktok_demo` 鐨勫樊鍒笉澶э紝鎵€浠ヤ綘瀛︿範杩佺Щ鎴愭湰浼氬緢浣庛€?
## 5. 妯″潡鍒掑垎

### 5.1 棣栭〉妯″潡 home

鑱岃矗锛?
- 棣栭〉娓叉煋
- 棣栭〉鍩虹淇℃伅鑱氬悎
- 杩斿洖绔欑偣浠嬬粛銆佸叆鍙ｄ俊鎭€佸彲閫夋渶杩戞枃绔?
瀵瑰簲鐩綍锛?
- `biz/handler/home`
- `biz/service/home`

### 5.2 鍐呭妯″潡 article

鑱岃矗锛?
- 鏂囩珷鍒楄〃
- 鏂囩珷璇︽儏
- 鍚庡彴鏂囩珷鍒涘缓
- 鍚庡彴鏂囩珷缂栬緫
- 鑽夌鍙戝竷鐘舵€佺鐞?- slug 鐢熸垚涓庡敮涓€鎬ф牎楠?
杩欐槸浣犵殑鏍稿績妯″潡锛岀涓€闃舵瑕佷紭鍏堝畬鎴愩€?
### 5.3 鏍囩妯″潡 tag

鑱岃矗锛?
- 鏍囩鍒涘缓
- 鏍囩鍒楄〃
- 鏂囩珷缁戝畾鏍囩
- 鏍囩绛涢€夋枃绔?
绗竴鐗堝彲浠ュ厛鍋氱畝鍗曞叧绯昏〃銆?
### 5.4 鍒嗙被妯″潡 category

鑱岃矗锛?
- 鍒嗙被鍒涘缓
- 鍒嗙被鍒楄〃
- 鏂囩珷鍏宠仈鍒嗙被

濡傛灉浣犳兂绠€鍖栵紝鍒嗙被鍙互鏅氫簬鏍囩鍋氥€?
### 5.5 璁よ瘉妯″潡 auth

鑱岃矗锛?
- 鐢ㄦ埛鐧诲綍
- Session 鐧诲綍鎬佸啓鍏ヤ笌鏍￠獙
- 绠＄悊鍚庡彴鏉冮檺淇濇姢

褰撳墠椤圭洰宸茬粡浠?JWT 璺嚎鍒囧埌 Session 璺嚎锛屽洜姝ょ涓€鐗堜互锛歚users` 琛ㄧ櫥褰曘€丼ession 缁存寔鐧诲綍鎬併€丷BAC 鎺у埗鍚庡彴鎺ュ彛鏉冮檺涓轰富锛屼笉鍐嶇户缁墿灞?JWT 涓婚摼銆?
### 5.6 涓婁紶妯″潡 upload

鑱岃矗锛?
- 鍥剧墖涓婁紶
- 杩斿洖鍙闂?URL
- 璁板綍鏂囦欢鍏冩暟鎹?
绗竴鐗堝缓璁湰鍦板瓨鍌紝鍚庣画鍐嶅垏 `MinIO`銆?
### 5.7 About 妯″潡 about

鑱岃矗锛?
- About 椤甸潰鍐呭杈撳嚭
- 鍙互鍏堥潤鎬侀厤缃?- 鍚庣画鍙浆鎴愬悗鍙板彲缂栬緫閰嶇疆椤?
## 6. 鏁版嵁搴撹璁″缓璁?
绗竴鐗堝缓璁繖浜涜〃锛?
### 6.1 users

```text
id
username
password_hash
nickname
created_at
updated_at
```

璇存槑锛?
- 绗竴鐗堝彧闇€瑕佷竴涓鐞嗗憳璐﹀彿
- 涓嶅仛澶嶆潅鐢ㄦ埛绯荤粺

### 6.2 articles

```text
id
title
slug
summary
content_md
content_html
status
cover_image
category_id
created_at
updated_at
published_at
```

璇存槑锛?
- `status` 寤鸿锛歚draft` / `published`
- `content_md` 淇濆瓨鍘熷 Markdown
- `content_html` 鍙€夌紦瀛樻覆鏌撶粨鏋?
### 6.3 tags

```text
id
name
slug
created_at
updated_at
```

### 6.4 article_tags

```text
article_id
tag_id
```

### 6.5 categories

```text
id
name
slug
created_at
updated_at
```

### 6.6 uploads

```text
id
file_name
file_path
file_url
mime_type
size
created_at
```

濡傛灉浣犳兂鍐嶇簿绠€锛岀涓€鐗堝彲浠ュ厛涓嶅仛 `categories`锛屽彧鍋?`tags`銆?
## 7. 璇锋眰閾捐矾鎬庝箞璁捐

鍙傝€?`tiktok_demo` 閲岀殑 handler 鍐欐硶锛屼綘鐨勮姹傚鐞嗗缓璁浐瀹氭垚涓嬮潰杩欎釜妯″紡锛?
### 7.1 handler 灞傝亴璐?
- 缁戝畾鍙傛暟
- 鍋氬熀纭€鏍￠獙
- 璋?service
- 缁熶竴杩斿洖 JSON 鎴栭〉闈㈡覆鏌?
涓嶈鎶婁笟鍔￠€昏緫濉炶繘 handler銆?
### 7.2 service 灞傝亴璐?
- 缂栨帓涓氬姟閫昏緫
- 璋冪敤 db 鏂规硶
- 澶勭悊 slug銆佸彂甯冩椂闂淬€佺姸鎬佸垏鎹€佹爣绛惧叧鑱?- 澶勭悊 Markdown 杞?HTML

### 7.3 db 灞傝亴璐?
- 鍙仛鏁版嵁搴撹闂?- 涓嶅啓椤甸潰鍜屼笟鍔￠€昏緫
- 涓€绫诲疄浣撲竴涓枃浠舵垨涓€缁勬枃浠?
瀵瑰簲浣犲彲浠ョ洿鎺ユā浠?`tiktok_demo` 鐨勬柟寮忥細

```text
handler: 澶勭悊璇锋眰
service: 澶勭悊涓氬姟
db: 鏌ヨ鍜屽啓搴?```

## 8. 璺敱瑙勫垝

### 8.1 绔欑偣椤甸潰璺敱

```text
GET  /
GET  /blog
GET  /blog/:slug
GET  /about
```

### 8.2 鍏紑 API

```text
GET  /api/articles
GET  /api/articles/:slug
GET  /api/tags
GET  /api/categories
```

### 8.3 绠＄悊鍚庡彴 API

```text
POST   /api/admin/login
GET    /api/admin/articles
POST   /api/admin/articles
PUT    /api/admin/articles/:id
DELETE /api/admin/articles/:id
GET    /api/admin/tags
POST   /api/admin/tags
GET    /api/admin/categories
POST   /api/admin/categories
POST   /api/admin/upload
```

### 8.4 璺敱鍒嗙粍寤鸿

```text
/site
/api
/api/admin
```

濡傛灉浣犺瀹屽叏璐磋繎 `tiktok_demo` 鐨勯鏍硷紝鍙互姣忎釜妯″潡鍗曠嫭寤?`router/blog`銆乣router/about`銆乣router/admin`銆?
## 9. 鍒濆鍖栭『搴?
鍙傝€?`tiktok_demo/main.go`锛屽缓璁綘鐨勫惎鍔ㄦ祦绋嬪啓鎴愶細

```text
InitConfig()
InitDB()
InitSession()
InitUploadStore()
InitTemplateRenderer()
RegisterRouter()
RunServer()
```

鏇存帴杩戜唬鐮佸眰闈㈢殑缁撴瀯锛?
```text
main.go
  -> Init()
      -> dal.Init()
      -> session.Init()
      -> upload.Init()
  -> register(h)
  -> h.Spin()
```

## 10. 涓棿浠惰鍒?
### 10.1 Session 璁よ瘉涓棿浠?
鐢ㄩ€旓細

- 淇濇姢鍚庡彴鎺ュ彛
- 鏍￠獙鐧诲綍鐘舵€?
### 10.2 鏃ュ織涓棿浠?
鐢ㄩ€旓細

- 璁板綍璇锋眰璺緞銆佽€楁椂銆佺姸鎬佺爜

### 10.3 鎭㈠涓棿浠?
鐢ㄩ€旓細

- 閬垮厤 panic 鐩存帴鎶婃湇鍔℃墦鎸?
### 10.4 鍙€夛細闄愭祦涓棿浠?
绗竴鐗堝彲浠ヤ笉鍋氥€?
## 11. 閰嶇疆瑙勫垝

寤鸿閰嶇疆鏂囦欢鑷冲皯鍖呭惈锛?
```text
server:
  port:

database:
  dsn:

jwt:
  secret:
  expire:

upload:
  dir:
  base_url:

site:
  name:
  github_url:
```

鏀惧湪锛?
```text
pkg/configs/
```

## 12. 閿欒鐮佷笌缁熶竴鍝嶅簲

杩欎釜閮ㄥ垎鍙互鐩存帴瀛?`tiktok_demo/pkg/errno` 鍜?`pkg/utils/resp` 鐨勬€濊矾銆?
寤鸿浣犱繚鐣欙細

- `pkg/errno`
- `pkg/response` 鎴?`pkg/utils/resp`

缁熶竴鍝嶅簲绀轰緥锛?
```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

椤甸潰娓叉煋鎺ュ彛鍒欏崟鐙繑鍥炴ā鏉裤€?
## 13. 绗竴闃舵寮€鍙戦『搴?
### 闃舵 1锛氳窇閫氶」鐩鏋?
鐩爣锛氶」鐩兘鍚姩銆佽兘璁块棶棣栭〉銆佽兘杩炴暟鎹簱

浠诲姟锛?
1. 鍒濆鍖?Hertz 椤圭洰
2. 鎼洰褰曠粨鏋?3. 鎺ュ叆 MySQL 鍜?GORM
4. 閰嶇疆鍔犺浇
5. 鍩虹璺敱娉ㄥ唽
6. 棣栭〉鍜?About 闈欐€侀〉杩斿洖

### 闃舵 2锛氬畬鎴愭枃绔犱富閾捐矾

鐩爣锛氳兘鍙戞枃绔犮€佽兘鐪嬫枃绔?
浠诲姟锛?
1. 寤?`articles` 琛?2. 鍋氭枃绔犲垪琛ㄦ帴鍙?3. 鍋氭枃绔犺鎯呮帴鍙?4. 鍋氬悗鍙版柊澧炴枃绔犳帴鍙?5. 鍋氬悗鍙扮紪杈戞枃绔犳帴鍙?6. 鍋氳崏绋?鍙戝竷鐘舵€?
杩欐槸鏈€鍏抽敭鐨勪竴闃舵銆?
### 闃舵 3锛氬畬鎴愬悗鍙扮櫥褰?
鐩爣锛氬悗鍙版帴鍙ｅ彈淇濇姢

浠诲姟锛?
1. 寤?`users` 琛?2. 鍒濆鍖栫鐞嗗憳璐﹀彿
3. 鐧诲綍鎺ュ彛
4. Session 涓棿浠?5. 鍚庡彴鎺ュ彛閴存潈

### 闃舵 4锛氳ˉ鍏ㄥ唴瀹圭鐞嗚兘鍔?
鐩爣锛氬悗鍙版洿濂界敤

浠诲姟锛?
1. 鏍囩绠＄悊
2. 鍒嗙被绠＄悊
3. 鍥剧墖涓婁紶
4. 鏂囩珷绛涢€夊拰鍒嗛〉
5. 缁熶竴閿欒鐮佷笌鏃ュ織

### 闃舵 5锛氬伐绋嬪寲瀹屽杽

鐩爣锛氶」鐩洿绋冲畾銆佸悗缁洿濂芥紨杩?
浠诲姟锛?
1. Docker 鍖?2. 鐜閰嶇疆鎷嗗垎
3. 娴嬭瘯琛ラ綈
4. 涓婁紶鑳藉姏鏇挎崲涓?`MinIO`
5. Redis 缂撳瓨鍙€夋帴鍏?
## 14. 鍜?Kitex 鐨勮鎺ユ柟寮?
铏界劧鐜板湪涓嶆媶寰湇鍔★紝浣嗕綘浠庣涓€澶╄捣灏卞彲浠ユ寜棰嗗煙鍐欐竻杈圭晫銆?
鍚庨潰鏈€瀹规槗鎷嗗嚭鍘荤殑妯″潡锛?
1. `article`
2. `auth`
3. `upload`

涔熷氨鏄锛屼綘鐜板湪鐨?service 杈圭晫瑕佸儚杩欐牱锛?
- `ArticleService`
- `AuthService`
- `UploadService`

浠ュ悗浠庘€滆繘绋嬪唴璋冪敤鈥濆彉鎴?鈥淩PC 璋冪敤鈥濇椂锛屾敼鍔ㄤ細灏忓緢澶氥€?
## 15. 鍜?Eino 鐨勮鎺ユ柟寮?
`Eino` 涓嶅簲璇ヨ繘鍏ョ涓€鐗堜富閾捐矾銆?
鏇撮€傚悎绗簩闃舵涔嬪悗鎺ュ叆鐨勫姛鑳斤細

1. 鑷姩鐢熸垚鏂囩珷鎽樿
2. 鑷姩鎺ㄨ崘鏍囩
3. 鏍规嵁鏂囩珷鐢熸垚瀛︿範鍗＄墖
4. 鍋氱珯鍐呴棶绛?
涔熷氨鏄锛宍Eino` 鏇村儚澧炲己灞傦紝涓嶆槸鍩虹璁炬柦灞傘€?
## 16. 浣犵幇鍦ㄦ渶璇ヤ紭鍏堝啓鐨勫嚑涓枃浠?
濡傛灉涓嬩竴姝ユ寮忓紑宸ワ紝鎴戝缓璁紭鍏堜粠杩欎簺鏂囦欢鍏ユ墜锛?
```text
main.go
biz/dal/init.go
biz/dal/db/init.go
biz/router/register.go
biz/handler/home/home_handler.go
biz/handler/blog/blog_handler.go
biz/handler/auth/auth_handler.go
biz/service/article/article_service.go
biz/service/auth/auth_service.go
pkg/configs/
pkg/errno/
```

## 17. 鏈€缁堝缓璁?
涓€鍙ヨ瘽鎬荤粨锛?
**浣犵殑涓汉灏忕珯鍚庣锛屾渶閫傚悎鎸?`tiktok_demo` 鐨勫垎灞傛柟寮忥紝鍋氭垚涓€涓竟鐣屾竻鏅扮殑 Hertz 鍗曚綋椤圭洰锛屽厛瀹屾垚鍐呭绠＄悊闂幆锛屽啀涓?`Kitex` 鍜?`Eino` 鐣欏嚭婕旇繘绌洪棿銆?*

鐜板湪涓嶈杩芥眰鎷嗘湇鍔★紝鍏堟妸锛?
- 鐧诲綍
- 鍙戞枃绔?- 鐪嬫枃绔?- 涓婁紶鍥剧墖
- 鏍囩绠＄悊

杩欎簺鏍稿績閾捐矾鍋氭墡瀹炪€?


