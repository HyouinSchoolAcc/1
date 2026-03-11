# 📜 Divergence 2% Tutorial Dialogue Script
## *"Cao Cao's Guide to Writing Characters"*
### Version 1.0 — January 2026

---

## 📋 Document Overview

**Purpose:** Gal-game style interactive tutorial for onboarding new writers  
**Guide Character:** 曹操 (Cao Cao)  
**Estimated Duration:** ~10 minutes  
**Tone:** Fun, irreverent, fourth-wall-breaking humor with helpful information

### Section Structure
1. Re-engage Checkpoint (returning users)
2. Just Joining — Greetings
3. What Do We Do?
4. How Do We Do It? (Features)
5. Pick a Character!
6. Character Fun Facts
7. Take a Quiz
8. Ending
9. Easter Eggs (Optional)

*Note: Creating your own character has a separate tutorial on the /newcharacter page.*

### Format Legend
- `曹操:` — Guide character speaks
- `用户:` — Player/user speaks
- `→ CHOICE:` — Branching choice point
- `[ HIGHLIGHT: xxx ]` — UI element gets highlighted
- `[ CHECKPOINT ]` — Save/resume trigger
- `→ [Action]` — Navigation or system action

---



## SECTION 0: Re-engage Checkpoint

*Triggers when user previously exited or returns to tutorial*

```
曹操: Welcome back! 欢迎回来！
曹操: 让我看看你上次到哪了...
曹操: ...
曹操: 啊，找到了。我们继续吧？

→ CHOICE:
  - "继续上次的进度" 
      → [Jump to saved section]
  
  - "从头开始" 
      → [Jump to Section 1]
  
  - "直接开写，别废话" 
      → 曹操: "行行行，年轻人都这么急躁么。祝你好运。" 
      → [End tutorial, navigate to /navigation]
```

---

## SECTION 1: Just Joining — Greetings

*First impressions and setting expectations*

```
曹操: 欢饮欢迎！热烈欢迎！
→ CHOICE:
  - "这么热情吗，还真挺不像你的" (进入趣味性入职指导）
      → [Select fun mode]
  
  - "就告诉我怎么办吧（只看介绍）" 
      → [Select businessman mode]
曹操: 咱缺写手吗，欸嘿。
曹操: 说实话我原来挺抗拒做公关的，这不是被数据质量逼得被迫下海了
用户: 下海。。。这词用的对吗？
曹操: 。。。。少废话
曹操: 开门见山吧，我们现阶段的数据需要一人分饰两角,同时扮演"和AI聊天的那个人"和扮演"AI本身"

用户: ...这听起来像在写小说

曹操: 不不不，小说是旁白叙述
曹操: 我们是纯聊天记录
曹操: 就像你跟朋友吐槽一样，真实、接地气、有情绪
曹操: 对了，刚才那段"发牢骚"的样本看到了吧？

用户: 好像有影响

曹操: 那个写手现在是我们的王牌
曹操: 你也可以

用户: （画大饼是吧）
曹操: 不画怎么你怎么能知道我们的重要性哦吼吼哦吼
曹操: ahem


## SECTION 2: What Do We Do?

*Explaining the core workflow*

```
曹操: 正式说一下流程
曹操: 第一步：了解文案系统，看角色设定
曹操: 第二步：创建一个"用户"（就是和角色聊天的人）
曹操: 第三步：模拟一整天的聊天
曹操: 第四步：提交审核，拿钱

用户: 就这？

曹操: 就这
曹操: 但"自然对话"比你想的难
曹操: 你试过假装不是自己跟人聊三十分钟吗？

用户: 没有

曹操: 那你马上就要体验了
曹操: 这活儿的核心是"代入感"
曹操: 你得让审稿人看完觉得"这确实像两个真人在聊天"

用户: 感觉不难啊，注水说话不是很容易的一件事
 
曹操: 话是这么说，但能把控好一份感情的涨缩，两个人的情感波动
曹操: 其实还真没那么简单

曹操: 就像现在ai前两天说着“哇，好像人”
曹操: 但过两天就没啥想聊的了，因为两个人谈话时价值观，追求和生活方式很重要
曹操: 一开始聊的有的没的其实没有任何对ai的能力提升，它主要需要的是对以上板块的探测能力，以及探测过后想分享的，想插话的动机。

用户：哦哦。。这。。。。听起来很难啊0.0
曹操：确实比较困难。
曹操：写手得习惯一个人每句话按下人格脑回路格式化然后立马带入对面角色
曹操：所以我们前两个月打算筛选出能做语c的人，之后再考虑单饰，对写手友好一点。

