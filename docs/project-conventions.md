# 涓汉灏忕珯椤圭洰瑙勮寖

## 1. 褰撳墠鐩爣

褰撳墠闃舵鐨勫敮涓€鐩爣鏄細

**瀹屾垚涓€涓熀浜?Hertz 鐨勫崟浣撲釜浜哄皬绔欙紝骞跺湪寮€鍙戣繃绋嬩腑鎸佺画璁板綍寮€鍙戞棩蹇楋紝淇濊瘉鍚庣画鏄撲簬 review銆佹槗浜庡洖椤俱€佹槗浜庢紨杩涖€?*

杩欎釜椤圭洰褰撳墠涓嶈拷姹備竴娆℃€у仛鍏紝鑰屾槸鎸夐樁娈垫帹杩涳細

1. 鍏堝畬鎴愬崟浣撻棴鐜?2. 鍐嶈ˉ瀹岀鐞嗚兘鍔?3. 鏈€鍚庝负 `Kitex` 鍜?`Eino` 棰勭暀婕旇繘绌洪棿

## 2. 寮€鍙戝師鍒?
### 2.1 椤圭洰褰㈡€?
- 褰撳墠鏄崟浣撻」鐩?- 鏍稿績鎺ュ彛閲囩敤 `thrift` 鍋?IDL
- `handler/service/dal` 杈圭晫蹇呴』娓呮櫚
- 椤甸潰娓叉煋鍜屽唴瀹规帴鍙ｅ彲浠ュ叡瀛?
### 2.2 涓嶅仛鐨勪簨

- 涓嶅仛杩囨棭寰湇鍔℃媶鍒?- 涓嶅仛杩囬噸鍚庡彴绯荤粺
- 涓嶆妸椤甸潰 view model 鍏ㄩ儴 IDL 鍖?- 涓嶆妸鎵€鏈夊姛鑳戒竴娆℃€у杩涚涓€鐗?
### 2.3 蹇呴』淇濇寔鐨勭害鏉?
- 鏍稿績鎺ュ彛浼樺厛璧?IDL
- 鐧诲綍閫昏緫鑷繁鍐?handler锛屼笉浣跨敤搴撳唴缃櫥褰?handler
- 褰撳墠涓荤櫥褰曟€佹柟妗堜负 Session
- 瀵嗙爜蹇呴』璧板搱甯屽瓨鍌ㄤ笌鏍￠獙锛屼笉鍋氭槑鏂囧瘑鐮佹寔涔呭寲
- 姣忓畬鎴愪竴涓樁娈甸兘琛ュ紑鍙戞棩蹇?- 姣忎釜鍔熻兘闂幆瀹屾垚鍚庡啀鎻愪氦鍚堝苟
- 鐢ㄦ埛璐熻矗浠ｇ爜寮€鍙戜笌鏈湴娴嬭瘯锛屾棩蹇楁暣鐞嗕笌 Git 娴佺▼鐢?Codex 鍗忓姪瀹屾垚

## 3. 鐩綍瑙勮寖

寤鸿鏈€缁堢洰褰曠粨鏋勶細

```text
personal_site/
鈹溾攢鈹€ idl/
鈹?  鈹溾攢鈹€ article.thrift
鈹?  鈹溾攢鈹€ auth.thrift
鈹?  鈹溾攢鈹€ upload.thrift
鈹?  鈹溾攢鈹€ tag.thrift
鈹?  鈹溾攢鈹€ category.thrift
鈹?  鈹斺攢鈹€ rbac.thrift
鈹溾攢鈹€ biz/
鈹?  鈹溾攢鈹€ handler/
鈹?  鈹溾攢鈹€ service/
鈹?  鈹溾攢鈹€ dal/
鈹?  鈹?  鈹斺攢鈹€ db/
鈹?  鈹溾攢鈹€ model/
鈹?  鈹溾攢鈹€ router/
鈹?  鈹斺攢鈹€ mw/
鈹溾攢鈹€ pkg/
鈹?  鈹溾攢鈹€ configs/
鈹?  鈹溾攢鈹€ constants/
鈹?  鈹溾攢鈹€ errno/
鈹?  鈹溾攢鈹€ response/
鈹?  鈹斺攢鈹€ utils/
鈹溾攢鈹€ templates/
鈹溾攢鈹€ static/
鈹溾攢鈹€ docs/
鈹?  鈹溾攢鈹€ devlog/
鈹?  鈹溾攢鈹€ personal-site-ui-spec.md
鈹?  鈹溾攢鈹€ personal-site-backend-plan.md
鈹?  鈹溾攢鈹€ personal-site-idl-plan.md
鈹?  鈹溾攢鈹€ personal-site-rbac-plan.md
鈹?  鈹斺攢鈹€ project-conventions.md
鈹斺攢鈹€ main.go
```

