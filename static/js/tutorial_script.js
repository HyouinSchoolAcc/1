// ========== TUTORIAL DIALOGUE SCRIPT ==========
// Based on docs/tutorial_dialogue_script.md

const TUTORIAL_SCRIPT = {
    // Section 0: Re-engage Checkpoint (returning users)
    section0: {
        title: "欢迎回来",
        dialogues: [
            { speaker: "曹操", text: "Welcome back! 欢迎回来！" },
            { speaker: "曹操", text: "让我看看你上次到哪了..." },
            { speaker: "曹操", text: "..." },
            { speaker: "曹操", text: "啊，找到了。我们继续吧？" }
        ],
        choices: [
            { text: "继续上次的进度", action: "resumeFromCheckpoint" },
            { text: "从头开始", action: "goToSection", target: "section1" },
            { 
                text: "直接开写，别废话", 
                action: "skipWithDialogue", 
                dialogue: [{ speaker: "曹操", text: "行行行，年轻人都这么急躁么。祝你好运。" }],
                then: "exitTutorial"
            }
        ]
    },

    // Section 1: Just Joining — Greetings
    section1: {
        title: "初见",
        page: "landing",
        dialogues: [
            { speaker: "曹操", text: "欢饮欢迎！热烈欢迎！" }
        ],
        choices: [
            { text: "这么热情吗，还真挺不像你的", action: "setMode", mode: "fun", then: "section1_fun" },
            { text: "就告诉我怎么办吧（只看介绍）", action: "setMode", mode: "business", then: "section1_business" }
        ]
    },

    section1_fun: {
        title: "初见",
        page: "landing",
        dialogues: [
            { speaker: "曹操", text: "咱缺写手吗，欸嘿。" },
            { speaker: "曹操", text: "说实话我原来挺抗拒做公关的，这不是被数据质量逼得被迫下海了" },
            { speaker: "用户", text: "下海。。。这词用的对吗？" },
            { speaker: "曹操", text: "。。。。少废话" },
            { speaker: "曹操", text: "开门见山吧，我们现阶段的数据需要一人分饰两角，同时扮演「和AI聊天的那个人」和扮演「AI本身」" },
            { speaker: "用户", text: "...这听起来像在写小说" },
            { speaker: "曹操", text: "不不不，小说是旁白叙述" },
            { speaker: "曹操", text: "我们是纯聊天记录" },
            { speaker: "曹操", text: "就像你跟朋友吐槽一样，真实、接地气、有情绪" },
            { speaker: "曹操", text: "对了，刚才那段「发牢骚」的样本看到了吧？" },
            { speaker: "用户", text: "好像有印象" },
            { speaker: "曹操", text: "那个写手现在是我们的王牌" },
            { speaker: "曹操", text: "你也可以" },
            { speaker: "用户", text: "（画大饼是吧）" },
            { speaker: "曹操", text: "不画怎么能让你知道我们的重要性哦吼吼吼" },
            { speaker: "曹操", text: "ahem" }
        ],
        next: "section2"
    },

    section1_business: {
        title: "初见",
        page: "landing",
        dialogues: [
            { speaker: "曹操", text: "好，简洁版：" },
            { speaker: "曹操", text: "我们需要你写AI角色和用户之间的对话数据" },
            { speaker: "曹操", text: "一人分饰两角，模拟真实聊天" },
            { speaker: "曹操", text: "质量过关就能拿稿费" }
        ],
        next: "section2"
    },

    // Section 2: What Do We Do?
    section2: {
        title: "我们做什么？",
        page: "landing",
        dialogues: [
            { speaker: "曹操", text: "正式说一下流程" },
            { speaker: "曹操", text: "第一步：了解文案系统，看角色设定" },
            { speaker: "曹操", text: "第二步：创建一个「用户」（就是和角色聊天的人）" },
            { speaker: "曹操", text: "第三步：模拟一整天的聊天" },
            { speaker: "曹操", text: "第四步：提交审核，拿钱" },
            { speaker: "用户", text: "就这？" },
            { speaker: "曹操", text: "就这" },
            { speaker: "曹操", text: "但「自然对话」比你想的难" },
            { speaker: "曹操", text: "你试过假装不是自己跟人聊三十分钟吗？" },
            { speaker: "用户", text: "没有" },
            { speaker: "曹操", text: "那你马上就要体验了" },
            { speaker: "曹操", text: "这活儿的核心是「代入感」" },
            { speaker: "曹操", text: "你得让审稿人看完觉得「这确实像两个真人在聊天」" },
            { speaker: "用户", text: "感觉不难啊，注水说话不是很容易的一件事" },
            { speaker: "曹操", text: "话是这么说，但能把控好一份感情的涨缩，两个人的情感波动" },
            { speaker: "曹操", text: "其实还真没那么简单" },
            { speaker: "曹操", text: "就像现在AI前两天说着「哇，好像人」" },
            { speaker: "曹操", text: "但过两天就没啥想聊的了，因为两个人谈话时价值观，追求和生活方式很重要" },
            { speaker: "曹操", text: "一开始聊的有的没的其实没有任何对AI的能力提升" },
            { speaker: "曹操", text: "它主要需要的是对以上板块的探测能力，以及探测过后想分享的，想插话的动机。" },
            { speaker: "用户", text: "哦哦。。这。。。。听起来很难啊0.0" },
            { speaker: "曹操", text: "确实比较困难。" },
            { speaker: "曹操", text: "写手得习惯一个人每句话按下人格脑回路格式化然后立马带入对面角色" },
            { speaker: "曹操", text: "所以我们前两个月打算筛选出能做语C的人，之后再考虑单饰，对写手友好一点。" }
        ],
        choices: [
            { text: "你话好多啊，说好的趣味性呢？", action: "goToSection", target: "section2_complaint" },
            { text: "哇，好有趣。有点跃跃欲试了", action: "goToSection", target: "section2_eager" }
        ]
    },

    section2_complaint: {
        title: "我们做什么？",
        dialogues: [
            { speaker: "曹操", text: "你好急啊，干着去投胎么？" },
            { speaker: "用户", text: "卧槽，你信不信老子不干了" },
            { speaker: "曹操", text: "啊啊别急啊，这后面有梗呢，你要听故事吗，还是要我求你" },
            { speaker: "用户", text: "求我" },
            { speaker: "曹操", text: "（嘟嘴）。。。。爱卿还请望步，老夫既所重望，还请再三而思" },
            { speaker: "用户", text: "不真诚啊" },
            { speaker: "曹操", text: "你妈的.....抱歉，对不起，我求你了，看看吧！" },
            { speaker: "用户", text: "这还差不多" }
        ],
        next: "section3"
    },

    section2_eager: {
        title: "我们做什么？",
        dialogues: [
            { speaker: "曹操", text: "是吧是吧，你看这样一来主观能动性有了，趣味性有了，尊重感有了。" },
            { speaker: "曹操", text: "而且还是那种你有问题了能及时关心的，在自己失败时还能扶你一把的人。" },
            { speaker: "曹操", text: "因为关心你所以能骂你让你看清方向。" },
            { speaker: "曹操", text: "来吧来吧，后续还有很多好玩的呢" }
        ],
        next: "section3"
    },

    // Section 3: How Do We Do It?
    section3: {
        title: "如何操作？",
        page: "landing",
        dialogues: [
            { speaker: "曹操", text: "接下来我跟你说说具体怎么操作吧" },
            { speaker: "曹操", text: "信息量可能比较大，但熟能生巧，准备好了吗？" }
        ],
        choices: [
            { text: "开始试炼", action: "goToSection", target: "section3_features" },
            { text: "我先上个厕所", action: "checkpoint", then: "section3_features" }
        ]
    },

    section3_features: {
        title: "如何操作？",
        page: "landing",
        dialogues: [
            { speaker: "曹操", text: "OK！" },
            { speaker: "曹操", text: "普通人可能在Word里写文案，对于我们的高质量创作其实非常麻烦。" },
            { speaker: "曹操", text: "首先，你随时可以回到主页" },
            { speaker: "曹操", text: "点上方的菜单，或者右下角那个圆圆的Home按钮都行" }
        ],
        highlight: { elements: ["hamburgerBtn", "floating-home-btn"], tooltip: "这两个地方都能回到主页" },
        next: "section3_nav"
    },

    section3_nav: {
        title: "如何操作？",
        page: "landing",
        dialogues: [
            { speaker: "曹操", text: "迷路了就按这个，懂？" },
            { speaker: "用户", text: "懂" },
            { speaker: "曹操", text: "现在我们去写作页面看看" }
        ],
        navigate: "/navigation"
    },

    // Section 3 continued on navigation page
    section3_writing: {
        title: "如何操作？",
        page: "navigation",
        dialogues: [
            { speaker: "曹操", text: "这是角色选择页面" },
            { speaker: "曹操", text: "在这里你可以看到其他人写的内容" },
            { speaker: "曹操", text: "写手们每天可以看到自己的文案以及每天更新一次的3天优秀文案" },
            { speaker: "曹操", text: "网站的数据并不公开（要不然会整天被其他大公司数据攻击）" },
            { speaker: "曹操", text: "写对话的时候有几个快捷键" },
            { speaker: "曹操", text: "Enter = 换行（同一个人继续说）" },
            { speaker: "曹操", text: "Shift + Enter = 切换说话人（轮到对方说了）" },
            { speaker: "用户", text: "为什么要这么设置？" },
            { speaker: "曹操", text: "因为大部分时候一个人会连续发好几条消息" },
            { speaker: "曹操", text: "你用微信也是这样吧，连发五六条" },
            { speaker: "曹操", text: "还有一个重要的：时间栏" },
            { speaker: "曹操", text: "你得在对话旁边标注时间" },
            { speaker: "曹操", text: "按下Enter换行之后，光标会自动跳到下一个时间输入框" },
            { speaker: "曹操", text: "小时不变，分钟你自己改" },
            { speaker: "曹操", text: "这样你就能模拟「随着时间推移的聊天节奏」" },
            { speaker: "用户", text: "这也太细了..." },
            { speaker: "曹操", text: "细节决定成败！" },
            { speaker: "曹操", text: "你看那些烂对话数据，AI一句人一句，完全机械" },
            { speaker: "曹操", text: "真实聊天是：AI发三句，人已读不回，过半小时人突然回一句「哈哈哈」" },
            { speaker: "用户", text: "太真实了" }
        ],
        next: "section3_lounge"
    },

    section3_lounge: {
        title: "如何操作？",
        page: "navigation",
        dialogues: [
            { speaker: "曹操", text: "写手论坛（Writers' Lounge）可以跟其他写手和开发组交流" }
        ],
        highlight: "feature-lounge",
        next: "section3_lounge2"
    },

    section3_lounge2: {
        title: "如何操作？",
        page: "navigation",
        dialogues: [
            { speaker: "曹操", text: "有问题就去那边问，我们运营的几个人都在" },
            { speaker: "曹操", text: "顺便一提，因为Discord中国访问不了，这个论坛是我们自己做的替代品" },
            { speaker: "曹操", text: "凑合用吧" },
            { speaker: "用户", text: "这么草台吗" },
            { speaker: "曹操", text: "没人干活啊（惨）" }
        ],
        next: "section4"
    },

    // Section 4: Pick a Character!
    section4: {
        title: "选择角色",
        page: "navigation",
        dialogues: [
            { speaker: "曹操", text: "好了，说说角色" },
            { speaker: "曹操", text: "我们目前有两个角色可以写：" },
            { speaker: "曹操", text: "一、牧濑红莉栖 (Kurisu) —— 《命运石之门》" },
            { speaker: "曹操", text: "住过物理/生物/计算机实验室的人写的" },
            { speaker: "曹操", text: "二、林路 —— 原创角色" },
            { speaker: "曹操", text: "文学系教授，彬彬有礼，有点喜欢聊社会有的没的，对文学要求很深" },
            { speaker: "曹操", text: "额。。。虽然你可以自己加角色，但现在新角色还在定。" },
            { speaker: "曹操", text: "哦对了，选角色的路非常曲折，想不想听？" }
        ],
        choices: [
            { text: "好啊", action: "goToSection", target: "section4_story" },
            { text: "Nah", action: "checkpoint", then: "section4_questions" }
        ]
    },

    section4_story: {
        title: "选择角色",
        dialogues: [
            { speaker: "曹操", text: "OK!" },
            { speaker: "曹操", text: "第一个角色选的时候纯粹创始人喜欢红莉栖，没得选" },
            { speaker: "曹操", text: "老是嚷嚷着什么「天底下一定要把对科研的情书带给全人类」就魔怔了没办法" },
            { speaker: "曹操", text: "第二个角色选的时候是因为我们招了个女编辑" },
            { speaker: "曹操", text: "一开始是想做女性化指导的，但她特别喜欢创作，就让她写了" },
            { speaker: "曹操", text: "我们一开始设定了个银行职员，但分分钟被毙了" },
            { speaker: "曹操", text: "「我靠我为什么会想跟银行职员谈恋爱，你其他角色都那么牛逼」" },
            { speaker: "曹操", text: "啊。。。虽然红莉栖就是个智商在线博士生" },
            { speaker: "曹操", text: "本来就是奔着现实中有原型的能接触到的人去写的" },
            { speaker: "曹操", text: "但！还是为了留住她，写了个生物教授的人设" },
            { speaker: "曹操", text: "虽然但是，后来我们发现她没生物背景，因为她是文科生" },
            { speaker: "曹操", text: "写基金本子，去开会，疏导学生干活那是一点不沾" },
            { speaker: "曹操", text: "但她本人认识非常多的文科生，见识很广，最后就选了个文学教授" }
        ],
        next: "section4_questions"
    },

    section4_questions: {
        title: "选择角色",
        dialogues: [
            { speaker: "曹操", text: "所以还有什么问题么？" }
        ],
        repeatableChoices: true,
        choices: [
            { 
                text: "角色设定是否追随原著？", 
                action: "answer", 
                dialogue: [
                    { speaker: "曹操", text: "直接扔了。" },
                    { speaker: "曹操", text: "我们发现很多设定如果加了现实根本不会去网上找陌生人聊天（有伴侣，在办实际性大事）" },
                    { speaker: "曹操", text: "而且为了追求现实生活中有可能遇见的感觉我们还加了很多2026年的设定进去" }
                ] 
            },
            { 
                text: "我能复刻一方但两个人一起就奇怪，怎么办？", 
                action: "answer", 
                dialogue: [
                    { speaker: "曹操", text: "去写手聊天室里吼一嗓子，试试有人和你对聊" },
                    { speaker: "曹操", text: "亲身经验有意思的" }
                ] 
            },
            { 
                text: "写作时灵感少了，想跟其他人交流可以吗？", 
                action: "answer", 
                dialogue: [
                    { speaker: "曹操", text: "基本每个人写的每一条我们都会看" },
                    { speaker: "曹操", text: "所以一直在我们排行榜上的人我们都推荐大家去聊聊天" },
                    { speaker: "曹操", text: "你也可以在聊天室里找人唠嗑。分享数据也就是一个键的事" }
                ] 
            },
            { text: "没什么问题了", action: "goToSection", target: "section4_end" }
        ]
    },

    section4_end: {
        title: "选择角色",
        dialogues: [
            { speaker: "曹操", text: "如果有的话可以去FAQ那边看看，问题的答案会持续增加的" },
            { speaker: "曹操", text: "那，没问题了你就先看看角色设定吧！" },
            { speaker: "曹操", text: "看完之后回来找我，我还有一些有趣的设定要告诉你" }
        ],
        choices: [
            { text: "去看角色设定", action: "navigate", target: "/descriptions" },
            { text: "我已经看过了，继续", action: "goToSection", target: "section5" }
        ]
    },

    // Section 5: Character Fun Facts
    section5: {
        title: "角色彩蛋",
        dialogues: [
            { speaker: "曹操", text: "欢迎回来！" },
            { speaker: "曹操", text: "看得怎么样？" }
        ],
        choices: [
            { text: "累死老子了，设定太多了", action: "goToSection", target: "section5_tired" },
            { text: "卧槽牛逼", action: "goToSection", target: "section5_excited" }
        ]
    },

    section5_tired: {
        title: "角色彩蛋",
        dialogues: [
            { speaker: "曹操", text: "理论上密度不应该太大啦" },
            { speaker: "曹操", text: "但是有一个地方不行就数据不能用了" }
        ],
        next: "section5_funfacts"
    },

    section5_excited: {
        title: "角色彩蛋",
        dialogues: [
            { speaker: "曹操", text: "hmhmmm~ 是吧是吧" }
        ],
        next: "section5_funfacts"
    },

    section5_funfacts: {
        title: "角色彩蛋",
        dialogues: [
            { speaker: "曹操", text: "对了，说几个有趣的设定" },
            { speaker: "曹操", text: "你可能在角色文档里看不到的那种" },
            { speaker: "曹操", text: "【Kurisu】" },
            { speaker: "曹操", text: "私底下喜欢幻想超能力" },
            { speaker: "曹操", text: "会在没人的时候想「如果我会念动力就好了」然后偷偷比划" },
            { speaker: "曹操", text: "被发现的话会超级尴尬" },
            { speaker: "曹操", text: "【林路】" },
            { speaker: "曹操", text: "这位教授有知识分子的通病" },
            { speaker: "曹操", text: "对很多物理世界的东西「想当然」" },
            { speaker: "曹操", text: "比如不理解为什么修电脑不能靠「重启一下」" },
            { speaker: "曹操", text: "或者困惑为什么学生觉得莎士比亚难读" },
            { speaker: "曹操", text: "他不是傲慢，是真的不理解" },
            { speaker: "曹操", text: "这种「知识诅咒」会让他有时候说话让人想笑" },
            { speaker: "用户", text: "听起来像我的某些教授" },
            { speaker: "曹操", text: "艺术来源于生活" },
            { speaker: "曹操", text: "哦对了，以后如果想创建自己的角色" },
            { speaker: "曹操", text: "新角色页面有单独的教程，到时候再看" },
            { speaker: "曹操", text: "现在先把这两个角色搞熟再说" }
        ],
        next: "section7"
    },

    // Section 7: Take a Quiz
    section7: {
        title: "小测验",
        dialogues: [
            { speaker: "曹操", text: "看完设定之后..." },
            { speaker: "曹操", text: "我们有个小测验" },
            { speaker: "曹操", text: "考你几个关于角色的问题" },
            { speaker: "曹操", text: "及格了才能正式开写" },
            { speaker: "用户", text: "还要考试？！" },
            { speaker: "曹操", text: "放心，开卷的" },
            { speaker: "曹操", text: "你可以一边看设定一边答题" },
            { speaker: "曹操", text: "主要是确保你真的看了" },
            { speaker: "曹操", text: "不然对话质量上不去，双方都浪费时间" },
            { speaker: "用户", text: "好吧，合理" }
        ],
        choices: [
            { text: "我准备好了，开始考试吧", action: "navigate", target: "/quiz" },
            { text: "让我再看看设定", action: "navigate", target: "/descriptions" },
            { 
                text: "我不需要考试，直接让我写", 
                action: "skipWithDialogue",
                dialogue: [{ speaker: "曹操", text: "自信满满啊，但规矩还是要守的，老老实实考一下吧" }],
                then: "goToQuiz"
            }
        ]
    },

    // Section 8: Ending
    section8: {
        title: "结语",
        dialogues: [
            { speaker: "曹操", text: "好了！流程大概就是这样" },
            { speaker: "曹操", text: "总结一下：" },
            { speaker: "曹操", text: "1️⃣ 选角色（或创建自己的）" },
            { speaker: "曹操", text: "2️⃣ 看设定（日程表很重要）" },
            { speaker: "曹操", text: "3️⃣ 创建「用户」身份" },
            { speaker: "曹操", text: "4️⃣ 写一整天的对话" },
            { speaker: "曹操", text: "5️⃣ 提交审核 → 拿钱！" },
            { speaker: "曹操", text: "有问题就去Writers' Lounge找我们" },
            { speaker: "曹操", text: "快捷键别忘了：Enter换行，Shift+Enter切人" },
            { speaker: "曹操", text: "还有那个时间栏...算了你用用就知道了" },
            { speaker: "用户", text: "懂了" },
            { speaker: "曹操", text: "那么..." },
            { speaker: "曹操", text: "欢迎来到奇点之门，研究员" },
            { speaker: "曹操", text: "去选个角色开始吧" },
            { speaker: "曹操", text: "我在这里随时待命" },
            { speaker: "曹操", text: "虽然穿着这身莫名其妙的衣服..." },
            { speaker: "用户", text: "（笑）" },
            { speaker: "曹操", text: "不准笑！！" }
        ],
        ending: true
    },

    // Easter eggs
    quizFailure: [
        { speaker: "曹操", text: "..." },
        { speaker: "曹操", text: "你真的看设定了吗？" },
        { speaker: "曹操", text: "回去再看看吧，这几道题都是送分题" },
        { speaker: "曹操", text: "我当年赤壁之战输了也没这么难看" },
        { speaker: "用户", text: "你还好意思提赤壁" },
        { speaker: "曹操", text: "...我收回刚才的话，你可以再考一次" }
    ]
};

// Export for use
if (typeof module !== 'undefined' && module.exports) {
    module.exports = TUTORIAL_SCRIPT;
}