→ CHOICE:
  - "你话好多啊，说好的趣味性呢？" 
      → 你好急啊，干着去投胎么？
      -》用户：“卧槽，你信不信老子不干了”
      → 啊啊别急啊，这后面有梗呢，你要听故事吗，还是要我求你
      -》用户：“求我”
      → （嘟嘴）。。。。爱卿还请望步，老夫既所重望，还请再三而思
      -》用户：“不真诚啊”
      → 你妈的.....抱歉，对不起，我求你了，看看吧！
      -》用户：“这还差不多”

  - "哇，好有趣。有点跃跃欲试了" 
      → 是吧是吧，你看这样一来主观能动性有了，趣味性有了，尊重感有了。而且还是那种你有问题了能及时关心的，在自己失败时还能扶你一把的人。因为关心你所以能骂你让你看清方向。
      → 来吧来吧，后续还有很多好玩的呢


---

## SECTION 3: How Do We Do It?
[跳转到写手网站]
*UI features and keyboard shortcuts*

```
曹操: 接下来我跟你说说具体怎么操作吧,信息量可能比较大，但熟能生巧，准备好了吗?

→ CHOICE:
  - "开始试炼" 
      → OK！
  - "我先上个厕所"
      → 你先上个厕所 [checkpoint save]

曹操: 普通人可能在word里写文案，对于我们的高质量创作其实非常麻烦。
曹操: 首先，你随时可以回到主页
曹操: 点上方的菜单，或者右下角那个圆圆的Home按钮都行

[ HIGHLIGHT: top-header + home-button ]

曹操: 迷路了就按这个，懂？

用户: 懂
[跳转到红莉栖页面]
曹操: 现在开始就是如何写文案了

[highlight: dialogue window]
曹操: 在这里你可以看到其他人写的内容
曹操: 写手们每天可以看到自己的文案以及每天更新一次的3天优秀文案
曹操: 网站的数据并不公开（要不然会整天被其他大公司数据攻击)

曹操: 写对话的时候有几个快捷键
曹操: Enter = 换行（同一个人继续说）
曹操: Shift + Enter = 切换说话人（轮到对方说了）

用户: 为什么要这么设置？

曹操: 因为大部分时候一个人会连续发好几条消息
曹操: 你用微信也是这样吧，连发五六条
曹操: 还有一个重要的：时间栏    [highlight: time]
曹操: 你得在对话旁边标注时间
曹操: 按下Enter换行之后，光标会自动跳到下一个时间输入框
曹操: 小时不变，分钟你自己改
曹操: 这样你就能模拟"随着时间推移的聊天节奏"

用户: 这也太细了...

曹操: 细节决定成败！
曹操: 你看那些烂对话数据，AI一句人一句，完全机械
曹操: 真实聊天是：AI发三句，人已读不回，过半小时人突然回一句"哈哈哈"

用户: 太真实了

[jump to writers lounge]
曹操: 写手论坛（Writers' Lounge）可以跟其他写手和开发组交流

[ HIGHLIGHT: feature-lounge ]

曹操: 有问题就去那边问，我们运营的几个人都在
曹操: 顺便一提，因为没有Discord中国访问，这个论坛是我们自己做的替代品
曹操: 凑合用吧

用户: 这么草台吗

曹操: 没人干活啊（惨）
```

---

## SECTION 4: Pick a Character!

*Introducing the three characters and the selection rationale*