## 4. Git 宸ヤ綔娴?
褰撳墠浠撳簱鐩爣锛?
- GitHub 浠撳簱锛歚Loe1210/personal-site`
- 浠撳簱鍦板潃锛歚https://github.com/Loe1210/personal-site.git`

鍚庣画缁熶竴浣跨敤 Git 鎺ㄨ繘锛屼笉鍐嶅仠鐣欏湪鈥滄湰鍦板彧鏀逛笉鎻愪氦鈥濈殑鐘舵€併€?
### 4.1 鏍囧噯闃舵娴佺▼

姣忓畬鎴愪竴涓樁娈碉紝鍥哄畾鎵ц浠ヤ笅娴佺▼锛?
1. 鏈湴瀹屾垚涓€涓樁娈?2. 鏇存柊瀵瑰簲 `devlog`
3. 鎻愪氦 `commit`
4. 鎺ㄩ€佸埌鍔熻兘鍒嗘敮
5. 鍚堝苟鍒颁富鍒嗘敮

杩欐潯娴佺▼鏄悗缁粯璁ゅ伐浣滄柟寮忋€?
### 4.2 鍒嗘敮绛栫暐

涓嶆寜鈥滄ā鍧椻€濆垎鏀紝鑰屾寜鈥滈樁娈?/ 鍔熻兘闂幆鈥濆垎鏀€?
鎺ㄨ崘鍒嗘敮锛?
- `main`
- `feat/project-bootstrap`
- `feat/idl-article-auth-upload`
- `feat/auth-session-users-rbac`
- `feat/article-crud`
- `feat/upload-image`
- `feat/tag-category`
- `feat/rbac-minimal`
- `feat/blog-pages`
- `feat/release-prep`

### 4.3 涓轰粈涔堜笉鎸夋ā鍧楀垎鏀?
鍥犱负涓€涓湡瀹炲姛鑳介€氬父浼氬悓鏃朵慨鏀癸細

- `idl`
- `handler`
- `service`
- `dal`
- `router`
- 椤甸潰
- 鏂囨。

濡傛灉鎸夋ā鍧楀垎鏀紝鍚庨潰鍚堝苟浼氬緢涔便€?
### 4.4 鍚堝苟鍘熷垯

- 涓€涓垎鏀彧鎵胯浇涓€涓槑纭樁娈电洰鏍?- 闃舵鏈畬鎴愬墠锛屼笉鎬ョ潃鍚堝苟
- 闃舵瀹屾垚鏃跺繀椤诲悓姝ユ洿鏂板紑鍙戞棩蹇?- 鍚堝苟鍓嶈嚦灏戜繚璇佹湰鍦拌嚜娴嬮€氳繃
- 鍚堝苟鍚庝富鍒嗘敮淇濇寔鍙户缁紑鍙戠姸鎬?
## 5. 鎻愪氦瑙勮寖

鎺ㄨ崘 commit 鏍煎紡锛?
```text
feat: add article thrift
feat: implement session login with users table
feat: add article admin crud in memory
refactor: simplify article service flow
fix: align thrift field names in article handlers
docs: update phase 02 devlog
```

## 6. 寮€鍙戞棩蹇楄鑼?
寮€鍙戞棩蹇楃粺涓€鏀惧湪锛?
```text
docs/devlog/
```

鏂囦欢鍛藉悕锛?
```text
phase-01.md
phase-02.md
phase-03.md
```

姣忕瘒寮€鍙戞棩蹇楀缓璁浐瀹氱粨鏋勶細

```md
# Phase X

## 鐩爣

## 瀹屾垚鍐呭

## 璁捐鍐崇瓥

## 閬囧埌鐨勯棶棰?
## 褰撳墠缁撴灉

## 涓嬩竴姝?```

## 7. IDL 瑙勮寖

### 7.1 姣忎釜棰嗗煙涓€涓?thrift

绗竴鎵?thrift锛?
- `article.thrift`
- `auth.thrift`
- `upload.thrift`

鍚庣画 thrift锛?
- `tag.thrift`
- `category.thrift`
- `rbac.thrift`

### 7.2 姣忎釜 thrift 鐙珛 package

渚嬪锛?
```thrift
namespace go article
```

### 7.3 涓嶅悓 thrift 涓嶇敓鎴愬埌鍚屼竴 Go 鍖?
杩欐槸涓轰簡閬垮厤浜掔浉瑕嗙洊鍜屽懡鍚嶅啿绐併€?
## 8. 绗竴闃舵鑼冨洿

绗竴闃舵鍙仛涓変欢浜嬶細

1. 纭畾鐩綍鍜岃鑼?2. 璧疯崏绗竴鐗?IDL
3. 寤虹珛绗竴绡囧紑鍙戞棩蹇?
涓嶈鍦ㄨ繖涓€闃舵鐩存帴鎵╁睍涓氬姟閫昏緫銆?
## 9. 涓嬩竴姝ラ『搴?
褰撳墠椤圭洰宸茬粡瀹屾垚锛?
- Phase 01锛氱洰褰曘€佽鑼冨拰鍒濈増 IDL
- Phase 02-05锛氶」鐩鏋躲€佹枃绔犱富閾捐矾銆佹暟鎹簱銆丼ession 鐧诲綍涓?RBAC 鏈€灏忛棴鐜?- Phase 06锛氭枃绔犲唴瀹瑰煙瑙勮寖鍖栵紝鏂囩珷鏍囩鍏崇郴鍒囨崲鍒?`article_tags`

鍚庣画寤鸿鎸夎繖涓『搴忕户缁仛锛?
1. 瀹炵幇鏈湴鏂囦欢瀛樺偍鐗?`upload` 妯″潡
2. 鎵撻€氭枃绔犲皝闈㈠浘涓婁紶涓庤闂?URL
3. 鍐嶈繘鍏ュ墠鍙伴椤点€丅log 鍒楄〃鍜岃鎯呴〉鑱旇皟
4. 鏈€鍚庤ˉ榻愭洿瀹屾暣鐨?RBAC 绠＄悊鎺ュ彛涓庡悗鍙板寮鸿兘鍔?
## 10. 褰撳墠闃舵鐨勬垚鍔熸爣鍑?
褰撳墠闃舵瀹屾垚鍚庯紝搴斿綋鍏峰锛?
- 鐩綍瑙勮寖鏄庣‘
- 寮€鍙戞棩蹇楄鑼冩槑纭?- 鏍稿績 IDL 宸插瓨鍦ㄥ苟鎸夐鍩熸媶鍒?- 鏂囩珷鍩熸暟鎹粨鏋勫凡缁忕ǔ瀹氬埌鍙户缁壙鎺ヤ笂浼犲拰鍓嶅彴鑱旇皟
- 鍚庣画寮€鍙戝彲浠ユ寜鏂囨。鍜?thrift 鐩存帴杩涘叆瀹炵幇
## 11. 上传演进规划

当前阶段的上传模块目标不是一次性做成生产级大文件系统，而是先完成**稳定版小文件上传**，满足文章封面图和站点资源的实际需求。

### 11.1 当前阶段要做的上传能力

- 支持图片上传
- 单文件大小限制，避免超大文件直接压垮服务
- MIME 白名单校验，只放行允许的图片类型
- 使用流式保存到磁盘，不把整文件一次性读入内存
- MySQL 保存上传元数据
- 返回可直接访问的文件 URL
- 失败时直接返回错误，不引入复杂状态机

### 11.2 当前阶段暂不做的能力

下面这些能力明确放到后续阶段，不在当前版本实现：

- 分片上传
- 断点续传
- 临时分块合并
- 秒传 / Hash 去重
- Bloom Filter 预判无效 file_id
- 上传状态补偿
- 对象存储接入
- 大文件专用网关限流和异步清理

### 11.3 后续什么时候再做

当项目出现以下场景时，再进入上传增强阶段：

- 上传文件从封面图扩展到大体积资源文件
- 单文件大小显著提升到几十 MB 甚至更大
- 上传失败重试和断点恢复成为真实需求
- 本地磁盘存储不再满足容量或部署要求

### 11.4 后续增强阶段建议顺序

1. 先补文件 Hash 与重复文件检测
2. 再补上传状态字段，如 `pending / success / failed`
3. 之后考虑对象存储
4. 最后再做分片上传、断点续传和状态补偿