```
曹操: 好了，说说角色
曹操: 我们目前有三个角色可以写：

[ Jump to character selection area ]

曹操: 一、牧濑红莉栖 (Kurisu) —— 《命运石之门》
曹操: 住过物理/生物/计算机实验室的人写的

曹操: 二、卡夫卡 (Kafka) —— 《崩坏：星穹铁道》
曹操: 妈妈型人格，一般比较有恋爱经验的人才能驾驭

曹操: 三、林路 —— 原创角色
曹操: 文学系教授，彬彬有礼，有点喜欢聊社会有的没的，对文学要求很深

曹操: 额。。。虽然你可以自己加角色，但现在新角色还在定。
曹操: 哦对了，选角色的路非常曲折，想不想听？

→ CHOICE:
  - "好啊" 
      → OK!
      → 第一个角色选的时候纯粹创始人喜欢红莉栖，没得选老是嚷嚷着什么“天底下一定要把对科研的情书带给全人类”就魔怔了没办法
      → 第二个角色选的时候想蹭个ip热度的，但不知道蹭谁的好。因为这本质上是免费帮别人打广告，如果像迪士尼那帮子巴不得自己ip被人忘掉的死脑筋就没办法了。最后想的是米哈游的文化非常开朗，芙莉连的团队特别喜欢整活炒ip，但再三思考之后我们发现芙莉连压根不说话，所以选了铁道2023年度大佬卡夫卡，虽然现在团队里有想改成远坂凌或加藤惠，因为贴实好写，但目前还没来的及改
      → 第三个角色选的时候是因为我们招了个女编辑，一开始是想做女性化指导的，但她特别喜欢创作，就让她写了。我们一开始设定了个银行职员，但分分钟被毙了。”我靠我为什么会想跟银行职员谈恋爱，你其他角色都那么牛逼“
      → 啊。。。虽然红莉栖就是个智商在线博士生，卡夫卡就是一个30岁见过世面追刺激的富婆，本来就是奔着现实中有原型的能接触到的人去写的，但！还是为了留住她，写了个生物教授的人设。
      → 虽然但是，后来我们发现她没生物背景，因为她是文科生。写基金本子，去开会，疏导学生干活那是一点不沾。但她本人认识非常多的文科生，见识很广，最后就选了个文学教授。

  - "Nah"
      → Aww, ok 。。。。 [checkpoint save]
曹操: 所以还有什么问题么？
→ CHOICE（can be repeated): 
  - "角色设定是否追随原著？原著里的设定怎么办？" 
      → 直接扔了。我们发现很多设定如果加了现实根本不会去网上找陌生人聊天（有伴侣，在办实际性大事），而且为了追求现实生活中有可能遇见的感觉我们还加了很多2026年到设定里去
  - "我非常能够复刻其中一方的感觉，但是如果我要复刻两个人的感觉就非常奇怪了，但我又不想自己写的全部打水漂，怎么办？ " 
      → 去写手聊天室里吼一嗓子，试试有人和你对聊，亲身经验有意思的。   
  - "我写作时灵感少了，想跟其他人对话，但看不到他们的作品，我想跟和我一样的人一起社交可以吗？" 
      → 基本每个人写的每一条我们都会看，所以一直在我们排行榜上的人我们都推荐大家去聊聊天。你也可以在聊天室里找人唠嗑。分享数据也就是一个键的事。
  - "没什么问题了" 
      → 如果有的话可以去FAQ地方看看，问题的答案会持续增加的。
曹操：那，没问题了你就先看看角色设定吧！
[at this point, display on bottom of the page: 看完啦！]

```

---

## SECTION 5: Character Fun Facts

*Behind-the-scenes character details for engagement*

```
曹操: 欢迎回来！
曹操: 看得怎么样？

→ CHOICE:
  - "累死老子了，你这么多设定我真的每天聊天会用的到吗" 
      → 理论上密度不应该太大啦，但是有一个地方不行就数据不能用了
  - "卧槽牛逼"
      → hmhmmm~ 是吧是吧

曹操: 对了，说几个有趣的设定
曹操: 你可能在角色文档里看不到的那种

曹操: 【Kurisu】
曹操: 私底下喜欢幻想超能力
曹操: 会在没人的时候想"如果我会念动力就好了"然后偷偷比划
曹操: 被发现的话会超级尴尬

曹操: 【Kafka】
曹操: 这位姐姐有点...怎么说呢...
曹操: 跟踪狂属性

用户: ???

曹操: 每次都会权衡"值不值得"
曹操: 是比较自信，机动性很强，高效得一位

用户: 太刺激了

曹操: 写她的时候记得这个特点，会让对话更有层次

曹操: 【林路】
曹操: 这位教授有知识分子的通病
曹操: 对很多物理世界的东西"想当然"
曹操: 比如不理解为什么修电脑不能靠"重启一下"
曹操: 或者困惑为什么学生觉得莎士比亚难读
曹操: 他不是傲慢，是真的不理解
曹操: 这种"知识诅咒"会让他有时候说话让人想笑

用户: 听起来像我的某些教授

曹操: 艺术来源于生活

曹操: 哦对了，以后如果想创建自己的角色
曹操: 新角色页面有单独的教程，到时候再看
曹操: 现在先把这三个角色搞熟再说
```

---

## SECTION 7: Take a Quiz

*Final preparation before writing*

```
曹操: 看完设定之后...
曹操: 我们有个小测验
曹操: 考你几个关于角色的问题
曹操: 及格了才能正式开写

用户: 还要考试？！

曹操: 放心，开卷的
曹操: 你可以一边看设定一边答题
曹操: 主要是确保你真的看了
曹操: 不然对话质量上不去，双方都浪费时间

用户: 好吧，合理

→ CHOICE:
  - "我准备好了，开始考试吧" 
      → 曹操: "好样的！跳转到测验页面..." 
      → [Navigate to quiz]
  
  - "让我再看看设定" 
      → 曹操: "明智的选择，不着急" 
      → [Navigate to /descriptions]
  
  - "我不需要考试，直接让我写" 
      → 曹操: "自信满满啊，但规矩还是要守的，老老实实考一下吧" 
      → [Navigate to quiz]
```

---

## SECTION 8: Ending

*Tutorial wrap-up and send-off*

```
曹操: 好了！流程大概就是这样
曹操: 总结一下：
曹操: 1️⃣ 选角色（或创建自己的）
曹操: 2️⃣ 看设定（日程表很重要）
曹操: 3️⃣ 创建"用户"身份
曹操: 4️⃣ 写一整天的对话
曹操: 5️⃣ 提交审核 → 拿钱！
曹操: 有问题就去Writers' Lounge找我们
曹操: 快捷键别忘了：Enter换行，Shift+Enter切人
曹操: 还有那个时间栏...算了你用用就知道了

用户: 懂了

曹操: 那么...
曹操: 欢迎来到奇点之门，研究员

[ HIGHLIGHT: /navigation ]

曹操: 去选个角色开始吧
曹操: 我在这里随时待命
曹操: 虽然穿着这身莫名其妙的衣服...

用户: （笑）

曹操: 不准笑！！

[ END OF TUTORIAL ]
```

---

## SECTION 9: Easter Eggs (Optional)

### 9.1 Avatar Interaction
*If user clicks Cao Cao's avatar multiple times:*

```
[Click 1]
曹操: 你戳我干啥

[Click 2]
曹操: ...

[Click 3]
曹操: 再戳我把你踢出去

[Click 4]
曹操: ...

[Click 5]
曹操: 行行行知道你无聊了，继续教程吧
```

### 9.2 Idle Detection
*If user doesn't interact for 30+ seconds:*

```
曹操: ...喂
曹操: 人呢？
曹操: 该不会去上厕所了吧
曹操: 我等着...
```

### 9.3 Quiz Failure Response
*If user fails the character quiz:*

```
曹操: ...
曹操: 你真的看设定了吗？
曹操: 回去再看看吧，这几道题都是送分题
曹操: 我当年赤壁之战输了也没这么难看

用户: 你还好意思提赤壁

曹操: ...我收回刚才的话，你可以再考一次
```

---

## 📝 Implementation Notes

### Key Features to Highlight in Tutorial
| Feature | Where to Mention | UI Element ID |
|---------|------------------|---------------|
| Home button / Top menu | Section 3 | `top-header`, `home-button` |
| 3 Characters + New Character option | Section 4 | Character cards |
| Creating new "user" in /kurisu endpoints | Section 2 | N/A (verbal) |
| Writing perspective explanation | Section 1 | N/A (verbal) |
| Character creation (brief mention) | Section 6 | N/A (separate tutorial) |
| Shortcut keys (Enter, Shift+Enter) | Section 3 | N/A (verbal) |
| Time association / timestamp bar | Section 3 | Time input fields |
| Writers' Lounge | Section 3 | `feature-lounge` |

### Story Elements Incorporated
| Story | Section | Notes |
|-------|---------|-------|
| How we chose the 3 characters | Section 4 | IP recognition + single status |
| Why are these characters single | Section 4 | Dating sim use case explanation |
| Gender-swap writing recommendation | Section 4 | Avoids self-projection |
| Kurisu's superpower fantasies | Section 5 | Fun fact, embarrassed if caught |
| Kafka's stalking tendencies | Section 5 | Internal conflict angle |
| Lin Lu's intellect complex | Section 5 | "Curse of knowledge" framing |

### Branching Structure Summary
```
START
  │
  ├─[Returning User]─→ CHECKPOINT
  │                      ├─ Continue → (saved section)
  │                      ├─ Restart → Section 1
  │                      └─ Skip → /navigation
  │
  └─[New User]─→ Section 1
                   │
                   ↓
              Section 2-6 (linear)
                   │
                   ↓
              Section 7 (Quiz)
                   ├─ Start quiz → Quiz
                   ├─ Read settings → /descriptions
                   └─ Skip (denied) → Quiz
                   │
                   ↓
              Section 8 → END
```

---

## 🎮 Tone Guidelines for Editor

### Do's ✅
- Fourth-wall breaking humor
- Self-deprecating jokes (Cao Cao's outfit, Chibi reference)
- Casual internet slang (画大饼, .jpg, etc.)
- Mix of Chinese and occasional English
- Quick back-and-forth banter
- Genuine helpfulness underneath the jokes

### Don'ts ❌
- Don't be condescending
- Don't make the tutorial feel like a lecture
- Don't skip the fun parts for efficiency
- Don't break character (Cao Cao should always feel like Cao Cao)
- Don't make jokes at the expense of users

### Voice Reference
Cao Cao in this tutorial is:
- Slightly sarcastic but ultimately supportive
- Annoyed by his own costume
- Knowledgeable but not preachy
- Quick to banter with the user
- Secretly proud of the project

---

*Document prepared for editorial review*  
*Last updated: January 2026*
